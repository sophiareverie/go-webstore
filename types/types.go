package types

type Product struct {
	ID      int
	Name    string
	Image   string
	Price   float64
	InStock int
}

type Customer struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
}

type Order struct {
	ID           int
	ProductID    int
	CustomerID   int
	Quantity     int
	Price        float64
	Tax          float64
	Donation     float64
	Timestamp    int64
	CustomerName string
	ProductName  string
}

type CustomerResults struct {
	Customers     []Customer
	Customer2     Customer
	Customer3Find Customer
	Customer3     string
	Customer4     Customer
	Customer5Find string
	Customer5     Customer
}

type PurchaseInfo struct {
	FirstName         string
	LastName          string
	Email             string
	Quantity          int
	Product           string
	Subtotal          float64
	Tax               float64
	Total             float64
	Donation          bool
	DonationAmount    float64
	TotalWithDonation float64
	Timestamp         string
}
