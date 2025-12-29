package dasht

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
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var base_url = "http://127.0.0.1:8080/api/"
var M map[string]string // map for updating env variables
var Debug = false

type LoginDetail struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
}
type AccDetail struct {
	GUID         string
	RefreshToken string
}

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

func Login(GUID string) (*LoginDetail, error) {
	url := "http://127.0.0.1:8080/api/Users/Login"
	method := "POST"

	data := fmt.Sprintf(`{
		"Guid": "%s"
		}`, GUID)
	payload := strings.NewReader(data)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Add("GenerationVersion", "102")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

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
	var detail LoginDetail
	err = json.Unmarshal(body, &detail)
	sync_db.KV_DB.Write(types.DASHT_ACCESS_TOKEN, []byte(detail.AccessToken))
	return &detail, err
}
func GetAcc() (AccDetail, error) {
	var loginDetail AccDetail
	err := sync_db.SQL_DB.QueryRow("SELECT GUID,RefreshToken from [FMK].[UserToken]").Scan(&loginDetail.GUID, &loginDetail.RefreshToken)
	return loginDetail, err
}
func RefreshToken() (*LoginDetail, error) {
	detail, err := GetAcc()
	if err != nil {
		log.Fatalf("Error while querying for refresh token %s", err)
	}
	url := "http://127.0.0.1:8080/api/Users/Refresh"
	method := "POST"

	data := fmt.Sprintf(`{ "RefreshToken": "%s"}`, detail.RefreshToken)
	payload := strings.NewReader(data)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Add("GenerationVersion", "102")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

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
	var login LoginDetail
	err = json.Unmarshal(body, &login)
	if err != nil && Debug {
		log.Printf("Error marshling refrest body %s", body)
	}
	sync_db.KV_DB.Write(types.DASHT_ACCESS_TOKEN, []byte(login.AccessToken))
	return &login, err
}

func GetSaleInvoise(accToken string) (string, error) {
	now := time.Now().AddDate(0, 0, -1).UTC()
	today := now.Format("2006-01-02T15:04:05.000Z")
	url := fmt.Sprintf("http://127.0.0.1:8080/api/SaleInvoices?fromDate=%s&limit=100&offset=0", today)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("GenerationVersion", "102")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accToken))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(body), nil
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
	if Debug {
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
func SyncPortalVariantWithDashtCode(token string, variantID string, portalVariant *types.PortalCSV, errch chan *types.PortalCSV, wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)
	item, err := GetItemDetailByCode(portalVariant, -1)
	if err != nil {
		log.Printf("Error while fetching detail for item %s Error: %s ", portalVariant.Sku, err)
		return
	}
	wg.Add(1)
	if item.Dasht != nil {
		go SyncPortal.Update_Variants(
			token,
			variantID, types.ToInt(item.Dasht.Quantity), // db qoh is in form of float
			item.Dasht.Code, types.ToInt(item.Dasht.Fee), wg,
		)
	}

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
func ListAllItemsChan(ch *chan *types.ItemDetail) {
	defer close(*ch)
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
		return
	}
	defer rows.Close()
	for rows.Next() {
		var item types.ItemDetail
		if err := rows.Scan(&item.Quantity, &item.Title, &item.Code, &item.ItemID, &item.Fee); err != nil {
			return
		}
		*ch <- &item
	}
	if err = rows.Err(); err != nil {
		return
	}
}
func ListAllItemsSync() ([]types.ItemDetail, error) {
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

// query automaticly without name matching
func GetItemDetailByCode(portal *types.PortalCSV, count int) (*types.DashtOrPortal, error) {
	var dashtOrPortal types.DashtOrPortal
	// sanity check
	if len(portal.Sku) == 0 {
		buff := make([]byte, 255)
		n := runtime.Stack(buff, false)
		message := fmt.Sprintf("Invalid Sku while checking %+v stack: %s", *portal, string(buff[:n]))
		buff = nil

		return &dashtOrPortal, &types.InvalidObject{Obj: portal, Message: message}
	}
	var item types.ItemDetail
	dashtOrPortal.Dasht = &item
	dashtOrPortal.Portal = portal
	query := fmt.Sprintf(
		`with cte as (SELECT  VI.VariantID,itemS.Quantity, i.Title,(select max(val) from (VALUES(max(IItem.Fee)),  (max(PurItem.Fee)) , (max(SaleItem.Price0))) as T(val) ) as Fee,i.Code,i.ItemID,IItem.SaleInvoiceNumber,
		ROW_NUMBER() OVER (
            PARTITION BY i.Code
            ORDER BY IItem.SaleInvoiceDate DESC )as rn  from [pos].Item as i 
left JOIN [POS].ItemStockSummary as itemS on itemS.ItemRef = i.ItemID
left join [POS].vwAllSaleInvoiceItem as IItem on IItem.ItemRef = i.ItemID
left join [POS].vwAllPurchaseInvoiceItem as PurItem on PurItem.ItemRef = i.ItemID
left join [pos].ItemSalePrice as SaleItem on SaleItem.ItemRef = i.ItemID  
left join VariantsItems as VI on VI.ItemID = i.ItemID
WHERE code  = N'%s' or BarCode = N'%s' and itemS.FiscalPeriodRef =( select TOP 1 FiscalPeriodID from [FMK].FiscalPeriod order by FiscalPeriodID desc )
	group by   itemS.Quantity,i.Title,  i.Code,i.ItemID, IItem.SaleInvoiceDate,IItem.SaleInvoiceItemID,IItem.SaleInvoiceNumber,VI.VariantID)
  SELECT  
	ISNULL(VariantID,0),
    Quantity,
    Title,
    Code,
    ItemID,
    Fee
FROM cte
 where rn=1
ORDER BY FEE,SaleInvoiceNumber DESC;`, portal.Sku, portal.Sku)

	if Debug {
		log.Printf("[INFO]: %s\n", query)
	}
	err := sync_db.SQL_DB.QueryRow(query).Scan(&item.VariantID, &item.Quantity, &item.Title, &item.Code, &item.ItemID, &item.Fee)
	if item.VariantID == 0 {
		variantID, err := strconv.ParseInt(portal.VariantID, 10, 64)
		if err != nil {
			return nil, err
		}
		item.VariantID = variantID
		go updateVariabtItemTable(portal.VariantID, item.Code)
	}
	if err != nil {
		return nil, err
	}
	return &dashtOrPortal, err

}
func GetItemDetailByPortalExactName(portal *types.PortalCSV, count int) (*types.DashtOrPortal, error) {
	var dashtOrPortal types.DashtOrPortal
	dashtOrPortal.Id = count
	query := fmt.Sprintf(`
	with cte as(SELECT VI.VariantID,itemS.Quantity, i.Title,i.Code,i.ItemID,
		(select max(val) from (VALUES(max(IItem.Fee)),  (max(PurItem.Fee)) , (max(SaleItem.Price0))) as T(val) )as Fee,
		IItem.SaleInvoiceItemID,IItem.SaleInvoiceNumber,ROW_NUMBER() OVER (
            PARTITION BY i.Code
            ORDER BY IItem.SaleInvoiceDate DESC )as rn from [Pos].Item as i
left JOIN [POS].ItemStockSummary as itemS on itemS.ItemRef = i.ItemID
left join [POS].vwAllSaleInvoiceItem as IItem on IItem.ItemRef = i.ItemID
left join [POS].vwAllPurchaseInvoiceItem as PurItem on PurItem.ItemRef = i.ItemID
left join [pos].ItemSalePrice as SaleItem on SaleItem.ItemRef = i.ItemID
left join VariantsItems as VI on VI.ItemID = CAST(i.ItemID as nvarchar(250))
  WHERE REPLACE(Title,N' ',N'') like REPLACE('%%%s%%',' ','') and
  IItem.SaleInvoiceDate>= dateadd(day,DATEDIFF(day,3,(select  max(IItem.SaleInvoiceDate) from  [POS].vwAllSaleInvoiceItem as IItem where IItem.ItemRef = i.ItemID) ),0 )
and itemS.FiscalPeriodRef =( SELECT TOP 1 FiscalPeriodID from [FMK].FiscalPeriod order by FiscalPeriodID desc )
group by itemS.Quantity,i.Title,  i.Code,i.ItemID, IItem.SaleInvoiceDate,IItem.SaleInvoiceItemID,IItem.SaleInvoiceNumber,VI.VariantID)
  SELECT  ISNULL(VariantID,0),
    Quantity,
    Title,
    Code,
    ItemID,
    Fee
	FROM cte
	where rn = 1
	ORDER BY  FEE,SaleInvoiceItemID,ItemID DESC;`, portal.Name)
	if Debug {
		log.Print(query)
	}
	var item types.ItemDetail
	err := sync_db.SQL_DB.QueryRow(query).Scan(&item.VariantID, &item.Quantity, &item.Title, &item.Code, &item.ItemID, &item.Fee)
	if err != nil {
		log.Printf("Error in Item Query item:%+v err:%s\n", *portal, err.Error())
		return nil, err
	}
	if item.ItemID == 0 {
		log.Fatalf("[FATAL] Error in Item Query item:%+v id is zero \n", *portal)
		return nil, err
	}

	dashtOrPortal.Dasht = &item
	dashtOrPortal.Portal = portal
	if len(portal.Sku) != 0 && item.VariantID != 0 && strconv.FormatInt(item.VariantID, 10) != portal.Sku {
		log.Printf("[FATAL]: Error incoming variant id doesn't match with itemId  portalCSV : %s  dasht Item: %s", portal.ToString(), item.ToString())
	}
	if item.ItemID != 0 && item.VariantID == 0 {
		item.VariantID, err = strconv.ParseInt(portal.VariantID, 10, 64)
		if err != nil {
			return nil, err
		}
		go updateVariabtItemTable(portal.VariantID, item.Code)
	}
	return &dashtOrPortal, nil
}
func GetItemDetailByPortalSimilarityName(portal *types.PortalCSV, count int, percents ...int) (*types.DashtOrPortal, error) {
	var dashtOrPortal types.DashtOrPortal
	dashtOrPortal.Id = count
	percent := 95
	if len(percents) != 0 {
		percent = percents[0]
	}
	query := fmt.Sprintf(`
	with cte as(
	select itemS.Quantity, i.Title,  i.Code,i.ItemID,max(IItem.Fee) as Fee,IItem.SaleInvoiceItemID,ROW_NUMBER() OVER (
            PARTITION BY i.Code
            ORDER BY IItem.SaleInvoiceDate DESC )as rn from [Pos].Item as i
left JOIN [POS].ItemStockSummary as itemS on itemS.ItemRef = i.ItemID
left join [POS].vwAllSaleInvoiceItem as IItem on IItem.ItemRef = i.ItemID
  WHERE Title EDIT_DISTANCE_SIMILARITY(%s)>%d and
  IItem.SaleInvoiceDate>= dateadd(day,DATEDIFF(day,3,(select  max(IItem.SaleInvoiceDate) from  [POS].vwAllSaleInvoiceItem as IItem where IItem.ItemRef = i.ItemID) ),0 )
group by  itemS.Quantity,i.Title,  i.Code,i.ItemID, IItem.SaleInvoiceDate,IItem.SaleInvoiceItemID)
  SELECT  
    Quantity,
    Title,
    Code,
    ItemID,
    Fee
FROM cte
where rn = 1
ORDER BY  SaleInvoiceItemID,ItemID DESC;`, portal.Name, percent)
	if Debug {
		log.Print(query)
	}
	var item types.ItemDetail
	err := sync_db.SQL_DB.QueryRow(query).Scan(&item.Quantity, &item.Title, &item.Code, &item.ItemID, &item.Fee)
	if err != nil {
		log.Printf("Error in Item Query item:%+v err:%s\n", *portal, err.Error())
		return nil, err
	}
	if item.ItemID == 0 {
		log.Fatalf("[FATAL] Error in Item Query item:%+v id is zero \n", *portal)
		return nil, err
	}
	item.VariantID, err = strconv.ParseInt(portal.VariantID, 10, 64)

	dashtOrPortal.Dasht = &item
	dashtOrPortal.Portal = portal
	if err != nil {
		return nil, err
	}
	if item.ItemID != 0 {
		go updateVariabtItemTable(portal.VariantID, item.Code)
	}
	return &dashtOrPortal, nil
}
func setLastOrderId(items []types.ItemDetail) {
	if len(items) != 0 {
		sync_db.KV_DB.Write(types.LAST_DASHT_PURCHASE, items[0].InvoiceId)
	}

}

func setLastCreatedId(items types.ItemDetail) {
	if items.ItemID != 0 {
		buff := make([]byte, 4)
		strconv.AppendUint(buff, uint64(items.ItemID), 10)
		err := sync_db.KV_DB.Write(types.LAST_DASHT_CREATED, buff)
		if err != nil {
			log.Printf("Error while setting last Created ItemID\n")
		}
		buff = nil
	}

}
func setLastOrderIdP(items *types.ItemDetail) {
	if items != nil {
		sync_db.KV_DB.Write(types.LAST_DASHT_PURCHASE, items.InvoiceId)
	}

}

func GetItemDetailByOrderId(ch *chan *types.ItemDetail) {
	var lastItem *types.ItemDetail
	defer close(*ch)
	defer setLastOrderIdP(lastItem)
	id, err := strconv.ParseInt(sync_db.KV_DB.ReadString(types.LAST_DASHT_PURCHASE), 10, 64)
	query, err := sync_db.SQL_DB.Query(`SELECT * 
		from [POS].vwAllSaleInvoiceItem where SaleInvoiceNumber > ? and 
	datediff(day,SaleInvoiceDate,GETDATE())<1
	order by SaleInvoiceDate DESC `, id)
	if err != nil {
		log.Printf("Error while query for last Order %s", err.Error())
		return
	}
	for query.Next() {
		var item types.ItemDetail
		if err := query.Scan(&item.Quantity, &item.Title, &item.Code, &item.ItemID, &item.Fee); err != nil {
			log.Printf("Error while Scanning Order %s", err.Error())
		}
		lastItem := &item
		*ch <- lastItem
	}

}
func SubmitOrder(accessToken string, order types.DashtOrder, retry bool) {
	url := base_url + "SaleInvoices"
	method := "POST"

	payloadB, err := json.Marshal(order)
	payload := bytes.NewReader(payloadB)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("GenerationVersion", "102")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	code := res.StatusCode
	if code != 200 {
		if retry {
			log.Fatal("[FATAL] Still not working after reloging in")
		}
		login, err := RefreshToken()
		if err != nil {
			acc, err := GetAcc()
			if err != nil {
				log.Printf("[FATAl] Error while get acc detail %s", err)
				return
			}
			login, err = Login(acc.GUID)
			if err != nil {
				log.Printf("[FATAl] Error while get acc detail %s", err)
				return
			}

		}
		SubmitOrder(login.AccessToken, order, true)
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Print(string(body))
}
func GetItnmDetailByOrderIdSync() ([]types.ItemDetail, error) {
	var items []types.ItemDetail
	defer setLastOrderId(items)

	id, err := strconv.ParseInt(sync_db.KV_DB.ReadString(types.LAST_DASHT_PURCHASE), 10, 64)
	query, err := sync_db.SQL_DB.Query(`select  * from [POS].vwAllSaleInvoiceItem where SaleInvoiceNumber >? and 
	datediff(day,SaleInvoiceDate,GETDATE())<1
	order by SaleInvoiceDate DESC `, id)
	if err != nil {
		return items, err
	}
	for query.Next() {
		var item types.ItemDetail
		if err := query.Scan(&item.Quantity, &item.Title, &item.Code, &item.ItemID, &item.Fee); err != nil {
			return items, err
		}
		items = append(items, item)
	}
	return items, err
}
func GetItemByCreationDate(date string, ch chan types.ItemDetail, update ...bool) {
	var item types.ItemDetail
	if len(update) > 0 && update[0] {
		defer setLastCreatedId(item)
	}
	defer close(ch)
	query := fmt.Sprintf(`
	with cte as(SELECT VI.VariantID,itemS.Quantity, i.Title,i.Code,i.ItemID,
		(select max(val) from (VALUES(max(IItem.Fee)),  (max(PurItem.Fee)) , (max(SaleItem.Price0))) as T(val) ) as Fee,
		IItem.SaleInvoiceItemID,IItem.SaleInvoiceNumber,ROW_NUMBER() OVER (
            PARTITION BY i.Code
            ORDER BY IItem.SaleInvoiceDate DESC )as rn from [Pos].Item as i
left JOIN [POS].ItemStockSummary as itemS on itemS.ItemRef = i.ItemID
left join [POS].vwAllSaleInvoiceItem as IItem on IItem.ItemRef = i.ItemID
left join [POS].vwAllPurchaseInvoiceItem as PurItem on PurItem.ItemRef = i.ItemID
left join [pos].ItemSalePrice as SaleItem on SaleItem.ItemRef = i.ItemID
left join VariantsItems as VI on VI.ItemID = CAST(i.ItemID as nvarchar(250))
  WHERE i.CreationDate > CONVERT(DATETIME, '%s', 111) and
  IItem.SaleInvoiceDate>= dateadd(day,DATEDIFF(day,3,(select  max(IItem.SaleInvoiceDate) from  [POS].vwAllSaleInvoiceItem as IItem where IItem.ItemRef = i.ItemID) ),0 )
and itemS.FiscalPeriodRef =( SELECT TOP 1 FiscalPeriodID from [FMK].FiscalPeriod order by FiscalPeriodID desc )
group by itemS.Quantity,i.Title,  i.Code,i.ItemID, IItem.SaleInvoiceDate,IItem.SaleInvoiceItemID,IItem.SaleInvoiceNumber,VI.VariantID)
  SELECT  ISNULL(VariantID,0),
    Quantity,
    Title,
    Code,
    ItemID,
    Fee
	FROM cte
	where rn = 1
	ORDER BY  FEE,SaleInvoiceItemID,ItemID DESC`, date)

	rows, err := sync_db.SQL_DB.Query(query)
	if types.Debug {
		log.Printf("[INFO]: %s", query)
	}
	if err != nil {
		log.Fatalf("Error While getting items %s", err.Error())
		return
	}
	for rows.Next() {
		if err := rows.Scan(&item.VariantID, &item.Quantity, &item.Title, &item.Code, &item.ItemID, &item.Fee); err != nil {
			log.Printf("Error while getting a item %s", err.Error())
			return
		}
		ch <- item
	}

}
func updateVariabtItemTable(variantId string, code string) {
	_, err := sync_db.SQL_DB.Exec(`Insert VariantsItems(VariantID,ItemID) SELECT ?,? WHERE NOT EXISTS (SELECT 1 from VariantsItems where VariantID = ? )`, variantId, code, variantId)
	if err != nil {
		log.Printf("Erorr while updateing VariantsItems(%s,%s) err: %s", variantId, code, err.Error())
	}

}
