package queries

const SchemaCreationStatement string = `
CREATE TABLE IF NOT EXISTS Users (
    username VARCHAR(255) PRIMARY KEY,
    role VARCHAR(255) NOT NULL,
    pw_hash VARCHAR(255) NOT NULL,
    salt BLOB NOT NULL
);`

const InsertUserQuery = `INSERT INTO Users (username, role, pw_hash, salt) VALUES (?, ?, ?, ?)`

const GetUserQuery = `SELECT role, pw_hash, salt FROM Users WHERE username = ?`

