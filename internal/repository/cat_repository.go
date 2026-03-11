package repository

import (
	"context"
	"database/sql"
	"encoding/json"

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

type postgresCatRepository struct {
	db *sql.DB
}

func NewPostgresCatRepository(db *sql.DB) CatRepository {
	return &postgresCatRepository{db: db}
}

func (r *postgresCatRepository) FindAll(ctx context.Context) ([]model.Cat, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, cat_id, url, width, height, breeds, api_used, created_at, updated_at FROM cats`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []model.Cat
	for rows.Next() {
		var c model.Cat
		var breedsRaw []byte
		var apiUsed sql.NullString
		if err := rows.Scan(&c.ID, &c.CatID, &c.URL, &c.Width, &c.Height, &breedsRaw, &apiUsed, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		if apiUsed.Valid {
			c.APIUsed = apiUsed.String
		}
		if len(breedsRaw) > 0 {
			c.Breeds = breedsRaw
		} else {
			c.Breeds = json.RawMessage("[]")
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

func (r *postgresCatRepository) FindByID(ctx context.Context, id int) (*model.Cat, error) {
	var c model.Cat
	var breedsRaw []byte
	var apiUsed sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, cat_id, url, width, height, breeds, api_used, created_at, updated_at FROM cats WHERE id = $1`, id,
	).Scan(&c.ID, &c.CatID, &c.URL, &c.Width, &c.Height, &breedsRaw, &apiUsed, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if apiUsed.Valid {
		c.APIUsed = apiUsed.String
	}
	if len(breedsRaw) > 0 {
		c.Breeds = breedsRaw
	} else {
		c.Breeds = json.RawMessage("[]")
	}
	return &c, nil
}

func (r *postgresCatRepository) FindByCatID(ctx context.Context, catID string) (*model.Cat, error) {
	var c model.Cat
	var breedsRaw []byte
	var apiUsed sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, cat_id, url, width, height, breeds, api_used, created_at, updated_at FROM cats WHERE cat_id = $1`, catID,
	).Scan(&c.ID, &c.CatID, &c.URL, &c.Width, &c.Height, &breedsRaw, &apiUsed, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if apiUsed.Valid {
		c.APIUsed = apiUsed.String
	}
	if len(breedsRaw) > 0 {
		c.Breeds = breedsRaw
	} else {
		c.Breeds = json.RawMessage("[]")
	}
	return &c, nil
}

func (r *postgresCatRepository) Create(ctx context.Context, input model.SaveCatInput, apiUsed string) (*model.Cat, error) {
	var c model.Cat
	var breedsRaw []byte
	if input.Breeds != nil {
		breedsRaw = input.Breeds
	} else {
		breedsRaw = []byte("[]")
	}

	err := r.db.QueryRowContext(ctx,
		`INSERT INTO cats (cat_id, url, width, height, breeds, api_used, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		 RETURNING id, cat_id, url, width, height, breeds, api_used, created_at, updated_at`,
		input.CatID, input.URL, input.Width, input.Height, breedsRaw, apiUsed,
	).Scan(&c.ID, &c.CatID, &c.URL, &c.Width, &c.Height, &breedsRaw, &c.APIUsed, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	c.Breeds = breedsRaw
	return &c, nil
}

func (r *postgresCatRepository) Update(ctx context.Context, input model.SaveCatInput, apiUsed string, id int) (*model.Cat, error) {
	var c model.Cat
	var breedsRaw []byte
	if input.Breeds != nil {
		breedsRaw = input.Breeds
	} else {
		breedsRaw = []byte("[]")
	}

	err := r.db.QueryRowContext(ctx,
		`UPDATE cats SET cat_id = $1, url = $2, width = $3, height = $4, breeds = $5, api_used = $6, updated_at = NOW()
		 WHERE id = $7
		 RETURNING id, cat_id, url, width, height, breeds, api_used, created_at, updated_at`,
		input.CatID, input.URL, input.Width, input.Height, breedsRaw, apiUsed, id,
	).Scan(&c.ID, &c.CatID, &c.URL, &c.Width, &c.Height, &breedsRaw, &c.APIUsed, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	c.Breeds = breedsRaw
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
