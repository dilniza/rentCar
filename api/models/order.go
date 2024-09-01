package models

type Order struct {
	Id        string `json:"id"`
	FromDate  string `json:"from_date"`
	ToDate    string `json:"to_date"`
	Status    string `json:"status"`
	Paid      bool   `json:"payment_status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateOrder struct {
	CarId      string `json:"car_id"`
	CustomerId string `json:"customer_id"`
	FromDate   string `json:"from_date"`
	ToDate     string `json:"to_date"`
	Status     string `json:"status"`
	Paid       bool   `json:"payment_status"`
}

type UpdateOrder struct {
	Id         string `json:"id"`
	CarId      string `json:"car_id"`
	CustomerId string `json:"customer_id"`
	FromDate   string `json:"from_date"`
	ToDate     string `json:"to_date"`
	Status     string `json:"status"`
	Paid       bool   `json:"payment_status"`
}

type GetOrderRequest struct {
	Id string `json:"id"`
}

type GetOrderResponse struct {
	Id        string      `json:"id"`
	Car       GetCar      `json:"car,omitempty"`
	Customer  GetCustomer `json:"customer,omitempty"`
	FromDate  string      `json:"from_date"`
	ToDate    string      `json:"to_date"`
	Status    string      `json:"status"`
	Paid      bool        `json:"payment_status"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

type GetAllOrdersRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllOrdersResponse struct {
	Orders []GetOrderResponse `json:"orders"`
	Count  int                `json:"count"`
}

type UpdateOrderStatus struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type UpdateStatus struct {
	OrderNumber    string  `json:"order_number"`
	ClientFullName string  `json:"client_full_name"`
	ClientPhone    string  `json:"client_phone"`
	Price          float64 `json:"price"`
	FromStatus     string  `json:"from_status"`
	ToStatus       string  `json:"to_status"`
	CarName        string  `json:"car_name"`
	FromDate       string  `json:"from_date"`
	ToDate         string  `json:"to_date"`
	Paid           bool    `json:"paid"`
}
