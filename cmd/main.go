package main

import (
	"context"
	"log"
	"time"

	"wallet/internal/app"
)

func main() {

	app, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println("Server is start on port", app.Config.Port)
		err := app.Server.ListenAndServe()
		if err != nil {
			log.Fatal("failed ListenAndServe: ", err)
		}
	}()

	sig := <-app.Sigint
	log.Printf("Received signal: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

}
