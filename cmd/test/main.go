package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigCh := make(chan os.Signal, 1)

	signal.Notify(
		sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	fmt.Println("app started")

	sig := <-sigCh

	fmt.Println("received signal:", sig)
	fmt.Println("shutting down")
}
