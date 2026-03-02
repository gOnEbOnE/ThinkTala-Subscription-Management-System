package main

import (
	"fmt"
	"log"
	"net"
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

// ensureCerts generates RSA keys ONCE from users, then copies to gateway & account.
// All services must share the same key pair: users signs, gateway/account verify.
func ensureCerts() {
	masterPrivate := "users/certs/private.pem"
	masterPublic := "users/certs/public.pem"

	// Step 1: Generate master key pair di users/certs/ jika belum ada
	_, errPriv := os.Stat(masterPrivate)
	_, errPub := os.Stat(masterPublic)
	if errPriv != nil || errPub != nil {
		fmt.Println("\n🔐 Generating master JWT key pair (users/certs/)...")
		os.MkdirAll("users/certs", 0755)
		if err := certs.GenerateAndSaveKeys(masterPrivate, masterPublic); err != nil {
			log.Fatalf("[FATAL] Failed to generate master certs: %v", err)
		}
		fmt.Println("✅ Master key pair generated")
	} else {
		fmt.Println("✅ Master key pair already exists (users/certs/)")
	}

	// Step 2: Baca master key
	privateKey, err := os.ReadFile(masterPrivate)
	if err != nil {
		log.Fatalf("[FATAL] Cannot read master private key: %v", err)
	}
	publicKey, err := os.ReadFile(masterPublic)
	if err != nil {
		log.Fatalf("[FATAL] Cannot read master public key: %v", err)
	}

	// Step 3: Copy ke semua service lain (gateway & account hanya butuh public key,
	// tapi kita copy private juga agar tidak ada error saat load)
	targets := []struct {
		name    string
		dir     string
	}{
		{"gateway", "gateway/certs"},
		{"account", "account/certs"},
	}

	for _, t := range targets {
		os.MkdirAll(t.dir, 0755)
		privPath := t.dir + "/private.pem"
		pubPath := t.dir + "/public.pem"

		// Cek apakah sudah sync
		existingPub, pubErr := os.ReadFile(pubPath)
		if pubErr == nil && string(existingPub) == string(publicKey) {
			fmt.Printf("✅ %s certificates already in sync\n", t.name)
			continue
		}

		if err := os.WriteFile(privPath, privateKey, 0600); err != nil {
			log.Fatalf("[FATAL] Failed to write %s private key: %v", t.name, err)
		}
		if err := os.WriteFile(pubPath, publicKey, 0644); err != nil {
			log.Fatalf("[FATAL] Failed to write %s public key: %v", t.name, err)
		}
		fmt.Printf("✅ %s certificates synced from master\n", t.name)
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

	// Start Redis (skip if already running, e.g. via WSL or Docker)
	var redisCmd *exec.Cmd
	if checkRedisConnectable() {
		fmt.Println("[REDIS] Redis already running on port 6379 ✅")
	} else {
		fmt.Println("[REDIS] Starting Redis server...")
		redisCmd = exec.Command("redis-server", "--port", "6379")
		redisCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if err := redisCmd.Start(); err != nil {
			log.Fatalf("[FATAL] Failed to start Redis: %v\nPastikan Redis sudah running (WSL: redis-server, atau Docker: docker run -d -p 6379:6379 redis:alpine)", err)
		}
		fmt.Printf("[+] Redis started with PID %d\n", redisCmd.Process.Pid)
		time.Sleep(500 * time.Millisecond)
	}

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

	// Kill Redis (only if we started it)
	if redisCmd != nil && redisCmd.Process != nil {
		redisCmd.Process.Kill()
		fmt.Println("[*] Terminated Redis")
	}

	fmt.Println("\n✅ All services stopped")
}

// checkGoInstalled checks if Go is installed
func checkGoInstalled() bool {
	var checkCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		checkCmd = exec.Command("cmd", "/c", "go version")
	} else {
		checkCmd = exec.Command("go", "version")
	}
	return checkCmd.Run() == nil
}

// checkRedisConnectable checks if Redis is reachable on localhost:6379 via TCP
func checkRedisConnectable() bool {
	conn, err := net.DialTimeout("tcp", "localhost:6379", 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// checkPrerequisites checks if required tools are installed
func checkPrerequisites() {
	fmt.Println("[*] Checking prerequisites...\n")

	allOk := true

	// Check Go
	if checkGoInstalled() {
		fmt.Printf("  ✅ Go - OK\n")
	} else {
		fmt.Printf("  ❌ Go - NOT INSTALLED\n")
		allOk = false
	}

	// Check Redis via TCP connection (supports WSL, Docker, or native Redis)
	if checkRedisConnectable() {
		fmt.Printf("  ✅ Redis - OK (connected to localhost:6379)\n")
	} else {
		fmt.Printf("  ❌ Redis - NOT REACHABLE (pastikan Redis sudah running di port 6379)\n")
		allOk = false
	}

	fmt.Printf("  ℹ️  PostgreSQL - (will be checked by services)\n")

	if !allOk {
		log.Fatal("\n[FATAL] Please install missing prerequisites")
	}

	fmt.Println()
}
