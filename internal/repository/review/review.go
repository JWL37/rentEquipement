package review

import "rentEquipement/internal/models"

type ReviewRepo interface {
	List(string) ([]models.Review, error)
}
