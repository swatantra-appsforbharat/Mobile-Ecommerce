// databasetup.go
package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect establishes a connection to the PostgreSQL database.
func Connect() {
	var err error
	connStr := "host=localhost port=5432 user=postgres password=yourpassword dbname=mobilecommerce sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Connected to PostgreSQL database")
	createTables()
	applyMigrations()
}

// createTables creates the required tables if they don't exist.
func createTables() {
	userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        first_name TEXT,
        last_name TEXT,
        email TEXT UNIQUE,
        password TEXT,
        role TEXT DEFAULT 'user'
    )`

	productTable := `
    CREATE TABLE IF NOT EXISTS products (
        id SERIAL PRIMARY KEY,
        name TEXT,
        description TEXT,
        price NUMERIC,
        quantity INT
    )`

	cartTable := `
    CREATE TABLE IF NOT EXISTS cart_items (
        id SERIAL PRIMARY KEY,
        user_id INT REFERENCES users(id) ON DELETE CASCADE,
        product_id INT REFERENCES products(id) ON DELETE CASCADE,
        quantity INT
    )`

	addressTable := `
    CREATE TABLE IF NOT EXISTS addresses (
        id SERIAL PRIMARY KEY,
        user_id INT REFERENCES users(id) ON DELETE CASCADE,
        street TEXT,
        city TEXT,
        state TEXT,
        zip_code TEXT,
        country TEXT
    )`

	orderTable := `
CREATE TABLE IF NOT EXISTS orders (
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(id) ON DELETE CASCADE,
	status TEXT
)`

	orderItemTable := `
CREATE TABLE IF NOT EXISTS order_items (
	id SERIAL PRIMARY KEY,
	order_id INT REFERENCES orders(id) ON DELETE CASCADE,
	product_id INT REFERENCES products(id),
	quantity INT
)`

	tableQueries := []string{userTable, productTable, cartTable, addressTable, orderTable, orderItemTable}

	for _, query := range tableQueries {
		_, err := DB.Exec(query)
		if err != nil {
			panic(fmt.Sprintf("Failed to create table: %v", err))
		}
	}

	fmt.Println("✅ Tables checked/created successfully")
}

func applyMigrations() {
	migrationFile := "database/migrations/001_create_indexes.sql"

	content, err := os.ReadFile(migrationFile)
	if err != nil {
		fmt.Printf("❌ Failed to read migration file: %v\n", err)
		return
	}

	_, err = DB.Exec(string(content))
	if err != nil {
		fmt.Printf("❌ Failed to execute migration: %v\n", err)
		return
	}

	fmt.Println("✅ Index migration applied")
}
