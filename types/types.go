package types

type Variant struct {
	ID           int      `json:"id"`
	ProductID    int      `json:"product_id"`
	Title        string   `json:"title"`
	Price        int      `json:"price"`
	ComparePrice int      `json:"compare_price,omitempty"`
	Tax          any      `json:"tax"`
	Shipping     any      `json:"shipping"`
	Weight       any      `json:"weight"`
	Length       any      `json:"length"`
	Width        any      `json:"width"`
	Height       any      `json:"height"`
	Stock        int      `json:"stock"`
	Minimum      any      `json:"minimum"`
	Maximum      any      `json:"maximum"`
	Sku          string   `json:"sku"`
	Image        any      `json:"image"`
	Type         string   `json:"type"`
	Status       []string `json:"status"`
	Files        any      `json:"files"`
}
type VariantResult struct {
	Success bool    `json:"success"`
	Variant Variant `json:"variant"`
}

type Person struct {
	PersonId   int    `db:"PersonID"`
	FirstName  string `db:"PersonFirstName"`
	LastName   string `db:"PersonLastName"`
	Address    string `db:"PersonAddress"`
	Mobile     string `db:"PersonMobile"`
	PostalCode string `db:"PostalCode"`
}

type GivItems struct {
	ItemPrice          float64 `db:"ItemCurrentSelPrice"`
	ItemQuantityOnHand float64 `db:"ItemQuantityOnHand"`
	ItemID             string  `db:"ItemID"`
	VariantId          int     `db:"VariantID"`
}
