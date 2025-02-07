package handlers

import (
	"net/http"
	"net/http/httptest"
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
		SwiftCode:     "ABCDEFGXXX",
		Name:          "Main Bank",
		Address:       "123 Main St",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
	}
	branch := models.SwiftCode{
		SwiftCode:     "ABCDEFG123",
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
	c.Params = gin.Params{gin.Param{Key: "swift-code", Value: "ABCDEFGXXX"}}

	GetSwiftCode(c)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedResponse := `{
		"address": "123 Main St",
		"bankName": "Main Bank",
		"countryISO2": "US",
		"countryName": "United States",
		"isHeadquarter": true,
		"swiftCode": "ABCDEFGXXX",
		"branches": [
			{
				"address": "456 Branch St",
				"bankName": "Branch Bank",
				"countryISO2": "US",
				"isHeadquarter": false,
				"swiftCode": "ABCDEFG123"
			}
		]
	}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}