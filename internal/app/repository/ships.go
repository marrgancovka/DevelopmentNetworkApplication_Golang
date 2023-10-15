package repository

import (
	"awesomeProject/internal/app/ds"
	"errors"
	"time"
)

// возвращает список космических кораблей
func (r *Repository) Select_ships(search string) (*[]ds.Ship, error) {
	var ships []ds.Ship
	if search != "" {
		res := r.db.Where("Is_delete = ?", "False").Where("Title LIKE ?", "%"+search+"%").Find(&ships)
		return &ships, res.Error
	}
	res := r.db.Where("Is_delete = ?", "False").Find(&ships)
	return &ships, res.Error
}

// возвращает информацию о космическом корабле по айди
func (r *Repository) Select_ship(id int) (*ds.Ship, error) {
	ship := &ds.Ship{}
	res := r.db.First(&ship, "id = ?", id)
	return ship, res.Error
}

// добавление нового космического корабля
func (r *Repository) Insert_ship(ship *ds.Ship) error {
	res := r.db.Create(&ship)
	return res.Error
}

// добавление космического корабля в заявку и создание заявки если ее не было
func (r *Repository) Insert_application(request *struct {
	Id_Ship            uint
	Id_Cosmodrom_Begin uint
	Id_cosmodrom_End   uint
	Id_user            uint
	Date_Flight        time.Time
}) error {
	var app ds.Application
	r.db.Where("id_user = ?", request.Id_user).Where("status = ?", "created").First(&app)

	if app.ID == 0 {
		newApp := ds.Application{
			Id_user:       request.Id_user,
			Id_admin:      1,
			Status:        "created",
			Date_creation: time.Now(),
		}
		res := r.db.Create(&newApp)
		if res.Error != nil {
			return res.Error
		}
		app = newApp
	}
	flight := ds.Flights{
		Id_Ship:            request.Id_Ship,
		Id_Application:     app.ID,
		Id_Cosmodrom_Begin: request.Id_Cosmodrom_Begin,
		Id_cosmodrom_End:   request.Id_cosmodrom_End,
		Date_Flight:        request.Date_Flight,
	}
	result := r.db.Create(&flight)
	return result.Error
}

// изменение информации о космическом корабле
func (r *Repository) Update_ship(updateShip *ds.Ship) error {
	var ship ds.Ship
	res := r.db.First(&ship, "id =?", updateShip.ID)
	if res.Error != nil {
		return res.Error
	}
	if updateShip.Is_delete != false {
		return errors.New("Нельзя менять статус услуги")
	}
	if updateShip.Title != "" {
		ship.Title = updateShip.Title
	}
	if updateShip.Type != "" {
		ship.Type = updateShip.Type
	}
	if updateShip.Description != "" {
		ship.Description = updateShip.Description
	}
	if updateShip.Image_url != "" {
		ship.Image_url = updateShip.Image_url
	}
	if updateShip.Rocket != "" {
		ship.Rocket = updateShip.Rocket
	}

	*updateShip = ship
	result := r.db.Save(updateShip)
	return result.Error
}

// логически удаляет космический корабль
func (r *Repository) Delete_ship(id int) error {
	var ship ds.Ship
	res := r.db.First(&ship, "id =?", id)
	if res.Error != nil {
		return res.Error
	}
	ship.Is_delete = true
	result := r.db.Save(ship)
	return result.Error
}