package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/igorpie1705/swift-codes-app/database"
	"github.com/igorpie1705/swift-codes-app/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}
	db.AutoMigrate(&models.SwiftCode{})
	return db
}

func TestGetSwiftCode_Headquarter(t *testing.T) {
	db := setupTestDB()
	database.SetDB(db)

	headquarter := models.SwiftCode{
		SwiftCode:     "ABCDEFGHXXX",
		Name:          "Main Bank",
		Address:       "123 Main St",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
	}
	branch := models.SwiftCode{
		SwiftCode:     "ABCDEFGH123",
		Name:          "Branch Bank",
		Address:       "456 Branch St",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: false,
	}
	db.Create(&headquarter)
	db.Create(&branch)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "swift-code", Value: "ABCDEFGHXXX"}}

	GetSwiftCode(c)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedResponse := `{
		"address": "123 Main St",
		"bankName": "Main Bank",
		"countryISO2": "US",
		"countryName": "United States",
		"isHeadquarter": true,
		"swiftCode": "ABCDEFGHXXX",
		"branches": [
			{
				"address": "456 Branch St",
				"bankName": "Branch Bank",
				"countryISO2": "US",
				"countryName": "United States",
				"isHeadquarter": false,
				"swiftCode": "ABCDEFGH123"
			}
		]
	}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestGetSwiftCodeByCountry_WithCodes(t *testing.T) {
	db := setupTestDB()
	database.SetDB(db)

	swiftCodes := []models.SwiftCode{
		{
			SwiftCode:   "ABCDEFGHXXX",
			Name:        "Main Bank",
			Address:     "123 Main St",
			CountryISO2: "US",
			CountryName: "United States",
		},
		{
			SwiftCode:   "ABCDEFGH123",
			Name:        "Branch Bank",
			Address:     "456 Branch St",
			CountryISO2: "US",
			CountryName: "United States",
		},
	}
	db.Create(&swiftCodes)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "countryISO2code", Value: "US"}}

	GetSwiftCodeByCountry(c)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedResponse := `{
		"countryISO2": "US",
		"countryName"	: "United States",
		"swiftCodes": [
			{
				"address": "123 Main St",
				"bankName": "Main Bank",
				"countryISO2": "US",
				"isHeadquarter": false,
				"swiftCode": "ABCDEFGHXXX"
			},
			{
				"address": "456 Branch St",
				"bankName": "Branch Bank",
				"countryISO2": "US",
				"isHeadquarter": false,
				"swiftCode": "ABCDEFGH123"
			}
		]
	}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestAddSwiftCode(t *testing.T) {
	db := setupTestDB()
	database.SetDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{
		"address": "789 New St",
		"bankName": "New Bank",
		"countryISO2": "GB",
		"countryName": "United Kingdom",
		"isHeadquarter": true,
		"swiftCode": "NEWBANKGBXXX"
	}`

	c.Request, _ = http.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	AddSwiftCode(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"message": "Swift code added successfully"}`, w.Body.String())
}

func TestDeleteSwiftCode(t *testing.T) {
	db := setupTestDB()
	database.SetDB(db)

	swiftCode := models.SwiftCode{
		SwiftCode:   "DELBANKKXXX",
		Name:        "Del Bank",
		Address:     "123 Delete St",
		CountryISO2: "DE",
		CountryName: "Germany",
	}
	db.Create(&swiftCode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "swift-code", Value: "DELBANKKXXX"}}

	c.Request, _ = http.NewRequest(http.MethodDelete, "/v1/swift-codes/DELBANKKXXX", nil)

	DeleteSwiftCode(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "Swift code deleted successfully"}`, w.Body.String())
}

func TestDeleteSwiftCode_NotFound(t *testing.T) {
	db := setupTestDB()
	database.SetDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "swift-code", Value: "NONEXISTXXX"}}

	c.Request, _ = http.NewRequest(http.MethodDelete, "/v1/swift-codes/NONEXISTXXX", nil)

	DeleteSwiftCode(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"error": "Swift code not found"}`, w.Body.String())
}
