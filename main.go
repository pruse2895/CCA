package main

import (
	"context"
	"log"
	"net/http"

	"cricketApp/db"
	"cricketApp/handlers"
	"cricketApp/router"
	"cricketApp/scheduler"
)

func main() {
	// Initialize MongoDB client and collections/indexes
	client, err := db.Init()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Create database instance using the initialized client
	dbName := "cricketApp"
	database := db.NewMongoDB(client, dbName)

	// Create handlers
	cricketerHandler := handlers.NewCricketerHandler(database)

	// Setup router with handlers and database instance
	r := router.SetupRouter(database, cricketerHandler)

	// Start the reminder scheduler
	reminderScheduler := scheduler.NewReminderScheduler(database)
	go reminderScheduler.Start()
	log.Println("Reminder scheduler started")

	// Start server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
