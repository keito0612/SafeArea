package response

import (
	"encoding/json"
	"net/http"
)

// ErrorCode はエラーコードを表す
type ErrorCode string

const (
	ErrorCodeValidation     ErrorCode = "VALIDATION_ERROR"
	ErrorCodeAddressNotFound ErrorCode = "ADDRESS_NOT_FOUND"
	ErrorCodeRiskDataNotFound ErrorCode = "RISK_DATA_NOT_FOUND"
	ErrorCodeExternalAPI    ErrorCode = "EXTERNAL_API_ERROR"
	ErrorCodeDatabase       ErrorCode = "DATABASE_ERROR"
	ErrorCodeInternal       ErrorCode = "INTERNAL_ERROR"
)

// ErrorDetail はバリデーションエラーの詳細を表す
type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorBody はエラーレスポンスのボディを表す
type ErrorBody struct {
	Code    ErrorCode     `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

// ErrorResponse はエラーレスポンスを表す
type ErrorResponse struct {
	Success bool      `json:"success"`
	Error   ErrorBody `json:"error"`
}

// SuccessResponse は成功レスポンスを表す
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// errorCodeToHTTPStatus はエラーコードをHTTPステータスコードに変換する
func errorCodeToHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrorCodeValidation:
		return http.StatusBadRequest
	case ErrorCodeAddressNotFound, ErrorCodeRiskDataNotFound:
		return http.StatusNotFound
	case ErrorCodeExternalAPI:
		return http.StatusBadGateway
	case ErrorCodeDatabase, ErrorCodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// WriteError はエラーレスポンスを書き込む
func WriteError(w http.ResponseWriter, code ErrorCode, message string, details []ErrorDetail) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCodeToHTTPStatus(code))

	resp := ErrorResponse{
		Success: false,
		Error: ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
	json.NewEncoder(w).Encode(resp)
}

// WriteValidationError はバリデーションエラーレスポンスを書き込む
func WriteValidationError(w http.ResponseWriter, details []ErrorDetail) {
	WriteError(w, ErrorCodeValidation, "入力内容に問題があります", details)
}

// WriteAddressNotFoundError は住所が見つからないエラーレスポンスを書き込む
func WriteAddressNotFoundError(w http.ResponseWriter) {
	WriteError(w, ErrorCodeAddressNotFound, "住所が見つかりませんでした", nil)
}

// WriteRiskDataNotFoundError はリスクデータが見つからないエラーレスポンスを書き込む
func WriteRiskDataNotFoundError(w http.ResponseWriter) {
	WriteError(w, ErrorCodeRiskDataNotFound, "リスクデータが見つかりませんでした", nil)
}

// WriteExternalAPIError は外部APIエラーレスポンスを書き込む
func WriteExternalAPIError(w http.ResponseWriter, message string) {
	WriteError(w, ErrorCodeExternalAPI, message, nil)
}

// WriteDatabaseError はデータベースエラーレスポンスを書き込む
func WriteDatabaseError(w http.ResponseWriter) {
	WriteError(w, ErrorCodeDatabase, "データベースエラーが発生しました", nil)
}

// WriteInternalError は内部エラーレスポンスを書き込む
func WriteInternalError(w http.ResponseWriter) {
	WriteError(w, ErrorCodeInternal, "内部エラーが発生しました", nil)
}

// WriteSuccess は成功レスポンスを書き込む
func WriteSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := SuccessResponse{
		Success: true,
		Data:    data,
	}
	json.NewEncoder(w).Encode(resp)
}
