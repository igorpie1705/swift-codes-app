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

func setupIntegrationTest() *gin.Engine {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}
	db.AutoMigrate(&models.SwiftCode{})
	database.SetDB(db)

	r := gin.Default()
	r.GET("/v1/swift-codes/:swift-code", handlers.GetSwiftCode)
	r.GET("/v1/swift-codes/country/:countryISO2code", handlers.GetSwiftCodeByCountry)
	r.POST("/v1/swift-codes", handlers.AddSwiftCode)
	r.DELETE("/v1/swift-codes/:swift-code", handlers.DeleteSwiftCode)

	return r
}

func TestIntegration_AddAndGetSwiftCode(t *testing.T) {
	r := setupIntegrationTest()

	w := httptest.NewRecorder()
	body := `{
		"address": "789 New St",
		"bankName": "New Bank",
		"countryISO2": "GB",
		"countryName": "United Kingdom",
		"isHeadquarter": true,
		"swiftCode": "NEWBANKGBXXX"
	}`
	req, _ := http.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/swift-codes/NEWBANKGBXXX", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_AddAndDeleteSwiftCode(t *testing.T) {
	r := setupIntegrationTest()

	w := httptest.NewRecorder()
	body := `{
		"address": "789 New St",
		"bankName": "New Bank",
		"countryISO2": "GB",
		"countryName": "United Kingdom",
		"isHeadquarter": true,
		"swiftCode": "NEWBANKGBAVXXX"
	}`
	req, _ := http.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, "/v1/swift-codes/NEWBANKGBAVXXX", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/swift-codes/NEWBANKGBAVXXX", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIntegration_GetSwiftCodesByCountry(t *testing.T) {
	r := setupIntegrationTest()

	codes := []string{
		`{"swiftCode": "BANKGB001", "bankName": "Bank 1", "address": "123 St", "countryISO2": "GB", "countryName": "United Kingdom", "isHeadquarter": false}`,
		`{"swiftCode": "BANKGB002", "bankName": "Bank 2", "address": "456 St", "countryISO2": "GB", "countryName": "United Kingdom", "isHeadquarter": false}`,
	}
	for _, code := range codes {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(code))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swift-codes/country/GB", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_AddDuplicateSwiftCode(t *testing.T) {
	r := setupIntegrationTest()

	w := httptest.NewRecorder()
	body := `{
		"address": "789 Dup St",
		"bankName": "Dup Bank",
		"countryISO2": "FR",
		"countryName": "France",
		"isHeadquarter": true,
		"swiftCode": "DUPLICATEXXX"
	}`
	req, _ := http.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}
