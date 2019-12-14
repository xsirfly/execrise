package utils

import (
	"net/http"
)

func ErrorResp(err error) (code int, resp interface{}) {
	return http.StatusInternalServerError, map[string]interface{}{
		"message": err,
	}
}

func SuccessResp(resp interface{}) (int, interface{}) {
	return http.StatusOK, resp
}