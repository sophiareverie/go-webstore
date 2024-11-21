package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"

	"go-store/db"
	"go-store/templates"
	"go-store/types"
	"strings"

	etag "github.com/pablor21/echo-etag/v4"
)

var CustomerResults types.CustomerResults
var OrdersFr []types.Order
var Orders []types.Order
var checker = ""
var Products []types.Product
var sessionuser types.SessionUser

var store = sessions.NewCookieStore([]byte("secret-key"))

func initialize(conn *sql.DB) {
	var err error
	//TABLE VARS
	CustomerResults.Customers, err = db.GetAllCustomers(conn)
	if err != nil {
	}
	Orders, err = db.GetAllOrders(conn)
	if err != nil || len(Orders) == 0 {
		checker = "No orders yet!"
	}
	OrdersFr, err = db.GetAllOrdersFriendly(conn)
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

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, _ := store.Get(c.Request(), "session")
			c.Set("session", session)
			return next(c)
		}
	})

	e.GET("/", func(ctx echo.Context) error {
		initialize(conn)
		errorMsg := ctx.QueryParam("error")

		return Render(ctx, http.StatusOK, templates.Base(sessionuser, templates.Login(errorMsg)))
	})

	e.POST("/logout", func(ctx echo.Context) error {
		sessionuser.First = ""
		sessionuser.Last = ""
		sessionuser.Username = ""
		sessionuser.Role = 0

		session, _ := store.Get(ctx.Request(), "session")
		session.Options.MaxAge = -1
		session.Values = make(map[interface{}]interface{})
		session.Save(ctx.Request(), ctx.Response())

		return ctx.Redirect(http.StatusFound, "/")
	})

	e.POST("/login", func(ctx echo.Context) error {
		initialize(conn)

		email := ctx.FormValue("email")
		password := ctx.FormValue("password")

		// Query the database for the user
		user, err := db.GetUserByEmailAndPassword(conn, email, password)
		if err != nil {
			return ctx.Redirect(http.StatusFound, "/?error=Invalid credentials")
		}
		sessionuser.First = user.FirstName
		sessionuser.Last = user.LastName
		sessionuser.Username = user.Email
		sessionuser.Role = user.Role

		session, _ := store.Get(ctx.Request(), "session")

		session.Values["userID"] = user.ID
		session.Values["email"] = user.Email
		session.Values["role"] = user.Role
		session.Save(ctx.Request(), ctx.Response())

		// Redirect based on user role
		if user.Role == 1 {
			return ctx.Redirect(http.StatusFound, "/order_entry")
		} else if user.Role == 2 {
			return ctx.Redirect(http.StatusFound, "/products")
		}

		return ctx.Redirect(http.StatusFound, "/?error=unknown_role")
	})

	e.POST("/guest", func(ctx echo.Context) error {

		session, err := store.Get(ctx.Request(), "session")
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to initialize session"})
		}

		// Set guest status in session
		session.Values["guest"] = true
		session.Values["role"] = 0

		if err := session.Save(ctx.Request(), ctx.Response()); err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to save session"})
		}

		// Redirect to order entry page
		err = ctx.Redirect(http.StatusFound, "/store")
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to redirect"})
		}

		return nil
	})

	e.GET("/store", func(ctx echo.Context) error {
		initialize(conn)

		return Render(ctx, http.StatusOK, templates.Base(sessionuser, templates.Store(Products)))
	})

	e.GET("/dbQueries", func(ctx echo.Context) error {
		initialize(conn)
		session, _ := store.Get(ctx.Request(), "session")
		role, ok := session.Values["role"].(int)

		if !ok || role == 0 {
			return ctx.Redirect(http.StatusFound, "/?error=Must log in first!")
		}

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

		return Render(ctx, http.StatusOK, templates.Base(sessionuser, templates.Queries(CustomerResults, Orders, Products, QuantityRemaining1, QuantityRemaining2)))
	})

	e.GET("/admin", func(ctx echo.Context) error {
		initialize(conn)
		session, _ := store.Get(ctx.Request(), "session")
		role, ok := session.Values["role"].(int)

		if !ok || role == 0 {
			return ctx.Redirect(http.StatusFound, "/?error=Must log in first!")
		}
		return Render(ctx, http.StatusOK, templates.Base(sessionuser, templates.Admin(CustomerResults, OrdersFr, Products, checker)))
	})

	e.GET("/products", func(ctx echo.Context) error {
		session, _ := store.Get(ctx.Request(), "session")
		role, ok := session.Values["role"].(int)

		if !ok || role == 0 {
			return ctx.Redirect(http.StatusFound, "/?error=Must log in first!")
		}
		if role < 2 {
			return ctx.Redirect(http.StatusFound, "/?error=You are not authorized for that page!")
		}

		return Render(ctx, http.StatusOK, templates.Base(sessionuser, templates.Products(Products)))
	})

	e.GET("/order_entry", func(ctx echo.Context) error {
		initialize(conn)
		session, _ := store.Get(ctx.Request(), "session")
		role, ok := session.Values["role"].(int)

		if !ok || role == 0 {
			return ctx.Redirect(http.StatusFound, "/?error=Must log in first!")
		}

		return Render(ctx, http.StatusOK, templates.Base(sessionuser, templates.OrderEntry(Products)))

	})

	e.GET("/get_product_quantity", func(ctx echo.Context) error {
		product := ctx.QueryParam("product")

		quantity, err := db.GetProductQuantity(conn, product)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Product not found"})
		}

		return ctx.JSON(http.StatusOK, map[string]int{"quantity": quantity})
	})

	e.GET("/get_customers", func(ctx echo.Context) error {
		initialize(conn)

		searchTerm := ctx.QueryParam("searchTerm")

		customers, err := db.GetAllCustomers(conn)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Error retrieving customers"})
		}

		var filteredCustomers []types.Customer
		if searchTerm != "" {
			for _, customer := range customers {
				if strings.Contains(strings.ToLower(customer.FirstName), strings.ToLower(searchTerm)) ||
					strings.Contains(strings.ToLower(customer.LastName), strings.ToLower(searchTerm)) {
					filteredCustomers = append(filteredCustomers, customer)
				}
			}
		} else {
			filteredCustomers = customers
		}

		if len(filteredCustomers) == 0 {
			return ctx.JSON(http.StatusOK, map[string]string{"message": "No matching customers found"})
		}

		return ctx.JSON(http.StatusOK, filteredCustomers)
	})

	e.POST("/product", func(ctx echo.Context) error {
		var formData struct {
			ProductID int     `json:"productID"`
			ItemName  string  `json:"itemName"`
			ItemImage string  `json:"itemImage"`
			Quantity  int     `json:"quantity"`
			Price     float64 `json:"price"`
			Inactive  int     `json:"inactive"`
			Action    string  `json:"action"`
		}

		if err := ctx.Bind(&formData); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("Error binding data: %v", err),
			})
		}
		if formData.Action == "" || formData.Action == "fetch" {
			Products, err := db.GetAllProducts(conn)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to retrieve products.",
				})
			}
			return ctx.JSON(http.StatusOK, map[string]interface{}{
				"products": Products,
			})
		}

		switch formData.Action {
		case "add":
			db.AddProduct(conn, formData.ItemName, formData.ItemImage, formData.Quantity, formData.Price, formData.Inactive)
		case "update":

			db.UpdateProduct(conn, formData.ProductID, formData.ItemName, formData.ItemImage, formData.Quantity, formData.Price, formData.Inactive)
		case "delete":
			db.DeleteProduct(conn, formData.ProductID)
		default:
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid action.",
			})
		}
		Products, err = db.GetAllProducts(conn)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Products failed",
			})
		}

		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"message":  "Product action successful",
			"products": Products,
		})
	})
	e.POST("/product-id", func(ctx echo.Context) error {

		var formData struct {
			ItemName  string  `json:"itemName"`
			ItemImage string  `json:"itemImage"`
			Quantity  int     `json:"quantity"`
			Price     float64 `json:"price"`
			Inactive  int     `json:"inactive"`
		}
		if err := ctx.Bind(&formData); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid data",
			})
		}

		var productID int
		err := conn.QueryRow(`
			SELECT id FROM product
			WHERE product_name = ? AND image_name = ? AND price = ? AND in_stock = ? AND inactive = ?
		`, formData.ItemName, formData.ItemImage, formData.Price, formData.Quantity, formData.Inactive).Scan(&productID)
		if err != nil {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "Product not found",
			})
		}
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"id": productID,
		})
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

		db.AddOrder(conn, productID, customerID, quantityInt, subtotal, tax, 0)

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
