package givsoft

import (
	"fmt"
	sync_db "giv/sync_db"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.env")
	sync_db.Init_kv_db()
	sync_db.InitSQL(true, false, "")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	code := m.Run()
	os.Exit(code)
}

func TestGetItemList(t *testing.T) {
	dirs, err := os.ReadDir("../portal_DB")
	if err != nil {
		log.Fatalf("error while opening porta_DB %s\n", err)
	}
	if len(dirs) <= 0 {
		t.Errorf("portal db empty")
	}

}
func TestGetItemDetail(t *testing.T) {
	dirs, err := os.ReadDir("../portal_DB")
	if err != nil {
		log.Fatalf("error while opening porta_DB %s\n", err)
	}

	idx := rand.Intn(len(dirs) - 5)
	itemId := dirs[idx].Name()
	itemDetail, err := GetItemDetail(itemId)
	if err != nil {
		log.Fatalf("error Fetching item detail %s :  %s\n", itemId, err)
	}
	log.Println(fmt.Sprintf("%s : %#v", itemDetail.ItemID, itemDetail))
}

// func TestMakeOrder(){
// }

// func  test_order(){
// }
