package types

import (
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"strconv"

	"github.com/google/uuid"
)

var LAST_DASHT_PURCHASE = "LAST_DASHT_PUARCHACE"
var LAST_DASHT_CREATED = "LAST_DASHT_CREATED"
var DASHT_ACCESS_TOKEN = "DASHT_ACCESS_TOKEN"
var LAST_PORTAL_PURCHASE = "LAST_PORTAL_PURCHASE"
var Debug = false

type PortalProductResult struct {
	Success bool          `json:"success"`
	Product PortalProduct `json:"product"`
}
type PortalProduct struct {
	ID                int      `json:"id"`
	Version           string   `json:"version"`
	Title             string   `json:"title"`
	Caption           string   `json:"caption"`
	Description       string   `json:"description"`
	Image             string   `json:"image"`
	Slug              string   `json:"slug"`
	URL               string   `json:"url"`
	Rate              any      `json:"rate"`
	RateCount         any      `json:"rate_count"`
	Password          any      `json:"password"`
	Layout            *string  `json:"layout"`
	CommentingEnabled bool     `json:"commenting_enabled"`
	MetaTitle         any      `json:"meta_title"`
	MetaDescription   any      `json:"meta_description"`
	MetaKeywords      any      `json:"meta_keywords"`
	MetaRobots        any      `json:"meta_robots"`
	CanonicalURL      any      `json:"canonical_url"`
	Redirect          any      `json:"redirect"`
	Stats             int      `json:"stats"`
	Comments          any      `json:"comments"`
	Position          int      `json:"position"`
	Status            []string `json:"status"`
	Contents          any      `json:"contents"`
	Fields            any      `json:"fields"`
	Images            []struct {
		Path  string `json:"path"`
		Title any    `json:"title"`
	} `json:"images"`
	Category   any `json:"category"`
	Categories []struct {
		ID     int    `json:"id"`
		Title  string `json:"title"`
		URL    string `json:"url"`
		Weight int    `json:"weight"`
	} `json:"categories"`
	Filters    []any `json:"filters"`
	Attributes any   `json:"attributes"`
	Variants   []struct {
		ID           int      `json:"id"`
		ProductID    int      `json:"product_id"`
		Title        string   `json:"title"`
		Price        int      `json:"price"`
		ComparePrice int      `json:"compare_price"`
		Tax          any      `json:"tax"`
		Shipping     any      `json:"shipping"`
		Weight       int      `json:"weight"`
		Length       any      `json:"length"`
		Width        any      `json:"width"`
		Height       any      `json:"height"`
		Stock        int      `json:"stock"`
		Minimum      int      `json:"minimum"`
		Maximum      int      `json:"maximum"`
		Sku          string   `json:"sku"`
		Image        string   `json:"image"`
		Type         string   `json:"type"`
		Status       []string `json:"status"`
		Files        any      `json:"files"`
	} `json:"variants"`
	Relates    any `json:"relates"`
	Expiration any `json:"expiration"`
	Published  struct {
		Year      string `json:"year"`
		Month     string `json:"month"`
		MonthName string `json:"month_name"`
		Day       string `json:"day"`
		Date      string `json:"date"`
		Time      string `json:"time"`
		Universal string `json:"universal"`
		Timestamp int    `json:"timestamp"`
		Subtract  string `json:"subtract"`
		Past      bool   `json:"past"`
	} `json:"published"`
	Created struct {
		Year      string `json:"year"`
		Month     string `json:"month"`
		MonthName string `json:"month_name"`
		Day       string `json:"day"`
		Date      string `json:"date"`
		Time      string `json:"time"`
		Universal string `json:"universal"`
		Timestamp int    `json:"timestamp"`
		Subtract  string `json:"subtract"`
		Past      bool   `json:"past"`
	} `json:"created"`
	Creator struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Nickname string `json:"nickname"`
		Avatar   any    `json:"avatar"`
	} `json:"creator"`
}

type OrderDetail struct {
	Success bool `json:"success"`
	Order   struct {
		ID             int      `json:"id"`
		Description    any      `json:"description"`
		Status         []string `json:"status"`
		Quantity       int      `json:"quantity"`
		Weight         int      `json:"weight"`
		Shipping       int      `json:"shipping"`
		Subtotal       int      `json:"subtotal"`
		Discount       int      `json:"discount"`
		Tax            int      `json:"tax"`
		Price          int      `json:"price"`
		RemainingPrice int      `json:"remaining_price"`
		IP             string   `json:"ip"`
		Contact        struct {
			Name    string `json:"name"`
			Mobile  string `json:"mobile"`
			Phone   any    `json:"phone"`
			Email   any    `json:"email"`
			Country struct {
				ID        int     `json:"id"`
				Name      string  `json:"name"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"country"`
			State struct {
				ID        int     `json:"id"`
				Name      string  `json:"name"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"state"`
			City struct {
				ID        int     `json:"id"`
				Name      string  `json:"name"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"city"`
			Zipcode   string  `json:"zipcode"`
			Address   string  `json:"address"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"contact"`
		Items []struct {
			Variant *struct {
				ID  int `json:"id"`
				Sku any `json:"sku"`
			} `json:"variant"`
			Product *struct {
				ID int `json:"id"`
			} `json:"product"`
			Title    string  `json:"title"`
			Price    int     `json:"price"`
			Quantity int     `json:"quantity"`
			Weight   any     `json:"weight"`
			Shipping any     `json:"shipping"`
			Discount any     `json:"discount"`
			Sku      *string `json:"sku"`
			Tax      int     `json:"tax"`
		} `json:"items"`
		Coupons   any `json:"coupons"`
		Shipments any `json:"shipments"`
		Payments  []*struct {
			ID          int      `json:"id"`
			Description any      `json:"description"`
			ReferenceID string   `json:"reference_id"`
			Type        string   `json:"type"`
			Status      []string `json:"status"`
			SubStatus   any      `json:"sub_status"`
			Amount      int      `json:"amount"`
			Created     struct {
				Year      string `json:"year"`
				Month     string `json:"month"`
				MonthName string `json:"month_name"`
				Day       string `json:"day"`
				Date      string `json:"date"`
				Time      string `json:"time"`
				Universal string `json:"universal"`
				Timestamp int    `json:"timestamp"`
				Subtract  string `json:"subtract"`
				Past      bool   `json:"past"`
			} `json:"created"`
			Gateway struct {
				ID    int    `json:"id"`
				Title string `json:"title"`
				Type  string `json:"type"`
				Owner string `json:"owner"`
			} `json:"gateway"`
		} `json:"payments"`
		User *struct {
			ID           int    `json:"id"`
			Username     string `json:"username"`
			Name         any    `json:"name"`
			Nickname     any    `json:"nickname"`
			NationalCode any    `json:"national_code"`
			Avatar       any    `json:"avatar"`
		} `json:"user"`
		Label         any `json:"label"`
		ShippingClass struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
			Type  string `json:"type"`
		} `json:"shipping_class"`
		Created struct {
			Year      string `json:"year"`
			Month     string `json:"month"`
			MonthName string `json:"month_name"`
			Day       string `json:"day"`
			Date      string `json:"date"`
			Time      string `json:"time"`
			Universal string `json:"universal"`
			Timestamp int    `json:"timestamp"`
			Subtract  string `json:"subtract"`
			Past      bool   `json:"past"`
		} `json:"created"`
		DueDate struct {
			Year      string `json:"year"`
			Month     string `json:"month"`
			MonthName string `json:"month_name"`
			Day       string `json:"day"`
			Date      string `json:"date"`
			Time      string `json:"time"`
			Universal string `json:"universal"`
			Timestamp int    `json:"timestamp"`
			Subtract  string `json:"subtract"`
			Past      bool   `json:"past"`
		} `json:"due_date"`
		Delivery any `json:"delivery"`
		Updated  any `json:"updated"`
	} `json:"order"`
}

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
type PortalCSV struct {
	VariantID    string
	Name         string
	Sku          string
	Price        string
	ComparePrice string
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

type ItemDetail struct {
	Quantity  []uint8
	Title     string
	Code      string
	ItemID    int64
	VariantID int64
	Fee       []uint8
	InvoiceId []uint8
}
type DashtOrPortal struct {
	Id     int
	Dasht  *ItemDetail
	Portal *PortalCSV
}

type CSVPath struct {
	Success     bool   `json:"success"`
	Description string `json:"description"`
	Path        string `json:"path"`
}

type DashtOrder struct {
	CustomerRef    int    `json:"CustomerRef"`
	StockRef       int    `json:"StockRef"`
	SaleType       int    `json:"SaleType"`
	SellerPartyRef int    `json:"SellerPartyRef"`
	Description    string `json:"Description"`
	GUID           string `json:"GUID"`
	AddressRef     int    `json:"AddressRef"`
	Items          []Item `json:"Items"`
}
type Item struct {
	RowNumber   int     `json:"RowNumber"`
	ItemRef     int     `json:"ItemRef"`
	UnitRef     int     `json:"UnitRef"`
	Quantity    float64 `json:"Quantity"`
	Fee         float64 `json:"Fee"`
	Price       float64 `json:"Price"`
	Discount    float64 `json:"Discount"`
	Tax         float64 `json:"Tax"`
	Duty        float64 `json:"Duty"`
	Description string  `json:"Description"`
	StockRef    int     `json:"StockRef"`
}

var PORTAL_BASE_URL string = "https://modernhyperindustry.com/"

func ToInt(b []byte) int {
	if len(b) == 0 {
		return 0
	}
	// Float Maybe SQL SERVER HAS . in it
	num, err := strconv.ParseFloat(string(b), 10)
	if err != nil {
		log.Fatalf("Error while parsing the stock/price %s", err.Error())
	}
	if num < 0. {
		return 0
	}
	return int(num)
}

func (item *ItemDetail) ToString() string {
	return fmt.Sprintf("Title:%s,Code:%s,Quantity:%d,ItemID:%d,VariantID:%d,Fee:%d,",
		item.Title,
		item.Code,
		int64(ToInt(item.Quantity)),
		item.ItemID,
		item.VariantID,
		int64(ToInt(item.Fee)),
	)
}
func (item *DashtOrPortal) ToString() string {
	if item.Portal == nil {
		return fmt.Sprintf("both dasht and portal are nil %+v  id: %d\n", item, item.Id)
	}
	if item.Dasht == nil {
		return fmt.Sprintf("dasht is nil %+v id: %d\n", item.Portal, item.Id)
	}

	return fmt.Sprintf("dasht %s %d \n", item.Dasht.ToString(), item.Id)

}
func (portalItem *PortalCSV) ToString() string {
	return fmt.Sprintf("%+v", *portalItem)
}

func (item *ItemDetail) ToStringList() []string {
	var names []string

	price := ToInt(item.Fee) / 10
	pseudoOff := rand.Float64()*16 + 5 // random number [5,20]
	comparePrice := float64(price) * (pseudoOff/100. + 1.)
	comparePrice = math.Round(comparePrice/1000) * 1000
	quanity := int64(ToInt(item.Quantity))
	names = append(names, item.Title)
	names = append(names, item.Code)
	names = append(names, strconv.FormatInt(quanity, 10))
	names = append(names, strconv.FormatInt(item.ItemID, 10))
	names = append(names, strconv.FormatInt(item.VariantID, 10))
	names = append(names, strconv.FormatInt(int64(price), 10))
	names = append(names, strconv.FormatInt(int64(comparePrice), 10))
	if quanity == 0 {
		names = append(names, "0")
	} else {
		names = append(names, "1")
	}
	names = append(names, strconv.FormatInt(max(quanity, quanity-5), 10))
	return names
}

func (o *OrderDetail) ToDasht() DashtOrder {
	order := o.Order

	d := DashtOrder{
		CustomerRef:    658591962,
		StockRef:       1,    // header-level stock usually resolved later
		SaleType:       1,    // online sale
		SellerPartyRef: 2178, // your seller ID
		Description:    "خرید  از پرتال",
		GUID:           uuid.New().String(),
	}

	items := make([]Item, 0, len(order.Items))

	for i, it := range order.Items {
		itemRef := 0
		if it.Product != nil {
			itemRef = it.Product.ID
		}

		price := float64(it.Price)
		qty := float64(it.Quantity)

		items = append(items, Item{
			RowNumber:   i + 1,
			ItemRef:     itemRef,
			UnitRef:     1,
			Quantity:    qty,
			Fee:         price,
			Price:       price * qty,
			Discount:    anyToFloat64(it.Discount),
			Tax:         float64(it.Tax),
			Duty:        0,
			Description: it.Title,
			StockRef:    d.StockRef,
		})
	}

	d.Items = items
	return d
}

type InvalidObject struct {
	Obj     *PortalCSV
	Message string
}

func (e *InvalidObject) Error() string {
	return e.Message
}

func anyToFloat64(v any) float64 {
	switch t := v.(type) {
	case int:
		return float64(t)
	case float64:
		return t
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	default:
		return 0
	}
}

func (item *PortalCSV) ToStringList(err error) []string {
	var names []string
	names = append(names, item.VariantID)
	names = append(names, item.Name)
	names = append(names, item.Sku)
	names = append(names, item.Price)
	names = append(names, item.ComparePrice)
	names = append(names, err.Error())
	return names

}
