package link

import (
	"context"
	"database/sql"
	"errors"
	"time"

	models "moss/go/internal/models/link"
	db "moss/go/internal/repository/db/sqlc"
)

var (
	ErrLinkNotFound = errors.New("link not found")
	ErrInvalidLink  = errors.New("invalid link")
)

type Repository interface {
	CreateEntryLink(ctx context.Context, l *models.Link) (*models.Link, error)
	DeleteEntryLink(ctx context.Context, sourceID, targetID string) error
	ListBySource(ctx context.Context, sourceID string) ([]*models.Link, error)
	ListByTarget(ctx context.Context, targetID string) ([]*models.Link, error)
	CountBySource(ctx context.Context, sourceID string) (int64, error)
	CountByTarget(ctx context.Context, targetID string) (int64, error)
}

type repository struct {
	queries *db.Queries
}

func NewRepository(dbConn *sql.DB) Repository {
	return &repository{
		queries: db.New(dbConn),
	}
}

func (r *repository) CreateEntryLink(ctx context.Context, l *models.Link) (*models.Link, error) {
	if err := l.Validate(); err != nil {
		return nil, ErrInvalidLink
	}
	now := time.Now().UTC()
	l.CreatedAt = now

	created, err := r.queries.CreateEntryLink(ctx, db.CreateEntryLinkParams{
		SourceEntryID: l.SourceEntryID,
		TargetEntryID: l.TargetEntryID,
		UserID:        l.UserID,
	})
	if err != nil {
		return nil, err
	}

	return fromDBLink(created), nil
}

func (r *repository) DeleteEntryLink(ctx context.Context, sourceID, targetID string) error {
	err := r.queries.DeleteEntryLink(ctx, db.DeleteEntryLinkParams{
		SourceEntryID: sourceID,
		TargetEntryID: targetID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrLinkNotFound
		}
		return err
	}
	return nil
}

func (r *repository) ListBySource(ctx context.Context, sourceID string) ([]*models.Link, error) {
	dbLinks, err := r.queries.ListLinksBySource(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	return fromDBLinks(dbLinks), nil
}

func (r *repository) ListByTarget(ctx context.Context, targetID string) ([]*models.Link, error) {
	dbLinks, err := r.queries.ListLinksByTarget(ctx, targetID)
	if err != nil {
		return nil, err
	}
	return fromDBLinks(dbLinks), nil
}

func (r *repository) CountBySource(ctx context.Context, sourceID string) (int64, error) {
	count, err := r.queries.CountLinksBySource(ctx, sourceID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) CountByTarget(ctx context.Context, targetID string) (int64, error) {
	count, err := r.queries.CountLinksByTarget(ctx, targetID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// fromDBLink converts a SQLC EntryLink row into a domain Link.
func fromDBLink(dbLink db.EntryLink) *models.Link {
	return &models.Link{
		SourceEntryID: dbLink.SourceEntryID,
		TargetEntryID: dbLink.TargetEntryID,
		UserID:        dbLink.UserID,
		CreatedAt:     dbLink.CreatedAt,
	}
}

// fromDBLinks converts a slice of SQLC EntryLink rows into []*models.Link.
func fromDBLinks(dbLinks []db.EntryLink) []*models.Link {
	result := make([]*models.Link, len(dbLinks))
	for i, l := range dbLinks {
		result[i] = fromDBLink(l)
	}
	return result
}
