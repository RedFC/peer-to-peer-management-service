package models

type ErrorResponse struct {
	Message string `json:"message"`
	Trace   string `json:"trace,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	TotalCount int64       `json:"total_count"`
	PageNo     string      `json:"page_no"`
	PageSize   string      `json:"page_size"`
	Message    string      `json:"message"`
}
