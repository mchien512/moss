package service

import (
	entrypb "moss/go/internal/genproto/entry"
	models "moss/go/internal/models/entry"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ConvertProtoToDomain converts a proto Entry to a domain Entry.
// (Note: we no longer map LinkedEntryIDs here.)
func ConvertProtoToDomain(protoEntry *entrypb.Entry) *models.Entry {
	return &models.Entry{
		ID:          protoEntry.Id,
		UserID:      protoEntry.UserId,
		Title:       protoEntry.Title,
		Content:     protoEntry.Content,
		GrowthStage: models.GrowthStage(protoEntry.GrowthStage.String()),
		CreatedAt:   protoEntry.CreatedAt.AsTime(),
		UpdatedAt:   protoEntry.UpdatedAt.AsTime(),
	}
}

// ConvertDomainToProto converts a domain Entry to a proto Entry.
// (Again, LinkedEntryIDs is omitted because it’s not on our proto)
func ConvertDomainToProto(domainEntry *models.Entry) *entrypb.Entry {
	return &entrypb.Entry{
		Id:          domainEntry.ID,
		UserId:      domainEntry.UserID,
		Title:       domainEntry.Title,
		Content:     domainEntry.Content,
		GrowthStage: entrypb.GrowthStage(entrypb.GrowthStage_value[string(domainEntry.GrowthStage)]),
		CreatedAt:   timestamppb.New(domainEntry.CreatedAt),
		UpdatedAt:   timestamppb.New(domainEntry.UpdatedAt),
	}
}
