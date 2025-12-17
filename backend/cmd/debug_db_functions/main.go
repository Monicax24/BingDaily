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

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func demoOperations(db *pgxpool.Pool) {
	fmt.Println("\n🚀 Demo Operations:")

	// Create a test community
	communityID := uuid.New().String()
	success, err := communities.CreateCommunity(db,
		communityID,
		"Daily Photos",                  // name
		"",                              // picture
		"Daily photo sharing community", // description
		"09:00",                         // postTime
		"Share your daily photo!")       // defaultPrompt
	if err != nil {
		log.Printf("Error creating community: %v", err)
	} else {
		fmt.Printf("✅ Status creating community (%s): %v\n", communityID, success)
	}

	// Test GetCommunity
	community, err := communities.GetCommunity(db, communityID)
	if err != nil {
		log.Printf("Failed to get community: %v", err)
	} else {
		fmt.Printf("✅ Retreived community: %v\n", community)
	}

	// Register test users with unique emails
	// user 1
	user1 := uuid.New().String()
	success, err = users.Register(db, user1, "Alice", "alice@bing.com", "")
	if err != nil {
		log.Printf("Error during User Registration: %v", err)
	} else {
		fmt.Printf("✅ Status registering Alice (%s): %v\n", user1, success)
	}
	// user 2
	user2 := uuid.New().String()
	success, err = users.Register(db, user2, "Bob", "bob@bing.com", "")
	if err != nil {
		log.Printf("Error during User Registration: %v", err)
	} else {
		fmt.Printf("✅ Status registering Bob (%s): %v\n", user2, success)
	}

	// Test in Community
	in, err := communities.UserInCommunity(db, communityID, user1)
	if err != nil {
		log.Printf("UserInCommunity check failed: %v", err)
	} else {
		fmt.Printf("✅ UserInCommunity check result before joining: %v\n", in)
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
	in, err = communities.UserInCommunity(db, communityID, user1)
	if err != nil {
		log.Printf("UserInCommunity check failed: %v", err)
	} else {
		fmt.Printf("✅ UserInCommunity check result after joining: %v\n", in)
	}

	// Create a daily post
	dailyID, err := dailies.CreateDaily(db, communityID, "", "Beautiful sunset today!", user1)
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
	dailyID2, err := dailies.CreateDaily(db, communityID, "", "Great hike today!", user1)
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

	// Check getting non-existent user
	fakeUser, err := users.GetUser(db, "not-a-real-user-id")
	if err != nil {
		fmt.Printf("Fetching nonexistent user failed: %v\n", err)
	} else {
		fmt.Printf("Value of nonexistent user: %v\n", fakeUser)
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

	fmt.Println("\n🎉 Demo completed successfully!")
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
