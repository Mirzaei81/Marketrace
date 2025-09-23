package digi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
)

type authJson struct {
	access_token  string
	refresh_token string
}

type refreshBody struct {
	status string   `json:"string"`
	data   authJson `json:"data"`
}
type Item struct {
	ID                      int       `json:"id"`
	ImageSrc                string    `json:"image_src"`
	SellerID                int       `json:"seller_id"`
	MainCategoryTitle       string    `json:"main_category_title"`
	CategoryID              int       `json:"category_id"`
	ProductID               int       `json:"product_id"`
	ProductURL              string    `json:"product_url"`
	ProductVariantID        int       `json:"product_variant_id"`
	SupplierCode            string    `json:"supplier_code"`
	ProductModerationStatus string    `json:"product_moderation_status"`
	Title                   string    `json:"title"`
	ProductTitle            string    `json:"product_title"`
	Active                  bool      `json:"active"`
	LeadTime                int       `json:"lead_time"`
	PriceList               int       `json:"price_list"`
	MarketPriceLastUpdate   string    `json:"market_price_last_update"`
	PriceType               string    `json:"price_type"`
	SellingChannelSite      string    `json:"selling_channel_site"`
	PriceSale               int       `json:"price_sale"`
	MarketplaceSellerStock  int       `json:"marketplace_seller_stock"`
	WarehouseStock          int       `json:"warehouse_stock"`
	OnTheWayStock           int       `json:"on_the_way_stock"`
	Reservation             int       `json:"reservation"`
	LeftConsumer            int       `json:"left_consumer"`
	MaximumPerOrder         int       `json:"maximum_per_order"`
	AllowedCount            int       `json:"allowed_count"`
	OvlSellingActive        bool      `json:"ovl_selling_active"`
	CreatedAt               time.Time `json:"created_at"`
	PolActive               bool      `json:"pol_active"`
	B2BParams               struct {
		SellerB2BActive bool `json:"seller_b2b_active"`
		IsOnlyB2B       bool `json:"is_only_b2b"`
		IsB2BActive     bool `json:"is_b2b_active"`
	} `json:"b2b_params"`
	MaxLeadTime      int    `json:"max_lead_time"`
	BuyBoxPrice      int    `json:"buy_box_price"`
	BuyBoxBadgeLabel string `json:"buy_box_badge_label"`
	IsBuyBoxWinner   bool   `json:"is_buy_box_winner"`
	IsSkuWinner      bool   `json:"is_sku_winner"`
	SkuConfig        struct {
		Color string `json:"color"`
		Size  string `json:"size"`
	} `json:"sku_config"`
	IsSellerBuyBoxWinner bool   `json:"is_seller_buy_box_winner"`
	IsInBuyBoxChallenge  bool   `json:"is_in_buy_box_challenge"`
	MinSellingPriceLimit int    `json:"min_selling_price_limit"`
	SuppressedUntil      string `json:"suppressed_until"`
	SuppressionReason    string `json:"suppression_reason"`
	ProductSellingChanel struct {
		ActiveDigikala  bool `json:"active_digikala"`
		ActiveDigistyle bool `json:"active_digistyle"`
	} `json:"product_selling_chanel"`
	VariantSellingChanel struct {
		ActiveDigikala  bool `json:"active_digikala"`
		ActiveDigistyle bool `json:"active_digistyle"`
	} `json:"variant_selling_chanel"`
	IsInIncrediblePromotion               bool    `json:"is_in_incredible_promotion"`
	IsInPeriodicPromotion                 bool    `json:"is_in_periodic_promotion"`
	IsInPromotion                         bool    `json:"is_in_promotion"`
	PromotionPrice                        int     `json:"promotion_price"`
	ShippingNatureID                      int     `json:"shipping_nature_id"`
	DefaultSellingChanelCode              int     `json:"default_selling_chanel_code"`
	Rating                                float64 `json:"rating"`
	IsPromotionManagementVisibleForSeller bool    `json:"is_promotion_management_visible_for_seller"`
	IsArchived                            bool    `json:"is_archived"`
	FulfilmentAndDeliveryCost             int     `json:"fulfilment_and_delivery_cost"`
	SellerReservation                     int     `json:"seller_reservation"`
	DigikalaReservation                   int     `json:"digikala_reservation"`
	SellerShippingLeadTime                int     `json:"seller_shipping_lead_time"`
	ShippingOptions                       struct {
		IsFbsAbilityEnable bool `json:"is_fbs_ability_enable"`
		IsFbdActive        bool `json:"is_fbd_active"`
		IsFbsActive        bool `json:"is_fbs_active"`
		IsNeededFbsSetting bool `json:"is_needed_fbs_setting"`
		IsSbsModuleActive  bool `json:"is_sbs_module_active"`
		OnlySbs            bool `json:"only_sbs"`
	} `json:"shipping_options"`
}
type kalaRequestResponse struct {
	Status string `json:"status"`
	Data   struct {
		SortData struct {
			SortColumn  string   `json:"sort_column"`
			SortOrder   string   `json:"sort_order"`
			SortColumns []string `json:"sort_columns"`
		} `json:"sort_data"`
		Pager struct {
			Page        int `json:"page"`
			ItemPerPage int `json:"item_per_page"`
			TotalPages  int `json:"total_pages"`
			TotalRows   int `json:"total_rows"`
		} `json:"pager"`
		FormData []interface{} `json:"form_data"`
		Items    []Item        `json:"items"`
		MetaData struct {
		} `json:"meta_data"`
	} `json:"data"`
}

const (
	clientCode      = "TklnazVDOVUrNjdOdWR0QWFpaXcwQT09"
	accessTokenPath = "access.json"
	TAG             = "digikala"
)

func main() {
	log.SetOutput(&lumberjack.Logger{
		Filename: "./digi.log",
		MaxSize:  5,
	})
}

func refreshToken() (string, error) {
	f, err := os.Open("auth.json")
	url := "https://seller.digikala.com/open-api/v1/auth/refresh-token"
	method := "POST"
	stat, err := f.Stat()
	if err != nil {
		log.Fatalf("%s : %s", TAG, err)
	}
	buffer := make([]byte, stat.Size())
	_, err = f.Read(buffer)
	if err != nil {
		log.Fatalf("%s : %s", TAG, err)
	}
	var tokens authJson
	err = json.Unmarshal(buffer, &tokens)
	if err != nil {
		log.Fatalf("%s : %s", TAG, err)
	}
	client := &http.Client{}
	payload := bytes.NewReader(buffer)
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Fatalf("%s : %s", TAG, err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cookie", "tracker_glob_new=3HT7sjs; tracker_session=7Hi80NT")

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("%s : %s", TAG, err)
	}
	defer res.Body.Close()

	var bodyStruct refreshBody
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(body, &bodyStruct)
	f.Seek(0, 0)
	empty := make([]byte, stat.Size())
	_, err = f.Write(empty)
	if err != nil {
		fmt.Println(err)
	}
	f.Seek(0, 0)
	f.Write(body)
	return bodyStruct.data.access_token, nil
}
func readToken() string {
	f, _ := os.Open(accessTokenPath)
	body, _ := io.ReadAll(f)
	return string(body)
}
func getDigiKala(searchParam string, accessToken string) []Item {
	client := &http.Client{}
	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://seller.digikala.com/open-api/v1/variants?page=1&size=50&search[active]=true&search[search_term]=%s", searchParam), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := client.Do(req)
	if resp.StatusCode == 401 {
		token, err := refreshToken()
		if err != nil {
			log.Fatal(TAG, err)
		}
		return getDigiKala(searchParam, token)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
	var response kalaRequestResponse
	err = json.Unmarshal(bodyText, &response)
	if err != nil {
		log.Fatal(err)
	}
	return response.Data.Items
}
