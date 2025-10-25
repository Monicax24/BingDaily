package main

import (
	"bingdaily/backend/internal/database"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func demoOperations(db *sql.DB) {
	fmt.Println("\n🚀 Demo Operations:")

	// Create a test community
	communityID, err := database.CreateCommunity(db, "community1.jpg", "Daily photo sharing community", "09:00:00", "Share your daily photo!")
	if err != nil {
		log.Printf("Community creation failed: %v", err)
	} else {
		fmt.Printf("✅ Created community with ID: %d\n", communityID)
	}

	// Register test users with unique emails
	user1, err := database.Register(db, "Alice", fmt.Sprintf("alice%d@bing.com", time.Now().Unix()), "alice.jpg")
	if err != nil {
		log.Printf("User registration failed: %v", err)
	} else {
		fmt.Printf("✅ Registered user Alice with ID: %d\n", user1)
	}

	user2, err := database.Register(db, "Bob", fmt.Sprintf("bob%d@bing.com", time.Now().Unix()+1), "bob.jpg")
	if err != nil {
		log.Printf("User registration failed: %v", err)
	} else {
		fmt.Printf("✅ Registered user Bob with ID: %d\n", user2)
	}

	// Users join community
	err = database.JoinCommunity(db, user1, communityID)
	if err != nil {
		log.Printf("Join community failed: %v", err)
	} else {
		fmt.Printf("✅ User %d joined community %d\n", user1, communityID)
	}

	err = database.JoinCommunity(db, user2, communityID)
	if err != nil {
		log.Printf("Join community failed: %v", err)
	} else {
		fmt.Printf("✅ User %d joined community %d\n", user2, communityID)
	}

	// Create a post
	postID, err := database.CreatePost(db, communityID, "sunset.jpg", "Beautiful sunset today!", user1)
	if err != nil {
		log.Printf("Post creation failed: %v", err)
	} else {
		fmt.Printf("✅ Created post with ID: %d\n", postID)
	}

	// Like the post
	err = database.LikePost(db, postID, user2)
	if err != nil {
		log.Printf("Like post failed: %v", err)
	} else {
		fmt.Printf("✅ User %d liked post %d\n", user2, postID)
	}

	// Check if user has posted today
	hasPosted, err := database.HasPostedToday(db, user1, communityID)
	if err != nil {
		log.Printf("Check post failed: %v", err)
	} else {
		fmt.Printf("✅ User %d has posted today: %t\n", user1, hasPosted)
	}

	fmt.Println("\n🎉 Demo completed successfully!")
}

func main() {
	// Debug info
	fmt.Println("=== DEBUG INFO ===")
	fmt.Println("PG_DSN:", os.Getenv("PG_DSN"))
	fmt.Println("DATABASE_URL:", os.Getenv("DATABASE_URL"))
	wd, _ := os.Getwd()
	fmt.Println("Working directory:", wd)
	fmt.Println("==================")

	// CHANGE THIS TO USE .envrc
	// Connect to database
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		log.Fatal("PG_DSN environment variable is not set")
	}

	fmt.Println("Connecting with DSN:", dsn)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println("✅ Connected to database!")

	// Verify the required tables exist
	fmt.Println("\n🔍 Checking database structure...")
	if err := database.VerifyDatabaseStructure(db); err != nil {
		log.Fatal("❌ Database structure issue:", err)
	}

	// Show current tables
	if err := database.ListAllTables(db); err != nil {
		log.Fatal("Failed to list tables:", err)
	}

	// Show table structures
	if err := database.ShowTableStructure(db, "communities"); err != nil {
		log.Fatal("Failed to show table structure:", err)
	}
	if err := database.ShowTableStructure(db, "posts"); err != nil {
		log.Fatal("Failed to show table structure:", err)
	}
	if err := database.ShowTableStructure(db, "users"); err != nil {
		log.Fatal("Failed to show table structure:", err)
	}

	// Your application logic starts here
	fmt.Println("\n🎉 Database is ready! Starting application...")

	// Demo some operations
	demoOperations(db)
}
