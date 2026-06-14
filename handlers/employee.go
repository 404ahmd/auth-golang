package handlers

import (
	"math"
	"net/http"
	"time"

	"auth-golang-jwt/config"
	"auth-golang-jwt/models"
	"auth-golang-jwt/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateEmployee creates a new employee record.
func CreateEmployee(c *gin.Context) {
	var req models.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		config.Log.Warn("CreateEmployee: Invalid request payload", zap.Error(err))
		utils.BadRequest(c, "Validation failed", err.Error())
		return
	}

	// Parse join_date (YYYY-MM-DD)
	joinDate, err := time.Parse("2006-01-02", req.JoinDate)
	if err != nil {
		config.Log.Warn("CreateEmployee: Invalid join_date format", zap.String("join_date", req.JoinDate), zap.Error(err))
		utils.BadRequest(c, "Invalid join_date format, must be YYYY-MM-DD", err.Error())
		return
	}

	// Check if NIK already exists
	var existing models.Employee
	if err := config.DB.Where("nik = ?", req.Nik).First(&existing).Error; err == nil {
		utils.BadRequest(c, "NIK already exists", nil)
		return
	}

	// Check if Email already exists
	if err := config.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		utils.BadRequest(c, "Email already exists", nil)
		return
	}

	employee := models.Employee{
		Nik:        req.Nik,
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Department: req.Department,
		Position:   req.Position,
		Salary:     req.Salary,
		JoinDate:   joinDate,
		IsActive:   true,
	}

	if err := config.DB.Create(&employee).Error; err != nil {
		config.Log.Error("CreateEmployee: Failed to save to database", zap.Error(err))
		utils.InternalError(c, "Failed to create employee", err.Error())
		return
	}

	utils.Created(c, "Employee created successfully", employee)
}

// GetEmployees retrieves all employee records with search, filter, and pagination.
func GetEmployees(c *gin.Context) {
	var params models.EmployeeQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		config.Log.Warn("GetEmployees: Invalid query parameters", zap.Error(err))
		utils.BadRequest(c, "Invalid query parameters", err.Error())
		return
	}

	query := config.DB.Model(&models.Employee{})

	// Apply Search filter (name, nik, email)
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("name LIKE ? OR nik LIKE ? OR email LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Apply Department filter
	if params.Department != "" {
		query = query.Where("department = ?", params.Department)
	}

	// Apply IsActive filter
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	// Count total rows
	var totalRows int64
	if err := query.Count(&totalRows).Error; err != nil {
		config.Log.Error("GetEmployees: Failed to count database records", zap.Error(err))
		utils.InternalError(c, "Failed to retrieve employee count", err.Error())
		return
	}

	// Pagination math
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	offset := (params.Page - 1) * params.Limit
	totalPages := int(math.Ceil(float64(totalRows) / float64(params.Limit)))

	// Fetch data
	var employees []models.Employee
	if err := query.Limit(params.Limit).Offset(offset).Order("created_at DESC").Find(&employees).Error; err != nil {
		config.Log.Error("GetEmployees: Failed to fetch database records", zap.Error(err))
		utils.InternalError(c, "Failed to retrieve employees list", err.Error())
		return
	}

	pagination := utils.Pagination{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}

	utils.OkPaginated(c, "Employees retrieved successfully", employees, pagination)
}

// GetEmployeeByID retrieves a single employee record by ID.
func GetEmployeeByID(c *gin.Context) {
	id := c.Param("id")

	var employee models.Employee
	if err := config.DB.First(&employee, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, utils.Response{
				Success: false,
				Message: "Employee not found",
			})
			return
		}
		config.Log.Error("GetEmployeeByID: Database error", zap.Error(err))
		utils.InternalError(c, "Database error", err.Error())
		return
	}

	utils.OK(c, "Employee retrieved successfully", employee)
}

// UpdateEmployee updates an existing employee record.
func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")

	var employee models.Employee
	if err := config.DB.First(&employee, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, utils.Response{
				Success: false,
				Message: "Employee not found",
			})
			return
		}
		config.Log.Error("UpdateEmployee: Database error", zap.Error(err))
		utils.InternalError(c, "Database error", err.Error())
		return
	}

	var req models.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		config.Log.Warn("UpdateEmployee: Invalid request payload", zap.Error(err))
		utils.BadRequest(c, "Validation failed", err.Error())
		return
	}

	// Check NIK uniqueness if updated and different
	if req.Nik != "" && req.Nik != employee.Nik {
		var existing models.Employee
		if err := config.DB.Where("nik = ?", req.Nik).First(&existing).Error; err == nil {
			utils.BadRequest(c, "NIK already exists", nil)
			return
		}
		employee.Nik = req.Nik
	}

	// Check Email uniqueness if updated and different
	if req.Email != "" && req.Email != employee.Email {
		var existing models.Employee
		if err := config.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
			utils.BadRequest(c, "Email already exists", nil)
			return
		}
		employee.Email = req.Email
	}

	// Parse and update join date if provided
	if req.JoinDate != "" {
		joinDate, err := time.Parse("2006-01-02", req.JoinDate)
		if err != nil {
			config.Log.Warn("UpdateEmployee: Invalid join_date format", zap.String("join_date", req.JoinDate), zap.Error(err))
			utils.BadRequest(c, "Invalid join_date format, must be YYYY-MM-DD", err.Error())
			return
		}
		employee.JoinDate = joinDate
	}

	// Update optional/other fields
	if req.Name != "" {
		employee.Name = req.Name
	}
	if req.Phone != "" {
		employee.Phone = req.Phone
	}
	if req.Department != "" {
		employee.Department = req.Department
	}
	if req.Position != "" {
		employee.Position = req.Position
	}
	if req.Salary > 0 {
		employee.Salary = req.Salary
	}
	if req.IsActive != nil {
		employee.IsActive = *req.IsActive
	}

	if err := config.DB.Save(&employee).Error; err != nil {
		config.Log.Error("UpdateEmployee: Failed to save to database", zap.Error(err))
		utils.InternalError(c, "Failed to update employee", err.Error())
		return
	}

	utils.OK(c, "Employee updated successfully", employee)
}

// DeleteEmployee deletes an employee record (soft delete).
func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	var employee models.Employee
	if err := config.DB.First(&employee, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, utils.Response{
				Success: false,
				Message: "Employee not found",
			})
			return
		}
		config.Log.Error("DeleteEmployee: Database error", zap.Error(err))
		utils.InternalError(c, "Database error", err.Error())
		return
	}

	if err := config.DB.Delete(&employee).Error; err != nil {
		config.Log.Error("DeleteEmployee: Failed to delete from database", zap.Error(err))
		utils.InternalError(c, "Failed to delete employee", err.Error())
		return
	}

	utils.OK(c, "Employee deleted successfully", nil)
}