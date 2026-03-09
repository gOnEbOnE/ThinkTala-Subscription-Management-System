
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"thinknalyze/certs"
)

func init() {
	ensureCerts()
}

// ================= CERTIFICATE SETUP =================

func ensureCerts() {
	masterPrivate := "users/certs/private.pem"
	masterPublic := "users/certs/public.pem"

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

	privateKey, err := os.ReadFile(masterPrivate)
	if err != nil {
		log.Fatalf("[FATAL] Cannot read master private key: %v", err)
	}

	publicKey, err := os.ReadFile(masterPublic)
	if err != nil {
		log.Fatalf("[FATAL] Cannot read master public key: %v", err)
	}

	targets := []struct {
		name string
		dir  string
	}{
		{"gateway", "gateway/certs"},
		{"account", "account/certs"},
	}

	for _, t := range targets {
		os.MkdirAll(t.dir, 0755)

		privPath := t.dir + "/private.pem"
		pubPath := t.dir + "/public.pem"

		existingPub, pubErr := os.ReadFile(pubPath)
		if pubErr == nil && string(existingPub) == string(publicKey) {
			fmt.Printf("✅ %s certificates already in sync\n", t.name)
			continue
		}

		os.WriteFile(privPath, privateKey, 0600)
		os.WriteFile(pubPath, publicKey, 0644)

		fmt.Printf("✅ %s certificates synced from master\n", t.name)
	}
}

// ================= SERVICE STRUCT =================

type Service struct {
	name string
	cmd  *exec.Cmd
	pid  int
}

// ================= MAIN =================

func main() {

	fmt.Println("=========================================")
	fmt.Println("  🚀 THINKNALYZE ORCHESTRATOR")
	fmt.Println("=========================================\n")

	checkPrerequisites()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var services []*Service
	var mu sync.Mutex

	serviceConfigs := []struct {
		name string
		dir  string
		port string
	}{
		{"account", "account", "2001"},
		{"gateway", "gateway", "2000"},
		{"users", "users", "2006"},
		{"notification", "notification", "5003"},
		{"operational", "operational", "5005"},
		{"subscription", "subscription", "5004"},
	}

	fmt.Println("\nMemulai semua service...\n")

	for _, cfg := range serviceConfigs {


		cmd := exec.CommandContext(ctx, "go", "run", "main.go")
		cmd.Dir = cfg.dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr


		if err := cmd.Start(); err != nil {
			log.Printf("[-] Failed to start %s: %v", cfg.name, err)
			continue
		}

		fmt.Printf("[+] Service %s berjalan dengan PID %d\n", cfg.name, cmd.Process.Pid)

		mu.Lock()
		services = append(services, &Service{
			name: cfg.name,
			cmd:  cmd,
			pid:  cmd.Process.Pid,
		})
		mu.Unlock()

		time.Sleep(400 * time.Millisecond)
	}

	printInfo()

	<-ctx.Done()

	fmt.Println("\n=========================================")
	fmt.Println("  🛑 SHUTDOWN SIGNAL RECEIVED")
	fmt.Println("=========================================")

	for _, svc := range services {

		if svc.cmd.Process != nil {

			fmt.Printf("[*] Terminating %s (PID %d)\n", svc.name, svc.pid)

			svc.cmd.Process.Signal(syscall.SIGTERM)

			done := make(chan error, 1)
			go func() {
				done <- svc.cmd.Wait()
			}()

			select {

			case <-done:

			case <-time.After(3 * time.Second):
				fmt.Printf("[!] Force killing %s\n", svc.name)
				svc.cmd.Process.Kill()
			}
		}
	}

	fmt.Println("\n✅ All services stopped")
}

// ================= UTILITIES =================

func checkPrerequisites() {

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
			parts := strings.Fields(cmd)
			checkCmd = exec.Command(parts[0], parts[1:]...)
		}

		if err := checkCmd.Run(); err != nil {

			if tool == "Redis" {
				fmt.Printf("  ⚠️  %s - not found locally (assuming Docker Redis)\n", tool)
			} else {
				fmt.Printf("  ❌ %s - NOT INSTALLED\n", tool)
				allOk = false
			}

		} else {
			fmt.Printf("  ✅ %s - OK\n", tool)
		}
	}

	fmt.Println("  ℹ️  PostgreSQL - checked by services")

	if !allOk {
		log.Fatal("\n[FATAL] Please install missing prerequisites")
	}

	fmt.Println()
}

func printInfo() {

	fmt.Println("\n=========================================")
	fmt.Println("  ✅ ALL SERVICES STARTED")
	fmt.Println("=========================================")

	fmt.Println("\n🌐 Gateway")
	fmt.Println("http://localhost:2000")

	fmt.Println("\n🔐 Login")
	fmt.Println("http://localhost:2000/account/login")

	fmt.Println("\n📊 Test Credentials")
	fmt.Println("Email: superadmin@thinktala.com")
	fmt.Println("Pass : Super123")

	fmt.Println("\n📌 Pages")
	fmt.Println("http://localhost:2000/ops/dashboard")
	fmt.Println("http://localhost:2000/ops/notifications")
	fmt.Println("http://localhost:2000/ops/notification-templates")

	fmt.Println("\n=========================================\n")
}