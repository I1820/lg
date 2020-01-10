package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/toskatok/lg/actions"
)

// ShutdownPeriod is a waiting period before forcing shutdown
const ShutdownPeriod = 5 * time.Second

func main() {
	fmt.Println("18.20 at Sep 07 2016 7:20 IR721")

	e := actions.App()

	go func() {
		if err := e.Start(":8080"); err != http.ErrServerClosed {
			log.Fatalf("API Service failed with %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("18.20 As always ... left me alone")

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownPeriod)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Printf("API Service failed on exit: %s", err)
	}
}
