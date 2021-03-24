package main

import (
	"blog/controllers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"blog/conf"
	"blog/databases"
)

func main() {
	// release resources
	defer func() {
		databases.Close()
	}()
	gin.ForceConsoleColor()

	serverURL := fmt.Sprintf("%s:%d", conf.Config.Server.Listen, conf.Config.Server.Port)
	srv := &http.Server{
		Addr:    serverURL,
		Handler: controllers.Router,
	}
	go func() {
		log.Printf("trying to listen on %s\n", serverURL)
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
