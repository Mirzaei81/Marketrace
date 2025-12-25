package update

import (
	"bytes"
	"encoding/json"
	"fmt"
	"giv/types"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var SetPrice bool

func Update_Variants(token string, variant_id string, stock int, sku int64, price int, wg *sync.WaitGroup) {
	defer wg.Done()
	url := fmt.Sprintf("https://modernhyperindustry.com/site/api/v1/manage/store/products/variants/%s", variant_id)
	method := "PUT"
	variant := getVariant(token, variant_id)
	if variant == nil {
		return
	}
	time.Sleep(time.Microsecond * 300)
	variant.Stock = stock
	if (price != 0 && variant.ComparePrice == 0 && variant.Price == 0) || SetPrice {
		variant.Price = int(price / 10) //TODO :  check
		variant.ComparePrice = int(price / 10)
	}
	variant.Sku = strconv.FormatInt(sku, 10)
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
	url := fmt.Sprintf("https://modernhyperindustry.com/site/api/v1/manage/store/products/variants/%s", variantid)
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

	dec := json.NewDecoder(res.Body)
	variant := new(types.VariantResult)
	err = dec.Decode(variant)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &variant.Variant
}

func updateProductSKU(token string, productId int, sku string) {
	url := fmt.Sprintf("https://modernhyperindustry.com/site/api/v1/manage/store/products/%d", productId)
	method := "PUT"

	payloadString := fmt.Sprintf(`{ "sku ":%s }`, sku)
	payload := strings.NewReader(payloadString)

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
