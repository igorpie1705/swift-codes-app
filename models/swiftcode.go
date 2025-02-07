package models

type SwiftCode struct {
	CountryISO2   string `csv:"COUNTRY ISO2 CODE" json:"countryISO2"`
	SwiftCode     string `csv:"SWIFT CODE" json:"swiftCode"`
	CodeType      string `csv:"CODE TYPE" json:"codeType"`
	Name          string `csv:"NAME" json:"bankName"`
	Address       string `csv:"ADDRESS" json:"address"`
	TownName      string `csv:"TOWN NAME" json:"townName"`
	CountryName   string `csv:"COUNTRY NAME" json:"countryName"`
	TimeZone      string `csv:"TIME ZONE" json:"timeZone"`
	IsHeadquarter bool   `json:"isHeadquarter"`
}
