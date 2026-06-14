package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"auth-golang-jwt/config"
	"auth-golang-jwt/models"
	"auth-golang-jwt/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to in-memory database: %v", err)
	}

	err = db.AutoMigrate(&models.Employee{})
	if err != nil {
		t.Fatalf("Failed to auto migrate models: %v", err)
	}

	config.DB = db
	return db
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestCreateEmployee(t *testing.T) {
	setupTestDB(t)
	r := setupTestRouter()
	r.POST("/api/employees", CreateEmployee)

	// Valid payload
	payload := models.CreateEmployeeRequest{
		Nik:        "EMP001",
		Name:       "John Doe",
		Email:      "john.doe@example.com",
		Phone:      "08123456789",
		Department: "Engineering",
		Position:   "Software Engineer",
		Salary:     15000000,
		JoinDate:   "2026-06-08",
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/employees", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d. Response: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var response utils.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response.Success {
		t.Errorf("Expected success to be true, got false")
	}

	// Test duplicate NIK
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for duplicate NIK, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetEmployees(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter()
	r.GET("/api/employees", GetEmployees)

	// Seed data
	employees := []models.Employee{
		{
			Nik:        "EMP001",
			Name:       "Alice",
			Email:      "alice@example.com",
			Department: "HR",
			Position:   "Manager",
			Salary:     12000000,
			JoinDate:   time.Now(),
			IsActive:   true,
		},
		{
			Nik:        "EMP002",
			Name:       "Bob",
			Email:      "bob@example.com",
			Department: "Engineering",
			Position:   "Developer",
			Salary:     10000000,
			JoinDate:   time.Now(),
			IsActive:   false,
		},
	}
	db.Create(&employees)

	// Get all
	req, _ := http.NewRequest("GET", "/api/employees", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Test search filter
	req, _ = http.NewRequest("GET", "/api/employees?search=Alice", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response utils.PaginatedResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Pagination.TotalRows != 1 {
		t.Errorf("Expected total rows to be 1, got %d", response.Pagination.TotalRows)
	}

	// Test department filter
	req, _ = http.NewRequest("GET", "/api/employees?department=Engineering", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Pagination.TotalRows != 1 {
		t.Errorf("Expected total rows to be 1, got %d", response.Pagination.TotalRows)
	}

	// Test active filter
	req, _ = http.NewRequest("GET", "/api/employees?is_active=false", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Pagination.TotalRows != 1 {
		t.Errorf("Expected total rows to be 1, got %d", response.Pagination.TotalRows)
	}
}

func TestGetEmployeeByID(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter()
	r.GET("/api/employees/:id", GetEmployeeByID)

	emp := models.Employee{
		Nik:        "EMP001",
		Name:       "Alice",
		Email:      "alice@example.com",
		Department: "HR",
		Position:   "Manager",
		Salary:     12000000,
		JoinDate:   time.Now(),
		IsActive:   true,
	}
	db.Create(&emp)

	// Test existing ID
	req, _ := http.NewRequest("GET", "/api/employees/"+strconv.FormatUint(emp.ID, 10), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Test non-existing ID
	req, _ = http.NewRequest("GET", "/api/employees/9999", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestUpdateEmployee(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter()
	r.PUT("/api/employees/:id", UpdateEmployee)

	emp := models.Employee{
		Nik:        "EMP001",
		Name:       "Alice",
		Email:      "alice@example.com",
		Department: "HR",
		Position:   "Manager",
		Salary:     12000000,
		JoinDate:   time.Now(),
		IsActive:   true,
	}
	db.Create(&emp)

	// Update payload
	payload := models.UpdateEmployeeRequest{
		Name:       "Alice Updated",
		Salary:     15000000,
		Department: "Engineering",
	}
	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", "/api/employees/"+strconv.FormatUint(emp.ID, 10), bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Response: %s", w.Code, w.Body.String())
	}

	var updated models.Employee
	db.First(&updated, emp.ID)
	if updated.Name != "Alice Updated" || updated.Salary != 15000000 || updated.Department != "Engineering" {
		t.Errorf("Employee fields were not updated correctly: %+v", updated)
	}
}

func TestDeleteEmployee(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter()
	r.DELETE("/api/employees/:id", DeleteEmployee)

	emp := models.Employee{
		Nik:        "EMP001",
		Name:       "Alice",
		Email:      "alice@example.com",
		Department: "HR",
		Position:   "Manager",
		Salary:     12000000,
		JoinDate:   time.Now(),
		IsActive:   true,
	}
	db.Create(&emp)

	req, _ := http.NewRequest("DELETE", "/api/employees/"+strconv.FormatUint(emp.ID, 10), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify soft delete
	var count int64
	db.Model(&models.Employee{}).Where("id = ?", emp.ID).Count(&count)
	if count != 0 {
		t.Errorf("Expected employee to be soft deleted, count is %d", count)
	}

	// Verify it still exists in the DB (unscoped)
	db.Unscoped().Model(&models.Employee{}).Where("id = ?", emp.ID).Count(&count)
	if count != 1 {
		t.Errorf("Expected employee to exist in DB unscoped, count is %d", count)
	}
}
