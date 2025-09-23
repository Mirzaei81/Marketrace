package main

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb"
)

type Product_resault struct {
	Success bool `json:"success"`
}
type csvPath struct {
	Success     bool   `json:"success"`
	Description string `json:"description"`
	Path        string `json:"path"`
}

var base_url = "https://batkap.com"

func abort(funcname string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcname, err))
}

// var (
// 	kernel32, _        = syscall.LoadLibrary("kernel32.dll")
// 	getModuleHandle, _ = syscall.GetProcAddress(kernel32, "GetModuleHandleW")
//
// 	user32, _     = syscall.LoadLibrary("user32.dll")
// 	messageBox, _ = syscall.GetProcAddress(user32, "MessageBoxW")
// )
//
// const (
// 	MB_OK                = 0x00000000
// 	MB_OKCANCEL          = 0x00000001
// 	MB_ABORTRETRYIGNORE  = 0x00000002
// 	MB_YESNOCANCEL       = 0x00000003
// 	MB_YESNO             = 0x00000004
// 	MB_RETRYCANCEL       = 0x00000005
// 	MB_CANCELTRYCONTINUE = 0x00000006
// 	MB_ICONHAND          = 0x00000010
// 	MB_ICONQUESTION      = 0x00000020
// 	MB_ICONEXCLAMATION   = 0x00000030
// 	MB_ICONASTERISK      = 0x00000040
// 	MB_USERICON          = 0x00000080
// 	MB_ICONWARNING       = MB_ICONEXCLAMATION
// 	MB_ICONERROR         = MB_ICONHAND
// 	MB_ICONINFORMATION   = MB_ICONASTERISK
// 	MB_ICONSTOP          = MB_ICONHAND
//
// 	MB_DEFBUTTON1 = 0x00000000
// 	MB_DEFBUTTON2 = 0x00000100
// 	MB_DEFBUTTON3 = 0x00000200
// 	MB_DEFBUTTON4 = 0x00000300
// )
//
// func DialogBox(caption, text string, style uintptr) (result int) {
// 	var nargs uintptr = 4
// 	ret, _, callErr := syscall.Syscall9(uintptr(messageBox),
// 		nargs,
// 		0,
// 		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
// 		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
// 		style,
// 		0,
// 		0,
// 		0,
// 		0,
// 		0)
// 	if callErr != 0 {
// 		abort("Call MessageBox", callErr)
// 	}
// 	result = int(ret)
// 	return
// }
//
// func GetModuleHandle() (handle uintptr) {
// 	var nargs uintptr = 0
// 	if ret, _, callErr := syscall.Syscall(uintptr(getModuleHandle), nargs, 0, 0, 0); callErr != 0 {
// 		abort("Call GetModuleHandle", callErr)
// 	} else {
// 		handle = ret
// 	}
// 	return
// }
//
// func main() {
// 	defer syscall.FreeLibrary(kernel32)
// 	defer syscall.FreeLibrary(user32)
// 	fmt.Printf("Input text was: %s\n")
// }

var pool *sql.DB

func main() {
	godotenv.Load("./.env")
	var windowsAuth = flag.Bool("auth", true, "should use windows authnication to connect to mssql")
	var debug = flag.Bool("debug", true, "should debug")
	flag.Parse()

	username, exists := os.LookupEnv("DB_USER")
	if !exists {
		log.Fatal("user not found")
	}
	password, exists := os.LookupEnv("PASS")
	if !exists {
		log.Fatal("pass not found")
	}

	host, exists := os.LookupEnv("HOST")
	if !exists {
		host = "localhost"
	}
	s_port, exists := os.LookupEnv("PORT")
	var port int = 1433
	if exists {
		var e error
		port, e = strconv.Atoi(s_port)
		if e != nil {
			log.Fatal("port format is invalid")
		}
	}
	db, exists := os.LookupEnv("DB")
	if !exists {
		db = "localhost"
	}
	if *debug {
		log.Print(host, "\t", db, "\t", port, "\t", username, "\t", password, "\t", *windowsAuth)
	}

	url := buildSQLServerURL(host, db, port, username, password, *windowsAuth)
	pool, err := sql.Open("mssql", url)
	if err != nil {
		log.Fatalf("Open connection failed %s", err)
	}
	stmt, err := pool.Prepare("SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE'")
	if err != nil {
		log.Fatalf("failed to prepare stmt %s for connection string %s", err, url)
	}
	defer stmt.Close()

	row := stmt.QueryRow()
	var catlog string
	var schema string
	var name string
	var tb_type string
	err = row.Scan(&catlog, &schema, &name, &tb_type)
	if err != nil {
		log.Fatal("Scan failed:", err.Error())
	}
	if err != nil {
		log.Fatal("unable to execute search query", err)
	}
	log.Printf(`
		catlog  = %s 
                schema  = %s 
                name = %s 
	        tb_type = %s`, catlog, schema, name, tb_type)
}

func buildSQLServerURL(host, database string, port int, user, password string, windowsAuth bool) string {
	query := url.Values{}
	query.Add("database", database)

	if windowsAuth {
		// Windows Authentication (integrated security)
		// Works only if running on Windows and driver supports it
		query.Add("trusted_connection", "yes")
		return fmt.Sprintf("sqlserver://%s:%d?%s", host, port, query.Encode())
	}

	// SQL Authentication with username + password
	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%d", host, port),
		// Path must begin with "/" if you want default DB pre-selected
		Path: "/" + database,
	}
	return u.String()
}

func Query() {

	defer pool.Close()
}

func readCsv(uri string, token string) {
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
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	reader := csv.NewReader(bytes.NewBuffer(body))
	_, err = reader.Read()
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
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
			fmt.Printf("Updating variant : %s with Sku Of %s with id of %d", line[2], line[8], itemId)
		}
	}
}
func readCsvbk(reader *csv.Reader, token string) {
	_, err := reader.Read()
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				log.Println(err)
				break
			}
		}
		if line[2] != "" {
			itemId, _ := strconv.ParseInt(line[0], 10, 64)
			fmt.Printf("Updating variant : %s with Sku Of %s with id of %d\n", line[1], line[2], itemId)
			//givsoft.QuantityOnhand_byitem(token, line[2], int(itemId), true)
		}
	}
}
