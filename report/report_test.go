package report

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	code := m.Run()
	os.Exit(code)
}
func TestReport(t *testing.T) {
	err := ReportFile("test", "imcomplete.csv")
	if err != nil {
		t.Errorf("Error while reporting incomplete file %s", err.Error())
	}

}
