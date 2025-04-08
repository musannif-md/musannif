package db

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/masroof-maindak/musannif/internal/db/queries"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

const dbname string = "experiments.db"

var db *sql.DB

func InitTestDb() error {
	var err error

	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		return fmt.Errorf("failed to open test database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to verify test database connection: %w", err)
	}

	pragmas := []string{
		"PRAGMA foreign_keys=ON",
		"PRAGMA temp_store=MEMORY",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return fmt.Errorf("failed to set pragma '%s': %w", pragma, err)
		}
	}

	_, err = db.Exec(queries.SchemaCreationStatement)
	if err != nil {
		return fmt.Errorf("failed to create test schema: %w", err)
	}

	return nil
}

func CleanupTestDb() error {
	if db != nil {
		if err := db.Close(); err != nil {
			return fmt.Errorf("failed to close test database: %w", err)
		}
		db = nil
	}
	return nil
}

func InitDb(dir string) error {
	path := filepath.Join(dir, dbname)
	var err error

	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to verify database connection: %w", err)
	}

	pragmas := []string{
		"PRAGMA foreign_keys=ON",
		"PRAGMA temp_store=MEMORY",
		"PRAGMA mmap_size=4000000000", // 4 GB
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return fmt.Errorf("failed to set pragma '%s': %w", pragma, err)
		}
	}

	_, err = db.Exec(queries.SchemaCreationStatement)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

func LoginUser(username, password string) (string, error) {
	var (
		hashedPassword string
		salt           []byte
		role           string
	)

	err := db.QueryRow(queries.GetUserQuery, username).Scan(&role, &hashedPassword, &salt)
	if err != nil {
		return "", fmt.Errorf("database error: %w", err)
	}

	saltedPassword := append([]byte(password), salt...)
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), saltedPassword); err != nil {
		return "", fmt.Errorf("invalid password")
	}

	return role, nil
}

func SignupUser(username, password, role string) error {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	saltedPassword := append([]byte(password), salt...)
	hashedPassword, err := bcrypt.GenerateFromPassword(saltedPassword, bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	_, err = db.Exec(
		queries.InsertUserQuery,
		username, role, string(hashedPassword), salt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func CleanupDb() error {
	if db != nil {
		if err := db.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
	}
	return nil
}
