package utils

// TODO: shift to protobufs ASAP

type NoteMetadata struct {
	Id           string `json:"note_id"`
	Name         string `json:"note_name"`
	CreatedAt    string `json:"created_at"`    // unix time
	LastModified string `json:"last_modified"` // unix time
}
