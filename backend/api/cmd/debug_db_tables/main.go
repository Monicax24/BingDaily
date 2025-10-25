package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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

	// List all tables and their columns
	fmt.Println("📊 ALL TABLES AND COLUMNS:")
	rows, err := db.Query(`
		SELECT table_name, column_name, data_type 
		FROM information_schema.columns 
		WHERE table_schema = 'public'
		ORDER BY table_name, ordinal_position;
	`)
	if err != nil {
		log.Fatal("Failed to get table info:", err)
	}
	defer rows.Close()

	var tableName, columnName, dataType string
	currentTable := ""
	for rows.Next() {
		if err := rows.Scan(&tableName, &columnName, &dataType); err != nil {
			log.Fatal("Failed to scan:", err)
		}

		if tableName != currentTable {
			fmt.Printf("\n%s:\n", tableName)
			currentTable = tableName
		}
		fmt.Printf("  - %s (%s)\n", columnName, dataType)
	}

	// Check if users table has user_id
	fmt.Println("\n🔍 CHECKING USERS TABLE COLUMNS:")
	userColumns, err := db.Query(`
		SELECT column_name FROM information_schema.columns 
		WHERE table_name = 'users' AND table_schema = 'public';
	`)
	if err != nil {
		log.Fatal("Failed to check users columns:", err)
	}
	defer userColumns.Close()

	fmt.Println("Users table columns:")
	for userColumns.Next() {
		var col string
		if err := userColumns.Scan(&col); err != nil {
			log.Fatal("Failed to scan column:", err)
		}
		fmt.Printf("  - %s\n", col)
	}
}
