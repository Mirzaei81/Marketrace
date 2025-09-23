package portal

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/joho/godotenv"
)

var token string

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
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
