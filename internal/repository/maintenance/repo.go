package maintenance

import (
	"database/sql"
	"fmt"
	"rentEquipement/internal/models"
	"time"
)

type Repository struct {
	DB *sql.DB
}

func NewRep(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (rep *Repository) List(equipmentID string) ([]models.Maintenance, error) {
	query := `
		SELECT id, equipmentID, date, description
		FROM maintenances
		WHERE equipmentID = $1
	`
	rows, err := rep.DB.Query(query, equipmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var maintenances []models.Maintenance

	for rows.Next() {
		var maintenance models.Maintenance
		var date string

		if err := rows.Scan(&maintenance.ID, &maintenance.EquipmentID, &date, &maintenance.Description); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Parse the date
		parsedDate, err := time.Parse(time.RFC3339, date)

		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		maintenance.Date = parsedDate

		maintenances = append(maintenances, maintenance)
	}

	return maintenances, nil
}
