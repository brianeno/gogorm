package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectToPostgreSQL() (*gorm.DB, error) {
	dsn := "user=postgres password=admin dbname=sessiondb host=localhost port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

type ChargeSession struct {
	Id   string `gorm:"primaryKey"`
	Watt uint   `gorm:"unique"`
	Vin  string
}

func createCs(db *gorm.DB, cs *ChargeSession) (*gorm.DB, error) {
	result := db.Create(cs)
	if result.Error != nil {
		return nil, result.Error
	}
	return result, nil
}

func getCsByID(db *gorm.DB, id string) (*gorm.DB, *ChargeSession, error) {
	var cs ChargeSession
	result := db.First(&cs, id)
	if result.Error != nil {
		return nil, nil, result.Error
	}
	return result, &cs, nil
}

func updateCs(db *gorm.DB, cs *ChargeSession) (*gorm.DB, error) {
	result := db.Save(cs)
	if result.Error != nil {
		return nil, result.Error
	}
	return result, nil
}

func deleteCs(db *gorm.DB, cs *ChargeSession) (*gorm.DB, error) {
	result := db.Delete(cs)
	if result.Error != nil {
		return nil, result.Error
	}
	return result, nil
}

func deleteAllCs(db *gorm.DB, cs *ChargeSession) (*gorm.DB, error) {
	result := db.Where("VIN NOT LIKE ?", "null").Delete(&ChargeSession{})
	if result.Error != nil {
		return nil, result.Error
	}
	return result, nil
}

func main() {
	db, err := connectToPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()

	// Perform database migration
	err = db.AutoMigrate(&ChargeSession{})
	if err != nil {
		log.Fatal(err)
	}

	// Create CS
	newCs := &ChargeSession{Id: "11111", Watt: 420, Vin: ""}
	res, err := createCs(db, newCs)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created CS:", newCs, "rows affected", res.RowsAffected)

	// Query CS by ID
	csID := newCs.Id
	res, cs, err := getCsByID(db, csID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("CS by ID:", cs, "rows affected", res.RowsAffected)

	// Update CS
	cs.Vin = "4Y1SL65848Z411439"
	res, err = updateCs(db, cs)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Updated CS:", cs, "rows affected", res.RowsAffected)

	// Delete CS
	res, err = deleteCs(db, cs)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Deleted CS:", cs, "rows affected", res.RowsAffected)

	res, err = deleteAllCs(db, cs)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Deleted All CS:", cs, "Rows affects", res.RowsAffected)
}
