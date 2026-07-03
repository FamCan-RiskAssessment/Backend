package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/crypto"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	postgresRepo "github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/seed"
)

func main() {
	phone := flag.String("phone", "", "super admin phone number (defaults to SUPER_ADMIN_PHONE)")
	password := flag.String("password", "", "new password (defaults to SUPER_ADMIN_PASSWORD)")
	flag.Parse()

	config := bootstrap.Run()

	if *phone == "" {
		*phone = config.Env.SuperAdmin.Phone
	}
	if *password == "" {
		*password = config.Env.SuperAdmin.Password
	}

	db := database.NewPostgresDatabase(&config.Env.Database, &config.Env.Server)
	userRepository := postgresRepo.NewUserRepository()
	passwordHasher := crypto.NewPasswordHasher()
	changer := seed.NewSuperAdminPasswordChanger(db, userRepository, passwordHasher)

	if err := changer.ChangePassword(*phone, *password); err != nil {
		log.Fatalf("change super admin password: %v", err)
	}

	fmt.Fprintf(os.Stdout, "super admin password updated for phone %s\n", *phone)
}
