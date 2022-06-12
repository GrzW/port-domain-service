package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/GrzW/port-domain-service/internal/handler"
	"github.com/GrzW/port-domain-service/internal/parser"
	"github.com/GrzW/port-domain-service/internal/storage"
)

const (
	serviceDefaultPort = "8080"
	servicePortEnvVar  = "SERVICE_PORT"
)

type PortsUpdatedMessage struct {
	FileURI string `json:"file_uri"`
}

func main() {
	portsDB := storage.NewMemoryDB()
	portsUpdatedHandler := handler.NewPortsUpdatedHandler(storage.NewPortsStore(portsDB), parser.NewJSONStreamParser())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		signalChanel := make(chan os.Signal, 1)
		signal.Notify(signalChanel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-signalChanel

		fmt.Printf("Gracefully stopping...")
		cancel()
	}()

	http.HandleFunc("/ports", func(w http.ResponseWriter, req *http.Request) {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
		defer func() { _ = req.Body.Close() }()

		var msg PortsUpdatedMessage
		if err = json.Unmarshal(b, &msg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		// because we expect it may take a while, let's respond with 202 status right away and run this in background
		go func() {
			if err = portsUpdatedHandler.HandlePortsUpdated(ctx, msg.FileURI); err != nil {
				fmt.Printf("Failed to handle ports updated message: %s", err.Error())
			}

			fmt.Printf("Finished handling %q", msg.FileURI)
		}()

		w.WriteHeader(http.StatusAccepted)
	})

	servicePort := os.Getenv(servicePortEnvVar)
	if servicePort == "" {
		servicePort = serviceDefaultPort
	}

	fmt.Printf("Service listening on port %s\n", servicePort)

	if err := http.ListenAndServe(":"+servicePort, nil); err != nil {
		log.Fatalf("starting HTTP server: %s", err.Error())
	}
}
