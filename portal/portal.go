package portal

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"giv/dasht"
	givsoft "giv/givsoft"
	sync_db "giv/sync_db"
	"giv/types"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/peterbourgon/diskv/v3"
	Jalaali "github.com/yaa110/go-persian-calendar"
)

var DB *diskv.Diskv

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
type VariantDetailResult struct {
	Success bool    `json:"success"`
	Variant Variant `json:"variant"`
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

type Order_result struct {
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
type csvPath struct {
	Success     bool   `json:"success"`
	Description string `json:"description"`
	Path        string `json:"path"`
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
func fetchOrders(token string, wg *sync.WaitGroup) (Orders, error) {
	defer wg.Done()
	todayJ := Jalaali.Now().AddDate(0, 0, -10).Format("yyy/MM/dd")
	status := []string{"paid", "cash_on_delivery"}
	var orders Orders
	for _, s := range status {
		wg.Add(1)
		url := base_url + fmt.Sprintf("/site/api/v1/manage/store/orders?page=1&size=20&status=%s&payment=&start=%s&end=&label_id=&user_id=&shipping_id=&ip=&keywords=", s, todayJ)
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
		wg.Done()
	}
	return orders, nil
}
func SyncGivByPortalOrders(token string, wg *sync.WaitGroup) {
	fetchWG := new(sync.WaitGroup)
	fetchWG.Add(1)
	lastPortalPurchase, _ := DB.Read("LAST_PORTAL_PURCHASE")
	lastPortalPurchaseValue := binary.LittleEndian.Uint32(lastPortalPurchase)
	orders, err := fetchOrders(token, fetchWG)
	fetchWG.Wait()

	if err != nil {
		return
	}
	procWG := new(sync.WaitGroup)

	for _, order := range orders.Orders {
		if uint32(order.ID) == lastPortalPurchaseValue {
			break
		}
		procWG.Add(1)
		go getOrderDetail(token, order.ID, procWG)
	}
	procWG.Wait()
	if orders.Count > 0 {
		buf := make([]byte, 4) // adjust size according to your int type
		binary.LittleEndian.PutUint32(buf, uint32(orders.Orders[0].ID))
		DB.Write("LAST_PORTAL_PURCHASE", buf)
	}
}
func getOrderDetail(token string, order_id int, wg *sync.WaitGroup) {
	defer wg.Done()
	url := fmt.Sprintf("%s/site/api/v1/manage/store/orders/%d", base_url, order_id)
	log.Printf("Getting order Detail %d\n", order_id)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "Application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	decder := json.NewDecoder(res.Body)
	var order_resault Order_result
	err = decder.Decode(&order_resault)
	if err != nil {
		body, _ := json.Marshal(order_resault)
		log.Println(body)
		log.Println(err)
		return
	}
	if order_resault.Success {
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
	ch := make(chan *csv.Reader)
	go GetVariants(token, &ch)

	reader := <-ch
	close(ch)
	syncDashtByCsv(token, reader)
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
func GetVariants(token string, ch *chan *csv.Reader) {
	url := base_url + "site/api/v1/manage/store/products/variants/export"
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
	csvPath := new(csvPath)
	decoder.Decode(csvPath)
	log.Printf("DEBUG: downloading csv path from %s \n", csvPath.Path)
	go GetAndParseCSV(csvPath.Path, token, ch)
}
func GetAndParseCSV(uri string, token string, ch *chan *csv.Reader) {
	url := base_url + uri
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	fmt.Println("getting new Csv File", url)

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
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	reader := csv.NewReader(bytes.NewBuffer(body))
	todayJ := Jalaali.Now().AddDate(0, 0, 0).Format("yyy-MM-dd")
	os.WriteFile(todayJ+".csv", body, 0777)
	_, err = reader.Read()
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
	reader.FieldsPerRecord = -1
	*ch <- reader
}
func syncDashtByCsv(token string, reader *csv.Reader) {
	wg := new(sync.WaitGroup)
	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				log.Println(err)
				break
			}
		}
		if line[8] != "" {
			itemId, _ := strconv.ParseInt(line[0], 10, 64)
			log.Printf("Updating variant : %s with Sku Of %s", line[2], line[8])
			// just insert into VariantsItem  and update  table accordingly
			// csv is in form ID ProductID Title Price ComparePrice Type Status Stock Sku
			dasht.SyncPortalVariantWithDashtCode(token, line[8], int(itemId), wg)
		}
	}
}

func dahstToStringList(item types.ItemDetail) []string {
	var names []string
	names = append(names, item.Title)
	names = append(names, item.Code)
	names = append(names, strconv.FormatInt(int64(intfrombytes(item.Quantity)), 10))
	names = append(names, strconv.FormatInt(item.ItemID, 10))
	names = append(names, strconv.FormatInt(item.VariantID, 10))
	names = append(names, strconv.FormatInt(int64(intfrombytes(item.Fee)), 10))
	return names
}
func poratlToString(item types.PortalCSV) []string {
	var names []string
	names = append(names, item.VariantID)
	names = append(names, item.Name)
	names = append(names, item.Sku)
	names = append(names, item.Price)
	names = append(names, item.ComparePrice)
	return names

}
func GetItemFromCsv(token string, path string) {
	f, err := os.Create("imcomplete.csv")
	compF, err := os.Create("finalized.csv")
	if err != nil {
		log.Fatalf("Error while opening out.csv for writing %s", err.Error())
	}
	_, err = f.Write([]byte{0xEF, 0xBB, 0xBF})
	if err != nil {
		log.Fatalf("Error  writing bom %s", err.Error())
	}
	defer f.Close()
	inCompleteWriter := csv.NewWriter(f)
	completeWrite := csv.NewWriter(compF)
	defer inCompleteWriter.Flush()
	err = inCompleteWriter.Write([]string{"ایدی", "نام", "SKU", "قیمت", "قیمت خط خورده"})
	if err != nil {
		log.Fatalf("Error while writing headers %s", err.Error())
	}
	err = completeWrite.Write([]string{"نام", "کد", "موجودی", "ایدی دشت", "ایدی پرتال", "قیمت"})
	if err != nil {
		log.Fatalf("Error while writing headers %s", err.Error())
	}
	portalCH := make(chan *types.PortalCSV)
	itemCH := make(chan *types.DashtOrPortal)
	go sync_db.GetItemFromCsv(path, portalCH)

	wg := new(sync.WaitGroup)
	defer wg.Wait()

	count := 0
	total := 0
	for portalItem := range portalCH {
		total += 1
		go dasht.GetItemDetailByPortalExactName(portalItem, itemCH)

	}
	for dashtItem := range itemCH {
		if dashtItem.Dasht == nil {
			count += 1
			inCompleteWriter.Write(poratlToString(*dashtItem.Portal))

		} else {
			wg.Add(1)
			time.Sleep(time.Millisecond * 300)
			completeWrite.Write(poratlToString(*dashtItem.Portal))
			// go updatePortalPoroductSKU(token, dashtItem.Dasht, wg)
		}
	}
	wg.Wait()
	log.Printf("finished setting sku for all the products missed count: %d", count)
}
func intfrombytes(b []uint8) int {
	num, err := strconv.ParseFloat(string(b), 10)
	if err != nil {
		log.Fatalf("Error while parsing the stock/price %s", err.Error())
	}
	if num < 0. {
		return 0
	}
	return int(num)
}

func updatePortalPoroductSKU(token string, detail *types.ItemDetail, wg *sync.WaitGroup) {
	defer wg.Done()
	url := fmt.Sprintf("%ssite/api/v1/manage/store/products/variants/%d", base_url, detail.VariantID)
	variant := GetVariant(token, detail.VariantID)
	method := "PUT"
	tomanPrice := intfrombytes(detail.Fee) / 10
	// Handling prices
	variant.Price = tomanPrice
	pseudoOff := rand.IntN(16) + 5 // random number [5,20]
	variant.ComparePrice = tomanPrice * (pseudoOff/100 + 1)

	//handle quantity and shippings
	variant.Stock = intfrombytes(detail.Quantity)
	variant.Minimum = 1
	variant.Maximum = max(variant.Stock, variant.Stock-5)

	payLoadB, err := json.Marshal(variant)
	if err != nil {
		log.Printf("Error while marshaling varaint for update %s", err.Error())
		return
	}
	payload := bytes.NewReader(payLoadB)

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
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}
