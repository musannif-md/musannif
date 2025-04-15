package queries

const SchemaCreationStatement string = `
CREATE TABLE IF NOT EXISTS Users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(255) NOT NULL,
    pw_hash blob NOT NULL,
    salt BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_username ON Users (username);

CREATE TABLE IF NOT EXISTS Notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at INTEGER DEFAULT (unixepoch()),
    last_modified INTEGER DEFAULT (unixepoch()),
	-- deleted BOOLEAN NOT NULL, -- TODO: For 'recently deleted' feature
    FOREIGN KEY (user_id) REFERENCES Users(id)
	UNIQUE (user_id, name)
);

CREATE INDEX IF NOT EXISTS idx_notes_user_id ON Notes (user_id);
`

const InsertUserQuery = `INSERT INTO Users (username, role, pw_hash, salt) VALUES (?, ?, ?, ?)`

const GetUserQuery = `SELECT role, pw_hash, salt FROM Users WHERE username = ?`

const InsertNoteQuery = `
INSERT INTO Notes (user_id, name) VALUES ((SELECT id FROM Users WHERE username = ?), ?)
`

// delete note entry given a username and note name
const DeleteNoteQuery = `
DELETE FROM Notes WHERE user_id = (SELECT id FROM Users WHERE username = ?) AND name = ?
`

const GetUsersNotesMetadata = `
SELECT n.id, n.name, n.created_at, n.last_modified from Notes n JOIN Users u on u.id = n.user_id
`

const UpdateNoteModificationTime = `
`
