package repository

import (
	"context"

	"monolith/internal/domains/geolocation/domain/entity"
)

type DestionationRepository interface {
	SetDestinationPoint(context.Context, *entity.FutureDestinationPoint)
}
