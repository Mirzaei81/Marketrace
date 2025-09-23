package portal

import (
	"encoding/csv"
	givsoft "giv/givsoft"
	"giv/sync_db"
	"giv/update"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/joho/godotenv"
)

var token string

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	sync_db.InitSQL(true, false, "")
	if err != nil {
		log.Fatal(err)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	token = Make_session()
	code := m.Run()
	os.Exit(code)
}

func TestFetchOrders(t *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	orders, _ := fetchOrders(token, wg)
	if orders.Count < 0 {
		t.Errorf("Oders are less then expected values %d", orders.Count)
	}
	wg.Wait()
	t.Logf("%#v", orders)
}
func TestVariantByPrice(t *testing.T) {
	update.SetPrice = true
	ch := make(chan *csv.Reader)
	GetVariants(token, &ch)
	reader := <-ch
	close(ch)
	_, err := reader.Read() // reading header
	if err != nil {
		t.Errorf("Error while reading header for variants csv %s", err)
	}
	line, err := reader.Read()
	wg := new(sync.WaitGroup)
	wg.Add(1)

	if line[8] == "" {
		t.Errorf("Failure : Invalid/Empty sku from variants in csv %s", err)
	}
	variantId, err := strconv.ParseInt(line[0], 10, 32)

	if line[8] == "" {
		t.Errorf("Failure : Failed while parsing variantID from csv %s", err)
	}
	initialVaraint := GetVariant(token, int(variantId))
	itemId, _ := strconv.ParseInt(line[0], 10, 64)
	log.Printf("Updating variant : %s with Sku Of %s", line[2], line[8])
	// just insert into VariantsItem  and update  table accordingly
	// csv is in form ID ProductID Title Price ComparePrice Type Status Stock Sku
	givsoft.SyncPortalVariantWithGivQOH(token, line[8], int(itemId), wg)
	wg.Wait()
	syncesVaraint := GetVariant(token, int(variantId))
	t.Logf("Initial %#v after update state is %#v", initialVaraint, syncesVaraint)

}
func TestFetchVariant(t *testing.T) {
	ch := make(chan *csv.Reader)
	GetVariants(token, &ch)

	reader := <-ch
	close(ch)
	count := readCsvbk(reader, token)
	if count <= 0 {
		t.Errorf("Incomplete date  while fetching variant from portal using csv export  current count %d", count)
	}
	t.Logf("Portal csv export  variant  count %d", count)

}
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
