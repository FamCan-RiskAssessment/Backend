package seed

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type LocationSeeder struct {
	db database.Database
}

func NewLocationSeeder(db database.Database) *LocationSeeder {
	return &LocationSeeder{db: db}
}

func (s *LocationSeeder) SeedLocations(provincesPath, citiesPath string) error {
	db := s.db.GetDB()

	if err := db.AutoMigrate(&entity.Province{}, &entity.City{}); err != nil {
		return fmt.Errorf("auto migrate locations: %w", err)
	}

	provinceMap, err := s.seedProvinces(db, provincesPath)
	if err != nil {
		return err
	}

	return s.seedCities(db, citiesPath, provinceMap)
}

func (s *LocationSeeder) seedProvinces(db *gorm.DB, path string) (map[uint]uint, error) {
	rows, err := readCSV(path)
	if err != nil {
		return nil, fmt.Errorf("read provinces csv: %w", err)
	}

	provinceMap := make(map[uint]uint, len(rows))
	for _, row := range rows {
		csvID, err := strconv.ParseUint(row["id"], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid province id %q: %w", row["id"], err)
		}

		name := row["name"]
		province := entity.Province{Name: name}
		result := db.Where("name = ?", name).FirstOrCreate(&province)
		if result.Error != nil {
			return nil, fmt.Errorf("create province %q: %w", name, result.Error)
		}

		provinceMap[uint(csvID)] = province.ID
	}

	return provinceMap, nil
}

func (s *LocationSeeder) seedCities(db *gorm.DB, path string, provinceMap map[uint]uint) error {
	rows, err := readCSV(path)
	if err != nil {
		return fmt.Errorf("read cities csv: %w", err)
	}

	cities := make([]entity.City, 0, len(rows))
	for _, row := range rows {
		csvProvinceID, err := strconv.ParseUint(row["province_id"], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid province_id %q for city %q: %w", row["province_id"], row["name"], err)
		}

		provinceID, ok := provinceMap[uint(csvProvinceID)]
		if !ok {
			return fmt.Errorf("city %q references unknown province_id %d", row["name"], csvProvinceID)
		}

		cities = append(cities, entity.City{
			Name:       row["name"],
			ProvinceID: provinceID,
		})
	}

	const batchSize = 100
	for i := 0; i < len(cities); i += batchSize {
		end := i + batchSize
		if end > len(cities) {
			end = len(cities)
		}
		batch := cities[i:end]

		for j := range batch {
			result := db.
				Where("name = ? AND province_id = ?", batch[j].Name, batch[j].ProvinceID).
				FirstOrCreate(&batch[j])
			if result.Error != nil {
				return fmt.Errorf("create city %q: %w", batch[j].Name, result.Error)
			}
		}
	}

	return nil
}

func readCSV(path string) ([]map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var rows []map[string]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		row := make(map[string]string, len(header))
		for i, col := range header {
			row[col] = record[i]
		}
		rows = append(rows, row)
	}

	return rows, nil
}
