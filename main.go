package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/sachanritik1/go-lang/internal/app"
	"github.com/sachanritik1/go-lang/internal/routes"
)

func main() {

	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	app.Logger.Println("Application started successfully")

	var port int
	flag.IntVar(&port, "port", 8080, "Port to run the server on")
	flag.Parse()

	r := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}
	app.Logger.Println("Server is running on port " + fmt.Sprintf("%d", port))
	err = server.ListenAndServe()

	if err != nil {
		app.Logger.Fatalf("Error starting server: %v", err)
	}

}
