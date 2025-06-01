package entry

import (
	"context"
	"errors"
	models "moss/go/internal/models/entry"
	entryRepo "moss/go/internal/repository/entry"
)

var (
	ErrInvalidEntry = errors.New("invalid entry")
	ErrUnauthorized = errors.New("unauthorized access")
)

type App interface {
	CreateEntry(ctx context.Context, entry *models.Entry) (*models.Entry, error)
	GetEntry(ctx context.Context, id string, userID string) (*models.Entry, error)
	UpdateEntry(ctx context.Context, entry *models.Entry) (*models.Entry, error)
	DeleteEntry(ctx context.Context, id string) error
	ListEntries(ctx context.Context, userID string) ([]*models.Entry, error)
}

type app struct {
	repo entryRepo.Repository
}

func NewApp(repo entryRepo.Repository) App {
	return &app{repo: repo}
}

func (a *app) CreateEntry(ctx context.Context, entry *models.Entry) (*models.Entry, error) {
	if err := entry.Validate(); err != nil {
		return nil, ErrInvalidEntry
	}

	return a.repo.Create(ctx, entry)
}

func (a *app) GetEntry(ctx context.Context, id string, userID string) (*models.Entry, error) {
	entry, err := a.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if entry.UserID != userID {
		return nil, ErrUnauthorized
	}

	return entry, nil
}

func (a *app) UpdateEntry(ctx context.Context, entry *models.Entry) (*models.Entry, error) {
	if err := entry.Validate(); err != nil {
		return nil, ErrInvalidEntry
	}

	existing, err := a.repo.GetByID(ctx, entry.ID)
	if err != nil {
		return nil, err
	}

	if existing.UserID != entry.UserID {
		return nil, ErrUnauthorized
	}

	return a.repo.Update(ctx, entry)
}

func (a *app) DeleteEntry(ctx context.Context, id string) error {
	_, err := a.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return a.repo.Delete(ctx, id)
}

func (a *app) ListEntries(ctx context.Context, userID string) ([]*models.Entry, error) {
	return a.repo.ListByUser(ctx, userID)
}
