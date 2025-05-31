package entry

import (
	"context"
	"database/sql"
	"errors"
	"time"

	db "lumo/internal/db/sqlc"
	models "lumo/internal/models/entry"
)

var ErrEntryNotFound = errors.New("entry not found")

type Repository interface {
	Create(ctx context.Context, e *models.Entry) (*models.Entry, error)
	GetByID(ctx context.Context, id string) (*models.Entry, error)
	ListByUser(ctx context.Context, userID string) ([]*models.Entry, error)
	Update(ctx context.Context, e *models.Entry) (*models.Entry, error)
	Delete(ctx context.Context, id string) error
	ListByUserSince(ctx context.Context, userID string, since time.Time) ([]*models.Entry, error)
}

type repository struct {
	queries *db.Queries
}

func NewRepository(dbConn *sql.DB) Repository {
	return &repository{
		queries: db.New(dbConn),
	}
}

func (r *repository) Create(ctx context.Context, e *models.Entry) (*models.Entry, error) {
	now := time.Now().UTC()
	e.CreatedAt = now
	e.UpdatedAt = now

	entry, err := r.queries.CreateEntry(ctx, db.CreateEntryParams{
		ID:        e.ID,
		UserID:    e.UserID,
		Title:     e.Title,
		Content:   e.Content,
		Mood:      e.Mood,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	})
	if err != nil {
		return nil, err
	}

	return fromDBEntry(entry), nil
}

func (r *repository) GetByID(ctx context.Context, id string) (*models.Entry, error) {
	entry, err := r.queries.GetEntryByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEntryNotFound
		}
		return nil, err
	}
	return fromDBEntry(entry), nil
}

func (r *repository) ListByUser(ctx context.Context, userID string) ([]*models.Entry, error) {
	entries, err := r.queries.ListEntriesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return fromDBEntries(entries), nil
}

func (r *repository) Update(ctx context.Context, e *models.Entry) (*models.Entry, error) {
	now := time.Now().UTC()
	e.UpdatedAt = now

	entry, err := r.queries.UpdateEntry(ctx, db.UpdateEntryParams{
		ID:        e.ID,
		Title:     e.Title,
		Content:   e.Content,
		Mood:      e.Mood,
		UpdatedAt: e.UpdatedAt,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEntryNotFound
		}
		return nil, err
	}
	return fromDBEntry(entry), nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.queries.DeleteEntry(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEntryNotFound
		}
		return err
	}
	return nil
}

func (r *repository) ListByUserSince(ctx context.Context, userID string, since time.Time) ([]*models.Entry, error) {
	entries, err := r.queries.ListEntriesByUserSince(ctx, db.ListEntriesByUserSinceParams{
		UserID:    userID,
		UpdatedAt: since,
	})
	if err != nil {
		return nil, err
	}
	return fromDBEntries(entries), nil
}

// Helper functions to convert between SQLC and domain models
func fromDBEntry(e db.Entry) *models.Entry {
	return &models.Entry{
		ID:        e.ID,
		UserID:    e.UserID,
		Title:     e.Title,
		Content:   e.Content,
		Mood:      e.Mood,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func fromDBEntries(entries []db.Entry) []*models.Entry {
	result := make([]*models.Entry, len(entries))
	for i, e := range entries {
		result[i] = fromDBEntry(e)
	}
	return result
}
