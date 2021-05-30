package http

import (
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/medivh13/evermos-test/internal/services"
	evermosConst "github.com/medivh13/evermos-test/pkg/common/const"
	"github.com/medivh13/evermos-test/pkg/common/crypto"
	"github.com/medivh13/evermos-test/pkg/dto"
	btbErrors "github.com/medivh13/evermos-test/pkg/errors"
	util "github.com/medivh13/evermos-test/pkg/utils"

	"github.com/labstack/echo"
)

type HttpHandler struct {
	service services.Services
}

func NewHttpHandler(e *echo.Echo, srv services.Services) {
	handler := &HttpHandler{
		srv,
	}

	e.GET("api/evermos-test/ping", handler.Ping)
	e.POST("api/evermos-test/register", handler.PostRegister)
	e.POST("api/evermos-test/login", handler.PostLogin)
	e.GET("api/evermos-test/products", handler.GetAllProducts)
	e.POST("api/evermos-test/order", handler.PostCheckOut)
	e.POST("api/evermos-test/expired", handler.CancelExpiredOrder) //this suppose to hit by CronJob
	e.POST("api/evermos-test/cancel", handler.PostCancelByID)

}

func (h *HttpHandler) Ping(c echo.Context) error {

	version := os.Getenv("VERSION")
	if version == "" {
		version = "pong"
	}

	data := version

	return c.JSON(http.StatusOK, data)

}

func (h *HttpHandler) PostRegister(c echo.Context) error {

	postDTO := dto.RegisterReqDTO{}

	if err := c.Bind(&postDTO); err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}
	postDTO.Password = crypto.EncodeSHA256HMAC(util.GetPrivKeySignature(), postDTO.Password)
	postDTO.Signature = c.Request().Header.Get("signature")

	err := h.service.Register(&postDTO)

	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: evermosConst.RegistrastionSuccess,
		Data:    nil,
	}

	return c.JSON(http.StatusOK, resp)

}

func (h *HttpHandler) PostLogin(c echo.Context) error {

	postDTO := dto.RegisterReqDTO{}

	if err := c.Bind(&postDTO); err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	postDTO.Password = crypto.EncodeSHA256HMAC(util.GetPrivKeySignature(), postDTO.Password)
	postDTO.Signature = c.Request().Header.Get("signature")

	data, err := h.service.Login(&postDTO)

	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: evermosConst.LoginSuccess,
		Data:    data,
	}

	return c.JSON(http.StatusOK, resp)

}

func (h *HttpHandler) GetAllProducts(c echo.Context) error {

	postDTO := dto.GetAllProductReqDTO{}

	if err := c.Bind(&postDTO); err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	postDTO.Tokens = c.Request().Header.Get("Authorization")
	postDTO.Signature = c.Request().Header.Get("signature")

	data, err := h.service.GetAllProducts(&postDTO)

	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: evermosConst.GetDataSuccess,
		Data:    data,
	}

	return c.JSON(http.StatusOK, resp)

}

func (h *HttpHandler) PostCheckOut(c echo.Context) error {

	postDTO := dto.CheckOutReqDTO{}

	if err := c.Bind(&postDTO); err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	postDTO.Tokens = c.Request().Header.Get("Authorization")
	postDTO.Signature = c.Request().Header.Get("signature")

	data, err := h.service.ChekOut(&postDTO)

	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: evermosConst.GetDataSuccess,
		Data:    data,
	}

	return c.JSON(http.StatusOK, resp)

}

func (h *HttpHandler) CancelExpiredOrder(c echo.Context) error {

	err := h.service.CancelExpiredOrder()

	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: evermosConst.UpdateStatusSuccess,
		Data:    nil,
	}

	return c.JSON(http.StatusOK, resp)

}

func (h *HttpHandler) PostCancelByID(c echo.Context) error {

	postDTO := dto.CancelByIDReqDTO{}

	if err := c.Bind(&postDTO); err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	postDTO.Tokens = c.Request().Header.Get("Authorization")
	postDTO.Signature = c.Request().Header.Get("signature")

	err := h.service.CancelByID(&postDTO)

	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: evermosConst.UpdateStatusSuccess,
		Data:    nil,
	}

	return c.JSON(http.StatusOK, resp)

}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case btbErrors.ErrInternalServerError:
		return http.StatusInternalServerError
	case btbErrors.ErrNotFound:
		return http.StatusNotFound
	case btbErrors.ErrConflict:
		return http.StatusConflict
	case btbErrors.ErrInvalidRequest:
		return http.StatusBadRequest
	case btbErrors.ErrFailAuth:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
