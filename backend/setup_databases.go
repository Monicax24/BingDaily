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

	// First, drop all tables if they exist (clean start)
	fmt.Println("🧹 Cleaning up existing tables...")
	cleanupStatements := []string{
		"DROP TABLE IF EXISTS posts CASCADE",
		"DROP TABLE IF EXISTS communities CASCADE",
		"DROP TABLE IF EXISTS users CASCADE",
	}

	for _, sql := range cleanupStatements {
		_, err = db.Exec(sql)
		if err != nil {
			log.Printf("Warning: Failed to drop table: %v", err)
		}
	}

	// Create tables in correct order with CORRECT column names
	fmt.Println("🏗️ Creating tables...")
	sqlStatements := []string{
		// 1. Users table
		`CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			profile_picture TEXT DEFAULT '',
			joined_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			communities JSONB DEFAULT '[]',
			friends JSONB DEFAULT '[]',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// 2. Communities table
		`CREATE TABLE communities (
			community_id SERIAL PRIMARY KEY,
			picture TEXT DEFAULT '',
			description TEXT,
			members JSONB DEFAULT '[]',
			moderators JSONB DEFAULT '[]',
			posts JSONB DEFAULT '[]',
			post_time TIME DEFAULT '09:00:00',
			default_prompt TEXT DEFAULT 'What did you do today?'
		)`,

		// 3. Posts table
		`CREATE TABLE posts (
			post_id SERIAL PRIMARY KEY,
			community_id INTEGER REFERENCES communities(community_id) ON DELETE CASCADE,
			picture TEXT NOT NULL,
			caption TEXT,
			author INTEGER REFERENCES users(id) ON DELETE CASCADE,
			time_posted TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			likes JSONB DEFAULT '[]'
		)`,

		// Indexes
		`CREATE INDEX idx_users_email ON users(email)`,
		`CREATE INDEX idx_posts_community ON posts(community_id)`,
		`CREATE INDEX idx_posts_author ON posts(author)`,
		`CREATE INDEX idx_posts_time ON posts(time_posted)`,
	}

	for i, sql := range sqlStatements {
		fmt.Printf("Executing statement %d...\n", i+1)
		_, err = db.Exec(sql)
		if err != nil {
			log.Fatalf("❌ Failed to execute statement %d: %v\nSQL: %s", i+1, err, sql)
		}
		fmt.Printf("✅ Created table/index %d\n", i+1)
	}

	fmt.Println("\n🎉 Database setup complete! All tables created successfully.")
}
