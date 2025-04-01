package response

import "net/http"

type Response struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

var (
	OK = "OK"

	PASSAGE_NOT_FOUND         = "PASSAGE_NOT_FOUND"
	PASSAGE_DATE_FORMAT_ERROR = "PASSAGE_DATE_FORMAT_ERROR"

	CLIENT_REQUEST_ERROR        = "CLIENT_REQUEST_ERROR"
	CLIENT_UNAUTHORIZED         = "CLIENT_UNAUTHORIZED"
	CLIENT_INVALID_ACCESS_TOKEN = "CLIENT_INVALID_ACCESS_TOKEN"

	USER_PERMISSION_DENIED = "USER_PERMISSION_DENIED"

	SERVER_ERROR     = "SERVER_ERROR"
	SERVER_NOT_FOUND = "SERVER_NOT_FOUND"
	SERVER_FORBIDDEN = "SERVER_FORBIDDEN"
)

func Default(httpCode int, success bool, code, msg string, data interface{}) (int, Response) {
	return httpCode, Response{Success: success, Code: code, Msg: msg, Data: data}
}
func Success(data interface{}) (int, Response) {
	return http.StatusOK, Response{Success: true, Code: OK, Msg: "success", Data: data}
}

func InternalServerError(data interface{}) (int, Response) {
	return http.StatusInternalServerError, Response{Success: false, Code: SERVER_ERROR, Msg: "Internal Server Error", Data: data}
}

func NotFound(data interface{}) (int, Response) {
	return http.StatusNotFound, Response{Success: false, Code: SERVER_NOT_FOUND, Msg: "Not Found", Data: data}
}

func ClientBadRequest(data interface{}) (int, Response) {
	return http.StatusBadRequest, Response{Success: false, Code: CLIENT_REQUEST_ERROR, Msg: "Bad Request", Data: data}
}

func Unauthorized(data interface{}) (int, Response) {
	return http.StatusUnauthorized, Response{Success: false, Code: CLIENT_UNAUTHORIZED, Msg: "Unauthorized", Data: data}
}

func Forbidden(data interface{}) (int, Response) {
	return http.StatusForbidden, Response{Success: false, Code: SERVER_FORBIDDEN, Msg: "Forbidden", Data: data}
}
