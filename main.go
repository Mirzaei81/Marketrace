package main

import (
	"flag"
	"giv/dasht"
	"giv/portal"
	"giv/types"
	"giv/update"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"sync"
	"syscall"
	"time"

	"giv/sync_db"

	godotenv "github.com/joho/godotenv"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var env_vars = [7]string{
	"WEB_TOKEN",
	"ITEM_DETAIL_ID",
	"PORTAL_PASS",
	"PORTAL_USER",
	"LAST_PORTAL_PURCHASE",
	"LAST_GIV_PURCHASE",
	"SUCCESS",
}

func gracefulShutdown() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s
		log.Println("Sutting down gracefully.")
		// clean up here
		os.Exit(0)
	}()
}

func main() {
	err := godotenv.Load()
	go func() { log.Println(http.ListenAndServe("127.0.0.1:8000", nil)) }()
	if err != nil {
		log.Fatalf("Error while Loading .env file  %s \n", err)
	}

	var windowsAuth = flag.Bool("auth", true, "should use windows authnication to connect to mssql")
	var shouldDebug = flag.Bool("debug", true, "should debug")
	var setPrice = flag.Bool("setPrice", false, "should update portal price while updating")
	var mode = flag.String("mode", "order", "At which mode does program run on? (order|stock|local|bootstrap|stats) ")
	var csv_path = flag.String("csv", "./bk.csv", "Path for csv to bulk insert date in table VariantItems")
	var memLimit = flag.String("mem", "5G", "Memory in Format of %d[G|M]")
	var shutdownHour = flag.Int("hour", 7, "Set time when app get's shut down:[0-24) -1 for always on")
	var logFile = flag.String("logName", "main.log", "name of ther output logfile")
	var fieldName = flag.String("fieldname", "cash_on_delivery", "field name to remove(cash_on_delivery)")
	if *shutdownHour != -1 {
		go func() {
			for range time.Tick(time.Hour*1 + 1*time.Minute) {
				if time.Now().Hour() == *shutdownHour {
					go gracefulShutdown()
				}
			}
		}()
	}
	if len(*memLimit) != 0 {
		n := len(*memLimit) - 1
		var unit int64
		if (*memLimit)[n] == 'G' {
			unit = 1e9
		} else {
			unit = 1e6
		}
		number, err := strconv.ParseInt((*memLimit)[:n], 10, 32)
		if err != nil {
			log.Fatalf("Invalid number for %s %d %d", *memLimit, unit, number)
		}
		debug.SetMemoryLimit(number * unit << 20)
	}
	flag.Parse()

	update.SetPrice = *setPrice
	log.SetOutput(&lumberjack.Logger{
		Filename:   *logFile,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     10,
		Compress:   true,
	})
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	dasht.Debug = *shouldDebug
	types.Debug = *shouldDebug

	sync_db.Init_kv_db()
	sync_db.InitSQL(*shouldDebug, *windowsAuth, *csv_path)
	token := portal.Make_session()
	wg := new(sync.WaitGroup)
	accTokenB, err := sync_db.KV_DB.Read(types.DASHT_ACCESS_TOKEN)
	var accToken string
	if err != nil {
		log.Printf("Error while fetching access token from kv %s", err)
		acc, err := dasht.GetAcc()
		if err != nil {
			log.Fatal(err)
		}
		detail, err := dasht.Login(acc.GUID)
		if err != nil {
			log.Fatal(err)
		}
		accToken = detail.AccessToken

	} else {
		accToken = string(accTokenB)
	}
	switch *mode {
	case "order":
		{
			go portal.SyncDashtByPortalOrders(accToken, token, wg)
			for range time.Tick(time.Minute * 5) {
				log.Printf("Syncing Begineing ...\n")
				wg.Add(1)
				//Syncing GIV Items via portal orders
				go portal.SyncDashtByPortalOrders(accToken, token, wg)
				//updating portal Product with  giv quantity on hand
				// go dasht.SyncPortalByDashtOrders(token, wg)
				wg.Wait()
			}
		}
	case "stock":
		{
			portal.SyncVariants(token)
			for range time.Tick(time.Minute * 35) {
				log.Println("Sync internal Database From new Portal Entries")
				portal.SyncVariants(token)
			}
		}
	case "local":
		f, err := os.Open(*csv_path)
		if err != nil {
			log.Fatalf("error While Opening the file %s\n", err.Error())
		}
		portal.GetAndUpdateItemFromCsv(token, f)
	case "bootstrap":
		{

			initialDate, found := os.LookupEnv("LAST_CREATED_DATE")
			if !found {
				initialDate = "2025/08/23"
			}
			ch := make(chan types.ItemDetail)
			go dasht.GetItemByCreationDate(initialDate, ch, true)
			for item := range ch {
				if item.ItemID != 0 {
					prod, err := portal.SearchByName(item.Title, token)
					if err != nil || prod == nil {
						// no Variant in portal
						if *shouldDebug {
							log.Printf("[INFO]: Creating Product %s", item.ToString())
						}
						go portal.CreateProduct(item, token)
					} else {
						go update.Update_Variants(token, strconv.FormatInt(int64(prod.ID), 10), types.ToInt(item.Quantity), strconv.FormatInt(item.ItemID, 10), types.ToInt(item.Fee),"", nil)
						if *shouldDebug {
							log.Printf("[INFO]: Updating  Variant %s", item.ToString())
						}
					}
				} else {
					log.Printf("Error invalid Item creationDate %s", item.ToString())
				}
			}
		}
	case "stats":{
		f,_ :=  os.Open(*csv_path)
		portal.UpdateVariantsField(token,f,*fieldName)
	
	}
	default:
		log.Fatalf("Mode %s is not support please choose (order,stock)", *mode)
	}

}
