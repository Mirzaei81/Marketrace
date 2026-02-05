package portal

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"giv/dasht"
	givsoft "giv/givsoft"
	"giv/report"
	sync_db "giv/sync_db"
	"giv/types"
	"giv/update"
	utils "giv/utils"
	"io"
	"log"
	"math"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	Jalaali "github.com/yaa110/go-persian-calendar"
)

var base_url string = "https://modernhyperindustry.com/"

type Update_Resault struct {
	Success  bool `json:"success"`
	Total    int  `json:"total"`
	Count    int  `json:"count"`
	Variants []struct {
		ID           int      `json:"id"`
		ProductID    int      `json:"product_id"`
		Title        string   `json:"title"`
		Price        int      `json:"price"`
		ComparePrice any      `json:"compare_price"`
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
	} `json:"variants"`
}
type Search_Resault struct {
	Success bool `json:"success"`
	Result  []struct {
		ID           int    `json:"id"`
		Title        string `json:"title"`
		Type         string `json:"type"`
		Description  any    `json:"description"`
		Image        any    `json:"image"`
		Price        int    `json:"price"`
		ComparePrice int    `json:"compare_price"`
		URL          string `json:"url"`
	} `json:"result"`
}
type VariantDetailResult struct {
	Success bool    `json:"success"`
	Variant Variant `json:"variant"`
}

type VariantListResult struct {
	Success  bool `json:"success"`
	Total    int  `json:"total"`
	Count    int  `json:"count"`
	Variants []Variant`json:"variants"`
}
type Variant struct {
	ID           int      `json:"id"`
	ProductID    int      `json:"product_id"`
	Title        string   `json:"title"`
	Price        int      `json:"price"`
	ComparePrice int      `json:"compare_price"`
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
type Order struct {
	ID          int      `json:"id"`
	Description any      `json:"description"`
	Status      []string `json:"status"`
	Quantity    int      `json:"quantity"`
	Weight      int      `json:"weight"`
	Shipping    int      `json:"shipping"`
	Subtotal    int      `json:"subtotal"`
	Discount    int      `json:"discount"`
	Tax         int      `json:"tax"`
	Price       int      `json:"price"`
	Payments    int      `json:"payments"`
	IP          string   `json:"ip"`
	User        struct {
		ID           int    `json:"id"`
		Username     string `json:"username"`
		Name         any    `json:"name"`
		Nickname     any    `json:"nickname"`
		NationalCode any    `json:"national_code"`
		Avatar       any    `json:"avatar"`
	} `json:"user"`
	Contact struct {
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
}
type Orders struct {
	Success bool     `json:"success"`
	Total   int      `json:"total"`
	Count   int      `json:"count"`
	Orders  []*Order `json:"orders"`
}
type session_resault struct {
	Success     bool   `json:"success"`
	Description string `json:"description"`
	Token       string `json:"token"`
}

type Product struct {
	ID                int      `json:"id"`
	Version           string   `json:"version"`
	Title             string   `json:"title"`
	Caption           any      `json:"caption"`
	Description       any      `json:"description"`
	Image             any      `json:"image"`
	Slug              string   `json:"slug"`
	URL               string   `json:"url"`
	Rate              any      `json:"rate"`
	RateCount         any      `json:"rate_count"`
	Password          any      `json:"password"`
	Layout            any      `json:"layout"`
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
	Images            any      `json:"images"`
	Category          any      `json:"category"`
	Categories        any      `json:"categories"`
	Filters           any      `json:"filters"`
	Attributes        any      `json:"attributes"`
	Variants          []struct {
		ID           int      `json:"id"`
		ProductID    int      `json:"product_id"`
		Title        string   `json:"title"`
		Price        int      `json:"price"`
		ComparePrice int      `json:"compare_price"`
		Tax          any      `json:"tax"`
		Shipping     any      `json:"shipping"`
		Weight       any      `json:"weight"`
		Length       any      `json:"length"`
		Width        any      `json:"width"`
		Height       any      `json:"height"`
		Stock        int      `json:"stock"`
		Minimum      int      `json:"minimum"`
		Maximum      int      `json:"maximum"`
		Sku          string   `json:"sku"`
		Image        any      `json:"image"`
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
		Name     any    `json:"name"`
		Nickname any    `json:"nickname"`
		Avatar   any    `json:"avatar"`
	} `json:"creator"`
}
type Product_resault struct {
	Success bool    `json:"success"`
	Product Product `json:"product"`
}

type ProductVariant struct {
	Status       []string `json:"status"`
	Price        int      `json:"price"`
	ComparePrice int      `json:"compare_price"`
	Stock        int      `json:"stock"`
	Sku          string   `json:"sku"`
	Minimum      any      `json:"minimum"`
	Maximum      any      `json:"maximum"`
	Weight       any      `json:"weight"`
	Width        any      `json:"width"`
	Length       any      `json:"length"`
	Height       any      `json:"height"`
	Title        string   `json:"title"`
	Type         string   `json:"type"`
}
type MakeProductBody struct {
	Title       string `json:"title"`
	Caption     string `json:"caption"`
	Description string `json:"description"`
	Contents    []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"contents"`
	Image             string   `json:"image"`
	Images            []string `json:"images"`
	CommentingEnabled bool     `json:"commenting_enabled"`
	Fields            []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"fields"`
	Variants        []ProductVariant `json:"variants"`
	Slug            string           `json:"slug"`
	Published       any              `json:"published"`
	Expiration      any              `json:"expiration"`
	Password        any              `json:"password"`
	MetaTitle       any              `json:"meta_title"`
	MetaDescription any              `json:"meta_description"`
	MetaRobots      any              `json:"meta_robots"`
	Redirect        any              `json:"redirect"`
	Filters         []int            `json:"filters"`
	Categories      []int            `json:"categories"`
	Status          []string         `json:"status"`
}

func Make_session() string {
	url := base_url + "site/api/v1/user/create-session"
	method := "POST"
	user, exists := os.LookupEnv("PORTAL_USER")
	if !exists {
		log.Fatal("user not exits in .env")
	}
	password, exists := os.LookupEnv("PORTAL_PASS")
	if !exists {
		log.Fatal("pass not exits in .env")
	}
	payload_string := fmt.Sprintf(`{
		"username": "%s",
		"password": "%s"
	}`, user, password)
	payload := strings.NewReader(payload_string)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	req.Close = true
	if err != nil {
		log.Fatalf("Error while creating request for Make_session %s", err.Error())
	}
	req.Header.Add("Content-Type", "Application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("error: %s payload: %s", err, payload_string)
	}
	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		log.Fatalf("Too many request error staus is %d body : %s with payloadString %s",
			res.StatusCode,
			string(body),
			payload_string)
	}
	defer res.Body.Close()
	dec := json.NewDecoder(res.Body)
	var Token_res session_resault
	err = dec.Decode(&Token_res)
	if err != nil {
		body, _ := io.ReadAll(res.Body)
		log.Println(res.StatusCode, "|-|", string(body))
		log.Printf("there was an error while decoding %+v\n", err)
		os.Exit(1)
	}
	if res.StatusCode != 200 {
		log.Println(url, payload_string)
		log.Printf("Token Not Found %s", res.Status)
		os.Exit(-1)
	} else {
		return Token_res.Token
	}
	return ""
}
func getOrders(token string) (Orders, error) {
	todayJ := Jalaali.Now().AddDate(0, 0, 0).Format("yyy/MM/dd")
	status := []string{"paid", "cash_on_delivery"}
	var orders Orders
	for _, s := range status {
		url := fmt.Sprintf("%s/site/api/v1/manage/store/orders?page=1&size=20&status=%s&payment=&start=%s&end=&label_id=&user_id=&shipping_id=&ip=&keywords=", base_url, s, todayJ)
		method := "GET"
		log.Print(url)
		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Add("Content-Type", "Application/json")
		req.Close = true
		if err != nil {
			log.Println(err)
			return Orders{}, err
		}
		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
			log.Printf("somthing is wrong with the Result %s\n", err)
			return Orders{}, err
		}
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		var current_res Orders
		err = decoder.Decode(&current_res)

		orders.Orders = append(orders.Orders, current_res.Orders...)
		orders.Count += current_res.Count

		if err != nil {
			log.Printf("There was an error while decoding orders %s\n", err)
			return Orders{}, err
		}
	}
	return orders, nil
}
func SyncGivByPortalOrders(token string, wg *sync.WaitGroup) {
	lastPortalPurchase, _ := sync_db.KV_DB.Read("LAST_PORTAL_PURCHASE")
	lastPortalPurchaseValue := binary.LittleEndian.Uint32(lastPortalPurchase)
	orders, err := getOrders(token)

	if err != nil {
		return
	}
	procWG := new(sync.WaitGroup)

	for _, order := range orders.Orders {
		if uint32(order.ID) == lastPortalPurchaseValue {
			break
		}
		procWG.Add(1)
		go getAndUpdateOrderDetail(token, order.ID, procWG)
	}
	procWG.Wait()
	if orders.Count > 0 {
		buf := make([]byte, 4) // adjust size according to your int type
		binary.LittleEndian.PutUint32(buf, uint32(orders.Orders[0].ID))
		sync_db.KV_DB.Write("LAST_PORTAL_PURCHASE", buf)
	}
}
func SyncDashtByPortalOrders(accessToken, portalToken string, wg *sync.WaitGroup) {
	lastPortalPurchase, _ := sync_db.KV_DB.Read(types.LAST_PORTAL_PURCHASE)
	lastPortalPurchaseValue := binary.LittleEndian.Uint32(lastPortalPurchase)
	orders, err := getOrders(portalToken)

	if err != nil {
		return
	}
	ch := make(chan *types.OrderDetail)

	for _, order := range orders.Orders {
		if uint32(order.ID) == lastPortalPurchaseValue {
			break
		}
		go getOrderDetail(portalToken, order.ID, ch)
	}
	select {
	case orderDetail := <-ch:
		{
			if orderDetail == nil {
				log.Printf("No OrderDetail found in \n")
			}
			dasht.SubmitOrder(accessToken, orderDetail.ToDasht(), false)
		}
	default:
		{
			log.Printf("No OrderDetail found in \n")
		}

	}
	if orders.Count > 0 {
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(orders.Orders[0].ID))
		sync_db.KV_DB.Write("LAST_PORTAL_PURCHASE", buf)
	}
}

func getOrderDetail(token string, order_id int, ch chan *types.OrderDetail) {
	url := fmt.Sprintf("%s/site/api/v1/manage/store/orders/%d", base_url, order_id)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
		ch <- nil
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "Application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		ch <- nil
	}
	defer res.Body.Close()
	decder := json.NewDecoder(res.Body)
	var order_resault types.OrderDetail
	err = decder.Decode(&order_resault)
	if err != nil {
		body, _ := json.Marshal(order_resault)
		log.Println(body)
		log.Println(err)
		ch <- &order_resault
	}
	ch <- &order_resault
}
func getAndUpdateOrderDetail(token string, order_id int, wg *sync.WaitGroup) {
	ch := make(chan *types.OrderDetail)
	go getOrderDetail(token, order_id, ch)

	order_resault := <-ch
	if order_resault != nil && order_resault.Success {
		var order_detail givsoft.Submit_Order_detail
		var ItemDetail []givsoft.ItemDetail
		var date_created_formated string
		if order_resault.Order.Payments != nil {
			date_created, _ := time.Parse("02/01/2006 15:04:05", order_resault.Order.Payments[0].Created.Universal)
			date_created_formated = date_created.Format("2006/01/02 15:04:05")
		} else {
			date_created_formated = time.Now().Format("2006/01/02 15:04:05")
		}
		date := Jalaali.Now().Format("yyyyMMdd")
		total_price := 0
		for _, item := range order_resault.Order.Items {
			if item.Sku != nil {
				total_price += item.Price
				itemId, _ := strconv.ParseInt(*item.Sku, 10, 32)
				ItemDetail = append(ItemDetail, givsoft.ItemDetail{
					ItemDetailID: int(itemId),
					ItemID:       itemId,
					OrderID:      order_id,
					ItemBarcode:  *item.Sku,
					Quantity:     float32(item.Quantity),
					Fee:          float32(item.Price),
				})
			}
		}
		order_detail = givsoft.Submit_Order_detail{
			OrderID:            -1,
			SourceID:           order_resault.Order.ID,
			Type:               "SALE",
			No:                 order_resault.Order.ID,
			Date:               date,
			EffectiveDate:      date,
			CouponCode:         "",
			TotalQuantity:      order_resault.Order.Quantity,
			TotalPrice:         total_price,
			TotalDiscount:      0,
			PackingCost:        0,
			TransferCost:       0,
			PostRefCode:        strconv.FormatInt(int64(order_resault.Order.ID), 10),
			ReceiverName:       order_resault.Order.Contact.Name,
			ReceiverCity:       order_resault.Order.Contact.City.Name,
			ReceiverAddress:    order_resault.Order.Contact.Address,
			ReceiverMobile:     order_resault.Order.Contact.Mobile,
			ReceiverPostalCode: order_resault.Order.Contact.Zipcode,
			PaymentType:        "",
			PaymentStatus:      "",
			DateCreated:        date_created_formated,
		}
		if len(order_resault.Order.Payments) > 0 {
			order_detail.PaymentBank = order_resault.Order.Payments[0].Gateway.Title
			order_detail.PaymentType = order_resault.Order.Payments[0].Gateway.Type
			order_detail.PaymentStatus = order_resault.Order.Payments[0].Status[0]
			order_detail.PaymentBankRefCode = order_resault.Order.Payments[0].ReferenceID
		} else {
			order_detail.PaymentType = "ONLINE"
			order_detail.PaymentStatus = "PAYMENT_STATUS_SUCCESSFUL"
		}
		order_detail.ItemDetail = ItemDetail
		wg.Add(1)
		go givsoft.GIVOrderHeader(&order_detail, wg)
	} else {
		log.Println(order_resault.Success)
	}
}

// Call For all variants
func SyncVariants(token string) {
	ch := make(chan *os.File)
	defer close(ch)
	go utils.GetVariants(token, ch)

	reader := <-ch
	go GetAndUpdateItemFromCsv(token, reader)
}

func UpdateVariantsField(token string,csvFile *os.File,fieldName string)  {
	lines,_ :=lineCounter(csvFile)
	csvFile.Seek(0, io.SeekStart)
	portalCH := make(chan *types.PortalCSV, lines)
	go sync_db.GetItemFromCsv(csvFile, portalCH,true)
	for portalItem := range portalCH {
		go update.Update_Variants(token,portalItem.VariantID,-1,"",0,fieldName,nil)
	}
}
func GetVariant(token string, variantID int64) Variant {
	vIdS := strconv.FormatInt(int64(variantID), 10)
	url := base_url + "/site/api/v1/manage/store/products/variants/" + vIdS
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	variantDetail := new(VariantDetailResult)
	decoder.Decode(variantDetail)
	return variantDetail.Variant
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
func GetAndUpdateItemFromCsv(token string, csvFile *os.File) {
	lines, err := lineCounter(csvFile)
	csvFile.Seek(0, io.SeekStart)

	f, err := os.Create("imcomplete.csv")
	compF, err := os.Create("finalized.csv")
	if err != nil {
		log.Fatalf("Error while opening out.csv for writing %s", err.Error())
	}
	//utf 8 bom
	_, err = f.Write([]byte{0xEF, 0xBB, 0xBF})
	if err != nil {
		log.Fatalf("Error  writing bom %s", err.Error())
	}
	_, err = compF.Write([]byte{0xEF, 0xBB, 0xBF})
	if err != nil {
		log.Fatalf("Error writing bom %s", err.Error())
	}


	defer f.Close()
	inCompleteWriter := csv.NewWriter(f)
	completeWrite := csv.NewWriter(compF)
	defer inCompleteWriter.Flush()
	defer completeWrite.Flush()
	err = inCompleteWriter.Write([]string{"ایدی", "نام", "SKU", "قیمت", "قیمت خط خورده", "خطا"})
	if err != nil {
		log.Fatalf("Error while writing headers %s", err.Error())
	}
	err = completeWrite.Write([]string{"نام", "کد", "موجودی", "ایدی دشت", "ایدی پرتال", "قیمت", "قیمت خط خورده", "کمترین", "بیشترین"})
	if err != nil {
		log.Fatalf("Error while writing headers %s", err.Error())
	}
	portalCH := make(chan *types.PortalCSV, lines)
	itemCH := make(chan *types.DashtOrPortal, lines)

	go sync_db.GetItemFromCsv(csvFile, portalCH,true)


	wg := new(sync.WaitGroup)
	dashConsumer := new(sync.WaitGroup)
	dashConsumer.Add(1)
	var zeros []*types.ItemDetail
	go func() {
		defer dashConsumer.Done()
		for dashtItem := range itemCH {
			log.Print(dashtItem.ToString())
			if dashtItem.Dasht != nil {
				time.Sleep(time.Millisecond * 333)
				if types.ToInt(dashtItem.Dasht.Quantity)==0{
					zeros =  append(zeros,dashtItem.Dasht)
				}
				completeWrite.Write(dashtItem.Dasht.ToStringList())
				completeWrite.Flush()
				dashConsumer.Add(1)
				go update.UpdatePortalVariantSKU(token, dashtItem.Dasht, dashConsumer)
			}
		}
	}()
	count := 0
	
	lastFiscal,err := dasht.GetLatestFiascal()
	if err!=nil{
		log.Fatal(err)
	}
	for portalItem := range portalCH {
		wg.Add(1)
		go func(item *types.PortalCSV) {
			defer wg.Done()
			var result *types.DashtOrPortal
			if len(item.Sku) == 0 {
				result, err = dasht.GetItemDetailByPortalExactName(item, lastFiscal)
			} else {
				result, err = dasht.GetItemDetailByCode(item,lastFiscal)
			}
			if err != nil {
				count++
				log.Printf("Error getting veriant detail %s", err)
				inCompleteWriter.Write(item.ToStringList(err))
				inCompleteWriter.Flush()
				return
			}
			itemCH <- result
		}(portalItem)
	}
	wg.Wait()
	close(itemCH)
	log.Printf("finished setting sku for all the products missed count: %d", count)
	caption := fmt.Sprintf("گزارش تاریخ %s", Jalaali.Now().Format("yyy/MM/dd,HH:mm:ss"))
	report.ReportFile(caption, f.Name())

	for _,z:=range zeros{
		dashConsumer.Add(1)
		go update.UpdatePortalVariantSKU(token, z, dashConsumer)
	}
	dashConsumer.Wait()
}
func CreateProduct(item types.ItemDetail, token string) {
	if len(item.Code) == 0 {
		log.Fatalf("Invalid item %+v", item)
	}

	url := fmt.Sprintf("%s/site/api/v1/manage/store/products", types.PORTAL_BASE_URL)
	method := "POST"

	variant := ProductVariant{
		Status:       []string{"approved", "online_payment", "bank_payment", "cash_on_delivery", "shipping_required"},
		Price:        0,
		ComparePrice: 0,
		Stock:        0,
		Sku:          item.Code,
		Minimum:      nil,
		Maximum:      nil,
		Weight:       nil,
		Width:        nil,
		Length:       nil,
		Height:       nil,
		Title:        "primary",
		Type:         "commodity",
	}
	tomanPrice := types.ToInt(item.Fee) / 10
	// Handling prices
	variant.Price = tomanPrice
	pseudoOff := rand.Float64()*16 + 5 // random number [5,20]
	if variant.Price != variant.ComparePrice {
		variant.ComparePrice = int(math.Round(float64(tomanPrice)*(pseudoOff/100+1)/1000) * 1000)
	}

	//handle quantity and shippings
	variant.Stock = types.ToInt(item.Quantity)
	if variant.Stock == 0 {
		variant.Minimum = 0
	} else {
		variant.Minimum = 1
	}
	if variant.Stock <= 5 {
		variant.Maximum = variant.Stock
	} else {
		variant.Maximum = max(5, variant.Stock-5)
	}
	var Product = MakeProductBody{
		Title:             item.Title,
		CommentingEnabled: false,
		Variants:          []ProductVariant{variant},
		Status:            []string{"pending", "available"},
	}
	data, err := json.Marshal(Product)
	payload := bytes.NewReader(data)
	if types.Debug {
		fmt.Printf("[INFO] posting variant Item %s", string(data))
	}
	if err != nil {
		log.Printf("error marshling %s\n", err.Error())
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("[INFO]: Product_submit: %s\n", string(body))

}
func SearchByName(name string, token string) (*Product, error) {
	url := fmt.Sprintf("%s/site/api/v1/search?q=%s", base_url, name)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var searchResault Search_Resault
	err = json.Unmarshal(body, &searchResault)
	if err != nil {
		log.Printf("Error while marshing the Search Result %s", err.Error())
		return nil, err
	}
	if len(searchResault.Result) == 1 {
		prod, err := getProduct(token, searchResault.Result[0].ID)
		if err != nil {
			log.Printf("Error while Getting Item: %s", err.Error())
			return nil, err
		}
		return prod, nil
	}
	return nil, err
}

func getProduct(token string, id int) (*Product, error) {
	url := fmt.Sprintf("%s/site/api/v1/manage/store/products/%d", base_url, id)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

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
	var product_result Product_resault
	err = json.Unmarshal(body, &product_result)
	if err != nil {
		log.Printf("Error while marshing the Product Result %s", err.Error())
		return nil, err
	}
	return &product_result.Product, nil
}
