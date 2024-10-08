package handler

import (
	"fmt"
	"net/http"
	"rent-car/api/models"
	"rent-car/pkg/check"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CreateCustomer godoc
// @Router      /customer [POST]
// @Summary     Create a customer
// @Description Create a new customer
// @Tags        customer
// @Accept      json
// @Produce 	json
// @Param 		customer body models.CreateCustomer true "customer"
// @Success 	200  {object}  string
// @Failure		400  {object}  models.Response
// @Failure		404  {object}  models.Response
// @Failure		500  {object}  models.Response
func (h Handler) CreateCustomer(c *gin.Context) {
	customer := models.CreateCustomer{}

	if err := c.ShouldBindJSON(&customer); err != nil {
		handleResponseLog(c, h.Log, "error while decoding request body", http.StatusBadRequest, err.Error())
		return
	}

	if _, err := check.ValidateEmail(customer.Email); err != nil {
		handleResponseLog(c, h.Log, "error while validating email"+customer.Email, http.StatusBadRequest, err.Error())
		return
	}

	if _, err := check.ValidatePhone(customer.Phone); err != nil {
		handleResponseLog(c, h.Log, "error while validating phone", http.StatusBadRequest, err.Error())
		return
	}

	if err := check.ValidatePassword(customer.Password); err != nil {
		handleResponseLog(c, h.Log, "error while validating password", http.StatusBadRequest, err.Error())
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(
		customer.Password,
	), bcrypt.DefaultCost)
	if err != nil {
		handleResponseLog(c, h.Log, "error while generating customer password", http.StatusInternalServerError, err.Error())
		return
	}
	customer.Password = string(hashedPass)

	id, err := h.Services.Customer().Create(c.Request.Context(), customer)
	if err != nil {
		handleResponseLog(c, h.Log, "error while creating customer", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "Customer was successfully created", http.StatusOK, id)
}

// UpdateCustomer godoc
// @Security ApiKeyAuth
// @Router		/customer/{id} [PUT]
// @Summary		update a customer
// @Description This api updates a customer by its id and returns its id
// @Tags		customer
// @Accept		json
// @Produce		json
// @Param 		id path string true "Customer ID"
// @Param		customer body models.UpdateCustomer true "customer"
// @Success		200  {object}  string
// @Failure		400  {object}  models.Response
// @Failure		404  {object}  models.Response
// @Failure		500  {object}  models.Response
func (h Handler) UpdateCustomer(c *gin.Context) {
	customer := models.UpdateCustomer{}

	if err := c.ShouldBindJSON(&customer); err != nil {
		handleResponseLog(c, h.Log, "error while decoding request body", http.StatusBadRequest, err.Error())
		return
	}

	id := c.Param("id")

	err := uuid.Validate(id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while validating id"+id, http.StatusBadRequest, err.Error())
		return
	}
	if _, err := check.ValidateEmail(customer.Email); err != nil {
		handleResponseLog(c, h.Log, "error while validating email"+customer.Email, http.StatusBadRequest, err.Error())
		return
	}

	if _, err := check.ValidatePhone(customer.Phone); err != nil {
		handleResponseLog(c, h.Log, "error while validating phone", http.StatusBadRequest, err.Error())
		return
	}

	ID, err := h.Services.Customer().Update(c.Request.Context(), customer, id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while updating customer", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "Customer was successfully updated", http.StatusOK, ID)
}

// GetCustomerById godoc
// @Security ApiKeyAuth
// @Router		/customer/{id} [GET]
// @Summary		get a customer by its id
// @Description This api gets a customer by its id and returns its info
// @Tags		customer
// @Accept		json
// @Produce		json
// @Param		id path string true "customer"
// @Success		200  {object}  models.Customer
// @Failure		400  {object}  models.Response
// @Failure		404  {object}  models.Response
// @Failure		500  {object}  models.Response
func (h Handler) GetCustomerByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		handleResponseLog(c, h.Log, "missing car ID", http.StatusBadRequest, id)
		return
	}

	customer, err := h.Services.Customer().GetByID(c.Request.Context(), id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while getting customer by ID", http.StatusBadRequest, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "Customer was successfully gotten by Id", http.StatusOK, customer)
}

// GetAllCustomers godoc
// @Security ApiKeyAuth
// @Router 			/customer [GET]
// @Summary 		Get all customers
// @Description		Retrieves information about all customers.
// @Tags 			customer
// @Accept 			json
// @Produce 		json
// @Param 			search query string true "customers"
// @Param 			page query uint64 false "page"
// @Param 			limit query uint64 false "limit"
// @Success 		200 {object} models.GetAllCustomersResponse
// @Failure 		400 {object} models.Response
// @Failure 		500 {object} models.Response
func (h Handler) GetAllCustomers(c *gin.Context) {
	var (
		req = models.GetAllCustomersRequest{}
	)

	req.Search = c.Query("search")

	page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		handleResponseLog(c, h.Log, "error while parsing page", http.StatusBadRequest, err.Error())
		return
	}

	limit, err := strconv.ParseUint(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		handleResponseLog(c, h.Log, "error while parsing limit", http.StatusBadRequest, err.Error())
		return
	}

	req.Page = page
	req.Limit = limit

	customers, err := h.Services.Customer().GetAll(c.Request.Context(), req)
	if err != nil {
		handleResponseLog(c, h.Log, "error while getting customers", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "Customers were successfully gotten by Id", http.StatusOK, customers)
}

// GetCustomerCars godoc
// @Security ApiKeyAuth
// @Router		/customer/cars [GET]
// @Summary		get customer's cars
// @Description This api gets customer cars and returns their info
// @Tags		customer
// @Accept		json
// @Produce		json
// @Param 		customerID query string true "Customer ID"
// @Param 		carName query string false "Car Name"
// @Success		200  {object}  models.GetCustomerCarsResponse
// @Failure		400  {object}  models.Response
// @Failure		404  {object}  models.Response
// @Failure		500  {object}  models.Response
func (h Handler) GetCustomerCars(c *gin.Context) {
	customerID := c.Query("customerID")
	carName := c.Query("carName")

	if customerID == "" && carName == "" {
		handleResponseLog(c, h.Log, "missing customerID or carName", http.StatusBadRequest, "")
		return
	}

	var customer models.GetCustomerCarsResponse
	var err error

	if customerID != "" && carName == "" {
		customer, err = h.Services.Customer().GetCustomerCars(c.Request.Context(), "", customerID, true)
		if err != nil {
			handleResponseLog(c, h.Log, "error while getting customer cars by Customer ID", http.StatusInternalServerError, err.Error())
			return
		}
	} else if carName != "" && customerID == "" {
		customer, err = h.Services.Customer().GetCustomerCars(c.Request.Context(), carName, "", false)
		if err != nil {
			handleResponseLog(c, h.Log, "error while getting customer cars by Car Name", http.StatusInternalServerError, err.Error())
			return
		}
	} else if carName != "" && customerID != "" {
		customer, err = h.Services.Customer().GetCustomerCars(c.Request.Context(), carName, customerID, false)
		if err != nil {
			handleResponseLog(c, h.Log, "error while getting customer cars by Car Name", http.StatusInternalServerError, err.Error())
			return
		}
	}

	handleResponseLog(c, h.Log, "Customer's cars were successfully gotten", http.StatusOK, customer)
}

// DeleteCustomer godoc
// @Security ApiKeyAuth
// @Router		/customer/{id} [DELETE]
// @Summary		delete a customer by its id
// @Description This api deletes a customer by its id and returns error or nil
// @Tags		customer
// @Accept		json
// @Produce		json
// @Param		id path string true "customer ID"
// @Success		200  {object}  nil
// @Failure		400  {object}  models.Response
// @Failure		404  {object}  models.Response
// @Failure		500  {object}  models.Response
func (h Handler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("id: ", id)

	if id == "" {
		handleResponseLog(c, h.Log, "missing car ID", http.StatusBadRequest, id)
		return
	}

	err := uuid.Validate(id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while validating id", http.StatusBadRequest, err.Error())
		return
	}

	err = h.Services.Customer().Delete(c.Request.Context(), id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while deleting customer", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "Customer was successfully deleted/updated by Id", http.StatusOK, id)
}
