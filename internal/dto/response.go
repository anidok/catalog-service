package dto

type ErrorObj struct {
	Code   string `json:"code"`
	Entity string `json:"entity"`
	Cause  string `json:"cause"`
}

type ServiceListResponse struct {
	Success bool             `json:"success"`
	Data    *ServiceListData `json:"data,omitempty"`
	Errors  []ErrorObj       `json:"errors,omitempty"`
}

type ServiceDetailResponse struct {
	Success bool        `json:"success"`
	Data    *ServiceDTO `json:"data,omitempty"`
	Errors  []ErrorObj  `json:"errors,omitempty"`
}
