package database

import (
	"database/sql"
	"fmt"
)

// VerifyDatabaseStructure checks all required tables exist
func VerifyDatabaseStructure(db *sql.DB) error {
	requiredTables := []string{"communities", "dailies", "users"}

	for _, table := range requiredTables {
		var exists bool
		err := db.QueryRow(`
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

// ListAllTables and ShowTableStructure functions remain here
func ListAllTables(db *sql.DB) error {
	fmt.Println("\n📊 Listing all tables in the database:")
	rows, err := db.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
		ORDER BY table_name;
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var tableName string
	tableCount := 0
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			return err
		}
		fmt.Printf(" - %s\n", tableName)
		tableCount++
	}

	if tableCount == 0 {
		fmt.Println("No tables found in the database.")
	} else {
		fmt.Printf("Total tables: %d\n", tableCount)
	}
	return nil
}

func ShowTableStructure(db *sql.DB, tableName string) error {
	fmt.Printf("\n📋 Structure of '%s' table:\n", tableName)
	columns, err := db.Query(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns 
		WHERE table_name = $1 AND table_schema = 'public'
		ORDER BY ordinal_position;
	`, tableName)
	if err != nil {
		return err
	}
	defer columns.Close()

	var colName, dataType, nullable string
	fmt.Println("Column Name     | Data Type     | Nullable")
	fmt.Println("----------------|---------------|----------")
	for columns.Next() {
		if err := columns.Scan(&colName, &dataType, &nullable); err != nil {
			return err
		}
		fmt.Printf("%-15s | %-13s | %-8s\n", colName, dataType, nullable)
	}
	return nil
}
