package db

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/musannif-md/musannif/internal/db/queries"
	"github.com/musannif-md/musannif/internal/utils"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const dbname string = "musannif.db"

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

func CleanupDb() error {
	if db != nil {
		if err := db.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
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

func CreateNote(username, notename string) (int64, error) {
	result, err := db.Exec(queries.InsertNoteQuery, username, notename)

	// TODO if err type is because a note w/ the same name exists already,
	// throw a different error

	if err != nil {
		return 0, fmt.Errorf("failed to create note: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last inserted id: %w", err)
	}

	return id, nil
}

func DeleteNote(username, notename string) error {
	_, err := db.Exec(queries.DeleteNoteQuery, username, notename)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	return nil
}

func GetUserNoteMetadata(username string) ([]utils.NoteMetadata, error) {
	rows, err := db.Query(queries.GetUsersNotesMetadata, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata for user's notes: %w", err)
	}
	defer rows.Close()

	var noteListMd []utils.NoteMetadata

	for rows.Next() {
		var noteId int64
		var notename string
		var createdAt int64
		var lastModified int64

		err = rows.Scan(&noteId, &notename, &createdAt, &lastModified)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row to NoteMetadata obj: %w", err)
		}

		// We're doing this to send strings over JSON as opposed to integers
		// Since JavaScript doesn't natively support 64-bit wide integers...

		md := utils.NoteMetadata{
			Id:           strconv.FormatInt(noteId, 10),
			Name:         notename,
			CreatedAt:    strconv.FormatInt(createdAt, 10),
			LastModified: strconv.FormatInt(lastModified, 10),
		}

		noteListMd = append(noteListMd, md)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %w", err)
	}

	return noteListMd, nil
}
