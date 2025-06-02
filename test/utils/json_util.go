package utils

import (
	"catalog-service/internal/models"
	"encoding/json"
	"os"
)

func UnmarshalServiceList(filePath string) ([]*models.Service, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var services []*models.Service
	if err = json.NewDecoder(f).Decode(&services); err != nil {
		return nil, err
	}
	return services, nil
}
