package main

import (
	"flag"
	"giv/givsoft"
	"giv/portal"
	"giv/update"
	"log"
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
	var debug = flag.Bool("debug", true, "should debug")
	var setPrice = flag.Bool("setPrice", false, "should update portal price while updating")
	var mode = flag.String("mode", "order", "At which mode does program run on? (order|stock) ")
	var csv_path = flag.String("csv", "./bk.csv", "Path for csv to bulk insert date in table VariantItems")
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
	if err != nil {
		log.Printf("Error:error while  loading .env %s\n", err)
	}

	sync_db.Init_kv_db()
	sync_db.InitSQL(*debug, *windowsAuth, *csv_path)
	token := portal.Make_session()
	wg := new(sync.WaitGroup)
	if *mode == "order" {

		for range time.Tick(time.Minute * 5) {
			log.Printf("Syncing Begineing ...\n")
			wg.Add(2)
			//Syncing GIV Items via portal orders
			go portal.SyncGivByPortalOrders(token, wg)
			//updating portal Product with  giv quantity on hand
			go givsoft.SyncPortalByGivOrders(token, wg)
			wg.Wait()
		}

	} else if *mode == "stock" {
		givsoft.SyncPortalWithGivQOH(token)
		portal.SyncVariants(token)
		for range time.Tick(time.Minute * 35) {
			log.Println("Sync internal Database From new Portal Entries")
			portal.SyncVariants(token)
		}
	} else {
		log.Fatalf("Mode %s is not support please choose (order,stock)", mode)
	}

}
