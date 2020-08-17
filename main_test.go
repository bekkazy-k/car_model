package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a.Init()
	code := m.Run()
	a.DB.Unscoped().Delete(&Car{})
	var seq sqlite_sequence
	a.DB.Delete(&seq)
	//delete from sqlite_sequence where name='your_table';
	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {

	req, _ := http.NewRequest("GET", "/cars", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentCar(t *testing.T) {

	req, _ := http.NewRequest("GET", "/car/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Car not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Car not found'. Got '%s'", m["error"])
	}
}

func TestCreateCar(t *testing.T) {

	var jsonStr = []byte(`{"Brand":"test brand", "Model":"test model", "Price": 20000, "Status":"В пути", "Mileage": 0}`)
	req, _ := http.NewRequest("POST", "/car", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["Brand"] != "test brand" {
		t.Errorf("Expected brand to be 'test brand'. Got '%v'", m["Brand"])
	}
	if m["Model"] != "test model" {
		t.Errorf("Expected model to be 'test model'. Got '%v'", m["Model"])
	}
	if m["Price"] != 20000 {
		t.Errorf("Expected price to be '20000'. Got '%v'", m["Brand"])
	}
	if m["Status"] != "В пути" {
		t.Errorf("Expected status to be 'В пути'. Got '%v'", m["Status"])
	}
	if m["Mileage"] != 0 {
		t.Errorf("Expected mileage to be '0'. Got '%v'", m["Mileage"])
	}
}

func TestGetCar(t *testing.T) {
	addCars(1)

	req, _ := http.NewRequest("GET", "/car/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateCar(t *testing.T) {

	addCars(1)

	req, _ := http.NewRequest("GET", "/car/1", nil)
	response := executeRequest(req)
	var responseCar map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &responseCar)

	var jsonStr = []byte(`{"Brand":"updated test", "Model":"updated test model", "Price": 40000, "Status":"На складе", "Mileage": 10}`)
	req, _ = http.NewRequest("PUT", "/car/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != responseCar["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", responseCar["id"], m["id"])
	}

	if m["Brand"] == responseCar["Brand"] {
		t.Errorf("Expected the brand to change from '%v' to '%v'. Got '%v'", responseCar["Brand"], m["Brand"], m["Brand"])
	}

	if m["Model"] == responseCar["Model"] {
		t.Errorf("Expected the model to change from '%v' to '%v'. Got '%v'", responseCar["Model"], m["Model"], m["Model"])
	}

	if m["Price"] == responseCar["Price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", responseCar["Price"], m["Price"], m["Price"])
	}

	if m["Status"] == responseCar["Status"] {
		t.Errorf("Expected the status to change from '%v' to '%v'. Got '%v'", responseCar["Status"], m["Status"], m["Status"])
	}

	if m["Mileage"] == responseCar["Mileage"] {
		t.Errorf("Expected the mileage to change from '%v' to '%v'. Got '%v'", responseCar["Mileage"], m["Mileage"], m["Mileage"])
	}
}

func TestDeleteCar(t *testing.T) {
	addCars(1)

	req, _ := http.NewRequest("GET", "/car/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/car/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/car/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func addCars(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		a.DB.Create(&Car{Brand: "test brand" + strconv.Itoa(i), Model: "test model" + strconv.Itoa(i), Price: 20000, Status: "В пути", Mileage: 0})
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
