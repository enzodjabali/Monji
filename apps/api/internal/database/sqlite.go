package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// DB is the global SQLite connection.
var DB *sql.DB

// InitSQLite initializes the SQLite database.
func InitSQLite(path string) {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}

	// Create users table.
	createUsersTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		company TEXT,
		password TEXT NOT NULL,
		role TEXT NOT NULL
	);`
	_, err = DB.Exec(createUsersTableSQL)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Create environments table.
	createEnvsTableSQL := `
	CREATE TABLE IF NOT EXISTS environments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		connection_string TEXT NOT NULL,
		created_by INTEGER NOT NULL
	);`
	_, err = DB.Exec(createEnvsTableSQL)
	if err != nil {
		log.Fatalf("Failed to create environments table: %v", err)
	}

	// Create user_env_permissions table:
	createUserEnvPerms := `
	CREATE TABLE IF NOT EXISTS user_env_permissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		environment_id INTEGER NOT NULL,
		permission TEXT NOT NULL, -- "none", "readOnly", "readAndWrite"
		UNIQUE (user_id, environment_id)
	);
	`
	_, err = DB.Exec(createUserEnvPerms)
	if err != nil {
		log.Fatalf("Failed to create user_env_permissions table: %v", err)
	}

	// Create user_db_permissions table:
	createUserDBPerms := `
	CREATE TABLE IF NOT EXISTS user_db_permissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		environment_id INTEGER NOT NULL,
		db_name TEXT NOT NULL,
		permission TEXT NOT NULL, -- "none", "readOnly", "readAndWrite"
		UNIQUE (user_id, environment_id, db_name)
	);
	`
	_, err = DB.Exec(createUserDBPerms)
	if err != nil {
		log.Fatalf("Failed to create user_db_permissions table: %v", err)
	}

	// Insert default admin user if none exist.
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to check users table: %v", err)
	}
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash default admin password: %v", err)
		}
		_, err = DB.Exec(
			`INSERT INTO users(first_name, last_name, email, company, password, role)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			"Admin", "User", "admin@example.com", "", string(hashedPassword), "admin")
		if err != nil {
			log.Fatalf("Failed to insert default admin user: %v", err)
		}
		log.Println("Default admin user created: email=admin@example.com, password=admin")
	}
}
