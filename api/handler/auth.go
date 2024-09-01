package handler

import (
	"fmt"
	"net/http"
	"rent-car/api/models"
	"rent-car/pkg/check"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CustomerLogin godoc
// @Router       /customer/login [POST]
// @Summary      Customer login
// @Description  Customer login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body models.CustomerLoginRequest true "login"
// @Success      201  {object}  models.CustomerLoginResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) LoginCustomer(c *gin.Context) {
	loginReq := models.CustomerLoginRequest{}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		handleResponseLog(c, h.Log, "error while binding body", http.StatusBadRequest, err)
		return
	}
	fmt.Println("loginReq: ", loginReq)

	if err := check.ValidatePassword(loginReq.Password); err != nil {
		handleResponseLog(c, h.Log, "error while validating password", http.StatusBadRequest, err.Error())
		return
	}

	loginResp, err := h.Services.Auth().CustomerLogin(c.Request.Context(), loginReq)
	if err != nil {
		handleResponseLog(c, h.Log, "unauthorized", http.StatusUnauthorized, err)
		return
	}

	handleResponseLog(c, h.Log, "Succes", http.StatusOK, loginResp)

}

// CustomerRegister godoc
// @Router       /customer/register [POST]
// @Summary      Customer register
// @Description  Customer register
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        register body models.CustomerRegisterRequest true "register"
// @Success      201  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) CustomerRegister(c *gin.Context) {
	loginReq := models.CustomerRegisterRequest{}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		handleResponseLog(c, h.Log, "error while binding body", http.StatusBadRequest, err)
		return
	}
	fmt.Println("loginReq: ", loginReq)

	if _, err := check.ValidateEmail(loginReq.Mail); err != nil {
		handleResponseLog(c, h.Log, "error while validating email"+loginReq.Mail, http.StatusBadRequest, err.Error())
		return
	}

	err := h.Services.Auth().CustomerRegister(c.Request.Context(), loginReq)
	if err != nil {
		handleResponseLog(c, h.Log, "", http.StatusInternalServerError, err)
		return
	}

	handleResponseLog(c, h.Log, "Otp sent successfull", http.StatusOK, "")
}

// CustomerRegister godoc
// @Router       /customer/register-confirm [POST]
// @Summary      Customer register
// @Description  Customer register
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        register body models.CustomerRegisterConfirm true "register"
// @Success      201  {object}  models.CustomerLoginResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) CustomerRegisterConfirm(c *gin.Context) {
	req := models.CustomerRegisterConfirm{}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponseLog(c, h.Log, "error while binding body", http.StatusBadRequest, err)
		return
	}
	fmt.Println("req: ", req)
	
	if _, err := check.ValidateEmail(req.Mail); err != nil {
		handleResponseLog(c, h.Log, "error while validating email"+req.Mail, http.StatusBadRequest, err.Error())
		return
	}
	
	if err := check.ValidatePassword(req.Customer.Password); err != nil {
		handleResponseLog(c, h.Log, "error while validating password", http.StatusBadRequest, err.Error())
		return
	}
	
	//login validation

	confResp, err := h.Services.Auth().CustomerRegisterConfirm(c.Request.Context(), req)
	if err != nil {
		handleResponseLog(c, h.Log, "error while confirming", http.StatusUnauthorized, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "Succes", http.StatusOK, confResp)

}

// CustomerChangePassword godoc
// @Security ApiKeyAuth
// @Router		/customer/ [PATCH]
// @Summary		customer change password
// @Description This api changes customer password by its login and password and returns message
// @Tags		customer
// @Accept		json
// @Produce		json
// @Param		customer body models.ChangePassword true "Change Customer Password"
// @Success		200  {object}  string
// @Failure		400  {object}  models.Response
// @Failure		404  {object}  models.Response
// @Failure		500  {object}  models.Response
func (h Handler) ChangePasswordCustomer(c *gin.Context) {
	pass := models.ChangePassword{}

	if err := c.ShouldBindJSON(&pass); err != nil {
		handleResponseLog(c, h.Log, "error while decoding request body", http.StatusBadRequest, err.Error())
		return
	}

	if err := check.ValidatePassword(pass.NewPassword); err != nil {
		handleResponseLog(c, h.Log, "error while validating new password", http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		handleResponseLog(c, h.Log, "error while hashing new password", http.StatusInternalServerError, err.Error())
		return
	}
	pass.NewPassword = string(hashedPassword)

	msg, err := h.Services.Auth().ChangePassword(c.Request.Context(), pass)
	if err != nil {
		handleResponseLog(c, h.Log, "error while updating customer", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "Customer password was successfully updated", http.StatusOK, msg)
}
