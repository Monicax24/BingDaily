package main

import (
	"bingdaily/backend/internal/database"
	"bingdaily/backend/internal/database/communities"
	"bingdaily/backend/internal/database/dailies"
	"bingdaily/backend/internal/database/users"
	"context"

	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func demoOperations(db *pgxpool.Pool) {
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

	// Test GetCommunity
	// Retrieve the community
	community, err := communities.GetCommunity(db, communityID)
	if err != nil {
		log.Printf("Failed to get community: %v", err)
	} else {
		fmt.Printf("✅ Retreived community: %v\n", community)
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

	// Test in Community
	in, err := communities.InCommunity(db, communityID, user1)
	if err != nil {
		log.Printf("InCommunity check failed: %v", err)
	} else {
		fmt.Printf("✅ InCommunity check result: %v\n", in)
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

	// Test in Community
	in, err = communities.InCommunity(db, communityID, user1)
	if err != nil {
		log.Printf("InCommunity check failed: %v", err)
	} else {
		fmt.Printf("✅ InCommunity check result: %v\n", in)
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

	// Create another daily to test bulk deletion
	dailyID2, err := dailies.CreateDaily(db, communityID, "mountain.jpg", "Great hike today!", user1)
	if err != nil {
		log.Printf("Second daily creation failed: %v", err)
	} else {
		fmt.Printf("✅ Created second daily with ID: %s\n", dailyID2)
	}

	// Fetch all daily posts
	dlies, err := dailies.FetchDailiesFromCommunity(db, communityID)
	if err != nil {
		log.Printf("Fetching dailies from community failed: %v", err)
	} else {
		fmt.Printf("✅ Dailies fetched: %v\n", dlies)
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

	fmt.Println("\n🎉 Demo completed successfully!")
}

func main() {
	// Debug info
	fmt.Println("=== DEBUG INFO ===")
	fmt.Println("PG_DSN:", os.Getenv("PG_DSN"))
	wd, _ := os.Getwd()
	fmt.Println("Working directory:", wd)
	fmt.Println("==================")

	// Connect to database
	db := database.InitializeDatabase()

	// Test connection
	if err := db.Ping(context.TODO()); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println("✅ Connected to database!")

	// Verify database structure
	if err := VerifyDatabaseStructure(db); err != nil {
		log.Fatal("❌ Database structure issue:", err)
	}
	demoOperations(db)
}

// VerifyDatabaseStructure checks all required tables exist
func VerifyDatabaseStructure(db *pgxpool.Pool) error {
	requiredTables := []string{"communities", "dailies", "users"}

	for _, table := range requiredTables {
		var exists bool
		err := db.QueryRow(context.TODO(), `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public'
				AND table_name = $1
			)`, table).Scan(&exists)

		if err != nil {
			return fmt.Errorf("failed to check table %s: %v", table, err)
		}

		if !exists {
			return fmt.Errorf("required table '%s' does not exist", table)
		}
		fmt.Printf("✅ Table '%s' exists\n", table)
	}
	return nil
}
