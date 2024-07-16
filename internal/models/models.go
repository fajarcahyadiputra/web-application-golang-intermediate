package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// DB model is the type for database connection values
type DBModel struct {
	DB *sql.DB
}

// Models is the wrapper for all models
type Models struct {
	DB DBModel
}

// return a model type with database connection pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{
			DB: db,
		},
	}
}

// Widget is the type for all widgets
type Widget struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	Image          string    `json:"image"`
	IsRecurring    bool      `json:"is_recurring"`
	PlanID         string    `json:"plan_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// type Order is the type for order
type Order struct {
	ID            int         `json:"id"`
	WidgetID      int         `json:"widget_id"`
	TransactionID int         `json:"transaction_id"`
	CustomerID    int         `json:"customer_id"`
	StatusID      int         `json:"status_id"`
	Quantity      int         `json:"quantity"`
	Amount        int         `json:"amount"`
	Widget        Widget      `json:"widget"`
	Transaction   Transaction `json:"transaction"`
	Customer      Customer    `json:"customer"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// Status for type for all statues
type Status struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TransactionStatus for type for all transaction status
type TransactionStatus struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// customer
type Customer struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Transaction for type for all transaction
type Transaction struct {
	ID                  int       `json:"id"`
	Amount              int       `json:"amount"`
	Currency            string    `json:"currency"`
	LastFour            string    `json:"last_four"`
	ExpiryMonth         int       `json:"expiry_month"`
	ExpiryYear          int       `json:"expiry_year"`
	BankReturnCode      string    `json:"bank_return_code"`
	TransactionStatusID int       `json:"transaction_status_id"`
	PaymentIntent       string    `json:"payment_intent"`
	PaymentMethod       string    `json:"payment_method"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// User for type for all transaction user
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	var widget Widget

	row := m.DB.QueryRowContext(ctx, `
	SELECT 
		id, name, description, inventory_level, price, coalesce(image, ''), is_recurring, plan_id, created_at, updated_at 
	FROM 
		widgets 
	WHERE id=?`, id)
	err := row.Scan(
		&widget.ID,
		&widget.Name,
		&widget.Description,
		&widget.InventoryLevel,
		&widget.Price,
		&widget.Image,
		&widget.IsRecurring,
		&widget.PlanID,
		&widget.CreatedAt,
		&widget.UpdatedAt,
	)
	if err != nil {
		return widget, err
	}
	return widget, nil
}

// insertTransaction insert new transaction and return transaction id
func (m *DBModel) InsertTransaction(txn Transaction) (int, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	stmt := `
		INSERT INTO transactions 
		(amount, currency, last_four, bank_return_code, expiry_month, expiry_year, transaction_status_id, payment_intent, payment_method, created_at, updated_at)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt, txn.Amount,
		txn.Currency,
		txn.LastFour,
		txn.BankReturnCode,
		txn.ExpiryMonth,
		txn.ExpiryYear,
		txn.TransactionStatusID,
		txn.PaymentIntent,
		txn.PaymentMethod,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// insert customer
func (m *DBModel) InsertCustomer(txn Customer) (int, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancle()

	stmt := `
		INSERT INTO customers 
		(first_name, last_name, email, created_at, updated_at)
		VALUES(?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt, txn.FirstName,
		txn.LastName,
		txn.Email,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// insert order
func (m *DBModel) InsertOrder(txn Order) (int, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	stmt := `
		INSERT INTO orders 
		(widget_id, status_id, transaction_id, customer_id, quantity, amount, created_at, updated_at)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt, txn.WidgetID,
		txn.StatusID,
		txn.TransactionID,
		txn.CustomerID,
		txn.Quantity,
		txn.Amount,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// get user by email
func (m *DBModel) GetUserByEmail(email string) (User, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	var user User
	email = strings.ToLower(email)

	row := m.DB.QueryRowContext(ctx, "SELECT id, first_name, last_name, email, password, created_at FROM users WHERE email=?", email)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
func (m *DBModel) UpdatePasswordForUser(u User, hash string) error {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	stmt := `UPDATE users SET password = ? WHERE id=?`
	_, err := m.DB.ExecContext(ctx, stmt, hash, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *DBModel) GetAllOrders() ([]*Order, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	var orders []*Order
	query := `
		SELECT o.id, o.widget_id, o.transaction_id, o.customer_id, o.status_id,
			o.quantity, o.amount, o.created_at, o.updated_at, w.id, w.name, t.id, t.amount,
			t.currency, t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
			t.bank_return_code, c.id, c.first_name, c.last_name, c.email
		FROM orders o
			LEFT JOIN transactions t ON (o.transaction_id=t.id)
			LEFT JOIN widgets w ON (o.widget_id=w.id)
			LEFT JOIN customers c ON (o.customer_id=c.id)
		WHERE w.is_recurring = 0
		ORDER BY o.created_at DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var o Order
		err = rows.Scan(
			&o.ID,
			&o.WidgetID,
			&o.TransactionID,
			&o.CustomerID,
			&o.StatusID,
			&o.Quantity,
			&o.Amount,
			&o.CreatedAt,
			&o.UpdatedAt,

			&o.Widget.ID,
			&o.Widget.Name,

			&o.Transaction.ID,
			&o.Transaction.Amount,
			&o.Transaction.Currency,
			&o.Transaction.LastFour,
			&o.Transaction.ExpiryMonth,
			&o.Transaction.ExpiryYear,
			&o.Transaction.BankReturnCode,
			&o.Transaction.PaymentIntent,

			&o.Customer.ID,
			&o.Customer.FirstName,
			&o.Customer.LastName,
			&o.Customer.Email,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &o)
	}

	defer rows.Close()
	return orders, nil
}
func (m *DBModel) GetAllOrdersPagination(pageSize, page int) ([]*Order, int, int, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	offset := (page - 1) * pageSize

	var orders []*Order
	query := `
		SELECT o.id, o.widget_id, o.transaction_id, o.customer_id, o.status_id,
			o.quantity, o.amount, o.created_at, o.updated_at, w.id, w.name, t.id, t.amount,
			t.currency, t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
			t.bank_return_code, c.id, c.first_name, c.last_name, c.email
		FROM orders o
			LEFT JOIN transactions t ON (o.transaction_id=t.id)
			LEFT JOIN widgets w ON (o.widget_id=w.id)
			LEFT JOIN customers c ON (o.customer_id=c.id)
		WHERE w.is_recurring = 0
		ORDER BY o.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := m.DB.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	for rows.Next() {
		fmt.Println(rows)
		var o Order
		err = rows.Scan(
			&o.ID,
			&o.WidgetID,
			&o.TransactionID,
			&o.CustomerID,
			&o.StatusID,
			&o.Quantity,
			&o.Amount,
			&o.CreatedAt,
			&o.UpdatedAt,

			&o.Widget.ID,
			&o.Widget.Name,

			&o.Transaction.ID,
			&o.Transaction.Amount,
			&o.Transaction.Currency,
			&o.Transaction.LastFour,
			&o.Transaction.ExpiryMonth,
			&o.Transaction.ExpiryYear,
			&o.Transaction.BankReturnCode,
			&o.Transaction.PaymentIntent,

			&o.Customer.ID,
			&o.Customer.FirstName,
			&o.Customer.LastName,
			&o.Customer.Email,
		)
		if err != nil {
			return nil, 0, 0, err
		}

		orders = append(orders, &o)
	}
	defer rows.Close()

	queryCount := `
		SELECT COUNT(o.id) FROM orders o 
		LEFT JOIN widgets w ON (o.widget_id=w.id)
		WHERE w.is_recurring = 0
	`
	var totalRecords int

	countRow := m.DB.QueryRowContext(ctx, queryCount)
	err = countRow.Scan(&totalRecords)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := totalRecords / pageSize

	return orders, lastPage, totalRecords, nil
}

func (m *DBModel) GetAllSubscription() ([]*Order, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	var orders []*Order
	query := `
		SELECT o.id, o.widget_id, o.transaction_id, o.customer_id, o.status_id,
			o.quantity, o.amount, o.created_at, o.updated_at, w.id, w.name, t.id, t.amount,
			t.currency, t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
			t.bank_return_code, c.id, c.first_name, c.last_name, c.email
		FROM orders o
			LEFT JOIN transactions t ON (o.transaction_id=t.id)
			LEFT JOIN widgets w ON (o.widget_id=w.id)
			LEFT JOIN customers c ON (o.customer_id=c.id)
		WHERE w.is_recurring = 1
		ORDER BY o.created_at DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var o Order
		err = rows.Scan(
			&o.ID,
			&o.WidgetID,
			&o.TransactionID,
			&o.CustomerID,
			&o.StatusID,
			&o.Quantity,
			&o.Amount,
			&o.CreatedAt,
			&o.UpdatedAt,

			&o.Widget.ID,
			&o.Widget.Name,

			&o.Transaction.ID,
			&o.Transaction.Amount,
			&o.Transaction.Currency,
			&o.Transaction.LastFour,
			&o.Transaction.ExpiryMonth,
			&o.Transaction.ExpiryYear,
			&o.Transaction.BankReturnCode,
			&o.Transaction.PaymentIntent,

			&o.Customer.ID,
			&o.Customer.FirstName,
			&o.Customer.LastName,
			&o.Customer.Email,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &o)
	}

	defer rows.Close()

	return orders, nil
}
func (m *DBModel) GetAllSubscriptionPagination(pageSize, page int) ([]*Order, int, int, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	offset := (page - 1) * pageSize

	var orders []*Order
	query := `
		SELECT o.id, o.widget_id, o.transaction_id, o.customer_id, o.status_id,
			o.quantity, o.amount, o.created_at, o.updated_at, w.id, w.name, t.id, t.amount,
			t.currency, t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
			t.bank_return_code, c.id, c.first_name, c.last_name, c.email
		FROM orders o
			LEFT JOIN transactions t ON (o.transaction_id=t.id)
			LEFT JOIN widgets w ON (o.widget_id=w.id)
			LEFT JOIN customers c ON (o.customer_id=c.id)
		WHERE w.is_recurring = 1
		ORDER BY o.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := m.DB.QueryContext(ctx, query, pageSize, offset)
	if err != nil {

		return nil, 0, 0, err
	}

	for rows.Next() {
		var o Order
		err = rows.Scan(
			&o.ID,
			&o.WidgetID,
			&o.TransactionID,
			&o.CustomerID,
			&o.StatusID,
			&o.Quantity,
			&o.Amount,
			&o.CreatedAt,
			&o.UpdatedAt,

			&o.Widget.ID,
			&o.Widget.Name,

			&o.Transaction.ID,
			&o.Transaction.Amount,
			&o.Transaction.Currency,
			&o.Transaction.LastFour,
			&o.Transaction.ExpiryMonth,
			&o.Transaction.ExpiryYear,
			&o.Transaction.BankReturnCode,
			&o.Transaction.PaymentIntent,

			&o.Customer.ID,
			&o.Customer.FirstName,
			&o.Customer.LastName,
			&o.Customer.Email,
		)
		if err != nil {

			return nil, 0, 0, err
		}

		orders = append(orders, &o)
	}

	defer rows.Close()

	queryCount := `
		SELECT COUNT(o.id) FROM orders o 
		LEFT JOIN widgets w ON (o.widget_id=w.id)
		WHERE w.is_recurring = 1
	`
	var totalRecords int
	countRow := m.DB.QueryRowContext(ctx, queryCount)
	err = countRow.Scan(&totalRecords)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := totalRecords / pageSize

	return orders, lastPage, totalRecords, nil
}

// GetOrderByID gets one order by id and returns the Order
func (m *DBModel) GetOrderByID(id int) (Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var o Order

	query := `
		select
			o.id, o.widget_id, o.transaction_id, o.customer_id,
			o.status_id, o.quantity, o.amount, o.created_at,
			o.updated_at, w.id, w.name, t.id, t.amount, t.currency,
			t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
			t.bank_return_code, c.id, c.first_name, c.last_name, c.email
		from
			orders o
			left join widgets w on (o.widget_id = w.id)
			left join transactions t on (o.transaction_id = t.id)
			left join customers c on (o.customer_id = c.id)
		where
			o.id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&o.ID,
		&o.WidgetID,
		&o.TransactionID,
		&o.CustomerID,
		&o.StatusID,
		&o.Quantity,
		&o.Amount,
		&o.CreatedAt,
		&o.UpdatedAt,
		&o.Widget.ID,
		&o.Widget.Name,
		&o.Transaction.ID,
		&o.Transaction.Amount,
		&o.Transaction.Currency,
		&o.Transaction.LastFour,
		&o.Transaction.ExpiryMonth,
		&o.Transaction.ExpiryYear,
		&o.Transaction.PaymentIntent,
		&o.Transaction.BankReturnCode,
		&o.Customer.ID,
		&o.Customer.FirstName,
		&o.Customer.LastName,
		&o.Customer.Email,
	)
	if err != nil {
		return o, err
	}

	return o, nil
}

func (m *DBModel) UpdateOrderStatus(id, statusID int) error {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*3)
	defer cancle()

	stmt := `UPDATE orders SET status_id = ? WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, stmt, statusID, id)
	if err != nil {
		return err
	}
	return nil
}
func (m *DBModel) GetAllUsers() ([]*User, error) {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*3)
	defer cancle()
	var users []*User

	query := `
		SELECT id, last_name, first_name, email, created_at, updated_at
			FROM users
		ORDER BY last_name, first_name
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.ID,
			&u.LastName,
			&u.FirstName,
			&u.Email,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	defer rows.Close()
	return users, nil
}
func (m *DBModel) GetOneUser(id int) (User, error) {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*3)
	defer cancle()
	var u User

	query := `
		SELECT id, last_name, first_name, email, created_at, updated_at
			FROM users
		WHERE id=?
		ORDER BY last_name, first_name
	`
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&u.ID,
		&u.LastName,
		&u.FirstName,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}
func (m *DBModel) Edituser(u User) error {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	stmt := `
		UPDATE users SET 
		   first_name=?,
		   last_name=?,
		   email=?,
		   updated_at=?
		WHERE id=?
	`

	_, err := m.DB.ExecContext(ctx, stmt, u.FirstName, u.LastName, u.Email, time.Now(), u.ID)
	if err != nil {
		return err
	}
	return nil

}

func (m *DBModel) Adduser(u User, hash string) error {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	stmt := `
		INSERT INTO users (first_name, last_name,email, password, created_at, updated_at)
		VALUES(?,?,?,?,?,?)
	`

	_, err := m.DB.ExecContext(ctx, stmt, u.FirstName, u.LastName, u.Email, hash, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil

}

func (m *DBModel) DeleteUser(id int) error {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	stmt := `DELETE FROM users WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	stmt = `DELETE FROM tokens WHERE user_id = ?`
	_, err = m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}
