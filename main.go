package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/igorpie1705/swift-codes-app/database"
	"github.com/igorpie1705/swift-codes-app/handlers"
	"github.com/igorpie1705/swift-codes-app/models"
)



func loadSwiftCodesFromFile(filePath string) ([]models.SwiftCode, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	var swiftCodes []models.SwiftCode

	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		swiftCode := models.SwiftCode{
			CountryISO2:   strings.ToUpper(record[0]),
			SwiftCode:     record[1],
			CodeType:      record[2],
			Name:          record[3],
			Address:       record[4],
			TownName:      record[5],
			CountryName:   record[6],
			TimeZone:      record[7],
		}

		swiftCodes = append(swiftCodes, swiftCode)
	}

	return swiftCodes, nil
}

func identifyHeadquartersAndBranches(swiftCodes []models.SwiftCode) []models.SwiftCode {
	for i := range swiftCodes {
		swiftCodes[i].IsHeadquarter = strings.HasSuffix(swiftCodes[i].SwiftCode, "XXX")
	}
	return swiftCodes
}

func main() {

	swiftCodes, err := loadSwiftCodesFromFile("swift_codes.csv")
    if err != nil {
        log.Fatalf("Failed to load swift codes: %v", err)
    }

    swiftCodes = identifyHeadquartersAndBranches(swiftCodes)

	db := database.InitDB()

	for _, code := range swiftCodes {
		result := db.Create(&code)
		if result.Error != nil {
			log.Printf("Failed to save swift code %s: %v", code.SwiftCode, result.Error)
			}
	}

	r := gin.Default()

	r.GET("/v1/swift-codes/:swift-code", handlers.GetSwiftCode)
	r.GET("/v1/swift-codes/country/:countryISO2code", handlers.GetSwiftCodeByCountry)
	r.POST("/v1/swift-codes", handlers.AddSwiftCode)
	r.DELETE("/v1/swift-codes/:swift-code", handlers.DeleteSwiftCode)

	r.Run(":8080")
}