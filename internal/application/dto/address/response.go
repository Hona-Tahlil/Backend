package address

type AddressInfoResponse struct {
	ID            uint    `json:"id"`
	ProvinceName  string  `json:"provinceName"`
	CityName      string  `json:"cityName"`
	StreetAddress string  `json:"streetAddress"`
	HouseNumber   uint    `json:"houseNumber"`
	Unit          uint    `json:"unit"`
	PostalCode    *string `json:"postalCode"`
}