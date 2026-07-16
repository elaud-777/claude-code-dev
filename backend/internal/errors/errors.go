package errors

import "fmt"

// ApiError represents the standard TaskFlow error response body:
// { error: { code, message, meta? } }
type ApiError struct {
	Status  int            `json:"-"`
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Meta    map[string]any `json:"meta,omitempty"`
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(status int, code, message string) *ApiError {
	return &ApiError{Status: status, Code: code, Message: message}
}

func ValidationError(message string) *ApiError {
	if message == "" {
		message = "올바른 형식이 아닙니다"
	}
	return New(400, "VALIDATION_ERROR", message)
}

func TooLong(limit, actual int) *ApiError {
	return &ApiError{
		Status:  400,
		Code:    "TOO_LONG",
		Message: fmt.Sprintf("메시지는 %d자 이내로 입력하세요", limit),
		Meta:    map[string]any{"limit": limit, "actual": actual},
	}
}

func InvalidCredentials() *ApiError {
	return New(401, "INVALID_CREDENTIALS", "이메일 또는 비밀번호가 일치하지 않습니다")
}

func TokenExpired() *ApiError {
	return New(401, "TOKEN_EXPIRED", "인증이 만료되었습니다")
}

func Forbidden() *ApiError {
	return New(403, "FORBIDDEN", "권한이 없습니다")
}

func NotOwner() *ApiError {
	return New(403, "NOT_OWNER", "본인의 메시지만 삭제할 수 있습니다")
}

func NotFound(message string) *ApiError {
	if message == "" {
		message = "해당 항목을 찾을 수 없습니다"
	}
	return New(404, "NOT_FOUND", message)
}

func EmailTaken() *ApiError {
	return New(409, "EMAIL_TAKEN", "이미 가입된 이메일입니다")
}

func Conflict(message string) *ApiError {
	if message == "" {
		message = "이미 다른 팀에 소속되어 있습니다"
	}
	return New(409, "CONFLICT", message)
}
