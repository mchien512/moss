package entry

import (
	"context"
	entryApp "lumo/go/internal/app/entry"
	entry2 "lumo/go/internal/genproto/entry"
	"lumo/go/internal/models/entry"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct {
	entry2.UnimplementedEntryServiceServer
	app entryApp.App
}

func NewService(app entryApp.App) entry2.EntryServiceServer {
	return &service{app: app}
}

func (s *service) CreateEntry(ctx context.Context, req *entry2.CreateEntryRequest) (*entry2.Entry, error) {
	domainEntry := &models.Entry{
		ID:        req.Id,
		UserID:    req.UserId,
		Title:     req.Title,
		Content:   req.Content,
		Mood:      req.Mood,
		CreatedAt: req.CreatedAt.AsTime(),
	}

	created, err := s.app.CreateEntry(ctx, domainEntry)
	if err != nil {
		if err == entryApp.ErrInvalidEntry {
			return nil, status.Error(codes.InvalidArgument, "invalid entry")
		}
		return nil, status.Errorf(codes.Internal, "failed to create entry: %v", err)
	}

	return &entry2.Entry{
		Id:        created.ID,
		UserId:    created.UserID,
		Title:     created.Title,
		Content:   created.Content,
		Mood:      created.Mood,
		CreatedAt: timestamppb.New(created.CreatedAt),
		UpdatedAt: timestamppb.New(created.UpdatedAt),
	}, nil
}

func (s *service) GetEntry(ctx context.Context, req *entry2.GetEntryRequest) (*entry2.Entry, error) {
	// In production, you would extract userID from context/auth token
	userID := "user-id-from-auth" // TODO: implement proper auth

	entry, err := s.app.GetEntry(ctx, req.Id, userID)
	if err != nil {
		switch err {
		case entryApp.ErrUnauthorized:
			return nil, status.Error(codes.PermissionDenied, "unauthorized access")
		default:
			return nil, status.Errorf(codes.Internal, "failed to get entry: %v", err)
		}
	}

	return &entry2.Entry{
		Id:        entry.ID,
		UserId:    entry.UserID,
		Title:     entry.Title,
		Content:   entry.Content,
		Mood:      entry.Mood,
		CreatedAt: timestamppb.New(entry.CreatedAt),
		UpdatedAt: timestamppb.New(entry.UpdatedAt),
	}, nil
}

func (s *service) UpdateEntry(ctx context.Context, req *entry2.UpdateEntryRequest) (*entry2.Entry, error) {
	domainEntry := &models.Entry{
		ID:        req.Id,
		UserID:    req.UserId,
		Title:     req.Title,
		Content:   req.Content,
		Mood:      req.Mood,
		UpdatedAt: req.UpdatedAt.AsTime(),
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

	return &entry2.Entry{
		Id:        updated.ID,
		UserId:    updated.UserID,
		Title:     updated.Title,
		Content:   updated.Content,
		Mood:      updated.Mood,
		CreatedAt: timestamppb.New(updated.CreatedAt),
		UpdatedAt: timestamppb.New(updated.UpdatedAt),
	}, nil
}

func (s *service) DeleteEntry(ctx context.Context, req *entry2.DeleteEntryRequest) (*emptypb.Empty, error) {
	err := s.app.DeleteEntry(ctx, req.Id, req.UserId)
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

func (s *service) ListEntries(ctx context.Context, req *entry2.ListEntriesRequest) (*entry2.ListEntriesResponse, error) {
	entries, err := s.app.ListEntries(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list entries: %v", err)
	}

	protoEntries := make([]*entry2.Entry, len(entries))
	for i, entry := range entries {
		protoEntries[i] = &entry2.Entry{
			Id:        entry.ID,
			UserId:    entry.UserID,
			Title:     entry.Title,
			Content:   entry.Content,
			Mood:      entry.Mood,
			CreatedAt: timestamppb.New(entry.CreatedAt),
			UpdatedAt: timestamppb.New(entry.UpdatedAt),
		}
	}

	return &entry2.ListEntriesResponse{
		Entries: protoEntries,
	}, nil
}

func (s *service) SyncEntries(ctx context.Context, req *entry2.SyncEntriesRequest) (*entry2.ListEntriesResponse, error) {
	entries, err := s.app.SyncEntries(ctx, req.UserId, req.Since.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to sync entries: %v", err)
	}

	protoEntries := make([]*entry2.Entry, len(entries))
	for i, entry := range entries {
		protoEntries[i] = &entry2.Entry{
			Id:        entry.ID,
			UserId:    entry.UserID,
			Title:     entry.Title,
			Content:   entry.Content,
			Mood:      entry.Mood,
			CreatedAt: timestamppb.New(entry.CreatedAt),
			UpdatedAt: timestamppb.New(entry.UpdatedAt),
		}
	}

	return &entry2.ListEntriesResponse{
		Entries: protoEntries,
	}, nil
}
