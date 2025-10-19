package types

import "time"

const ApiVersionV1 = "v1"

// APIResponse 统一的 API 响应格式
type APIResponse struct {
	Success    bool        `json:"success"`
	Timestamp  int64       `json:"timestamp"`
	ApiVersion string      `json:"api_version"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Error      *APIError   `json:"error,omitempty"`
}

// 分页信息
type Pagination struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// APIError API 错误信息
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// 成功响应
func Success(apiVersion string, data interface{}, pagination *Pagination) *APIResponse {
	res := &APIResponse{
		Success:    true,
		Timestamp:  time.Now().UnixMilli(),
		ApiVersion: apiVersion,
		Data:       data,
	}
	if pagination != nil {
		res.Pagination = pagination
	}
	return res
}

// 错误响应
func Error(code string, message, details string) *APIResponse {
	return &APIResponse{
		Success:   false,
		Timestamp: time.Now().UnixMilli(),
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}
