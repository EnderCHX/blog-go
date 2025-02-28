package response

import "net/http"

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) (int, Response) {
	return http.StatusOK, Response{Code: http.StatusOK, Msg: "success", Data: data}
}

func InternalServerError(data interface{}) (int, Response) {
	return http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Msg: "Internal Server Error", Data: data}
}

func NotFound(data interface{}) (int, Response) {
	return http.StatusNotFound, Response{Code: http.StatusNotFound, Msg: "Not Found", Data: data}
}

func BadRequest(data interface{}) (int, Response) {
	return http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: "Bad Request", Data: data}
}

func Unauthorized(data interface{}) (int, Response) {
	return http.StatusUnauthorized, Response{Code: http.StatusUnauthorized, Msg: "Unauthorized", Data: data}
}

func Forbidden(data interface{}) (int, Response) {
	return http.StatusForbidden, Response{Code: http.StatusForbidden, Msg: "Forbidden", Data: data}
}
