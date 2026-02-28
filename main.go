package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// Daftar folder service yang ingin dijalankan
	services := []string{"account", "gateway", "users"}

	var wg sync.WaitGroup
	var cmds []*exec.Cmd

	// Membuat channel untuk mendengarkan sinyal interupsi (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Memulai semua service...")

	for _, service := range services {
		wg.Add(1)

		// Menyiapkan command 'go run main.go'
		cmd := exec.Command("go", "run", "main.go")

		// Mengatur direktori kerja ke folder service masing-masing
		cmd.Dir = service

		// Meneruskan output service ke terminal root
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Memulai proses tanpa memblokir
		err := cmd.Start()
		if err != nil {
			log.Printf("Gagal menjalankan service %s: %v\n", service, err)
			wg.Done()
			continue
		}

		cmds = append(cmds, cmd)
		fmt.Printf("[+] Service %s berjalan dengan PID %d\n", service, cmd.Process.Pid)

		// Goroutine untuk menunggu proses selesai / error
		go func(c *exec.Cmd, name string) {
			defer wg.Done()
			err := c.Wait()
			if err != nil {
				// Akan log error jika service mati paksa
				log.Printf("[-] Service %s berhenti: %v\n", name, err)
			} else {
				fmt.Printf("[-] Service %s berhenti dengan normal\n", name)
			}
		}(cmd, service)
	}

	// Menunggu sinyal interupsi (Ctrl+C) dari user
	<-quit
	fmt.Println("\nSinyal terminasi diterima. Mematikan semua service...")

	// Mematikan semua child process saat root dimatikan
	for _, cmd := range cmds {
		if cmd.Process != nil {
			// Mengirim sinyal kill ke setiap process
			cmd.Process.Signal(syscall.SIGTERM)
		}
	}

	// Menunggu semua goroutine selesai memastikan proses mati
	wg.Wait()
	fmt.Println("Semua service telah dimatikan. Keluar dari program.")
}
