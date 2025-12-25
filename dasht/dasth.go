package dasht

// "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjEzLCJleHBpcmVUaW1lIjoiXC9EYXRlKDE3NjYzMTg0NDM2MDApXC8ifQ.HsbhdZvxQpnSsPAVryUP_aeL5zZP5_-8F5HtXv7segY"

import (
	"bytes"
	"encoding/json"
	"fmt"
	sync_db "giv/sync_db"
	"giv/types"
	SyncPortal "giv/update"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var base_url = "http://127.0.0.1:8080/api/"
var M map[string]string // map for updating env variables
var debug = false

type CustomersBody struct {
	GUID         string    `json:"GUID"`
	PhoneNumber  string    `json:"PhoneNumber"`
	CustomerType int64     `json:"CustomerType"`
	Name         string    `json:"Name"`
	LastName     string    `json:"LastName"`
	BirthDate    string    `json:"BirthDate"`
	NationalID   string    `json:"NationalID"`
	EconomicCode string    `json:"EconomicCode"`
	Addresses    []Address `json:"Addresses"`
}
type CustomerResp struct {
	CustomerID           int64     `json:"CustomerID"`
	GUID                 string    `json:"GUID"`
	Title                string    `json:"Title"`
	Code                 string    `json:"Code"`
	DLCode               string    `json:"DLCode"`
	HasCredit            bool      `json:"HasCredit"`
	CustomerCredit       int64     `json:"CustomerCredit"`
	AcceptCustomerCheque bool      `json:"AcceptCustomerCheque"`
	PhoneNumber          string    `json:"PhoneNumber"`
	CustomerType         int64     `json:"CustomerType"`
	Name                 string    `json:"Name"`
	LastName             string    `json:"LastName"`
	CustomerRemaining    int64     `json:"CustomerRemaining"`
	BirthDate            string    `json:"BirthDate"`
	NationalID           string    `json:"NationalID"`
	EconomicCode         string    `json:"EconomicCode"`
	Version              int64     `json:"Version"`
	PartyGroupRef        int64     `json:"PartyGroupRef"`
	Addresses            []Address `json:"Addresses"`
}

type Address struct {
	Title   string `json:"Title"`
	IsMain  bool   `json:"IsMain"`
	CityRef int64  `json:"CityRef"`
	Address string `json:"Address"`
	ZipCode string `json:"ZipCode"`
	GUID    string `json:"GUID"`
}

type Orders struct {
	SaleInvoiceID  int64       `json:"SaleInvoiceID"`
	Date           string      `json:"Date"`
	Number         int64       `json:"Number"`
	CustomerRef    int64       `json:"CustomerRef"`
	StockRef       int64       `json:"StockRef"`
	SaleType       int64       `json:"SaleType"`
	TotalPrice     float64     `json:"TotalPrice"`
	TotalDiscount  float64     `json:"TotalDiscount"`
	TotalTax       float64     `json:"TotalTax"`
	TotalDuty      float64     `json:"TotalDuty"`
	TotalNetPrice  float64     `json:"TotalNetPrice"`
	RemainingPrice float64     `json:"RemainingPrice"`
	State          int64       `json:"State"`
	AddressRef     *int64      `json:"AddressRef"`
	GUID           interface{} `json:"GUID"`
	Items          interface{} `json:"Items"`
}

func GetOrCreateCustomer(customer CustomersBody, accessToken string) (*CustomerResp, error) {
	url := base_url + "api/Customers"
	method := "POST"
	client := &http.Client{}
	bodyByte, err := json.Marshal(customer)
	payload := bytes.NewReader(bodyByte)
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("GenerationVersion", "102")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if debug {
		log.Printf("creating customer %s", string(body))
	}
	var customerResp CustomerResp
	err = json.Unmarshal(body, &customerResp)
	if err != nil {
		log.Printf("Error while parsing Customer resp %s \n", err)
		return nil, err

	}
	return &customerResp, nil
}
func SyncPortalVariantWithDashtList(token string) {

}
func SyncPortalVariantWithDashtCode(token string, variantID string, itemCode int, wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)
	item, err := getItemDetailByCode(itemCode)
	if err != nil {
		log.Printf("Error while fetching detail for item %d Error: %s ", itemCode, err)
		return
	}
	SyncPortal.Update_Variants(
		token,
		variantID, int(item.Quantity[0]), // db qoh is in form of float
		item.ItemID, int(item.Fee[0]), wg,
	)

}
func SyncDashtWithPortalOrders(token string, wg *sync.WaitGroup) {
	defer wg.Done()

}

func SyncPortalByDashtOrders(token string, wg *sync.WaitGroup) {
	defer wg.Done()

	// lastGivOrderBuff, _ := sync_db.KV_DB.Read("LASTDASHTODER")
	// var lastGivOrder uint32

	// orders := getNewOrders(lastGivOrder)
	// var lastToken uint32 = 0
	// for _, order := range orders {
	// 	if order.SentNo > lastToken {
	// 		lastToken = order.SentNo
	// 	}
	// 	wg.Add(1)
	// 	SyncPortal.Update_Variants(token, order.VariantID, order.ItemQuantityOnHand,
	// 		order.ItemID, float64(order.ItemFee), wg)
	// }
	// buff := make([]byte, 4)
	// binary.LittleEndian.AppendUint32(buff, lastToken)
	// sync_db.KV_DB.Write("LASTGIVODER", buff)
}

func getNewOrders(lastGivOrder uint32) {
}

func ListAllItems() ([]types.ItemDetail, error) {
	rows, err := sync_db.SQL_DB.Query(
		`with cte as (select itemS.Quantity, i.Title,  i.Code,i.ItemID,IItem.Fee,IItem.PurchaseInvoiceItemID,
		ROW_NUMBER() OVER (
            PARTITION BY i.Code
            ORDER BY IItem.PurchaseInvoiceItemID DESC
        ) AS rn from [POS].vwPurchaseInvoiceItem as IItem
inner join [Pos].Item as i on i.ItemID = ItemRef 
inner JOIN [POS].ItemStockSummary as itemS on itemS.ItemRef = IItem.ItemRef
  WHERE itemS.FiscalPeriodRef =(SELECT TOP 1 FiscalPeriodID from [FMK].FiscalPeriod order by FiscalPeriodID desc ) )
  SELECT 
    Quantity,
    Title,
    Code,
    ItemID,
    Fee
FROM cte
 where rn=1
ORDER BY PurchaseInvoiceItemID ,ItemID DESC;`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []types.ItemDetail

	for rows.Next() {
		var item types.ItemDetail
		if err := rows.Scan(&item.Quantity, &item.Title, &item.Code, &item.ItemID, &item.Fee); err != nil {
			return items, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return items, err
	}
	return items, err
}
func getItemDetailByCode(code int) (*types.ItemDetail, error) {
	var item types.ItemDetail
	err := sync_db.SQL_DB.QueryRow(
		`with cte as (select itemS.Quantity, i.Title,  i.Code,i.ItemID,IItem.Fee,IItem.PurchaseInvoiceItemID,
		ROW_NUMBER() OVER (
            PARTITION BY i.Code
            ORDER BY IItem.PurchaseInvoiceItemID DESC
        ) AS rn from [POS].vwPurchaseInvoiceItem as IItem
inner join [Pos].Item as i on i.ItemID = ItemRef 
inner JOIN [POS].ItemStockSummary as itemS on itemS.ItemRef = IItem.ItemRef
  WHERE code  = ? and itemS.FiscalPeriodRef =( select TOP 1 FiscalPeriodID from [FMK].FiscalPeriod order by FiscalPeriodID desc ) )
  SELECT  
    Quantity,
    Title,
    Code,
    ItemID,
    Fee
FROM cte
 where rn=1
ORDER BY PurchaseInvoiceItemID ,ItemID DESC;`, code,
	).Scan(&item.Quantity, &item.Title, &item.Code, &item.ItemID, &item.Fee)
	if err != nil {
		return nil, err
	}
	return &item, err

}
func GetItemDetailByPortalExactName(portal *types.PortalCSV, ch chan *types.DashtOrPortal) {
	var dashtOrPortal types.DashtOrPortal
	var item types.ItemDetail
	err := sync_db.SQL_DB.QueryRow(
		`with cte as (select itemS.Quantity, i.Title,  i.Code,i.ItemID,IItem.Fee,IItem.PurchaseInvoiceItemID,
		ROW_NUMBER() OVER (
            PARTITION BY i.Code
            ORDER BY IItem.PurchaseInvoiceItemID DESC
        ) AS rn from [POS].vwPurchaseInvoiceItem as IItem
inner join [Pos].Item as i on i.ItemID = ItemRef 
inner JOIN [POS].ItemStockSummary as itemS on itemS.ItemRef = IItem.ItemRef
  WHERE Title = ? and itemS.FiscalPeriodRef =( select TOP 1 FiscalPeriodID from [FMK].FiscalPeriod order by FiscalPeriodID desc ) )
  SELECT  
    Quantity,
    Title,
    Code,
    ItemID,
    Fee
FROM cte
ORDER BY PurchaseInvoiceItemID ,ItemID DESC;`, portal.Name,
	).Scan(&item.Quantity, &item.Title, &item.Code, &item.ItemID, &item.Fee)
	if err != nil {
		log.Printf("Erorr in Item Query %s\n", err.Error())
		ch <- &dashtOrPortal
	}
	item.VariantID, err = strconv.ParseInt(portal.VariantID, 10, 64)

	dashtOrPortal.Dasht = &item
	dashtOrPortal.Portal = portal
	if err != nil {
		ch <- &dashtOrPortal
	}
	if item.ItemID != 0 {
		go updateVariabtItemTable(portal.VariantID, item.Code)
	}
	ch <- &dashtOrPortal
}
func updateVariabtItemTable(variantId string, itemID string) {
	_, err := sync_db.SQL_DB.Exec(`Insert VariantsItems(VariantID,ItemID) SELECT ?,? WHERE NOT EXISTS (SELECT 1 from VariantsItems where VariantID = ? )`, variantId, itemID, variantId)
	if err != nil {
		log.Printf("Erorr while updateing VariantsItems(%s,%d) err: %s", variantId, itemID, err.Error())
	}

}
