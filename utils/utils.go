package utils

import (
	"encoding/json"
	"fmt"
	"giv/types"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	Jalaali "github.com/yaa110/go-persian-calendar"
)

func GetVariants(token string, ch chan *os.File) {
	url := types.PORTAL_BASE_URL + "site/api/v1/manage/store/products/variants/export"
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
	csvPath := new(types.CSVPath)
	decoder.Decode(csvPath)
	log.Printf("DEBUG: downloading csv path from %s \n", csvPath.Path)
	go getCSV(csvPath.Path, token, ch)
}
func getCSV(uri string, token string, ch chan *os.File) {
	url := types.PORTAL_BASE_URL + uri
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
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	todayJ := Jalaali.Now().AddDate(0, 0, 0).Format("yyy-MM-dd")
	err = os.WriteFile(todayJ+".csv", body, 0777)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open(todayJ + ".csv")
	if err != nil {
		log.Fatal(err)
	}
	ch <- f
}

type Throttler struct {
	invocations int
	Duration    time.Duration
	m           sync.Mutex
	Max         int
}

func (t *Throttler) Throttle() {
	t.m.Lock()
	defer t.m.Unlock()
	t.invocations += 1
	if t.invocations >= t.Max {
		<-time.After(t.Duration) // This block until the time period expires
		t.invocations = 0
	}
}
