package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"serviceNest/model"
)

type HouseholderRepository struct {
	filePath string
}

// NewHouseholderRepository initializes a new HouseholderRepository
func NewHouseholderRepository(filePath string) *HouseholderRepository {
	return &HouseholderRepository{filePath: filePath}
}

// SaveHouseholder saves a new householder to the file
func (repo *HouseholderRepository) SaveHouseholder(householder model.Householder) error {
	householders, err := repo.loadHouseholders()
	if err != nil {
		return err
	}

	householders = append(householders, householder)

	return repo.saveHouseholders(householders)
}

// GetHouseholderByID retrieves a householder by their ID
func (repo *HouseholderRepository) GetHouseholderByID(id string) (*model.Householder, error) {
	householders, err := repo.loadHouseholders()
	if err != nil {
		return nil, err
	}

	for _, householder := range householders {
		if householder.ID == id {
			return &householder, nil
		}
	}

	return nil, fmt.Errorf("householder not found")
}

// Private helper methods for loading and saving householders and service requests
func (repo *HouseholderRepository) loadHouseholders() ([]model.Householder, error) {
	var householders []model.Householder

	file, err := ioutil.ReadFile(repo.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return householders, nil
		}
		return nil, err
	}

	err = json.Unmarshal(file, &householders)
	if err != nil {
		return nil, err
	}

	return householders, nil
}

func (repo *HouseholderRepository) saveHouseholders(householders []model.Householder) error {
	data, err := json.MarshalIndent(householders, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(repo.filePath, data, 0644)
}
