package user

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"rentEquipement/internal/models"
	"strconv"
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

func (rep *Repository) CreateUser(username string, hashedPassword []byte, email, phoneNumber, address string) (*models.User, error) {
	tx, err := rep.DB.Begin()
	if err != nil {
		log.Println("Failed to start transaction:", err)
		return nil, err
	}

	var userID int

	err = tx.QueryRow("INSERT INTO users (Username, PasswordHash, Role) VALUES ($1, $2, $3) RETURNING id",
		username, string(hashedPassword), "customer").Scan(&userID)
	if err != nil {
		log.Println("Error inserting into users table:", err)
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO customers (userID, email, phoneNumber, address) VALUES ($1, $2, $3, $4)",
		userID, email, phoneNumber, address)
	if err != nil {
		log.Println("Error inserting into customers table:", err)
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Println("Transaction commit failed:", err)
		return nil, err
	}

	user := models.User{
		ID:       strconv.Itoa(userID),
		Username: username,
		Role:     "customer",
	}

	return &user, nil
}

func (rep *Repository) Authorize(username string) (*models.User, error) {

	row := rep.DB.QueryRow("SELECT id, username, passwordHash, role FROM Users WHERE username = $1", username)

	user := models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {

		log.Println(err)
		return nil, err
	}

	return &user, nil

}

func (rep *Repository) GetCustomerByUsername(username string) (*customerInfo, error) {
	query := `
		SELECT email, phoneNumber, address
		FROM customers join users on customers.userID=users.id
		WHERE username = $1
	`

	var customer customerInfo
	err := rep.DB.QueryRow(query, username).Scan(
		&customer.Email,
		&customer.PhoneNumber,
		&customer.Address,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("customer not found for username %s", username)
		}
		return nil, fmt.Errorf("error fetching customer: %w", err)
	}

	return &customer, nil
}

func (rep *Repository) GetLogs() ([]models.Log, error) {
	// Запрос для получения всех записей из таблицы logs
	rows, err := rep.DB.Query("SELECT id, event_time, event_type, user_id, message FROM logs")
	if err != nil {
		log.Println("Error fetching logs:", err)
		return nil, err
	}
	defer rows.Close() // Закрыть rows после завершения работы с ними

	var logs []models.Log

	// Перебор всех строк и добавление их в срез logs
	for rows.Next() {
		var logEntry models.Log
		err := rows.Scan(&logEntry.ID, &logEntry.EventTime, &logEntry.EventType, &logEntry.UserID, &logEntry.Message)
		if err != nil {
			log.Println("Error scanning log row:", err)
			return nil, err
		}
		logs = append(logs, logEntry)
	}

	// Проверка на наличие ошибок после завершения перебора строк
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}

	return logs, nil
}

type OrderView struct {
	OrderID       int       `json:"order_id"`
	CustomerID    int       `json:"customer_id"`
	EquipmentName string    `json:"equipment_name"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	TotalCost     float64   `json:"total_cost"`
}

func (rep *Repository) GetListOrdersForCustomer(userID string) ([]OrderView, error) {
	// SQL query to fetch the orders for a specific customer
	query := `
    WITH user_customer_mapping AS (
        SELECT
            u.id AS user_id,
            u.username,
            u.role,
            c.id AS customer_id,
            c.email,
            c.phoneNumber,
            c.address
        FROM
            users u
        INNER JOIN
            customers c
            ON c.userID = u.id
    )
    SELECT 
        v.order_id,
        v.customerID,
        v.equipment_name,
        v.start_date,
        v.endDate,
        v.totalCost
    FROM 
        orders_view v
    WHERE 
        v.customerID = (SELECT customer_id FROM user_customer_mapping WHERE user_id = $1);
    `

	// Initialize the result slice to hold the fetched orders
	var orders []OrderView

	// Execute the query
	rows, err := rep.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and populate the orders slice
	for rows.Next() {
		var order OrderView
		if err := rows.Scan(&order.OrderID, &order.CustomerID, &order.EquipmentName, &order.StartDate, &order.EndDate, &order.TotalCost); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	// Check if there was an error during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
