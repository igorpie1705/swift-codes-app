package handlers

import (
	"net/http"

	"github.com/igorpie1705/swift-codes-app/database"
	"github.com/igorpie1705/swift-codes-app/models"

	"github.com/gin-gonic/gin"
)

func GetSwiftCode(c *gin.Context) {
	swiftCode := c.Param("swift-code")


	var swiftCodeData models.SwiftCode
	db := database.GetDB()
	result := db.Where("swift_code = ?", swiftCode).First(&swiftCodeData)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Swift code not found"})
		return
	}

	if swiftCodeData.IsHeadquarter {
		var branches []models.SwiftCode
		db.Where("SUBSTRING(swift_code, 1, 8) = ? AND is_headquarter = ?", swiftCodeData.SwiftCode[:8], false).Find(&branches)

	response := gin.H{
		"address": swiftCodeData.Address,
		"bankName": swiftCodeData.Name,
		"countryISO2": swiftCodeData.CountryISO2,
		"countryName": swiftCodeData.CountryName,
		"isHeadquarter": true,
		"swiftCode": swiftCodeData.SwiftCode,
		"branches": branches,
	}
	c.JSON(http.StatusOK, response)
} else {
	response := gin.H{
		"address": swiftCodeData.Address,
		"bankName": swiftCodeData.Name,
		"countryISO2": swiftCodeData.CountryISO2,
		"countryName": swiftCodeData.CountryName,
		"isHeadquarter": false,
		"swiftCode": swiftCodeData.SwiftCode,
	}
	c.JSON(http.StatusOK, response)
}
}

func GetSwiftCodeByCountry(c *gin.Context) {
	countryISO2 := c.Param("countryISO2code")

	var swiftCodes []models.SwiftCode
	db := database.GetDB()
	db.Where("country_iso2 = ?", countryISO2).Find(&swiftCodes)

	response := gin.H{
		"countryISO2": countryISO2,
		"countryName": swiftCodes[0].CountryName,
		"swiftCodes": swiftCodes,
	}
	c.JSON(http.StatusOK, response)
}


func AddSwiftCode(c *gin.Context) {
	var swiftCode models.SwiftCode

	if err := c.ShouldBindJSON(&swiftCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	result := db.Create(&swiftCode)
	if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add swift code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Swift code added successfully"})
}

func DeleteSwiftCode(c *gin.Context) {
	swiftCode := c.Param("swift-code")

	db := database.GetDB()
	result := db.Where("swift_code = ?", swiftCode).Delete(&models.SwiftCode{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete swift code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Swift code deleted successfully"})
}