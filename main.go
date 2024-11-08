package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"go-store/db"
	"go-store/templates"
	"go-store/types"
	"strings"

	etag "github.com/pablor21/echo-etag/v4"
)

var conn *sql.DB
var CustomerResults types.CustomerResults
var Orders []types.Order
var checker = ""
var Products []types.Product

func initialize(conn *sql.DB) {
	var err error
	//TABLE VARS
	CustomerResults.Customers, err = db.GetAllCustomers(conn)
	if err != nil {
	}
	Orders, err = db.GetAllOrdersFriendly(conn)
	if err != nil || len(Orders) == 0 {
		checker = "No orders yet!"
	}
	Products, err = db.GetAllProducts(conn)
	if err != nil {
	}

}

func main() {
	var err error

	cfg := mysql.Config{
		User:   "mulholland",
		Passwd: "226466",
		DBName: "mulholland",
	}

	conn, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()

	initialize(conn)

	e := echo.New()
	e.Use(etag.Etag())
	e.Static("assets", "./assets")

	// TODO: Render your base store page here
	e.GET("/store", func(ctx echo.Context) error {
		initialize(conn)

		return Render(ctx, http.StatusOK, templates.Base(templates.Store()))
	})

	e.GET("/dbQueries", func(ctx echo.Context) error {
		initialize(conn)

		CustomerResults.Customer2, err = db.GetCustomerByID(conn, 2)
		CustomerResults.Customer3Find, err = db.GetCustomerByID(conn, 3)
		if err != nil {
			CustomerResults.Customer3 = "Customer 2 not found!"
		}
		CustomerResults.Customer4, err = db.GetCustomerByEmail(conn, "dt@mines.edu")
		CustomerResults.Customer5, err = db.GetCustomerByEmail(conn, "nobody@gmail.com")
		if err != nil {
			CustomerResults.Customer5Find = "Customer nobody@gmail.com not found... adding!"
		}
		db.AddCustomer(conn, "No", "Body", "nobody@gmail.com")
		CustomerResults.Customer5, err = db.GetCustomerByEmail(conn, "nobody@gmail.com")
		db.AddOrder(conn, 1, 1, 1, 1.00, 0.07, 0.0)
		var QuantityRemaining1 int
		QuantityRemaining1, err = db.SellProduct(conn, 3, 5)
		var QuantityRemaining2 int
		QuantityRemaining2, err = db.SellProduct(conn, 3, 10)

		return Render(ctx, http.StatusOK, templates.Base(templates.Queries(CustomerResults, Orders, Products, QuantityRemaining1, QuantityRemaining2)))
	})

	e.GET("/admin", func(ctx echo.Context) error {
		initialize(conn)

		return Render(ctx, http.StatusOK, templates.Base(templates.Admin(CustomerResults, Orders, Products, checker)))
	})
	e.GET("/get_product_quantity", func(ctx echo.Context) error {
		// Retrieve the product name from the query parameters
		product := ctx.QueryParam("product")

		// Mock function to get quantity based on the product name
		// Replace this with your actual function for retrieving product quantity
		quantity, err := db.GetProductQuantity(conn, product)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Product not found"})
		}

		// Return the quantity as a JSON response
		return ctx.JSON(http.StatusOK, map[string]int{"quantity": quantity})
	})
	e.GET("/get_customers", func(ctx echo.Context) error {
		// Initialize the database connection
		initialize(conn)

		// Retrieve the search term from the query parameters
		searchTerm := ctx.QueryParam("searchTerm")

		// Call GetAllCustomers to get all customers from the database
		customers, err := db.GetAllCustomers(conn)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Error retrieving customers"})
		}

		// If searchTerm is provided, filter the customers
		var filteredCustomers []types.Customer
		if searchTerm != "" {
			for _, customer := range customers {
				// Search customers based on first or last name (case insensitive)
				if strings.Contains(strings.ToLower(customer.FirstName), strings.ToLower(searchTerm)) ||
					strings.Contains(strings.ToLower(customer.LastName), strings.ToLower(searchTerm)) {
					filteredCustomers = append(filteredCustomers, customer)
				}
			}
		} else {
			// If no searchTerm is provided, return all customers
			filteredCustomers = customers
		}

		// If no customers match the search term, return a message
		if len(filteredCustomers) == 0 {
			return ctx.JSON(http.StatusOK, map[string]string{"message": "No matching customers found"})
		}

		// Return the filtered customers as JSON
		return ctx.JSON(http.StatusOK, filteredCustomers)
	})

	e.GET("/order_entry", func(ctx echo.Context) error {
		initialize(conn)

		return Render(ctx, http.StatusOK, templates.Base(templates.OrderEntry()))

	})

	e.POST("/purchase", func(ctx echo.Context) error {

		firstName := ctx.FormValue("firstName")
		lastName := ctx.FormValue("lastName")
		email := ctx.FormValue("email")
		customerID, _ := db.AddCustomer(conn, firstName, lastName, email)
		product := ctx.FormValue("product")

		var productObj types.Product
		productObj, _ = db.GetProductByName(conn, product)
		var productID = productObj.ID

		quantity := ctx.FormValue("quantity")
		timestamp := ctx.FormValue("timestamp")
		quantityInt, err := strconv.Atoi(quantity)
		if err != nil {
			return ctx.String(http.StatusBadRequest, "Invalid quantity")
		}

		var price = productObj.Price
		subtotal := price * float64(quantityInt)
		tax := subtotal * 0.06
		total := subtotal + tax
		donation := ctx.FormValue("donation") == "yes"
		donationAmount := 3.0 //3 dolla donation
		totalWithDonation := total

		if donation {
			totalWithDonation += donationAmount
		}

		purchaseInfo := types.PurchaseInfo{
			FirstName:         firstName,
			LastName:          lastName,
			Email:             email,
			Quantity:          quantityInt,
			Product:           product,
			Subtotal:          subtotal,
			Tax:               tax,
			Total:             total,
			Donation:          donation,
			DonationAmount:    donationAmount,
			TotalWithDonation: totalWithDonation,
			Timestamp:         timestamp,
		}

		db.AddOrder(conn, productID, customerID, quantityInt, subtotal, tax, totalWithDonation-total)
		return Render(ctx, http.StatusOK, templates.Base(templates.PurchaseConfirmation(purchaseInfo)))
	})

	e.POST("/purchasebrief", func(ctx echo.Context) error {
		firstName := ctx.FormValue("firstName")
		lastName := ctx.FormValue("lastName")
		email := ctx.FormValue("email")
		customerID, _ := db.AddCustomer(conn, firstName, lastName, email)
		product := ctx.FormValue("product")

		var productObj types.Product
		productObj, _ = db.GetProductByName(conn, product)
		var productID = productObj.ID

		quantity := ctx.FormValue("quantity")
		quantityInt, err := strconv.Atoi(quantity)
		if err != nil {
			return ctx.String(http.StatusBadRequest, "Invalid quantity")
		}

		var price = productObj.Price
		subtotal := price * float64(quantityInt)
		tax := subtotal * 0.06
		total := subtotal + tax
		donation := ctx.FormValue("donation") == "yes"
		donationAmount := 3.0 // 3 dolla donation
		totalWithDonation := total

		if donation {
			totalWithDonation += donationAmount
		}

		// Add the order to the database
		db.AddOrder(conn, productID, customerID, quantityInt, subtotal, tax, totalWithDonation-total)

		confirmationMessage := fmt.Sprintf(`
			<h3>Purchase Confirmation</h3>
			<p>Order submitted for: %s %s %d %s %.2f</p>
		`, firstName, lastName, quantityInt, product, total)

		return ctx.HTML(http.StatusOK, confirmationMessage)
	})

	e.Logger.Fatal(e.Start(":8000"))
}

// INFO: This is a simplified render method that replaces `echo`'s with a custom
// one. This should simplify rendering out of an echo route.
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
