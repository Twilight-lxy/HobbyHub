package models

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type SuccessResponse struct {
	SuccessMessage string `json:"successMessage"`
}
