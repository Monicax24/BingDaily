package main

import (
	"bingdaily/backend/internal/database"
	"bingdaily/backend/internal/database/communities"
	"bingdaily/backend/internal/database/users"
	"context"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5"
)

func main() {
	db := database.InitializeDatabase()
	defer db.Close()

	// Clean up existing tables
	fmt.Println("🧹 Cleaning up existing tables...")
	cleanupStatements := []string{
		"DROP TABLE IF EXISTS dailies CASCADE",
		"DROP TABLE IF EXISTS communities CASCADE",
		"DROP TABLE IF EXISTS users CASCADE",
	}

	for _, sql := range cleanupStatements {
		_, err := db.Exec(context.TODO(), sql)
		if err != nil {
			log.Printf("Warning: Failed to drop table: %v", err)
		}
	}

	// Create tables with string IDs
	fmt.Println("🏗️ Creating tables with string IDs...")
	sqlStatements := []string{
		// Users table with string IDs
		`CREATE TABLE users (
			user_id VARCHAR(50) PRIMARY KEY,
			name VARCHAR(100) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			profile_picture TEXT,
			joined_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			communities TEXT[] DEFAULT '{}'
		)`,

		// Communities table with string IDs
		`CREATE TABLE communities (
			community_id VARCHAR(50) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			picture TEXT,
			description TEXT,
			members TEXT[] DEFAULT '{}',
			posts TEXT[] DEFAULT '{}',
			post_time VARCHAR(20) DEFAULT '09:00',
			prompt VARCHAR(255) DEFAULT 'What did you do today?'
		)`,

		// Dailies table with string IDs
		`CREATE TABLE dailies (
			post_id VARCHAR(50) PRIMARY KEY,
			community_id VARCHAR(50) REFERENCES communities(community_id) ON DELETE CASCADE,
			picture TEXT NOT NULL,
			caption TEXT,
			author VARCHAR(50) REFERENCES users(user_id) ON DELETE CASCADE,
			time_posted TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			likes TEXT[] DEFAULT '{}'
		)`,

		// Indexes
		`CREATE INDEX idx_users_email ON users(email)`,
		`CREATE INDEX idx_users_name ON users(name)`,
		`CREATE INDEX idx_dailies_community ON dailies(community_id)`,
		`CREATE INDEX idx_dailies_author ON dailies(author)`,
		`CREATE INDEX idx_dailies_time ON dailies(time_posted)`,
		`CREATE INDEX idx_communities_name ON communities(name)`,
	}

	for i, sql := range sqlStatements {
		fmt.Printf("Executing statement %d...\n", i+1)
		_, err := db.Exec(context.TODO(), sql)
		if err != nil {
			log.Fatalf("❌ Failed to execute statement %d: %v\nSQL: %s", i+1, err, sql)
		}
		fmt.Printf("✅ Created table/index %d\n", i+1)
	}

	fmt.Println("\n🎉 Database setup complete! All tables created with string IDs.")

	// default community
	communities.CreateCommunity(
		db,
		"ACM Project Team 4",
		"nopicture.jpg",
		"default community for testing",
		"09:00",
		"Upload a cool photo!",
	)
	db.Exec(context.TODO(),
		`UPDATE public.communities
		SET community_id = '6a6a671e-2543-4fad-ba82-dedc37338f14'
		WHERE name = 'ACM Project Team 4'`,
	)

	// WARNING: remove this after testing
	// default user
	users.Register(
		db,
		"testuser",
		"test@test.com",
		"nopicture.jpg",
	)
	// set id for default user
	db.Exec(context.TODO(),
		`UPDATE public.users
		SET user_id = '697b8a69-0c01-4ccb-aabc-6dccd6a22fa3'
		WHERE name = 'testuser'`,
	)
}
