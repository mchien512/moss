package models

import (
	"errors"
	"time"
)

// Entry represents a single journal entry in our application.
type Entry struct {
	ID        string    // UUID, unique identifier
	UserID    string    // UUID of the user who owns this entry
	Title     string    // Non-empty title of the entry
	Content   string    // The main text/body of the entry
	Mood      string    // User’s mood annotation (e.g. "happy", "sad")
	CreatedAt time.Time // Timestamp when the entry was first created
	UpdatedAt time.Time // Timestamp when the entry was last modified
}

// ErrInvalidEntry is returned when an Entry does not meet basic domain rules.
var ErrInvalidEntry = errors.New("invalid entry: missing required fields")

// Validate checks basic invariants on the Entry.
// For example, we require that UserID, Title, and Content are non-empty.
func (e *Entry) Validate() error {
	if e.UserID == "" {
		return errors.New("entry must have a UserID")
	}
	if e.Title == "" {
		return errors.New("entry must have a Title")
	}
	if e.Content == "" {
		return errors.New("entry must have Content")
	}
	return nil
}
