package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/seed"
)

func main() {
	provincesPath := flag.String("provinces", "assets/provinces.csv", "path to provinces CSV file")
	citiesPath := flag.String("cities", "assets/cities.csv", "path to cities CSV file")
	flag.Parse()

	config := bootstrap.Run()
	db := database.NewPostgresDatabase(&config.Env.Database)
	locationSeeder := seed.NewLocationSeeder(db)

	if err := locationSeeder.SeedLocations(*provincesPath, *citiesPath); err != nil {
		log.Fatalf("seed locations: %v", err)
	}

	fmt.Println("provinces and cities seeded successfully")
}
