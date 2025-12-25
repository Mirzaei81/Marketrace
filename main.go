package main

import (
	"flag"
	"giv/dasht"
	"giv/portal"
	"giv/update"
	"log"
	"os"
	"strconv"
	"sync"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error while Loading .env file  %s \n", err)
	}
	var windowsAuth = flag.Bool("auth", true, "should use windows authnication to connect to mssql")
	var shouldDebug = flag.Bool("debug", true, "should debug")
	var setPrice = flag.Bool("setPrice", false, "should update portal price while updating")
	var mode = flag.String("mode", "order", "At which mode does program run on? (order|stock) ")
	var csv_path = flag.String("csv", "./bk.csv", "Path for csv to bulk insert date in table VariantItems")
	var memLimit = flag.String("mem", "5G", "Memory in Format of %d[G|M]")
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
		// debug.SetMemoryLimit(number * unit << 20)
	}
	flag.Parse()

	update.SetPrice = *setPrice
	log.SetOutput(&lumberjack.Logger{
		Filename:   "./main.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     10,
		Compress:   true,
	})
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	dasht.Debug = *shouldDebug

	sync_db.Init_kv_db()
	sync_db.InitSQL(*shouldDebug, *windowsAuth, *csv_path)
	token := portal.Make_session()
	wg := new(sync.WaitGroup)
	switch *mode {
	case "order":
		for range time.Tick(time.Minute * 5) {
			log.Printf("Syncing Begineing ...\n")
			wg.Add(2)
			//Syncing GIV Items via portal orders
			go portal.SyncGivByPortalOrders(token, wg)
			//updating portal Product with  giv quantity on hand
			go dasht.SyncPortalByDashtOrders(token, wg)
			wg.Wait()
		}

	case "stock":
		portal.SyncVariants(token)
		for range time.Tick(time.Minute * 35) {
			log.Println("Sync internal Database From new Portal Entries")
			portal.SyncVariants(token)
		}
	case "sku":
		f, err := os.Open(*csv_path)
		if err != nil {
			log.Fatalf("error While Opening the file %s\n", err.Error())
		}
		portal.GetItemFromCsv(token, f)

	default:
		log.Fatalf("Mode %s is not support please choose (order,stock)", *mode)
	}

}
