package sync_db

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	types "giv/types"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.env")
	Init_kv_db()
	InitSQL(true, false, "/VAR/OPT/MSSQL/DATA/batkap.csv")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	code := m.Run()
	os.Exit(code)
}
func TestSQLX_EXISTS(t *testing.T) {
	dirs, err := os.ReadDir("../portal_DB/")
	if err != nil {
		log.Fatalf("Error reading portal db %s \n", err)
	}
	n := rand.Intn(len(dirs) - 4) // last 4 dirs are for meta data
	itemId, err := strconv.ParseInt(dirs[n].Name(), 10, 0)
	if err != nil {
		log.Fatalf("Error while parsing itemId to int  %s \n", err)
	}

	givItem := getItem(int(itemId))
	t.Logf("%#v", givItem)
}
func getItem(itemId int) types.GivItems {
	var itemDetail types.GivItems
	SQL_DB.Get(&itemDetail, `
		SELECT qoh.ItemQuantityOnHand,p.ItemCurrentSelPrice,i.ItemID,i.VariantID
		from [dbo].QuantityOnHand as qoh join 
		[dbo].ItemParent p on p.ItemParentID = qoh.ItemID/10000000 
		JOIN VariantsItems i on i.ItemID =qoh.ItemID
		where StockID = 4 AND i.ItemID = $1
		`, itemId)
	return itemDetail
}
