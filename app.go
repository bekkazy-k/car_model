package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// App - Структура приложения
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Init - Инициализация подключения к БД
func (a *App) Init() {
	var err error
	a.DB, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal("failed to connect database")
	}
	// defer a.DB.Close()
	a.DB.Debug().AutoMigrate(&Car{})

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

// Run - Запуск приложения
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/car/{id:[0-9]+}", a.getCar).Methods("GET")
	a.Router.HandleFunc("/cars", a.getCars).Methods("GET")
	a.Router.HandleFunc("/car", a.createCar).Methods("POST")
	a.Router.HandleFunc("/car/{id:[0-9]+}", a.updateCar).Methods("PUT")
	a.Router.HandleFunc("/car/{id:[0-9]+}", a.deleteCar).Methods("DELETE")
}

func (a *App) getCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid car ID")
		return
	}

	c := Car{ID: uint(id)}
	if err := c.getCar(a.DB); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			respondWithError(w, http.StatusNotFound, "Car not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) getCars(w http.ResponseWriter, r *http.Request) {

	cars, err := getCars(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, cars)
}

func (a *App) createCar(w http.ResponseWriter, r *http.Request) {
	var c Car
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := c.createCar(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, c)
}

func (a *App) updateCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid car ID")
		return
	}

	var c Car
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	c.ID = uint(id)

	if err := c.updateCar(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) deleteCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid car ID")
		return
	}

	c := Car{ID: uint(id)}
	if err := c.deleteCar(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
