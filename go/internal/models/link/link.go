package link

import (
	"errors"
	"time"
)

var (
	ErrInvalidLink = errors.New("invalid link: missing required fields")
)

// Link represents a directional connection between two Entries,
// as stored in the `entry_links` table.
type Link struct {
	SourceEntryID string    // UUID of the entry that points to TargetEntryID
	TargetEntryID string    // UUID of the entry being pointed to
	UserID        string    // UUID of the user who created/owns this link
	CreatedAt     time.Time // Timestamp when the link was created
}

// Validate ensures the Link has all required fields.
// In particular, SourceEntryID, TargetEntryID, and UserID must be non-empty.
func (l *Link) Validate() error {
	if l.SourceEntryID == "" {
		return errors.New("link must have a SourceEntryID")
	}
	if l.TargetEntryID == "" {
		return errors.New("link must have a TargetEntryID")
	}
	if l.UserID == "" {
		return errors.New("link must have a UserID")
	}
	return nil
}
