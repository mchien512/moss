package link

import (
	"context"
	"errors"

	models "moss/go/internal/models/link"
	linkRepo "moss/go/internal/repository/link"
)

var (
	ErrInvalidLink  = errors.New("invalid link")
	ErrUnauthorized = errors.New("unauthorized access")
)

type App interface {
	CreateLink(ctx context.Context, l *models.Link) (*models.Link, error)
	DeleteLink(ctx context.Context, sourceID string, targetID string) error
	ListLinksBySource(ctx context.Context, sourceID string) ([]*models.Link, error)
	ListLinksByTarget(ctx context.Context, targetID string) ([]*models.Link, error)
	CountLinksBySource(ctx context.Context, sourceID string) (int64, error)
	CountLinksByTarget(ctx context.Context, targetID string) (int64, error)
}

type app struct {
	repo linkRepo.Repository
}

func NewApp(repo linkRepo.Repository) App {
	return &app{repo: repo}
}

// CreateLink validates and creates a new Link.
func (a *app) CreateLink(ctx context.Context, l *models.Link) (*models.Link, error) {
	if err := l.Validate(); err != nil {
		return nil, ErrInvalidLink
	}
	return a.repo.CreateEntryLink(ctx, l)
}

// DeleteLink checks ownership then deletes the Link.
func (a *app) DeleteLink(ctx context.Context, sourceID string, targetID string) error {
	// Ensure that the link exists and that userID matches.
	_, err := a.repo.ListBySource(ctx, sourceID)
	if err != nil {
		return err
	}
	return a.repo.DeleteEntryLink(ctx, sourceID, targetID)
}

func (a *app) ListLinksBySource(ctx context.Context, sourceID string) ([]*models.Link, error) {
	return a.repo.ListBySource(ctx, sourceID)
}

func (a *app) ListLinksByTarget(ctx context.Context, targetID string) ([]*models.Link, error) {
	return a.repo.ListByTarget(ctx, targetID)
}

func (a *app) CountLinksBySource(ctx context.Context, sourceID string) (int64, error) {
	return a.repo.CountBySource(ctx, sourceID)
}

func (a *app) CountLinksByTarget(ctx context.Context, targetID string) (int64, error) {
	return a.repo.CountByTarget(ctx, targetID)
}
