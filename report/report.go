package report

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var eitaToken string
var chatID string

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal(err)
	}
	var found bool
	eitaToken, found = os.LookupEnv("EIYTA_TOKEN")
	if !found {
		log.Fatal("EIYTA Token not found ")
	}
	chatID, found = os.LookupEnv("CHAT_ID")
	if !found {
		log.Fatal("ChatID not found ")
	}

}
func ReportFile(caption string, filename string) error {
	url := fmt.Sprintf("https://eitaayar.ir/api/%s/sendFile", eitaToken)
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open(filename)
	if errFile1 != nil {
		log.Println(errFile1)
		return errFile1
	}
	defer file.Close()
	part1,
		errFile1 := writer.CreateFormFile("file", filename)
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		log.Println(errFile1)
		return errFile1
	}
	_ = writer.WriteField("chat_id", chatID)
	_ = writer.WriteField("caption", caption)
	_ = writer.WriteField("date", "0")
	_ = writer.WriteField("parse_mode", "")
	_ = writer.WriteField("pin", "on")
	_ = writer.WriteField("viewCountForDelete", "")
	err := writer.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("origin", "https://eitaayar.ir")
	req.Header.Add("priority", "u=0, i")
	req.Header.Add("referer", "https://eitaayar.ir/testApi")
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"143\", \"Chromium\";v=\"143\", \"Not A(Brand\";v=\"24\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return nil
}
