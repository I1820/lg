package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/toskatok/lg/handler"
)

// ShutdownPeriod is a waiting period before forcing shutdown
const ShutdownPeriod = 5 * time.Second

func main() {
	e := echo.New()

	api := e.Group("api")
	handler.NewInstance().Register(api)

	go func() {
		if err := e.Start(":1378"); err != http.ErrServerClosed {
			log.Fatalf("API Service failed with %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownPeriod)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Printf("API Service failed on exit: %s", err)
	}
}

func Register(root *cobra.Command) {
	root.AddCommand(&cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	})
}
