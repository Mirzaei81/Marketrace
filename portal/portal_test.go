package portal

import (
	"encoding/csv"
	"giv/dasht"
	"giv/sync_db"
	"giv/types"
	"io"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var token string

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}
	types.Debug = true
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	sync_db.InitSQL(true, true, "")
	sync_db.Init_kv_db()
	token = Make_session()
	code := m.Run()
	os.Exit(code)
}

func TestFetchOrders(t *testing.T) {
	orders, _ := getOrders(token)
	if orders.Count < 0 {
		t.Errorf("Oders are less then expected values %d", orders.Count)
	}
	t.Logf("%#v", orders)
}

//	func TestVariantByPrice(t *testing.T) {
//		update.SetPrice = true
//		ch := make(chan *os.File)
//		utils.GetVariants(token, ch)
//		file := <-ch
//		close(ch)
//		reader := csv.NewReader(file)
//		_, err := reader.Read() // reading header
//		if err != nil {
//			t.Errorf("Error while reading header for variants csv %s", err)
//		}
//		line, err := reader.Read()
//		wg := new(sync.WaitGroup)
//		wg.Add(1)
//
//		if line[8] == "" {
//			t.Errorf("Failure : Invalid/Empty sku from variants in csv %s", err)
//		}
//		variantId, err := strconv.ParseInt(line[0], 10, 32)
//		initialVaraint := GetVariant(token, variantId)
//		itemId, _ := strconv.ParseInt(line[0], 10, 64)
//		log.Printf("Updating variant : %s with Sku Of %s", line[2], line[8])
//		// just insert into VariantsItem  and update  table accordingly
//		// csv is in form ID ProductID Title Price ComparePrice Type Status Stock Sku
//		givsoft.SyncPortalVariantWithGivQOH(token, line[8], int(itemId), wg)
//		wg.Wait()
//		syncesVaraint := GetVariant(token, variantId)
//		t.Logf("Initial %#v after update state is %#v", initialVaraint, syncesVaraint)
//	}
//
//	func TestFetchVariant(t *testing.T) {
//		ch := make(chan *os.File)
//		utils.GetVariants(token, ch)
//
//		f := <-ch
//		close(ch)
//		reader := csv.NewReader(f)
//		count := readCsvbk(reader, token)
//		if count <= 0 {
//			t.Errorf("Incomplete date  while fetching variant from portal using csv export  current count %d", count)
//		}
//		t.Logf("Portal csv export  variant  count %d", count)
//	}
func readCsvbk(reader *csv.Reader, token string) int {
	count := 0
	_, err := reader.Read()
	if err != nil {
		log.Print(err)
		os.Exit(-1)
	}
	for {
		_, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		count++
	}
	return count
}
func TestCreateItemOne(t *testing.T) {
	ch := make(chan types.ItemDetail)
	testDate := "2023/01/01"

	go dasht.GetItemByCreationDate(testDate, ch, true)

	count := 0
	for item := range ch {
		count++
		if item.ItemID == 0 {
			t.Errorf("Expected ItemID to be populated, got %d", item.ItemID)
		}
		if item.Title == "" {
			t.Error("Expected Title to be populated, got empty string")
		}
		if count == 1 {
			CreateProduct(item, token)
		}
	}

	t.Logf("Successfully processed %d items from database", count)
}
