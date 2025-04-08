package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"cricketApp/models"
)

// Removed global Client: var Client *mongo.Client

// Init connects to MongoDB, initializes collections/indexes, and returns the client.
func Init() (*mongo.Client, error) { // Changed signature to return client and error
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding without it") // Changed log level
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Println("MONGODB_URI environment variable not set, using default mongodb://localhost:27017")
		mongoURI = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions) // Use local client variable
	if err != nil {
		log.Printf("Error connecting to MongoDB: %v", err)
		return nil, err // Return error
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Error pinging MongoDB: %v", err)
		return nil, err // Return error
	}
	log.Println("Connected to MongoDB!")

	// Initialize collections and indexes using the local client
	err = initCollections(client) // Pass client
	if err != nil {
		// Disconnect if initialization fails?
		_ = client.Disconnect(context.Background()) // Attempt disconnect
		return nil, err                             // Return initialization error
	}

	return client, nil // Return connected client
}

// initCollections initializes all necessary collections and indexes.
func initCollections(client *mongo.Client) error { // Accept client
	dbName := "cricketApp" // Or get from env
	if err := initCricketersCollection(client, dbName); err != nil {
		return err
	}
	if err := initAnnouncementsCollection(client, dbName); err != nil {
		return err
	}
	if err := initAdminsCollection(client, dbName); err != nil {
		return err
	}
	log.Println("Collections and indexes created successfully")
	return nil
}

// initCricketersCollection creates indexes for the cricketers collection.
func initCricketersCollection(client *mongo.Client, dbName string) error { // Accept client and dbName
	ctx := context.Background()
	cricketersCollection := client.Database(dbName).Collection("cricketers")

	// Create unique indexes for email and mobile
	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	mobileIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "mobile", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := cricketersCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{emailIndex, mobileIndex})
	if err != nil {
		log.Printf("Error creating cricketers indexes: %v", err)
		return err
	}
	return nil
}

// initAnnouncementsCollection creates indexes for the announcements collection.
func initAnnouncementsCollection(client *mongo.Client, dbName string) error { // Accept client and dbName
	ctx := context.Background()
	announcementsCollection := client.Database(dbName).Collection("announcements")

	// Create index for sorting by creation date
	createdAtIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "createdAt", Value: -1}}, // -1 for descending order
	}

	_, err := announcementsCollection.Indexes().CreateOne(ctx, createdAtIndex)
	if err != nil {
		log.Printf("Error creating announcements index: %v", err)
		return err
	}
	return nil
}

// initAdminsCollection creates index and default admin for the admins collection.
func initAdminsCollection(client *mongo.Client, dbName string) error { // Accept client and dbName
	ctx := context.Background()
	adminsCollection := client.Database(dbName).Collection("admins")

	// Create unique index for email
	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := adminsCollection.Indexes().CreateOne(ctx, emailIndex)
	if err != nil {
		// Don't treat duplicate index error as fatal for this setup
		if !mongo.IsDuplicateKeyError(err) && !isIndexAlreadyExistsError(err) {
			log.Printf("Error creating admins index: %v", err)
			return err
		}
		log.Println("Admin email index already exists or non-fatal error occurred.")
	}

	// Hash the default admin password
	defaultPassword := os.Getenv("DEFAULT_ADMIN_PASSWORD")
	if defaultPassword == "" {
		defaultPassword = "admin123"
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing admin password: %v", err)
		return err
	}

	// Create default admin if not exists
	defaultAdminEmail := os.Getenv("DEFAULT_ADMIN_EMAIL")
	if defaultAdminEmail == "" {
		defaultAdminEmail = "admin@example.com"
	}
	defaultAdmin := models.Admin{
		// ID: primitive.NewObjectID(), // Let MongoDB generate ID
		Name:     "Admin",
		Email:    defaultAdminEmail,
		Password: string(hashedPassword),
	}

	// Try to insert default admin
	_, err = adminsCollection.InsertOne(ctx, defaultAdmin)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Println("Default admin already exists")
		} else {
			log.Printf("Error creating default admin: %v", err)
			return err
		}
	} else {
		log.Println("Default admin created successfully")
	}
	return nil
}

// Helper function to check for index already exists errors (example structure)
func isIndexAlreadyExistsError(err error) bool {
	// MongoDB driver errors might not have a specific type for this,
	// often requires checking error message strings. This is brittle.
	// Example: return strings.Contains(err.Error(), "index already exists")
	// Or check for specific command error codes if available.
	return false // Placeholder - adjust based on actual error inspection
}

// Removed global GetCollection function

// Removed Cricketer operations (moved to cricketer.go)

// Removed Admin operations (moved to admin.go)

// Removed Announcement operations (moved to announcement.go)
