package maintenance

import "rentEquipement/internal/models"

type MaintenanceRepo interface {
	List(string) ([]models.Maintenance, error)
}
