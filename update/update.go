package update

import (
	"bytes"
	"encoding/json"
	"fmt"
	"giv/types"
	"io"
	"log"
	"math"
	"math/rand/v2"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var SetPrice bool

func Update_Variants(token string, variant_id string, stock int, sku string, price int, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	url := fmt.Sprintf("https://modernhyperindustry.com/site/api/v1/manage/store/products/variants/%s", variant_id)
	method := "PUT"
	variant := getVariant(token, variant_id)
	if variant == nil {
		return
	}
	variant.Stock = stock
	if (price != 0 && variant.ComparePrice == 0 && variant.Price == 0) || SetPrice {
		variant.Price = int(price / 10) //TODO :  check
		variant.ComparePrice = int(price / 10)
	}

	variant.Sku = sku
	product_byte, err := json.Marshal(variant)
	log.Printf("Updating Variant %s\n", string(product_byte))
	fmt.Printf("Updating Variant %s\n", string(product_byte))
	if err != nil {
		log.Println(string(product_byte))
	}
	payload := bytes.NewReader(product_byte)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Println(err)
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
}

func getVariant(token string, variantid string) *types.Variant {
	url := fmt.Sprintf("%ssite/api/v1/manage/store/products/variants/%s", types.PORTAL_BASE_URL, variantid)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	variant := new(types.VariantResult)
	err = json.Unmarshal(body, &variant)
	if err != nil {
		log.Println(err)
		return nil
	}
	if !variant.Success {
		log.Printf("Erro while fethcing variant %s", string(body))
		return nil
	}
	return &variant.Variant
}

func UpdatePortalVariantSKU(token string, detail *types.ItemDetail, wg *sync.WaitGroup) {
	defer wg.Done()
	url := fmt.Sprintf("%ssite/api/v1/manage/store/products/variants/%d", types.PORTAL_BASE_URL, detail.VariantID)
	variant := getVariant(token, strconv.FormatInt(detail.VariantID, 10))
	if variant == nil {
		return
	}
	time.Sleep(time.Millisecond * 333) // sleep for the next request
	method := "PUT"
	tomanPrice := types.ToInt(detail.Fee) / 10
	// Handling prices
	variant.Price = tomanPrice
	pseudoOff := rand.Float64()*16 + 5 // random number [5,20]
	if variant.Price != variant.ComparePrice {
		variant.ComparePrice = int(math.Round(float64(tomanPrice)*(pseudoOff/100+1)/1000) * 1000)
	}

	//handle quantity and shippings
	variant.Stock = types.ToInt(detail.Quantity)
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
	log.Printf("%s \n%s\n", string(payLoadB), string(body))
	fmt.Printf("%s \n%s\n", string(payLoadB), string(body))

}
