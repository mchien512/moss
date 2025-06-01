package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"moss/go/internal/entry/repository/db/sqlc"
	models "moss/go/internal/models/entry"
)

var ErrEntryNotFound = errors.New("entry not found")

type Repository interface {
	Create(ctx context.Context, e *models.Entry) (*models.Entry, error)
	GetByID(ctx context.Context, id string) (*models.Entry, error)
	ListByUser(ctx context.Context, userID string) ([]*models.Entry, error)
	ListByUserSince(ctx context.Context, userID string, since time.Time) ([]*models.Entry, error)
	Update(ctx context.Context, e *models.Entry) (*models.Entry, error)
	Delete(ctx context.Context, id string) error
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
		ID:          e.ID,
		UserID:      e.UserID,
		Title:       e.Title,
		Content:     e.Content,
		GrowthStage: string(e.GrowthStage),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	})
	if err != nil {
		return nil, err
	}

	return fromDBEntry(entry), nil
}

func (r *repository) GetByID(ctx context.Context, id string) (*models.Entry, error) {
	dbEntry, err := r.queries.GetEntryByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEntryNotFound
		}
		return nil, err
	}

	return fromDBEntry(dbEntry), nil
}

func (r *repository) ListByUser(ctx context.Context, userID string) ([]*models.Entry, error) {
	dbEntries, err := r.queries.ListEntriesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return fromDBEntries(dbEntries), nil
}

func (r *repository) ListByUserSince(ctx context.Context, userID string, since time.Time) ([]*models.Entry, error) {
	dbEntries, err := r.queries.ListEntriesByUserSince(ctx, db.ListEntriesByUserSinceParams{
		UserID:    userID,
		UpdatedAt: since,
	})
	if err != nil {
		return nil, err
	}

	return fromDBEntries(dbEntries), nil
}

func (r *repository) Update(ctx context.Context, e *models.Entry) (*models.Entry, error) {
	e.UpdatedAt = time.Now().UTC()

	entry, err := r.queries.UpdateEntry(ctx, db.UpdateEntryParams{
		ID:          e.ID,
		Title:       e.Title,
		Content:     e.Content,
		GrowthStage: string(e.GrowthStage),
		UpdatedAt:   e.UpdatedAt,
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

func fromDBEntry(dbEntry db.Entry) *models.Entry {
	return &models.Entry{
		ID:          dbEntry.ID,
		UserID:      dbEntry.UserID,
		Title:       dbEntry.Title,
		Content:     dbEntry.Content,
		GrowthStage: models.GrowthStage(dbEntry.GrowthStage),
		CreatedAt:   dbEntry.CreatedAt,
		UpdatedAt:   dbEntry.UpdatedAt,
	}
}

func fromDBEntries(dbEntries []db.Entry) []*models.Entry {
	entries := make([]*models.Entry, len(dbEntries))
	for i, dbEntry := range dbEntries {
		entries[i] = fromDBEntry(dbEntry)
	}
	return entries
}
