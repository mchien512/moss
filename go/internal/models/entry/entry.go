package models

import (
	"errors"
	"time"
)

type GrowthStage string

const (
	GrowthStageSeed      GrowthStage = "seed"
	GrowthStageSprout    GrowthStage = "sprout"
	GrowthStageBloom     GrowthStage = "bloom"
	GrowthStageEvergreen GrowthStage = "evergreen"
)

// Entry represents a single digital garden entry in the Moss app.
type Entry struct {
	ID          string      // UUID, unique identifier
	UserID      string      // UUID of the user who owns this entry
	Title       string      // Required title of the entry
	Content     string      // Markdown content
	GrowthStage GrowthStage // Lifecycle stage of the entry
	CreatedAt   time.Time   // Timestamp of creation
	UpdatedAt   time.Time   // Timestamp of last update
}

var ErrInvalidEntry = errors.New("invalid entry: missing required fields")

func (e *Entry) Validate() error {
	if e.Title == "" {
		return errors.New("entry must have a Title")
	}
	if e.Content == "" {
		return errors.New("entry must have Content")
	}
	switch e.GrowthStage {
	case GrowthStageSeed, GrowthStageSprout, GrowthStageBloom, GrowthStageEvergreen:
		// valid
	default:
		return errors.New("entry has invalid GrowthStage")
	}
	return nil
}
