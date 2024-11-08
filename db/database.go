package db

import (
	"database/sql"
	"errors"
	"fmt"
	"go-store/types"
)

// Define your data structs

// GetAllProducts retrieves all products from the database
func GetAllProducts(conn *sql.DB) ([]types.Product, error) {
	rows, err := conn.Query("SELECT id, product_name, image_name, price, in_stock FROM product")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []types.Product
	for rows.Next() {
		var product types.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Image, &product.Price, &product.InStock); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func SellProduct(conn *sql.DB, productID int, quantity int) (int, error) {
	var inStock int

	err := conn.QueryRow("SELECT in_stock FROM product WHERE id = ?", productID).Scan(&inStock)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("product not found")
		}
		return 0, fmt.Errorf("error retrieving stock: %v", err)
	}

	if quantity > inStock {
		return 0, fmt.Errorf("not enough stock available")
	}

	_, err = conn.Exec("UPDATE product SET in_stock = in_stock - ? WHERE id = ?", quantity, productID)
	if err != nil {
		return 0, fmt.Errorf("error updating stock: %v", err)
	}

	remainingStock := inStock - quantity

	return remainingStock, nil
}

// GetAllCustomers retrieves all customers from the database
func GetAllCustomers(conn *sql.DB) ([]types.Customer, error) {
	rows, err := conn.Query("SELECT id, first_name, last_name, email FROM customer")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []types.Customer
	for rows.Next() {
		var customer types.Customer
		if err := rows.Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Email); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

// GetAllOrders retrieves all orders from the database
func GetAllOrders(conn *sql.DB) ([]types.Order, error) {
	rows, err := conn.Query("SELECT id, product_id, customer_id, quantity, price, tax, donation, timestamp FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []types.Order
	for rows.Next() {
		var order types.Order
		if err := rows.Scan(&order.ID, &order.ProductID, &order.CustomerID, &order.Quantity, &order.Price, &order.Tax, &order.Donation, &order.Timestamp); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func GetAllOrdersFriendly(conn *sql.DB) ([]types.Order, error) {
	query := `
        SELECT 
            o.id, 
            o.quantity, 
            o.price, 
            o.tax, 
            o.donation, 
            o.timestamp,
            c.first_name, 
            c.last_name, 
            p.product_name
        FROM 
            orders o
        JOIN 
            customer c ON o.customer_id = c.id
        JOIN 
            product p ON o.product_id = p.id
    `

	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []types.Order
	for rows.Next() {
		var order types.Order
		var firstName, lastName, productName string

		if err := rows.Scan(
			&order.ID,
			&order.Quantity,
			&order.Price,
			&order.Tax,
			&order.Donation,
			&order.Timestamp,
			&firstName,
			&lastName,
			&productName,
		); err != nil {
			return nil, err
		}

		order.CustomerName = firstName + " " + lastName
		order.ProductName = productName

		orders = append(orders, order)
	}

	return orders, nil
}

func AddOrder(conn *sql.DB, productID, customerID, quantity int, price, tax, donation float64) (int64, error) {
	result, err := conn.Exec("INSERT INTO orders (product_id, customer_id, quantity, price, tax, donation, timestamp) VALUES (?, ?, ?, ?, ?, ?, NOW())", productID, customerID, quantity, price, tax, donation)
	if err != nil {
		return 0, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	_, err = SellProduct(conn, productID, quantity)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// GetCustomerByID retrieves a customer by their ID
func GetCustomerByID(conn *sql.DB, id int) (types.Customer, error) {
	var customer types.Customer
	err := conn.QueryRow("SELECT id, first_name, last_name, email FROM customer WHERE id = ?", id).Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return customer, errors.New("customer not found")
		}
		return customer, err
	}
	return customer, nil
}

// GetCustomerByEmail retrieves a customer by their email
func GetCustomerByEmail(conn *sql.DB, email string) (types.Customer, error) {
	var customer types.Customer
	err := conn.QueryRow("SELECT id, first_name, last_name, email FROM customer WHERE email = ?", email).Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return customer, errors.New("customer not found")
		}
		return customer, err
	}
	return customer, nil
}

func AddCustomer(conn *sql.DB, first string, last string, email string) (int, error) {

	var customerID int
	err := conn.QueryRow(
		"SELECT id FROM customer WHERE first_name = ? AND last_name = ? AND email = ?",
		first, last, email,
	).Scan(&customerID)

	if err != nil {
		if err == sql.ErrNoRows {

			result, err := conn.Exec("INSERT INTO customer (first_name, last_name, email) VALUES (?, ?, ?)", first, last, email)
			if err != nil {
				return 0, err
			}

			newID, err := result.LastInsertId()
			if err != nil {
				return 0, err
			}
			customerID = int(newID)
		} else {
			return 0, err
		}
	} else {
	}

	return customerID, nil
}

func GetProductByName(conn *sql.DB, productName string) (types.Product, error) {
	products, err := GetAllProducts(conn)
	if err != nil {
		return types.Product{}, err
	}

	for _, product := range products {
		if product.Name == productName {
			return product, nil
		}
	}

	return types.Product{}, fmt.Errorf("product with name '%s' not found", productName)
}

func GetProductQuantity(conn *sql.DB, productName string) (int, error) {
	var err error
	query := "SELECT in_stock FROM product WHERE product_name = ?"
	var inStock int

	// Execute the query
	err = conn.QueryRow(query, productName).Scan(&inStock)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("product not found: %s", productName)
		}
		return 0, fmt.Errorf("query failed: %v", err)
	}

	return inStock, nil
}
