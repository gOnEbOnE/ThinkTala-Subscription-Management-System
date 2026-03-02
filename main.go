package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
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

func main() {
	fmt.Println("=========================================")
	fmt.Println("  🚀 THINKNALYZE ORCHESTRATOR")
	fmt.Println("=========================================\n")

	// Check prerequisites
	checkPrerequisites()

	var wg sync.WaitGroup
	services := []*Service{}
	var mu sync.Mutex

	// Start Redis
	fmt.Println("[REDIS] Starting Redis server...")
	redisCmd := exec.Command("redis-server", "--port", "6379")
	redisCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := redisCmd.Start(); err != nil {
		log.Fatalf("[FATAL] Failed to start Redis: %v", err)
	}
	fmt.Printf("[+] Redis started with PID %d\n", redisCmd.Process.Pid)
	time.Sleep(500 * time.Millisecond)

	// Service configurations
	serviceConfigs := []struct {
		name string
		dir  string
		port string
	}{
		{"account", "account", "2001"},
		{"gateway", "gateway", "2000"},
		{"users", "users", "2006"},
	}

	fmt.Println("\nMemulai semua service...\n")

	// Start services
	for _, cfg := range serviceConfigs {
		wg.Add(1)
		go func(config struct {
			name string
			dir  string
			port string
		}) {
			defer wg.Done()

			cmd := exec.Command("go", "run", "main.go")
			cmd.Dir = config.dir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}

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

	fmt.Println("\n=========================================")
	fmt.Println("  ✅ ALL SERVICES STARTED")
	fmt.Println("=========================================")
	fmt.Println("\n  🌐 Gateway (Port 2000):")
	fmt.Println("     http://localhost:2000")
	fmt.Println("\n  🔐 Login:")
	fmt.Println("     http://localhost:2000/account/login")
	fmt.Println("\n  📊 Test Credentials:")
	fmt.Println("     Email: superadmin@thinktala.com")
	fmt.Println("     Pass:  Super123")
	fmt.Println("\n=========================================\n")

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

	// Kill Redis
	if redisCmd.Process != nil {
		redisCmd.Process.Kill()
		fmt.Println("[*] Terminated Redis")
	}

	fmt.Println("\n✅ All services stopped")
}

// checkPrerequisites checks if required tools are installed
func checkPrerequisites() {
	// ⭐ MODIFIED: Only check Go & Redis, skip PostgreSQL
	tools := map[string]string{
		"Go":    "go version",
		"Redis": "redis-cli --version",
	}

	fmt.Println("[*] Checking prerequisites...\n")

	allOk := true
	for tool, cmd := range tools {
		var checkCmd *exec.Cmd
		if runtime.GOOS == "windows" {
			checkCmd = exec.Command("cmd", "/c", cmd)
		} else {
			parts := []string{cmd}
			checkCmd = exec.Command(parts[0])
		}

		if err := checkCmd.Run(); err != nil {
			fmt.Printf("  ❌ %s - NOT INSTALLED\n", tool)
			allOk = false
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
