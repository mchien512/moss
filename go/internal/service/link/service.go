package link

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	linkApp "moss/go/internal/app/link"
	linkpb "moss/go/internal/genproto/protobuf/link"
	models "moss/go/internal/models/link"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	app linkApp.App
}

func NewService(app linkApp.App) *Service {
	return &Service{app: app}
}

// CreateLink implements the LinkServiceHandler interface
func (s *Service) CreateLink(ctx context.Context, req *connect.Request[linkpb.CreateLinkRequest]) (*connect.Response[linkpb.CreateLinkResponse], error) {
	domainLink := &models.Link{
		SourceEntryID: req.Msg.SourceEntryId,
		TargetEntryID: req.Msg.TargetEntryId,
		UserID:        req.Msg.UserId,
		// CreatedAt is set by the repository/app, so we ignore req.CreatedAt if present.
	}

	created, err := s.app.CreateLink(ctx, domainLink)
	if err != nil {
		switch err {
		case linkApp.ErrInvalidLink:
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid link"))
		default:
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create link: %w", err))
		}
	}

	return connect.NewResponse(&linkpb.CreateLinkResponse{
		Link: toProtoLink(created),
	}), nil
}

// DeleteLink implements the LinkServiceHandler interface
func (s *Service) DeleteLink(ctx context.Context, req *connect.Request[linkpb.DeleteLinkRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.app.DeleteLink(ctx, req.Msg.SourceEntryId, req.Msg.TargetEntryId)
	if err != nil {
		switch err {
		case linkApp.ErrUnauthorized:
			return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("unauthorized access"))
		case linkApp.ErrInvalidLink:
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("link not found"))
		default:
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to delete link: %w", err))
		}
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

// ListLinksBySource implements the LinkServiceHandler interface
func (s *Service) ListLinksBySource(ctx context.Context, req *connect.Request[linkpb.ListLinksBySourceRequest]) (*connect.Response[linkpb.ListLinksBySourceResponse], error) {
	links, err := s.app.ListLinksBySource(ctx, req.Msg.SourceEntryId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to list links by source: %w", err))
	}

	protoLinks := make([]*linkpb.Link, len(links))
	for i, l := range links {
		protoLinks[i] = toProtoLink(l)
	}
	return connect.NewResponse(&linkpb.ListLinksBySourceResponse{Links: protoLinks}), nil
}

// ListLinksByTarget implements the LinkServiceHandler interface
func (s *Service) ListLinksByTarget(ctx context.Context, req *connect.Request[linkpb.ListLinksByTargetRequest]) (*connect.Response[linkpb.ListLinksByTargetResponse], error) {
	links, err := s.app.ListLinksByTarget(ctx, req.Msg.TargetEntryId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to list links by target: %w", err))
	}

	protoLinks := make([]*linkpb.Link, len(links))
	for i, l := range links {
		protoLinks[i] = toProtoLink(l)
	}
	return connect.NewResponse(&linkpb.ListLinksByTargetResponse{Links: protoLinks}), nil
}

// CountLinksBySource implements the LinkServiceHandler interface
func (s *Service) CountLinksBySource(ctx context.Context, req *connect.Request[linkpb.CountLinksBySourceRequest]) (*connect.Response[linkpb.CountLinksBySourceResponse], error) {
	count, err := s.app.CountLinksBySource(ctx, req.Msg.SourceEntryId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to count links by source: %w", err))
	}
	return connect.NewResponse(&linkpb.CountLinksBySourceResponse{Count: count}), nil
}

// CountLinksByTarget implements the LinkServiceHandler interface
func (s *Service) CountLinksByTarget(ctx context.Context, req *connect.Request[linkpb.CountLinksByTargetRequest]) (*connect.Response[linkpb.CountLinksByTargetResponse], error) {
	count, err := s.app.CountLinksByTarget(ctx, req.Msg.TargetEntryId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to count links by target: %w", err))
	}
	return connect.NewResponse(&linkpb.CountLinksByTargetResponse{Count: count}), nil
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
