package models

type SwiftCode struct {
	CountryISO2   string `csv:"COUNTRY ISO2 CODE"`
	SwiftCode     string `csv:"SWIFT CODE"`
	CodeType      string `csv:"CODE TYPE"`
	Name          string `csv:"NAME"`
	Address       string `csv:"ADDRESS"`
	TownName      string `csv:"TOWN NAME"`
	CountryName   string `csv:"COUNTRY NAME"`
	TimeZone      string `csv:"TIME ZONE"`
	IsHeadquarter bool
}