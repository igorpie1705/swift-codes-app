package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/igorpie1705/swift-codes-app/database"
	"github.com/igorpie1705/swift-codes-app/handlers"
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
		SwiftCode:     "ABCDEFGJXXX",
		Name:          "Main Bank",
		Address:       "123 Main St",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
	}
	branch := models.SwiftCode{
		SwiftCode:     "ABCDEFGJ123",
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
	c.Params = gin.Params{gin.Param{Key: "swift-code", Value: "ABCDEFGJXXX"}}

	handlers.GetSwiftCode(c)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedResponse := `{
		"address": "123 Main St",
		"bankName": "Main Bank",
		"countryISO2": "US",
		"countryName": "United States",
		"isHeadquarter": true,
		"swiftCode": "ABCDEFGJXXX",
		"branches": [
			{
				"address": "456 Branch St",
				"bankName": "Branch Bank",
				"countryISO2": "US",
				"countryName": "United States",
				"isHeadquarter": false,
				"swiftCode": "ABCDEFGJ123"
			}
		]
	}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestGetSwiftCode_NotFound(t *testing.T) {
	db := setupTestDB()
	database.SetDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "swift-code", Value: "INVALID123"}}

	handlers.GetSwiftCode(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetSwiftCodeByCountry_WithCodes(t *testing.T) {
	db := setupTestDB()
	database.SetDB(db)
	db.Where("1 = 1").Delete(&models.SwiftCode{})

	swiftCodes := []models.SwiftCode{
		{
			SwiftCode:   "ABCDEFGAXXX",
			Name:        "Main Bank",
			Address:     "123 Main St",
			CountryISO2: "US",
			CountryName: "United States",
		},
		{
			SwiftCode:   "ABCDEFGA123",
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

	handlers.GetSwiftCodeByCountry(c)

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
				"swiftCode": "ABCDEFGAXXX"
			},
			{
				"address": "456 Branch St",
				"bankName": "Branch Bank",
				"countryISO2": "US",
				"isHeadquarter": false,
				"swiftCode": "ABCDEFGA123"
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

	handlers.AddSwiftCode(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAddSwiftCode_Duplicate(t *testing.T) {
	db := setupTestDB()
	database.SetDB(db)

	swiftCode := models.SwiftCode{
		SwiftCode:     "DUPLICATEXXX",
		Name:          "Duplicate Bank",
		Address:       "1 Dup St",
		CountryISO2:   "FR",
		CountryName:   "France",
		IsHeadquarter: true,
	}
	db.Create(&swiftCode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{
        "address": "1 Dup St",
        "bankName": "Duplicate Bank",
        "countryISO2": "FR",
        "countryName": "France",
        "isHeadquarter": true,
        "swiftCode": "DUPLICATEXXX"
    }`

	c.Request, _ = http.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handlers.AddSwiftCode(c)

	assert.Equal(t, http.StatusConflict, w.Code)
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

	handlers.DeleteSwiftCode(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteSwiftCode_NotFound(t *testing.T) {
	db := setupTestDB()
	database.SetDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "swift-code", Value: "NONEXISTXXX"}}

	c.Request, _ = http.NewRequest(http.MethodDelete, "/v1/swift-codes/NONEXISTXXX", nil)

	handlers.DeleteSwiftCode(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
