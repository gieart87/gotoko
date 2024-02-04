package models

type ProvinceResponse struct {
	ProvinceData ProvinceData `json:"rajaongkir"`
}

type ProvinceData struct {
	Results []Province `json:"results"`
}

type Province struct {
	ID   string `json:"province_id"`
	Name string `json:"province"`
}

type CityResponse struct {
	CityData CityData `json:"rajaongkir"`
}

type CityData struct {
	Results []City `json:"results"`
}

type City struct {
	ID         string `json:"city_id"`
	Name       string `json:"city_name"`
	PostalCode string `json:"postal_code"`
	ProvinceID string `json:"province_id"`
}

type OngkirResponse struct {
	OngkirData OngkirData `json:"rajaongkir"`
}

type ShippingFeeParams struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Weight      int    `json:"weight"`
	Courier     string `json:"courier"`
}

type OngkirData struct {
	OriginDetails      OriginDetails      `json:"origin_details"`
	DestinationDetails DestinationDetails `json:"destination_details"`
	Results            []OngkirResult     `json:"results"`
}

type OriginDetails struct {
	CityID   string `json:"city_id"`
	CityName string `json:"city_name"`
}

type DestinationDetails struct {
	CityID   string `json:"city_id"`
	CityName string `json:"city_name"`
}

type OngkirResult struct {
	Code  string       `json:"code"`
	Name  string       `json:"name"`
	Costs []OngkirCost `json:"costs"`
}

type OngkirCost struct {
	Service     string       `json:"service"`
	Description string       `json:"description"`
	Cost        []CostDetail `json:"cost"`
}

type CostDetail struct {
	Value int64  `json:"value"`
	Etd   string `json:"etd"`
	Note  string `json:"note"`
}

type ShippingFeeOption struct {
	Service string `json:"service"`
	Fee     int64  `json:"fee"`
}
