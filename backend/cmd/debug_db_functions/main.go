package main

import (
	"bingdaily/backend/internal/database"
	"bingdaily/backend/internal/database/communities"
	"bingdaily/backend/internal/database/dailies"
	"bingdaily/backend/internal/database/users"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func demoOperations(db *sql.DB) {
	fmt.Println("\n🚀 Demo Operations:")

	// Create a test community
	communityID, err := communities.CreateCommunity(db,
		"Daily Photos",                  // name
		"community1.jpg",                // picture
		"Daily photo sharing community", // description
		"09:00",                         // postTime
		"Share your daily photo!")       // defaultPrompt
	if err != nil {
		log.Printf("Community creation failed: %v", err)
	} else {
		fmt.Printf("✅ Created community with ID: %s\n", communityID)
	}

	// Register test users with unique emails
	user1, err := users.Register(db, "Alice", fmt.Sprintf("alice%d@bing.com", time.Now().Unix()), "alice.jpg")
	if err != nil {
		log.Printf("User registration failed: %v", err)
	} else {
		fmt.Printf("✅ Registered user Alice with ID: %s\n", user1)
	}

	user2, err := users.Register(db, "Bob", fmt.Sprintf("bob%d@bing.com", time.Now().Unix()+1), "bob.jpg")
	if err != nil {
		log.Printf("User registration failed: %v", err)
	} else {
		fmt.Printf("✅ Registered user Bob with ID: %s\n", user2)
	}

	// Users join community
	err = communities.JoinCommunity(db, user1, communityID)
	if err != nil {
		log.Printf("Join community failed: %v", err)
	} else {
		fmt.Printf("✅ User %s joined community %s\n", user1, communityID)
	}

	err = communities.JoinCommunity(db, user2, communityID)
	if err != nil {
		log.Printf("Join community failed: %v", err)
	} else {
		fmt.Printf("✅ User %s joined community %s\n", user2, communityID)
	}

	// Create a daily post
	dailyID, err := dailies.CreateDaily(db, communityID, "sunset.jpg", "Beautiful sunset today!", user1)
	if err != nil {
		log.Printf("Daily creation failed: %v", err)
	} else {
		fmt.Printf("✅ Created daily with ID: %s\n", dailyID)
	}

	// Like the daily post
	err = dailies.LikeDaily(db, dailyID, user2)
	if err != nil {
		log.Printf("Like daily failed: %v", err)
	} else {
		fmt.Printf("✅ User %s liked daily %s\n", user2, dailyID)
	}

	// Check if user has posted today
	hasPosted, err := dailies.HasPostedToday(db, user1, communityID)
	if err != nil {
		log.Printf("Check daily post failed: %v", err)
	} else {
		fmt.Printf("✅ User %s has posted today: %t\n", user1, hasPosted)
	}

	// Like the daily post
	err = dailies.LikeDaily(db, dailyID, user2)
	if err != nil {
		log.Printf("Like daily failed: %v", err)
	} else {
		fmt.Printf("✅ User %s liked daily %s\n", user2, dailyID)
	}

	// Demo: Delete the daily post
	fmt.Println("\n🗑️ Testing delete operations...")
	err = dailies.DeleteDaily(db, dailyID)
	if err != nil {
		log.Printf("Delete daily failed: %v", err)
	} else {
		fmt.Printf("✅ Deleted daily with ID: %s\n", dailyID)
	}

	// Verify the daily was deleted
	_, err = dailies.GetDaily(db, dailyID)
	if err != nil {
		fmt.Printf("✅ Confirmed daily %s no longer exists: %v\n", dailyID, err)
	} else {
		fmt.Printf("❌ Daily %s still exists after deletion\n", dailyID)
	}

	// Create another daily to test bulk deletion
	dailyID2, err := dailies.CreateDaily(db, communityID, "mountain.jpg", "Great hike today!", user1)
	if err != nil {
		log.Printf("Second daily creation failed: %v", err)
	} else {
		fmt.Printf("✅ Created second daily with ID: %s\n", dailyID2)
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
	if err := database.ShowTableStructure(db, "dailies"); err != nil {
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
