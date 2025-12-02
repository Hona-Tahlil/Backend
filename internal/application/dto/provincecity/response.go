package provincecity

import "hona/backend/internal/domain/enums"

type ProvinceResponse struct {
	Num  enums.Province `json:"num"`
	Name string         `json:"name"`
}

type CityResponse struct {
	Num  enums.City `json:"num"`
	Name string     `json:"name"`
}
