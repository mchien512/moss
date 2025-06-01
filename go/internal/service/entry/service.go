package entry

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	entryApp "moss/go/internal/app/entry"
	entrypb "moss/go/internal/genproto/protobuf/entry"
	models "moss/go/internal/models/entry"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Service implements the EntryServiceHandler interface
type Service struct {
	app entryApp.App
}

// NewService constructs a new Connect service for entries.
func NewService(app entryApp.App) *Service {
	return &Service{app: app}
}

// CreateEntry implements the EntryServiceHandler interface
func (s *Service) CreateEntry(ctx context.Context, req *connect.Request[entrypb.CreateEntryRequest]) (*connect.Response[entrypb.CreateEntryResponse], error) {
	domainEntry := &models.Entry{
		UserID:      req.Msg.UserId,
		Title:       req.Msg.Title,
		Content:     req.Msg.Content,
		GrowthStage: models.GrowthStage(req.Msg.GrowthStage.String()),
	}

	created, err := s.app.CreateEntry(ctx, domainEntry)
	if err != nil {
		if errors.Is(entryApp.ErrInvalidEntry, err) {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid entry"))
		}
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create entry: %w", err))
	}

	return connect.NewResponse(&entrypb.CreateEntryResponse{
		Entry: toProtoEntry(created, 0),
	}), nil
}

// GetEntry implements the EntryServiceHandler interface
func (s *Service) GetEntry(ctx context.Context, req *connect.Request[entrypb.GetEntryRequest]) (*connect.Response[entrypb.GetEntryResponse], error) {
	// TODO: Replace with real userID from auth context
	userID := "user-id-from-auth"

	domainEntry, err := s.app.GetEntry(ctx, req.Msg.EntryId, userID)
	if err != nil {
		switch err {
		case entryApp.ErrUnauthorized:
			return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("unauthorized access"))
		default:
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get entry: %w", err))
		}
	}

	return connect.NewResponse(&entrypb.GetEntryResponse{
		Entry: toProtoEntry(domainEntry, 0),
	}), nil
}

// UpdateEntry implements the EntryServiceHandler interface
func (s *Service) UpdateEntry(ctx context.Context, req *connect.Request[entrypb.UpdateEntryRequest]) (*connect.Response[entrypb.UpdateEntryResponse], error) {
	domainEntry := &models.Entry{
		ID:          req.Msg.EntryId,
		Title:       req.Msg.Title,
		Content:     req.Msg.Content,
		GrowthStage: models.GrowthStage(req.Msg.GrowthStage.String()),
	}

	updated, err := s.app.UpdateEntry(ctx, domainEntry)
	if err != nil {
		switch err {
		case entryApp.ErrInvalidEntry:
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid entry"))
		case entryApp.ErrUnauthorized:
			return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("unauthorized access"))
		default:
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update entry: %w", err))
		}
	}

	return connect.NewResponse(&entrypb.UpdateEntryResponse{
		Entry: toProtoEntry(updated, 0),
	}), nil
}

// DeleteEntry implements the EntryServiceHandler interface
func (s *Service) DeleteEntry(ctx context.Context, req *connect.Request[entrypb.DeleteEntryRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.app.DeleteEntry(ctx, req.Msg.EntryId)
	if err != nil {
		switch err {
		case entryApp.ErrUnauthorized:
			return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("unauthorized access"))
		default:
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to delete entry: %w", err))
		}
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

// ListEntries implements the EntryServiceHandler interface
func (s *Service) ListEntries(ctx context.Context, req *connect.Request[entrypb.ListEntriesRequest]) (*connect.Response[entrypb.ListEntriesResponse], error) {
	domainEntries, err := s.app.ListEntries(ctx, req.Msg.UserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to list entries: %w", err))
	}

	protoEntries := make([]*entrypb.Entry, len(domainEntries))
	for i, e := range domainEntries {
		protoEntries[i] = toProtoEntry(e, 0)
	}

	return connect.NewResponse(&entrypb.ListEntriesResponse{
		Entries: protoEntries,
	}), nil
}

// toProtoEntry converts a domain Entry into a proto Entry, injecting linkCount separately.
// Replace hardcoded 0 with actual lookup when you add entry_links support.
func toProtoEntry(domain *models.Entry, linkCount int) *entrypb.Entry {
	return &entrypb.Entry{
		Id:          domain.ID,
		UserId:      domain.UserID,
		Title:       domain.Title,
		Content:     domain.Content,
		CreatedAt:   timestamppb.New(domain.CreatedAt),
		UpdatedAt:   timestamppb.New(domain.UpdatedAt),
		GrowthStage: entrypb.GrowthStage(entrypb.GrowthStage_value[string(domain.GrowthStage)]),
		LinkCount:   int32(linkCount),
	}
}
