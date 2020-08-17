package main

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type sqlite_sequence struct {
	name string
	seg  int64
}

// Car - Модель автомобиля
type Car struct {
	ID      uint   `gorm:"primary_key"`
	Brand   string `gorm:"size:50"` // Согласно ссылке https://1gai.ru/publ/515183-spisok-vseh-avtomobilnyh-marok-mira.html, максимальной длины 50 хватит с головой
	Model   string `gorm:"size:255"`
	Price   uint32 // от 0 до 4294967295
	Status  string `gorm:"size:15"` // В пути, На складе, Продан, Снят с продажи
	Mileage uint32
}

func (c *Car) validate() error {
	if c.Brand == "" {
		return errors.New("Brand is empty")
	}
	if c.Model == "" {
		return errors.New("Model is empty")
	}
	if c.Price < 0 {
		return errors.New("Price cannot be less than zero")
	}
	if c.Status == "" {
		return errors.New("Status is empty")
	}
	if !isValidStatus(c.Status) {
		return errors.New("Status can be one of the following 'В пути, На складе, Продан, Снят с продажи'")
	}
	if c.Mileage < 0 {
		return errors.New("Mileage cannot be less than zero")
	}
	return nil
}

func isValidStatus(status string) bool {
	switch status {
	case
		"В пути",
		"На складе",
		"Продан",
		"Снят с продажи":
		return true
	}
	return false
}

func (c *Car) getCar(db *gorm.DB) error {
	return db.First(&c, c.ID).Error
}

func (c *Car) createCar(db *gorm.DB) error {
	if err := c.validate(); err != nil {
		return err
	}
	return db.Create(&c).Error
}

func (c *Car) updateCar(db *gorm.DB) error {
	if err := c.validate(); err != nil {
		return err
	}
	return db.Save(&c).Error
}

func (c *Car) deleteCar(db *gorm.DB) error {
	return db.Unscoped().Delete(&c).Error
}

func getCars(db *gorm.DB) ([]Car, error) {
	var car []Car
	err := db.Find(&car).Error
	return car, err
}
