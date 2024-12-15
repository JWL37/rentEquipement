package equipment

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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

const (
	querryListAvailableEquipments         = `SELECT id,name,type,description,pricePerDay,isAvailable from equipments WHERE isAvailable=$1`
	querryEquipmentInfo                   = `SELECT id,name,type,description,pricePerDay,isAvailable FROM equipments WHERE id=$1`
	querryPriceAndIsAvailableForEquipment = `SELECT pricePerDay,isAvailable FROM equipments WHERE id=$1`
)

func (rep *Repository) ListAvailableEquipments() ([]models.Equipment, error) {
	rows, err := rep.DB.Query(querryListAvailableEquipments, "TRUE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var equipments []models.Equipment

	for rows.Next() {
		var eq models.Equipment
		err := rows.Scan(&eq.ID, &eq.Name, &eq.Type, &eq.Description, &eq.PricePerDay, &eq.IsAvailable)
		if err != nil {
			return nil, err

		}
		equipments = append(equipments, eq)
	}
	return equipments, nil
}

func (rep *Repository) ListNotAvailableEquipments() ([]models.Equipment, error) {
	rows, err := rep.DB.Query(querryListAvailableEquipments, "FALSE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var equipments []models.Equipment

	for rows.Next() {
		var eq models.Equipment
		err := rows.Scan(&eq.ID, &eq.Name, &eq.Type, &eq.Description, &eq.PricePerDay, &eq.IsAvailable)
		if err != nil {
			return nil, err

		}
		equipments = append(equipments, eq)
	}
	return equipments, nil
}

func (rep *Repository) EquipmentInfo(equipmentID string) (*models.Equipment, error) {
	row := rep.DB.QueryRow(querryEquipmentInfo, equipmentID)
	eq := models.Equipment{}
	err := row.Scan(&eq.ID, &eq.Name, &eq.Type, &eq.Description, &eq.PricePerDay, &eq.IsAvailable)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &eq, nil
}

func (rep *Repository) PriceAndIsAvailableForEquipment(equipmentID string) (float64, bool, error) {
	row := rep.DB.QueryRow(querryPriceAndIsAvailableForEquipment, equipmentID)
	pricePerDay := 0.
	isAvailable := false
	err := row.Scan(&pricePerDay, &isAvailable)
	if err != nil {
		log.Println(err)
		return pricePerDay, isAvailable, err
	}
	return pricePerDay, isAvailable, nil
}

func (rep *Repository) CreateOrder(equipmentID, userID string, start time.Time, countDays int) (*OrderInfo, error) {
	tx, err := rep.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return nil, err
	}

	pricePerDay, isAvailable, err := rep.PriceAndIsAvailableForEquipment(equipmentID)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return nil, err
	}
	if !isAvailable {
		tx.Rollback()
		return nil, errors.New("this equipment is not available")
	}
	query := `
		SELECT id FROM customers 
		WHERE userID = $1
	`
	customerID := 0
	err = rep.DB.QueryRow(query, userID).Scan(
		&customerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("customer not found for userID %s", userID)
		}
		return nil, fmt.Errorf("error fetching customer: %w", err)
	}
	endDate := start.AddDate(0, 0, countDays)
	priceForRent := pricePerDay * float64(countDays)

	_, err = tx.Exec(`
	INSERT INTO orders (customerid, equipmentID, start_date, endDate, totalCost)
	VALUES ($1, $2, $3, $4, $5)`,
		customerID, equipmentID, start, endDate, priceForRent)

	if err != nil {
		tx.Rollback()
		log.Println("Error inserting order:", err)
		return nil, err
	}

	_, err = tx.Exec(`
		UPDATE equipments
		SET isAvailable = FALSE
		WHERE id = $1`, equipmentID)
	if err != nil {
		tx.Rollback()
		log.Println("Error updating equipment availability:", err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return nil, err
	}

	return &OrderInfo{
		StartDate: start,
		EndDate:   endDate,
		TotalCost: priceForRent,
	}, nil
}

func (rep *Repository) DoAvailable(equipmentID string) error {
	query := `UPDATE equipments SET isAvailable = TRUE WHERE id = $1`

	_, err := rep.DB.Exec(query, equipmentID)
	if err != nil {
		return fmt.Errorf("error updating availability: %w", err)
	}

	return nil
}
