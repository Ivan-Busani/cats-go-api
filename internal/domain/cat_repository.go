package domain

import (
	"context"

	"cats-go-api/internal/model"
)

type CatRepository interface {
	FindAll(ctx context.Context) ([]model.Cat, error)
	FindByID(ctx context.Context, id int) (*model.Cat, error)
	FindByCatID(ctx context.Context, catID string) (*model.Cat, error)
	Create(ctx context.Context, input model.SaveCatInput, apiUsed string) (*model.Cat, error)
	Update(ctx context.Context, input model.SaveCatInput, apiUsed string, id int) (*model.Cat, error)
	Delete(ctx context.Context, id int) error
}
