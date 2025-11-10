package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		log.Fatal("PG_DSN environment variable is not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer db.Close()

	// Clean up existing tables
	fmt.Println("🧹 Cleaning up existing tables...")
	cleanupStatements := []string{
		"DROP TABLE IF EXISTS dailies CASCADE",
		"DROP TABLE IF EXISTS communities CASCADE",
		"DROP TABLE IF EXISTS users CASCADE",
	}

	for _, sql := range cleanupStatements {
		_, err = db.Exec(sql)
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
			name VARCHAR(100) NOT NULL,
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
			default_prompt VARCHAR(255) DEFAULT 'What did you do today?'
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
		_, err = db.Exec(sql)
		if err != nil {
			log.Fatalf("❌ Failed to execute statement %d: %v\nSQL: %s", i+1, err, sql)
		}
		fmt.Printf("✅ Created table/index %d\n", i+1)
	}

	fmt.Println("\n🎉 Database setup complete! All tables created with string IDs.")
}
