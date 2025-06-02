package dto

type ErrorObj struct {
	Code   string `json:"code"`
	Entity string `json:"entity"`
	Cause  string `json:"cause"`
}

type ServiceListResponse struct {
	Data   *ServiceListData `json:"data,omitempty"`
	Errors []ErrorObj       `json:"errors,omitempty"`
}
