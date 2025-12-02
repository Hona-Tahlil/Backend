package address

import "hona/backend/internal/domain/enums"

type AddressInfo struct {
	// Refer         uint
	// Type          string
	ProvinceName  enums.Province
	CityName      enums.City
	StreetAddress string
	HouseNumber   uint
	Unit          uint
	PostalCode    *string
}
