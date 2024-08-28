package interfaces

import "serviceNest/model"

type HouseholderRepository interface {
	SaveHouseholder(householder model.Householder) error
	GetHouseholderByID(id string) (*model.Householder, error)
}
