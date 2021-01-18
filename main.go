package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/k3nnyM/OpenGroupTestIssue/model"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	queue := model.NewQueue()
	go queue.Run()

	router := new(mux.Router)

	router.HandleFunc("/add_file", queue.UploadFile)
	router.HandleFunc("/check_status", queue.CheckStatus)

	httpServer := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		fmt.Println("Server shutdown")
		return
	}

	return
}
