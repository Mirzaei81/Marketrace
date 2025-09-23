package givsoft

import (
	"encoding/binary"
	e "giv/error"
	sync_db "giv/sync_db"
	givTypes "giv/types"
	SyncPortal "giv/update"
	"log"
	"strings"
	"sync"
)

var base_url = "http://91.92.214.97:8201"
var M map[string]string // map for updating env variables

type ItemDetail struct {
	ItemDetailID int     `db:"ItemID"`
	OrderID      int     `db:"ID"`
	RowID        int     `db:"RowID"` //idx
	ItemID       int64   `db:"ItemID"`
	ItemBarcode  string  `db:"ItemBarcode"`
	Quantity     float32 `db:"Quantity"`
	Fee          float32 `db:"Fee"`
}

type NewGivOrder struct {
	SentNo             uint32 `db:"SentNo"`
	ItemQuantityOnHand int    `db:"ItemQuantityOnHand"`
	ItemFee            int    `db:"ItemFee"`
	ItemID             string `db:"ItemID"`
	VariantID          int    `db:"VariantID"`
}
type Submit_Order_detail struct {
	OrderID            int    `json:"OrderID" db:"ID"`
	SourceID           int    `json:"SourceID" db:"SourceID"`
	Type               string `db:"Type"`
	No                 int    `db:"No"`
	Date               string `json:"Date" db:"Date"` //JALALI format
	EffectiveDate      string `json:"EffectiveDate"`
	PersonID           string `json:"PersonID"`   // create customer beforehand
	CouponCode         string `json:"CouponCode"` // inv maybe
	TotalQuantity      int    `json:"TotalQuantity" db:"TotalQuantity"`
	TotalPrice         int    `json:"TotalPrice" db:"TotalPrice"`
	TotalDiscount      int    `json:"TotalDiscount" db:"TotalDiscount"`
	PackingCost        int    `json:"PackingCost" db:"PackingCost" `
	TransferCost       int    `json:"TransferCost" db:"TransferCost"`
	PostRefCode        string `json:"PostRefCode" db:"PostRefCode"`
	ReceiverName       string `json:"ReceiverName" db:"ReceiverName"`
	ReceiverProvinceID int    `json:"ReceiverProvinceID" db:"ReceiverProvinceID"`
	ReceiverCity       string `db:"ReceiverCity"`
	ReceiverAddress    string `db:"ReceiverAddress"`
	ReceiverTel        string `db:"ReceiverTel"`
	ReceiverMobile     string `db:"ReceiverMobile"`
	ReceiverPostalCode string `db:"ReceiverPostalCode"`
	PaymentBank        string `db:"PaymentBank"`
	PaymentType        string `json:"PaymentType" db:"PaymentType"`
	PaymentStatus      string `json:"PaymentStatus" db:"PaymentStatus"`
	PaymentBankRefCode string `json:"PaymentBankRefCode" db:"PaymentBankRefCode"`
	DateCreated        string `json:"DateCreated" db:"DateCreated"`
	ItemDetail         []ItemDetail
}

type QuantityOnhand struct {
	StockID            int     `db:"StockID"`
	ItemID             int     `db:"ItemID"`
	ItemQuantityOnHand float64 `db:"ItemQuantityOnHand"`
}

func getOrCreateUserID(id string) int {
	var personId int
	sync_db.SQL_DB.QueryRow(`
		IF EXISTS(select 1 from Person where PersonID = $1)
		BEGIN 
		select -1 
		END
		ELSE
		BEGIN
		select MAX(PersonID)+1 from Person 
		END
	`, id).Scan(&personId)

	return personId

}

func createOrderID() int {
	var order_id int
	sync_db.SQL_DB.QueryRow("select MAX(ID)+1 from OrderImported").Scan(&order_id)
	return order_id
}
func Create_customer(person *givTypes.Person) int {
	sync_db.SQL_DB.MustExec(`
		INSERT INTO Person(PersonTypeID,PersonCategoryCode,PersonFirstName,PersonLastName,
		PersonID,PersonMobile,PersonIsActive)
		) VALUES(0,1514,:PersonFirstName,:PersonLastName,:PersonID,:PersonMobile,1)
	`, &person)
	sync_db.SQL_DB.MustExec(`
		INSERT INTO PersonAddress(
		PersonID,PersonTypeName,PersonIsReal,PersonFirstName,
		PersonLastName,ContactName,ProvinceName,PostalCode,PersonMobile)
		VALUES(:PersonID,'خانم/آقا',1,:PersonFirstName,:PersonLastName,'','',:PostalCode,:PersonMobile)
	`)

	return person.PersonId
}

func GIVOrderHeader(order_datail *Submit_Order_detail, wg *sync.WaitGroup) {
	defer wg.Done()
	personID := getOrCreateUserID(order_datail.PersonID)
	if personID == -1 {
		parts := strings.Split(order_datail.ReceiverName, " ")
		firstname := parts[0]
		lastaname := ""
		if len(parts) > 1 {
			lastaname = parts[1]
		}
		Create_customer(&givTypes.Person{
			PersonId:   personID,
			FirstName:  firstname,
			LastName:   lastaname,
			Address:    order_datail.ReceiverAddress,
			Mobile:     order_datail.ReceiverMobile,
			PostalCode: order_datail.ReceiverPostalCode,
		})
	}

	order_datail.OrderID = createOrderID()
	sync_db.SQL_DB.MustExec(`
		INSERT INTO Orders (
			ID,
			ClientID,
			SourceID,
			SourceType,
			Type,
			Status,
			No,
			Date,
			EffectiveDate,
			PersonID,
			Description,
			CouponCode,
			TotalQuantity,
			TotalPrice,
			TotalDiscount,
			PackingCost,
			TransferCost,
			PaymentType,
			PaymentStatus,
			PaymentBank,
			PaymentBankRefCode,
			PostRefCode,
			ReceiverName,
			ReceiverProvinceID,
			ReceiverCity,
			ReceiverAddress,
			ReceiverTel,
			ReceiverMobile,
			ReceiverPostalCode,
			DateCreated,
			DateChanged,
			DataReceiptDate,
			CreditUsed,
			MoneyReceiptID,
			MoneyReceiptClientID,
			UserChanged,
			WebsiteID,
			ShippingMethod,
			ShippingTitle
		)
		VALUES (
			1,                        -- ID
			99,                      -- ClientID
			:SourceID,                     -- SourceID
			'WEB',                -- SourceType
			'SALE',                 -- Type
			'NEW',                -- Status
			:SourceID,      -- No (order number)
			:Date,                -- Date
			:DATE,                     -- EffectiveDate
			:PersonID,                     -- PersonID
			'',      -- Description
			':CouponCode',              -- CouponCode
			:TotalQuantity,                        -- TotalQuantity
			:TotalPrice,                -- TotalPrice
			:TotalDiscount,                 -- TotalDiscount
			:PackingCost,                  -- PackingCost
			:TransferCost,                 -- TransferCost
			'Online',             -- PaymentType
			'PAYMENT_STATUS_SUCCESSFUL',        -- PaymentStatus
			':PaymentBank',             -- PaymentBank
			':PaymentBankRefCode',              -- PaymentBankRefCode
			':PostRefCode',             -- PostRefCode
			N':ReceiverName',            -- ReceiverName
			-1,                       -- ReceiverProvinceID
			N':ReceiverCity',                -- ReceiverCity
			N':ReceiverAddress',-- ReceiverAddress
			':ReceiverTel',            -- ReceiverTel
			':ReceiverMobile',            -- ReceiverMobile
			':ReceiverPostalCode',             -- ReceiverPostalCode
			:DateCreated,                -- DateCreated
			NULL,                     -- DateChanged
			NULL,                     -- DataReceiptDate
			NULL,                     -- CreditUsed
			NULL,                     -- MoneyReceiptID
			NULL,                     -- MoneyReceiptClientID
			NULL,                     -- UserChanged
			1,                        -- WebsiteID
			NULL,                   -- ShippingMethod
			NULL      -- ShippingTitle
		);`, order_datail)
	log.Printf("%#v\n", order_datail)
	for idx, item := range order_datail.ItemDetail {
		ItemDetail := ItemDetail{
			ItemDetailID: item.OrderID + idx,
			OrderID:      order_datail.OrderID,
			RowID:        idx,
			ItemID:       item.ItemID,
			ItemBarcode:  item.ItemBarcode,
			Quantity:     item.Quantity,
			Fee:          item.Fee,
		}
		wg.Add(1)
		go orderRow(ItemDetail, idx, wg)
	}
}

func orderRow(itemDetail ItemDetail, idx int, wg *sync.WaitGroup) {
	defer wg.Done()
	sync_db.SQL_DB.MustExec(`
		INSERT INTO OrderItems (
		ID,
		ClientID,
		RowID,
		ItemID,
		ItemBarcode,
		Quantity,
		Fee,
		Price,
		RowDiscount,
		TotalPrice,
		TotalDiscount,
		VatBaseValue,
		VatValue,
		DateCreated,
		DateChanged,
		DataReceiptDate,
		SuggestedStockID,
		SuggestedStockQOH,
		ConfirmedStockID,
		ConfirmedStockQOH,
		ItemToSendID,
		ItemToSendClientID,
		ItemToSendTypeID,
		ItemReceiptID,
		ItemReceiptClientID,
		ItemSentID,
		ItemSentClientID,
		Status,
		UserChanged,
		ItemPackID
		)
		VALUES (
		$1,                  -- ID (line item id)
		99,                -- ClientID
		$2,                  -- RowID (order row number)
		5001,               -- ItemID
		'$3',    -- ItemBarcode
		$4,                  -- Quantity
		$5,           -- Fee (shipping/handling fee per row, if used)
		$5,          -- Price (per item price)
		0,            -- RowDiscount
		$5,          -- TotalPrice 
		0,            -- TotalDiscount
		$5,          -- VatBaseValue (value before tax)
		0,           -- VatValue
		GETDATE(),          -- DateCreated
		GETDATE(),               -- DateChanged
		NULL,               -- DataReceiptDate
		NULL,                 -- SuggestedStockID
		NULL,                 -- SuggestedStockQOH (quantity on hand at suggested stock)
		NULL,                 -- ConfirmedStockID
		NULL,                 -- ConfirmedStockQOH
		NULL,               -- ItemToSendID
		NULL,               -- ItemToSendClientID
		NULL,               -- ItemToSendTypeID
		NULL,               -- ItemReceiptID
		NULL,               -- ItemReceiptClientID
		NULL,               -- ItemSentID
		NULL,               -- ItemSentClientID
		'NEW',			-- Status
		NULL,               -- UserChanged
		NULL                -- ItemPackID
		);`)
	log.Printf("Item details are %+v", itemDetail)
}

func ListItems() []givTypes.GivItems {
	var itemDetail []givTypes.GivItems
	sync_db.SQL_DB.Select(&itemDetail, `
		SELECT qoh.ItemQuantityOnHand,p.ItemCurrentSelPrice,i.ItemID,i.VariantID
		from [dbo].QuantityOnHand as qoh join 
		[dbo].ItemParent p on p.ItemParentID = qoh.ItemID/10000000 
		JOIN VariantsItems i on i.ItemID =qoh.ItemID
		where StockID = 4 
		`)
	return itemDetail
}
func SyncPortalVariantWithGivQOH(token, sku string, variantID int, wg *sync.WaitGroup) {
	defer wg.Done()
	item, err := GetItemDetail(sku)
	if err != nil {
		log.Printf("Error Getting item detail for sku of csv %s", err)
		return
	}

	wg.Add(1)
	SyncPortal.Update_Variants(token,
		item.VariantId, int(item.ItemQuantityOnHand), // db qoh is in form of float
		item.ItemID, item.ItemPrice, wg)

}
func SyncPortalWithGivQOH(token string, wg *sync.WaitGroup) {
	defer wg.Done()
	givItems := ListItems()

	for _, item := range givItems {
		wg.Add(1)
		SyncPortal.Update_Variants(token,
			item.VariantId, int(item.ItemQuantityOnHand), // db qoh is in form of float
			item.ItemID, item.ItemPrice, wg)

	}
}

func SyncPortalByGivOrders(token string, wg *sync.WaitGroup) {
	defer wg.Done()
	orders := getNewOrders()
	var lastToken uint32 = 0
	for _, order := range orders {
		if order.SentNo > lastToken {
			lastToken = order.SentNo
		}
		wg.Add(1)
		SyncPortal.Update_Variants(token, order.VariantID, order.ItemQuantityOnHand,
			order.ItemID, float64(order.ItemFee), wg)
	}
	buff := make([]byte, 4)
	binary.LittleEndian.AppendUint32(buff, lastToken)
	sync_db.KV_DB.Write("LASTGIVODER", buff)
}

func getNewOrders() []NewGivOrder {
	lastGivOrderBuff, _ := sync_db.KV_DB.Read("LASTGIVODER")
	var lastGivOrder uint32
	binary.LittleEndian.PutUint32(lastGivOrderBuff, lastGivOrder)
	var orders []NewGivOrder
	sync_db.SQL_DB.Select(&orders, `Select s.SentNo,qoh.ItemQuantityOnHand,r.ItemFee,v.ItemID,v.VariantID 
		from ItemSent s join ItemSentRow r on s.ItemSentID = r.ItemSentID
		LEFT join QuantityOnHand qoh on qoh.ItemID = r.ItemID and qoh.StockID = 4
		inner join VariantsItems v on v.ItemID = r.ItemSentID or v.ItemID = r.ItemBarCode 
		WHERE s.SentNo >$1
		order by DateCreated DESC
	`, lastGivOrder)

	return orders
}

func GetItemDetail(itemId string) (givTypes.GivItems, error) {
	var itemDetail *givTypes.GivItems = &givTypes.GivItems{}
	if sync_db.SQL_DB == nil {
		return *itemDetail, e.New("SQL_DB is not initlized please call INIT_SQL beforehand")

	}
	noZeroLeadingID := itemId
	if itemId[:2] == "00" {
		noZeroLeadingID = itemId[2:]
	}
	err := sync_db.SQL_DB.Get(itemDetail, `
		SELECT qoh.ItemQuantityOnHand,p.ItemCurrentSelPrice,i.ItemID,i.VariantID
		from QuantityOnHand as qoh join 
		ItemParent p on p.ItemParentID = qoh.ItemID/10000000 
		JOIN VariantsItems i on i.ItemID =qoh.ItemID
		where StockID = 4 AND (i.ItemID = $1 or i.ItemID = $2)
		`, itemId, noZeroLeadingID)
	return *itemDetail, err
}
