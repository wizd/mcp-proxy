package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type MiddlewareFunc func(http.Handler) http.Handler

func chainMiddleware(h http.Handler, middlewares ...MiddlewareFunc) http.Handler {
	for _, mw := range middlewares {
		h = mw(h)
	}
	return h
}

func newAuthMiddleware(tokens []string) MiddlewareFunc {
	tokenSet := make(map[string]struct{}, len(tokens))
	for _, token := range tokens {
		tokenSet[token] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(tokens) != 0 {
				token := r.Header.Get("Authorization")
				if strings.HasPrefix(token, "Bearer ") {
					token = strings.TrimPrefix(token, "Bearer ")
				}
				if token == "" {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				if _, ok := tokenSet[token]; !ok {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func startHTTPServer(config *Config) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var errorGroup errgroup.Group
	httpMux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    config.Server.Addr,
		Handler: httpMux,
	}
	info := mcp.Implementation{
		Name:    config.Server.Name,
		Version: config.Server.Version,
	}

	for name, clientConfig := range config.Clients {
		mcpClient, err := newMCPClient(name, clientConfig)
		if err != nil {
			log.Fatalf("<%s> Failed to create client: %v", name, err)
		}
		server := newMCPServer(name, &config.Server, &clientConfig)
		errorGroup.Go(func() error {
			log.Printf("<%s> Connecting", name)
			addErr := mcpClient.addToMCPServer(ctx, info, server.mcpServer)
			if addErr != nil {
				log.Printf("<%s> Failed to add client to server: %v", name, addErr)
				if clientConfig.PanicIfInvalid {
					return addErr
				}
				return nil
			}
			log.Printf("<%s> Connected", name)
			httpMux.Handle(fmt.Sprintf("/%s/", name), chainMiddleware(server.sseServer, newAuthMiddleware(server.tokens)))
			httpServer.RegisterOnShutdown(func() {
				log.Printf("<%s> Shutting down", name)
				_ = mcpClient.Close()
			})
			return nil
		})
	}

	go func() {
		err := errorGroup.Wait()
		if err != nil {
			log.Fatalf("Failed to add clients: %v", err)
		}
		log.Printf("All clients initialized")
	}()

	go func() {
		log.Printf("Starting SSE server")
		log.Printf("SSE server listening on %s", config.Server.Addr)
		hErr := httpServer.ListenAndServe()
		if hErr != nil && !errors.Is(hErr, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", hErr)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	err := httpServer.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}
