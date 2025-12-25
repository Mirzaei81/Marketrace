package givsoft

import (
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
	log.Printf("%s : %#v\n", itemDetail.ItemID, itemDetail)
}
func TesGetOrCreateUserIDExist(t *testing.T) {
	userId := getOrCreateUserID("1")
	if userId != -1 {
		t.Errorf("Failed : UserId 1 doesn't exists in table Persion")
	}
}

func TesGetOrCreateUserIDNotExist(t *testing.T) {
	userId := getOrCreateUserID("-1")
	if userId == -1 {
		t.Errorf("Failed : UserId -1 exists in table Persion")
	}
}

func TestMakeOrder(t *testing.T) {
	orders := getNewOrders(23835)
	t.Logf("orders after 23835 %#v\n", orders)
}

// func  test_order(){
// }
