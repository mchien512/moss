package service

import (
	"context"
	"errors"
	entryApp "moss/go/internal/app"
	entrypb "moss/go/internal/genproto/entry"
	models "moss/go/internal/models/entry"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// service implements entrypb.EntryServiceServer
type service struct {
	entrypb.UnimplementedEntryServiceServer
	app entryApp.App
}

// NewService constructs a new gRPC service for entries.
func NewService(app entryApp.App) entrypb.EntryServiceServer {
	return &service{app: app}
}

// CreateEntry RPC must return *CreateEntryResponse
func (s *service) CreateEntry(ctx context.Context, req *entrypb.CreateEntryRequest) (*entrypb.CreateEntryResponse, error) {
	domainEntry := &models.Entry{
		UserID:      req.UserId,
		Title:       req.Title,
		Content:     req.Content,
		GrowthStage: models.GrowthStage(req.GrowthStage.String()),
	}

	created, err := s.app.CreateEntry(ctx, domainEntry)
	if err != nil {
		if errors.Is(entryApp.ErrInvalidEntry, err) {
			return nil, status.Error(codes.InvalidArgument, "invalid entry")
		}
		return nil, status.Errorf(codes.Internal, "failed to create entry: %v", err)
	}

	return &entrypb.CreateEntryResponse{
		Entry: toProtoEntry(created, 0),
	}, nil
}

// GetEntry RPC must return *GetEntryResponse
func (s *service) GetEntry(ctx context.Context, req *entrypb.GetEntryRequest) (*entrypb.GetEntryResponse, error) {
	// TODO: Replace with real userID from auth context
	userID := "user-id-from-auth"

	domainEntry, err := s.app.GetEntry(ctx, req.EntryId, userID)
	if err != nil {
		switch err {
		case entryApp.ErrUnauthorized:
			return nil, status.Error(codes.PermissionDenied, "unauthorized access")
		default:
			return nil, status.Errorf(codes.Internal, "failed to get entry: %v", err)
		}
	}

	return &entrypb.GetEntryResponse{
		Entry: toProtoEntry(domainEntry, 0),
	}, nil
}

// UpdateEntry RPC must return *UpdateEntryResponse
func (s *service) UpdateEntry(ctx context.Context, req *entrypb.UpdateEntryRequest) (*entrypb.UpdateEntryResponse, error) {
	domainEntry := &models.Entry{
		ID:          req.EntryId,
		Title:       req.Title,
		Content:     req.Content,
		GrowthStage: models.GrowthStage(req.GrowthStage.String()),
	}

	updated, err := s.app.UpdateEntry(ctx, domainEntry)
	if err != nil {
		switch err {
		case entryApp.ErrInvalidEntry:
			return nil, status.Error(codes.InvalidArgument, "invalid entry")
		case entryApp.ErrUnauthorized:
			return nil, status.Error(codes.PermissionDenied, "unauthorized access")
		default:
			return nil, status.Errorf(codes.Internal, "failed to update entry: %v", err)
		}
	}

	return &entrypb.UpdateEntryResponse{
		Entry: toProtoEntry(updated, 0),
	}, nil
}

// DeleteEntry RPC already returns emptypb.Empty
func (s *service) DeleteEntry(ctx context.Context, req *entrypb.DeleteEntryRequest) (*emptypb.Empty, error) {
	err := s.app.DeleteEntry(ctx, req.EntryId)
	if err != nil {
		switch err {
		case entryApp.ErrUnauthorized:
			return nil, status.Error(codes.PermissionDenied, "unauthorized access")
		default:
			return nil, status.Errorf(codes.Internal, "failed to delete entry: %v", err)
		}
	}
	return &emptypb.Empty{}, nil
}

// ListEntries RPC must return *ListEntriesResponse
func (s *service) ListEntries(ctx context.Context, req *entrypb.ListEntriesRequest) (*entrypb.ListEntriesResponse, error) {
	domainEntries, err := s.app.ListEntries(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list entries: %v", err)
	}

	protoEntries := make([]*entrypb.Entry, len(domainEntries))
	for i, e := range domainEntries {
		protoEntries[i] = toProtoEntry(e, 0)
	}

	return &entrypb.ListEntriesResponse{
		Entries: protoEntries,
	}, nil
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
