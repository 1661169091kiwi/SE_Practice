package repo

import (
	"database/sql"
	"se_practice/backend/internal/db"
	"se_practice/backend/internal/model"
)

type CarouselRepo struct{}

func NewCarouselRepo() *CarouselRepo {
	return &CarouselRepo{}
}

func (r *CarouselRepo) Create(img *model.CarouselImage) error {
	query := `INSERT INTO carousel_images (image_url, title, link_url, sort_order, is_active) VALUES (?, ?, ?, ?, ?)`
	result, err := db.Exec(query, img.ImageURL, img.Title, img.LinkURL, img.SortOrder, img.IsActive)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	img.ID = id
	return nil
}

func (r *CarouselRepo) GetAll(activeOnly bool) ([]model.CarouselImage, error) {
	query := `SELECT id, image_url, title, link_url, sort_order, is_active, created_at, updated_at FROM carousel_images`
	if activeOnly {
		query += ` WHERE is_active = TRUE`
	}
	query += ` ORDER BY sort_order ASC, created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []model.CarouselImage
	for rows.Next() {
		var img model.CarouselImage
		var linkUrl sql.NullString
		var title sql.NullString
		if err := rows.Scan(&img.ID, &img.ImageURL, &title, &linkUrl, &img.SortOrder, &img.IsActive, &img.CreatedAt, &img.UpdatedAt); err != nil {
			return nil, err
		}
		if title.Valid {
			img.Title = title.String
		}
		if linkUrl.Valid {
			img.LinkURL = linkUrl.String
		}
		images = append(images, img)
	}
	// 如果没有数据，返回空切片而不是 nil
	if images == nil {
		images = []model.CarouselImage{}
	}
	return images, nil
}

func (r *CarouselRepo) Update(img *model.CarouselImage) error {
	query := `UPDATE carousel_images SET image_url=?, title=?, link_url=?, sort_order=?, is_active=? WHERE id=?`
	_, err := db.Exec(query, img.ImageURL, img.Title, img.LinkURL, img.SortOrder, img.IsActive, img.ID)
	return err
}

func (r *CarouselRepo) UpdateSortOrder(id int64, sortOrder int) error {
	_, err := db.Exec(`UPDATE carousel_images SET sort_order=? WHERE id=?`, sortOrder, id)
	return err
}

func (r *CarouselRepo) Delete(id int64) error {
	_, err := db.Exec(`DELETE FROM carousel_images WHERE id=?`, id)
	return err
}
