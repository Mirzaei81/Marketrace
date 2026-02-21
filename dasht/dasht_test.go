package dasht

import (
	sync_db "giv/sync_db"
	"giv/types"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/natefinch/lumberjack"
)

func TestMain(m *testing.M) {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "./main.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     10,
		Compress:   true,
	})
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	godotenv.Load("../.env")
	sync_db.Init_kv_db()
	sync_db.InitSQL(true, true, "")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	code := m.Run()
	os.Exit(code)
}
func TestGetRefresh_OK(t *testing.T) {
	acc, err := GetAcc()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(acc.GUID) == 0 {
		t.Fatalf("Invalid GUID")
	}
	if len(acc.RefreshToken) == 0 {
		t.Fatalf("Invalid Refresh Token")
	}
	t.Log(acc)
}

func TestLogin_OK(t *testing.T) {
	acc, err := GetAcc()
	if err != nil {
		t.Fatalf("Error Getting acc %s", err)
	}
	resp, err := Login(acc.GUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body, err := GetSaleInvoise(resp.AccessToken)
	if err != nil {
		t.Fatalf("unexpected error Orders: %v", err)
	}

	t.Logf("accToken %s,orders: %s", resp.AccessToken, body)
}
func TestListAllItems(t *testing.T) {
	fiscal, err := GetLatestFiascal()
	if err != nil {
		t.Fatal(err)
	}
	allItems, err := ListAllItemsSync(int(fiscal))
	if err != nil {
		t.Errorf("Error while ListingElements %s", err.Error())
	}
	if len(allItems) == 0 {
		t.Error("No Item Where Found")
	}
	t.Logf("Item List All Item %+v", allItems[:10])
}
func TestGetItemByCreationDate(t *testing.T) {
	ch := make(chan types.ItemDetail)
	testDate := "2023/01/01"

	lastFiscal, _ := GetLatestFiascal()

	go GetItemByCreationDate(testDate, ch, int(lastFiscal), true)

	count := 0
	for item := range ch {
		count++
		if item.ItemID == 0 {
			t.Errorf("Expected ItemID to be populated, got %d", item.ItemID)
		}
		if item.Title == "" {
			t.Error("Expected Title to be populated, got empty string")
		}
	}

	t.Logf("Successfully processed %d items from database", count)
}
