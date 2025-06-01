package link

import (
	"context"

	linkApp "moss/go/internal/app/link"
	linkpb "moss/go/internal/genproto/link"
	models "moss/go/internal/models/link"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct {
	linkpb.UnimplementedLinkServiceServer
	app linkApp.App
}

func NewService(app linkApp.App) linkpb.LinkServiceServer {
	return &service{app: app}
}

// CreateLink RPC → converts request → domain model → calls app → returns proto Link.
func (s *service) CreateLink(ctx context.Context, req *linkpb.CreateLinkRequest) (*linkpb.CreateLinkResponse, error) {
	domainLink := &models.Link{
		SourceEntryID: req.SourceEntryId,
		TargetEntryID: req.TargetEntryId,
		UserID:        req.UserId,
		// CreatedAt is set by the repository/app, so we ignore req.CreatedAt if present.
	}

	created, err := s.app.CreateLink(ctx, domainLink)
	if err != nil {
		switch err {
		case linkApp.ErrInvalidLink:
			return nil, status.Error(codes.InvalidArgument, "invalid link")
		default:
			return nil, status.Errorf(codes.Internal, "failed to create link: %v", err)
		}
	}

	return &linkpb.CreateLinkResponse{
		Link: toProtoLink(created),
	}, nil
}

// DeleteLink RPC → calls app.DeleteLink and wraps errors.
func (s *service) DeleteLink(ctx context.Context, req *linkpb.DeleteLinkRequest) (*emptypb.Empty, error) {
	err := s.app.DeleteLink(ctx, req.SourceEntryId, req.TargetEntryId)
	if err != nil {
		switch err {
		case linkApp.ErrUnauthorized:
			return nil, status.Error(codes.PermissionDenied, "unauthorized access")
		case linkApp.ErrInvalidLink:
			return nil, status.Error(codes.NotFound, "link not found")
		default:
			return nil, status.Errorf(codes.Internal, "failed to delete link: %v", err)
		}
	}
	return &emptypb.Empty{}, nil
}

// ListLinksBySource RPC → returns a list of Link messages.
func (s *service) ListLinksBySource(ctx context.Context, req *linkpb.ListLinksBySourceRequest) (*linkpb.ListLinksBySourceResponse, error) {
	links, err := s.app.ListLinksBySource(ctx, req.SourceEntryId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list links by source: %v", err)
	}

	protoLinks := make([]*linkpb.Link, len(links))
	for i, l := range links {
		protoLinks[i] = toProtoLink(l)
	}
	return &linkpb.ListLinksBySourceResponse{Links: protoLinks}, nil
}

// ListLinksByTarget RPC → returns a list of Link messages.
func (s *service) ListLinksByTarget(ctx context.Context, req *linkpb.ListLinksByTargetRequest) (*linkpb.ListLinksByTargetResponse, error) {
	links, err := s.app.ListLinksByTarget(ctx, req.TargetEntryId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list links by target: %v", err)
	}

	protoLinks := make([]*linkpb.Link, len(links))
	for i, l := range links {
		protoLinks[i] = toProtoLink(l)
	}
	return &linkpb.ListLinksByTargetResponse{Links: protoLinks}, nil
}

// CountLinksBySource RPC → returns the count of outgoing links.
func (s *service) CountLinksBySource(ctx context.Context, req *linkpb.CountLinksBySourceRequest) (*linkpb.CountLinksBySourceResponse, error) {
	count, err := s.app.CountLinksBySource(ctx, req.SourceEntryId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count links by source: %v", err)
	}
	return &linkpb.CountLinksBySourceResponse{Count: count}, nil
}

// CountLinksByTarget RPC → returns the count of incoming links.
func (s *service) CountLinksByTarget(ctx context.Context, req *linkpb.CountLinksByTargetRequest) (*linkpb.CountLinksByTargetResponse, error) {
	count, err := s.app.CountLinksByTarget(ctx, req.TargetEntryId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count links by target: %v", err)
	}
	return &linkpb.CountLinksByTargetResponse{Count: count}, nil
}

// toProtoLink converts a domain Link into a proto Link.
func toProtoLink(domain *models.Link) *linkpb.Link {
	return &linkpb.Link{
		SourceEntryId: domain.SourceEntryID,
		TargetEntryId: domain.TargetEntryID,
		UserId:        domain.UserID,
		CreatedAt:     timestamppb.New(domain.CreatedAt),
	}
}
