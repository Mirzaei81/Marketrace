package utils

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"giv/types"
	"io"
	"log"
	"net/http"
	"os"

	Jalaali "github.com/yaa110/go-persian-calendar"
)

func GetVariants(token string, ch *chan *csv.Reader) {
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
func getCSV(uri string, token string, ch *chan *csv.Reader) {
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
