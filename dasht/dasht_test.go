package dasht

import (
	sync_db "giv/sync_db"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// "AccessToken":  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjEzLCJleHBpcmVUaW1lIjoiXC9EYXRlKDE3NjYzMzgzOTk3MjUpXC8ifQ.dimsvFPt0aTayyCu9nnEYpYSW7ayRnYDbcOOu_qhLXU",
// "RefreshToken":  "ggvNZhKPXCdcekDRGlQa2RS5A6KaocQCv/alKIli2EyQ1HKWjHeg86mqEpctdfcUE7blUgRPJPsh85fzEkEcZw=="
var accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjEzLCJleHBpcmVUaW1lIjoiXC9EYXRlKDE3NjYzMzgzOTk3MjUpXC8ifQ.dimsvFPt0aTayyCu9nnEYpYSW7ayRnYDbcOOu_qhLXU"
var refreshToken = "hMOjOUFgMrxw4TpUXbzxAngSvTirHlq71+oJK65bHewLFYpwqA6qAhYxdR5yg/Xy4sjhTRsoa/I/JLFrlaMGSg=="

func TestMain(m *testing.M) {
	godotenv.Load("../.env")
	sync_db.Init_kv_db()
	sync_db.InitSQL(true, false, "")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	code := m.Run()
	os.Exit(code)
}

func TestListAllItems(t *testing.T) {
	allItems, err := ListAllItems()
	if err != nil {
		t.Errorf("Error while ListingElements %s", err.Error())
	}
	if len(allItems) == 0 {
		t.Error("No Item Where Found")
	}
	t.Logf("Item List All Item %+v", allItems[:10])

}
