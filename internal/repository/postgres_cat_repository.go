package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/jmoiron/sqlx"

	"cats-go-api/internal/domain"
	"cats-go-api/internal/model"
)

type postgresCatRepository struct {
	db *sqlx.DB
}

func NewPostgresCatRepository(db *sqlx.DB) domain.CatRepository {
	return &postgresCatRepository{db: db}
}

func (r *postgresCatRepository) FindAll(ctx context.Context) ([]model.Cat, error) {
	cats := []model.Cat{}
	err := r.db.SelectContext(ctx, &cats,
		`SELECT id, cat_id, url, width, height, breeds, api_used, created_at, updated_at FROM cats ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	for i := range cats {
		if cats[i].Breeds == nil {
			cats[i].Breeds = json.RawMessage("[]")
		}
	}
	return cats, nil
}

func (r *postgresCatRepository) FindByID(ctx context.Context, id int) (*model.Cat, error) {
	var c model.Cat
	err := r.db.GetContext(ctx, &c,
		`SELECT id, cat_id, url, width, height, breeds, api_used, created_at, updated_at FROM cats WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if c.Breeds == nil {
		c.Breeds = json.RawMessage("[]")
	}
	return &c, nil
}

func (r *postgresCatRepository) FindByCatID(ctx context.Context, catID string) (*model.Cat, error) {
	var c model.Cat
	err := r.db.GetContext(ctx, &c,
		`SELECT id, cat_id, url, width, height, breeds, api_used, created_at, updated_at FROM cats WHERE cat_id = $1`, catID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if c.Breeds == nil {
		c.Breeds = json.RawMessage("[]")
	}
	return &c, nil
}

func (r *postgresCatRepository) Create(ctx context.Context, input model.SaveCatInput, apiUsed string) (*model.Cat, error) {
	breedsRaw := input.Breeds
	if breedsRaw == nil {
		breedsRaw = json.RawMessage("[]")
	}

	var c model.Cat
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO cats (cat_id, url, width, height, breeds, api_used, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		 RETURNING id, cat_id, url, width, height, breeds, api_used, created_at, updated_at`,
		input.CatID, input.URL, input.Width, input.Height, breedsRaw, apiUsed,
	).StructScan(&c)
	if err != nil {
		return nil, err
	}
	if c.Breeds == nil {
		c.Breeds = json.RawMessage("[]")
	}
	return &c, nil
}

func (r *postgresCatRepository) Update(ctx context.Context, input model.SaveCatInput, apiUsed string, id int) (*model.Cat, error) {
	breedsRaw := input.Breeds
	if breedsRaw == nil {
		breedsRaw = json.RawMessage("[]")
	}

	var c model.Cat
	err := r.db.QueryRowxContext(ctx,
		`UPDATE cats SET cat_id = $1, url = $2, width = $3, height = $4, breeds = $5, api_used = $6, updated_at = NOW()
		 WHERE id = $7
		 RETURNING id, cat_id, url, width, height, breeds, api_used, created_at, updated_at`,
		input.CatID, input.URL, input.Width, input.Height, breedsRaw, apiUsed, id,
	).StructScan(&c)
	if err != nil {
		return nil, err
	}
	if c.Breeds == nil {
		c.Breeds = json.RawMessage("[]")
	}
	return &c, nil
}

func (r *postgresCatRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM cats WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
