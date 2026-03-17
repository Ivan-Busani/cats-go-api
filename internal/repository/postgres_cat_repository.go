package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"cats-go-api/internal/domain"
	"cats-go-api/internal/model"
)

type postgresCatRepository struct {
	db *gorm.DB
}

func NewPostgresCatRepository(db *gorm.DB) domain.CatRepository {
	return &postgresCatRepository{db: db}
}

func (r *postgresCatRepository) FindAll(ctx context.Context) ([]model.Cat, error) {
	cats := []model.Cat{}
	err := r.db.WithContext(ctx).Order("id DESC").Find(&cats).Error
	if err != nil {
		return nil, err
	}
	for i := range cats {
		if cats[i].Breeds == nil {
			cats[i].Breeds = datatypes.JSON(json.RawMessage("[]"))
		}
	}
	return cats, nil
}

func (r *postgresCatRepository) FindByID(ctx context.Context, id int) (*model.Cat, error) {
	var c model.Cat
	err := r.db.WithContext(ctx).First(&c, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if c.Breeds == nil {
		c.Breeds = datatypes.JSON(json.RawMessage("[]"))
	}
	return &c, nil
}

func (r *postgresCatRepository) FindByCatID(ctx context.Context, catID string) (*model.Cat, error) {
	var c model.Cat
	err := r.db.WithContext(ctx).Where("cat_id = ?", catID).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if c.Breeds == nil {
		c.Breeds = datatypes.JSON(json.RawMessage("[]"))
	}
	return &c, nil
}

func (r *postgresCatRepository) Create(ctx context.Context, input model.SaveCatInput, apiUsed string) (*model.Cat, error) {
	breedsRaw := input.Breeds
	if breedsRaw == nil {
		breedsRaw = json.RawMessage("[]")
	}

	cat := model.Cat{
		CatID:   input.CatID,
		URL:     input.URL,
		Width:   input.Width,
		Height:  input.Height,
		Breeds:  datatypes.JSON(breedsRaw),
		APIUsed: &apiUsed,
	}
	err := r.db.WithContext(ctx).Create(&cat).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *postgresCatRepository) Update(ctx context.Context, input model.SaveCatInput, apiUsed string, id int) (*model.Cat, error) {
	breedsRaw := input.Breeds
	if breedsRaw == nil {
		breedsRaw = json.RawMessage("[]")
	}

	result := r.db.WithContext(ctx).Model(&model.Cat{}).Where("id = ?", id).Updates(model.Cat{
		CatID:   input.CatID,
		URL:     input.URL,
		Width:   input.Width,
		Height:  input.Height,
		Breeds:  datatypes.JSON(breedsRaw),
		APIUsed: &apiUsed,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	var c model.Cat
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}
	if c.Breeds == nil {
		c.Breeds = datatypes.JSON(json.RawMessage("[]"))
	}
	return &c, nil
}

func (r *postgresCatRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&model.Cat{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
