package entry

import (
	"lumo/internal/genproto/entry"
	models "lumo/internal/models/entry"
)

// ConvertProtoToDomain converts a proto entry to a domain model entry
func ConvertProtoToDomain(protoEntry *entry.Entry) *models.Entry {
	return &models.Entry{
		ID:      protoEntry.Id,
		UserID:  protoEntry.UserId,
		Title:   protoEntry.Title,
		Content: protoEntry.Content,
		Mood:    protoEntry.Mood,
	}
}

// ConvertDomainToProto converts a domain model entry to a proto entry
func ConvertDomainToProto(domainEntry *models.Entry) *entry.Entry {
	return &entry.Entry{
		Id:      domainEntry.ID,
		UserId:  domainEntry.UserID,
		Title:   domainEntry.Title,
		Content: domainEntry.Content,
		Mood:    domainEntry.Mood,
	}
}
