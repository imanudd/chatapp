package helper

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseHandler struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func ErrBadRequest(c echo.Context, err error) error {
	resp := ResponseHandler{
		StatusCode: http.StatusBadRequest,
		Message:    err.Error(),
		Data:       nil,
	}

	return c.JSON(http.StatusBadRequest, resp)
}

func ErrNotFound(c echo.Context, err error) error {
	resp := ResponseHandler{
		StatusCode: http.StatusNotFound,
		Message:    err.Error(),
		Data:       nil,
	}

	return c.JSON(http.StatusNotFound, resp)
}

func ErrInternalServer(c echo.Context, err error) error {
	resp := ResponseHandler{
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
		Data:       nil,
	}

	return c.JSON(http.StatusInternalServerError, resp)
}

func Success(c echo.Context, data interface{}) error {
	resp := ResponseHandler{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data:       data,
	}

	return c.JSON(resp.StatusCode, resp)
}
