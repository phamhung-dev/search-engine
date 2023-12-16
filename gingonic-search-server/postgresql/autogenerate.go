package postgresql

import (
	"fmt"
	"gingonic-search-server/models"
	"log"

	"gorm.io/gorm"
)

func autogenerate(db *gorm.DB) {
	fmt.Println("Autogenerating...")

	if err := db.AutoMigrate(
		&models.User{},
		&models.File{},
	); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Autogenerate successfully!")
}
