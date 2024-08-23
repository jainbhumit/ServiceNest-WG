package repository

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"serviceNest/model"
)

type AdminRepository struct {
	filePath string
}

func NewAdminRepository(filePath string) *AdminRepository {
	return &AdminRepository{filePath: filePath}
}

// GetAdmin retrieves the admin details from the file
func (r *AdminRepository) GetAdmin() (*model.Admin, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var admin model.Admin
	err = json.NewDecoder(file).Decode(&admin)
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

// SaveAdmin saves the admin details to the file
func (r *AdminRepository) SaveAdmin(admin *model.Admin) error {
	data, err := json.MarshalIndent(admin, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(r.filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
