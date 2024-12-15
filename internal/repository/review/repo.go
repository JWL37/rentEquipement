package review

import (
	"database/sql"
	"fmt"
	"rentEquipement/internal/models"
)

type Repository struct {
	DB *sql.DB
}

func NewRep(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (rep *Repository) List(equipmentID string) ([]models.Review, error) {
	query := `
        SELECT id, customerID, equipmentID, rating, comment, reviewDate
        FROM reviews
        WHERE equipmentID = $1
    `
	rows, err := rep.DB.Query(query, equipmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query reviews: %v", err)
	}
	defer rows.Close()

	var reviews []models.Review

	for rows.Next() {
		var review models.Review
		err := rows.Scan(&review.ID, &review.CustomerID, &review.EquipmentID, &review.Rating, &review.Comment, &review.ReviewDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %v", err)
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}
