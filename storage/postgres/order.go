package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"rent-car/api/models"
	"rent-car/pkg"
	"rent-car/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepo struct {
	db     *pgxpool.Pool
	logger logger.ILogger
}

func NewOrderRepo(db *pgxpool.Pool, log logger.ILogger) OrderRepo {
	return OrderRepo{
		db:     db,
		logger: log,
	}
}

func (o *OrderRepo) Create(ctx context.Context, order models.CreateOrder) (string, error) {
	id := uuid.New().String()

	maxQuery := `SELECT COALESCE(MAX(CAST(SUBSTRING(order_number, 4) AS INTEGER)), 0) FROM orders`
	var max int
	err := o.db.QueryRow(ctx, maxQuery).Scan(&max)
	if err != nil {
		o.logger.Error("failed to get max order number", logger.Error(err))
		return "", err
	}
	orderNumber := pkg.GetSerialId(&max)

	query := `INSERT INTO orders (
		id,
		order_number,
		car_id,
		customer_id,
		from_date,
		to_date,
		status,
		payment_status,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err = o.db.Exec(ctx, query,
		id,
		orderNumber,
		order.CarId,
		order.CustomerId,
		order.FromDate,
		order.ToDate,
		order.Status,
		order.Paid,
	)

	if err != nil {
		o.logger.Error("failed to create order in database", logger.Error(err))
		return "", err
	}

	return id, nil
}

func (o *OrderRepo) Update(ctx context.Context, order models.UpdateOrder) (string, error) {
	query := `UPDATE orders SET
		car_id = $1,
		customer_id = $2,
		from_date = $3,
		to_date = $4,
		status = $5,
		payment_status = $6,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $7 AND deleted_at = 0`

	_, err := o.db.Exec(ctx, query,
		order.CarId,
		order.CustomerId,
		order.FromDate,
		order.ToDate,
		order.Status,
		order.Paid,
		order.Id,
	)

	if err != nil {
		o.logger.Error("failed to update order in database", logger.Error(err))
		return "", err
	}

	return order.Id, nil
}

func (o *OrderRepo) UpdateStatus(ctx context.Context, status models.UpdateOrderStatus) (models.UpdateStatus, error) {
	var (
		updatedOrder   models.UpdateStatus
		orderNumber    sql.NullString
		clientFullName sql.NullString
		clientPhone    sql.NullString
		price          sql.NullFloat64
		carName        sql.NullString
		fromStatus     sql.NullString
		toStatus       sql.NullString
		fromDate       sql.NullString
		toDate         sql.NullString
		paid           sql.NullBool
	)

	query := `SELECT status
		FROM orders
		WHERE id = $1`

	err := o.db.QueryRow(ctx, query, status.Id).Scan(
		&fromStatus,
	)

	if err != nil {
		o.logger.Error("failed to get order current status in database", logger.Error(err))
		return models.UpdateStatus{}, err
	}

	query = `UPDATE orders SET
		status = $2,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $1`

	_, err = o.db.Exec(ctx, query,
		status.Id,
		status.Status,
	)

	if err != nil {
		o.logger.Error("failed to update order STATUS in database", logger.Error(err))
		return models.UpdateStatus{}, err
	}

	query = `SELECT order_number, 
                 (SELECT first_name || ' ' || last_name FROM customers WHERE id = orders.customer_id) AS client_full_name,
                 (SELECT phone FROM customers WHERE id = orders.customer_id) AS client_phone,
                 (SELECT price FROM cars WHERE id = orders.car_id) AS price, 
                 status,
                 (SELECT name FROM cars WHERE id = orders.car_id) AS car_name,
                 from_date, 
                 to_date, 
                 payment_status AS paid
             FROM orders 
             WHERE id = $1`

	err = o.db.QueryRow(ctx, query, status.Id).Scan(
		&orderNumber,
		&clientFullName,
		&clientPhone,
		&price,
		&toStatus,
		&carName,
		&fromDate,
		&toDate,
		&paid,
	)

	if err != nil {
		o.logger.Error("failed to get updated order", logger.Error(err))
		return models.UpdateStatus{}, err
	}

	updatedOrder.OrderNumber = orderNumber.String
	updatedOrder.ClientFullName = clientFullName.String
	updatedOrder.ClientPhone = clientPhone.String
	updatedOrder.Price = price.Float64
	updatedOrder.FromStatus = fromStatus.String
	updatedOrder.ToStatus = toStatus.String
	updatedOrder.CarName = carName.String
	updatedOrder.FromDate = fromDate.String
	updatedOrder.ToDate = toDate.String
	updatedOrder.Paid = paid.Bool

	return updatedOrder, nil
}

func (o *OrderRepo) GetByID(ctx context.Context, id string) (models.GetOrderResponse, error) {
	var (
		order             = models.GetOrderResponse{}
		carName           sql.NullString
		carBrand          sql.NullString
		customerFirstName sql.NullString
		customerLastName  sql.NullString
		customerEmail     sql.NullString
		customerPhone     sql.NullString
		customerAddress   sql.NullString
		fromDate          sql.NullString
		toDate            sql.NullString
		status            sql.NullString
		paid              sql.NullBool
		createdAt         sql.NullString
		updatedAt         sql.NullString
	)

	query := `SELECT
		o.id,
		c.id AS car_id,
		c.name AS car_name,
		c.brand AS car_brand,
		cu.id AS customer_id,
		cu.first_name AS customer_first_name,
		cu.last_name AS customer_last_name,
		cu.email AS customer_email,
		cu.phone AS customer_phone,
		cu.address AS customer_address,
		o.from_date,
		o.to_date,
		o.status,
		o.payment_status,
		o.created_at,
		o.updated_at
	FROM orders o
	JOIN cars c ON o.car_id = c.id
	JOIN customers cu ON o.customer_id = cu.id
	WHERE o.id = $1 AND o.deleted_at = 0`

	row := o.db.QueryRow(ctx, query, id)

	err := row.Scan(
		&order.Id,
		&order.Car.ID,
		&carName,
		&carBrand,
		&order.Customer.ID,
		&customerFirstName,
		&customerLastName,
		&customerEmail,
		&customerPhone,
		&customerAddress,
		&fromDate,
		&toDate,
		&status,
		&paid,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		o.logger.Error("failed to get order by ID from database", logger.Error(err))
		return models.GetOrderResponse{}, err
	}

	order.Car = models.GetCar{
		ID:    order.Car.ID,
		Name:  carName.String,
		Brand: carBrand.String,
	}

	order.Customer = models.GetCustomer{
		ID:        order.Customer.ID,
		FirstName: customerFirstName.String,
		LastName:  customerLastName.String,
		Email:     customerEmail.String,
		Phone:     customerPhone.String,
		Address:   customerAddress.String,
	}

	order.FromDate = fromDate.String
	order.ToDate = toDate.String
	order.Status = status.String
	order.Paid = paid.Bool
	order.CreatedAt = createdAt.String
	order.UpdatedAt = updatedAt.String

	return order, nil
}

func (o *OrderRepo) GetAll(ctx context.Context, req models.GetAllOrdersRequest) (models.GetAllOrdersResponse, error) {
	var (
		resp   = models.GetAllOrdersResponse{}
		filter string
		count  sql.NullInt64
	)

	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter = fmt.Sprintf(` AND (c.name ILIKE '%%%v%%' OR cu.first_name ILIKE '%%%v%%' OR cu.last_name ILIKE '%%%v%%')`, req.Search, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)
	//fmt.Println("filter: ", filter)

	query := `SELECT
		o.id,
		c.id AS car_id,
		c.name AS car_name,
		c.brand AS car_brand,
		cu.id AS customer_id,
		cu.first_name AS customer_first_name,
		cu.last_name AS customer_last_name,
		cu.email AS customer_email,
		cu.phone AS customer_phone,
		cu.address AS customer_address,
		o.from_date,
		o.to_date,
		o.status,
		o.payment_status,
		o.created_at,
		o.updated_at
		FROM orders o
		JOIN cars c ON o.car_id = c.id
		JOIN customers cu ON o.customer_id = cu.id
		WHERE o.deleted_at = 0` + filter

	//.Println(query)

	rows, err := o.db.Query(ctx, query)
	if err != nil {
		o.logger.Error("failed to get all orders from database", logger.Error(err))
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		order := models.GetOrderResponse{
			Car:      models.GetCar{},
			Customer: models.GetCustomer{},
		}

		var (
			carName           sql.NullString
			carBrand          sql.NullString
			customerFirstName sql.NullString
			customerLastName  sql.NullString
			customerEmail     sql.NullString
			customerPhone     sql.NullString
			customerAddress   sql.NullString
			fromDate          sql.NullString
			toDate            sql.NullString
			status            sql.NullString
			paid              sql.NullBool
			createdAt         sql.NullString
			updatedAt         sql.NullString
		)

		err := rows.Scan(
			&order.Id,
			&order.Car.ID,
			&carName,
			&carBrand,
			&order.Customer.ID,
			&customerFirstName,
			&customerLastName,
			&customerEmail,
			&customerPhone,
			&customerAddress,
			&fromDate,
			&toDate,
			&status,
			&paid,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			o.logger.Error("failed to scan all orders from database", logger.Error(err))
			return resp, err
		}

		order.Car.Name = carName.String
		order.Car.Brand = carBrand.String
		order.Customer.FirstName = customerFirstName.String
		order.Customer.LastName = customerLastName.String
		order.Customer.Email = customerEmail.String
		order.Customer.Phone = customerPhone.String
		order.Customer.Address = customerAddress.String
		order.FromDate = fromDate.String
		order.ToDate = toDate.String
		order.Status = status.String
		order.Paid = paid.Bool
		order.CreatedAt = createdAt.String
		order.UpdatedAt = updatedAt.String

		resp.Orders = append(resp.Orders, order)
	}

	if err = rows.Err(); err != nil {
		o.logger.Error("failed to get all orders from database", logger.Error(err))
		return resp, err
	}

	countQuery := `SELECT COUNT(id) FROM orders`
	err = o.db.QueryRow(ctx, countQuery).Scan(&count)
	resp.Count = int(count.Int64)
	if err != nil {
		o.logger.Error("failed to get count of orders from database", logger.Error(err))
		return resp, err
	}

	return resp, nil
}

func (o *OrderRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE orders SET deleted_at = date_part('epoch', CURRENT_TIMESTAMP)::int WHERE id = $1 AND deleted_at = 0`

	_, err := o.db.Exec(ctx, query, id)
	if err != nil {
		o.logger.Error("failed to delete order from database", logger.Error(err))
		return err
	}

	return nil
}

func (o *OrderRepo) DeleteHard(ctx context.Context, id string) error {
	query := `DELETE FROM orders WHERE id = $1`

	_, err := o.db.Exec(ctx, query, id)
	if err != nil {
		o.logger.Error("failed to hard delete order from database", logger.Error(err))
		return err
	}

	return nil
}
