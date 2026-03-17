package service

import (
	"context"
	"database/sql"
	"errors"

	"gorm.io/gorm"

	"cats-go-api/internal/domain"
	"cats-go-api/internal/model"
)

var ErrDuplicateCat = errors.New("duplicate cat")

type CatService struct {
	repo domain.CatRepository
}

func NewCatService(repo domain.CatRepository) *CatService {
	return &CatService{repo: repo}
}

func (s *CatService) List(ctx context.Context) ([]model.Cat, error) {
	return s.repo.FindAll(ctx)
}

func (s *CatService) GetByID(ctx context.Context, id int) (*model.Cat, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CatService) GetByCatID(ctx context.Context, catID string) (*model.Cat, error) {
	return s.repo.FindByCatID(ctx, catID)
}

func (s *CatService) Save(ctx context.Context, input model.SaveCatInput) (*model.Cat, error) {
	cat, err := s.repo.Create(ctx, input, "go")
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrDuplicateCat
		}
		return nil, err
	}
	return cat, nil
}

func (s *CatService) Update(ctx context.Context, input model.SaveCatInput, id int) (*model.Cat, error) {
	return s.repo.Update(ctx, input, "go", id)
}

func (s *CatService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err == sql.ErrNoRows {
		return sql.ErrNoRows
	}
	return err
}
