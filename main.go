package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"thinknalyze/certs"
)

func init() {
	// Ensure certificates exist
	ensureCerts()
}

// ensureCerts generates RSA keys if they don't exist
func ensureCerts() {
	certDirs := []struct {
		name    string
		private string
		public  string
	}{
		{
			name:    "gateway",
			private: "gateway/certs/private.pem",
			public:  "gateway/certs/public.pem",
		},
		{
			name:    "users",
			private: "users/certs/private.pem",
			public:  "users/certs/public.pem",
		},
		{
			name:    "account",
			private: "account/certs/private.pem",
			public:  "account/certs/public.pem",
		},
		{
			name:    "subscription",
			private: "subscription/certs/private.pem",
			public:  "subscription/certs/public.pem",
		},
	}

	for _, cd := range certDirs {
		// Skip if both files exist
		if _, errPriv := os.Stat(cd.private); errPriv == nil {
			if _, errPub := os.Stat(cd.public); errPub == nil {
				fmt.Printf("✅ %s certificates already exist\n", cd.name)
				continue
			}
		}

		fmt.Printf("\n🔐 Generating certificates for %s...\n", cd.name)
		if err := certs.GenerateAndSaveKeys(cd.private, cd.public); err != nil {
			log.Fatalf("[FATAL] Failed to generate certs for %s: %v", cd.name, err)
		}
	}
}

// Service represents a running service
type Service struct {
	name string
	cmd  *exec.Cmd
	pid  int
}

type ServiceConfig struct {
	name string
	dir  string
	port string
}

// isPortAvailable returns true when the TCP port can be bound.
func isPortAvailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

func findListeningPIDs(port string) ([]string, error) {
	if runtime.GOOS == "windows" {
		intPIDs, err := getListeningPIDsWindows(port)
		if err != nil {
			return nil, err
		}
		selfPID := os.Getpid()
		var pids []string
		for _, pid := range intPIDs {
			if pid == selfPID {
				continue
			}
			pids = append(pids, strconv.Itoa(pid))
		}
		return pids, nil
	}

	// Linux/macOS
	cmd := exec.Command("lsof", "-tiTCP:"+port, "-sTCP:LISTEN")
	out, err := cmd.Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return []string{}, nil
		}
		return nil, err
	}

	selfPID := fmt.Sprintf("%d", os.Getpid())
	var pids []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		pid := strings.TrimSpace(line)
		if pid == "" || pid == selfPID {
			continue
		}
		pids = append(pids, pid)
	}

	return pids, nil
}

func waitUntilPortFree(port string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if isPortAvailable(port) {
			return true
		}
		time.Sleep(150 * time.Millisecond)
	}
	return isPortAvailable(port)
}

func autoFreePort(serviceName, port string) bool {
	if isPortAvailable(port) {
		return true
	}

	pids, err := findListeningPIDs(port)
	if err != nil {
		log.Printf("[!] Gagal cek konflik port %s untuk %s: %v", port, serviceName, err)
		return false
	}
	if len(pids) == 0 {
		log.Printf("[!] Port %s bentrok untuk %s, tapi PID listener tidak terdeteksi", port, serviceName)
		return false
	}

	fmt.Printf("[!] Port %s bentrok untuk %s. Auto-kill PID: %s\n", port, serviceName, strings.Join(pids, ", "))
	for _, pid := range pids {
		_ = exec.Command("kill", "-TERM", pid).Run()
	}

	if waitUntilPortFree(port, 2*time.Second) {
		return true
	}

	fmt.Printf("[!] Port %s masih dipakai, paksa kill PID: %s\n", port, strings.Join(pids, ", "))
	for _, pid := range pids {
		_ = exec.Command("kill", "-KILL", pid).Run()
	}

	if waitUntilPortFree(port, 2*time.Second) {
		return true
	}

	log.Printf("[-] Port %s tetap tidak bisa dibebaskan untuk service %s", port, serviceName)
	return false
}

func main() {
	fmt.Println("=========================================")
	fmt.Println("  🚀 THINKNALYZE ORCHESTRATOR")
	fmt.Println("=========================================")
	fmt.Println()

	// Check prerequisites
	checkPrerequisites()

	var wg sync.WaitGroup
	services := []*Service{}
	var mu sync.Mutex

	// Skip auto-start Redis — assumed to be running via Docker
	fmt.Println("[REDIS] Using existing Redis (Docker or local)...")
	var redisCmd *exec.Cmd // placeholder agar teardown tetap aman

	// Service configurations
	serviceConfigs := []ServiceConfig{
		{"account", "account", "2001"},
		{"gateway", "gateway", "2000"},
		{"users", "users", "2006"},
		{"tickets", "tickets", "2004"},
		{"notification", "notification", "5003"},
		{"operational", "operational", "5005"},
		{"subscription", "subscription", "5004"},
		{"management", "management", "5006"},
	}

	fmt.Println()
	fmt.Println("Memulai semua service...")
	fmt.Println()

	startedCount := 0
	skippedCount := 0

	// Start services
	for _, cfg := range serviceConfigs {
		if !autoFreePort(cfg.name, cfg.port) {
			fmt.Printf("[!] Skip service %s: port %s tidak bisa dibebaskan\n", cfg.name, cfg.port)
			skippedCount++
			continue
		}

		wg.Add(1)
		startedCount++
		go func(config ServiceConfig) {
			defer wg.Done()

			cmd := exec.Command("go", "run", "main.go")
			cmd.Dir = config.dir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				log.Printf("[-] Failed to start %s: %v", config.name, err)
				return
			}

			fmt.Printf("[+] Service %s berjalan dengan PID %d\n", config.name, cmd.Process.Pid)

			mu.Lock()
			services = append(services, &Service{
				name: config.name,
				cmd:  cmd,
				pid:  cmd.Process.Pid,
			})
			mu.Unlock()

			// Wait for process to finish
			if err := cmd.Wait(); err != nil {
				fmt.Printf("[-] Service %s berhenti: %v\n", config.name, err)
			}
		}(cfg)

		time.Sleep(500 * time.Millisecond)
	}

	if startedCount == 0 {
		fmt.Println("\n[!] Tidak ada service baru yang dijalankan oleh orchestrator.")
		if skippedCount > 0 {
			fmt.Printf("[i] %d service di-skip karena port sudah dipakai.\n", skippedCount)
		}
		fmt.Println("[i] Jalankan service per modul atau hentikan proses lama terlebih dahulu.")
		return
	}

	fmt.Println("\n=========================================")
	fmt.Printf("  ✅ ORCHESTRATOR STARTED %d SERVICES\n", startedCount)
	if skippedCount > 0 {
		fmt.Printf("  ℹ️  SKIPPED %d SERVICES (PORT IN USE)\n", skippedCount)
	}
	fmt.Println("=========================================")
	fmt.Println("\n  🌐 Gateway (Port 2000):")
	fmt.Println("     http://localhost:2000")
	fmt.Println("\n  🔐 Login:")
	fmt.Println("     http://localhost:2000/account/login")
	fmt.Println("\n  📊 Test Credentials:")
	fmt.Println("     Email: superadmin@thinktala.com")
	fmt.Println("     Pass:  Super123")
	fmt.Println("\n=========================================")
	fmt.Println()

	// Wait for all services
	wg.Wait()

	// Cleanup
	fmt.Println("\n[*] Shutting down services...")
	for _, svc := range services {
		if svc.cmd.Process != nil {
			svc.cmd.Process.Kill()
			fmt.Printf("[*] Terminated %s (PID %d)\n", svc.name, svc.pid)
		}
	}

	// Kill Redis (only if we started it)
	if redisCmd != nil && redisCmd.Process != nil {
		redisCmd.Process.Kill()
		fmt.Println("[*] Terminated Redis")
	}

	fmt.Println("\n✅ All services stopped")
}

// checkPrerequisites checks if required tools are installed
func checkPrerequisites() {
	tools := map[string]string{
		"Go":    "go version",
		"Redis": "redis-cli --version",
	}

	fmt.Println("[*] Checking prerequisites...")
	fmt.Println()

	allOk := true
	for tool, cmd := range tools {
		var checkCmd *exec.Cmd
		if runtime.GOOS == "windows" {
			checkCmd = exec.Command("cmd", "/c", cmd)
		} else {
			parts := strings.Fields(cmd)
			checkCmd = exec.Command(parts[0], parts[1:]...)
		}

		if err := checkCmd.Run(); err != nil {
			if tool == "Redis" {
				fmt.Printf("  ⚠️  %s - not found locally (assuming Docker Redis is running)\n", tool)
			} else {
				fmt.Printf("  ❌ %s - NOT INSTALLED\n", tool)
				allOk = false
			}
		} else {
			fmt.Printf("  ✅ %s - OK\n", tool)
		}
	}

	fmt.Printf("  ℹ️  PostgreSQL - (will be checked by services)\n")

	if !allOk {
		log.Fatal("\n[FATAL] Please install missing prerequisites")
	}

	fmt.Println()
}

func ensureServicePortsAvailable(configs []ServiceConfig) error {
	fmt.Println("[*] Checking service ports...\n")

	for _, cfg := range configs {
		if !isPortInUse(cfg.port) {
			fmt.Printf("  ✅ Port %s (%s) - available\n", cfg.port, cfg.name)
			continue
		}

		if runtime.GOOS != "windows" {
			return fmt.Errorf("port %s (%s) is already in use", cfg.port, cfg.name)
		}

		pids, err := getListeningPIDsWindows(cfg.port)
		if err != nil {
			return fmt.Errorf("failed to inspect port %s (%s): %w", cfg.port, cfg.name, err)
		}

		releasedByCleanup := false
		for _, pid := range pids {
			processName, err := getProcessNameWindows(pid)
			if err != nil {
				continue
			}

			if strings.EqualFold(processName, "main.exe") {
				fmt.Printf("  ⚠️  Port %s (%s) used by stale %s (PID %d). Terminating...\n", cfg.port, cfg.name, processName, pid)
				if err := killPIDWindows(pid); err != nil {
					return fmt.Errorf("failed to terminate stale process PID %d on port %s: %w", pid, cfg.port, err)
				}
				releasedByCleanup = true
			}
		}

		if releasedByCleanup {
			time.Sleep(300 * time.Millisecond)
		}

		if isPortInUse(cfg.port) {
			if releasedByCleanup {
				return fmt.Errorf("port %s (%s) still in use after cleanup", cfg.port, cfg.name)
			}
			return fmt.Errorf("port %s (%s) is in use by non-orchestrator process", cfg.port, cfg.name)
		}

		fmt.Printf("  ✅ Port %s (%s) - released\n", cfg.port, cfg.name)
	}

	fmt.Println()
	return nil
}

func isPortInUse(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return true
	}
	_ = ln.Close()
	return false
}

func getListeningPIDsWindows(port string) ([]int, error) {
	cmd := exec.Command("cmd", "/c", fmt.Sprintf("netstat -ano | findstr :%s | findstr LISTENING", port))
	out, err := cmd.Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return []int{}, nil
		}
		return nil, err
	}

	pidSet := make(map[int]bool)
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 5 {
			continue
		}
		pid, err := strconv.Atoi(fields[len(fields)-1])
		if err != nil {
			continue
		}
		pidSet[pid] = true
	}

	result := make([]int, 0, len(pidSet))
	for pid := range pidSet {
		result = append(result, pid)
	}

	return result, nil
}

func getProcessNameWindows(pid int) (string, error) {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/FO", "CSV", "/NH")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	line := strings.TrimSpace(string(out))
	if line == "" || strings.HasPrefix(line, "INFO:") {
		return "", fmt.Errorf("no process found for pid %d", pid)
	}

	r := csv.NewReader(strings.NewReader(line))
	record, err := r.Read()
	if err != nil || len(record) == 0 {
		return "", fmt.Errorf("unable to parse process info for pid %d", pid)
	}

	return strings.TrimSpace(record[0]), nil
}

func killPIDWindows(pid int) error {
	cmd := exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/F")
	return cmd.Run()
}
