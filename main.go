package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var cmds []*exec.Cmd

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// ===== AUTO START REDIS DULU =====
	fmt.Println("[REDIS] Starting Redis server...")
	redisCmd := exec.Command("redis-server")

	// Suppress redis output (optional, biar ga berisik)
	// redisCmd.Stdout = nil
	// redisCmd.Stderr = nil

	if err := redisCmd.Start(); err != nil {
		log.Printf("[WARNING] Redis gagal start: %v (App will run in no-cache mode)", err)
	} else {
		cmds = append(cmds, redisCmd)
		fmt.Printf("[+] Redis started with PID %d\n", redisCmd.Process.Pid)

		// Tunggu 2 detik supaya Redis ready
		time.Sleep(2 * time.Second)
	}

	// ===== START SERVICES =====
	services := []string{"account", "gateway", "users"}

	fmt.Println("Memulai semua service...")

	for _, service := range services {
		wg.Add(1)

		cmd := exec.Command("go", "run", "main.go")
		cmd.Dir = service
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Start()
		if err != nil {
			log.Printf("Gagal menjalankan service %s: %v\n", service, err)
			wg.Done()
			continue
		}

		cmds = append(cmds, cmd)
		fmt.Printf("[+] Service %s berjalan dengan PID %d\n", service, cmd.Process.Pid)

		go func(c *exec.Cmd, name string) {
			defer wg.Done()
			err := c.Wait()
			if err != nil {
				log.Printf("[-] Service %s berhenti: %v\n", name, err)
			} else {
				fmt.Printf("[-] Service %s berhenti dengan normal\n", name)
			}
		}(cmd, service)
	}

	<-quit
	fmt.Println("\nSinyal terminasi diterima. Mematikan semua service...")

	for _, cmd := range cmds {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}

	wg.Wait()
	fmt.Println("Semua service telah dimatikan. Keluar dari program.")
}
