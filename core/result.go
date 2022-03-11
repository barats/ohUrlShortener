package core

import (
	"net/http"
	"time"
)

type ResultJson struct {
	Code    int
	Message string
	Result  interface{}
	Date    time.Time
}

func ResultJsonSuccess() ResultJson {
	return ResultJson{
		Code:    http.StatusOK,
		Message: "success",
		Result:  nil,
		Date:    time.Now(),
	}
}

func ResultJsonSuccessWithData(data interface{}) ResultJson {
	return ResultJson{
		Code:    http.StatusOK,
		Message: "success",
		Result:  data,
		Date:    time.Now(),
	}
}

func ResultJsonError(message string) ResultJson {
	return ResultJson{
		Code:    http.StatusInternalServerError,
		Message: message,
		Result:  nil,
		Date:    time.Now(),
	}
}

func ResultJsonBadRequest(message string) ResultJson {
	return ResultJson{
		Code:    http.StatusBadRequest,
		Message: message,
		Result:  nil,
		Date:    time.Now(),
	}
}
