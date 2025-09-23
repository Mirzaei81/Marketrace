package sync_db

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"
	"github.com/peterbourgon/diskv/v3"
)

var SQL_DB *sqlx.DB
var KV_DB *diskv.Diskv

func InitSQL(debug bool, windowsAuth bool, csvPath string) {

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
		db = "GivKohancharm04"
	}

	url := buildSQLServerURL(host, db, port, username, password, windowsAuth)
	var err error
	SQL_DB, err = sqlx.Connect("mssql", url)
	if err != nil {
		log.Fatalf("Open connection failed %s", err)
	}

	var dbName string
	SQL_DB.Get(&dbName, "SELECT DB_NAME() AS [Current Database];")

	if debug {
		log.Print("Datbase DSN : ", host, "\t", db, "\t", port, "\t", username, "\t", password, "\t", windowsAuth,
			" current selected is ", dbName)
	}
	bootStrapTables(csvPath)
}
func bootStrapTables(csvPath string) {
	dropProcStmt := `IF EXISTS (
	SELECT type_desc, type
	FROM sys.procedures WITH(NOLOCK)
	WHERE NAME = 'BootVarToItemTable'
	AND type = 'P'
	)
	DROP PROCEDURE dbo.BootVarToItemTable;
	`
	SQL_DB.MustExec(dropProcStmt)
	createProcStmt := `CREATE PROCEDURE [dbo].BootVarToItemTable
	(
	@CSV_PATH nvarchar(255) = NULL
	)
	AS 
	BEGIN
	-- BootVarToItemTable
	IF OBJECT_ID('VariantsItems') IS NOT NULL
	DROP TABLE VariantsItems ; 
	IF NOT EXISTS (
	SELECT 1
	FROM sys.key_constraints  where name = 'uq_ItemBarCode_Barcode'
	)
	BEGIN
	ALTER TABLE ItemBarCode
	ADD CONSTRAINT uq_ItemBarCode_Barcode UNIQUE (ItemBarCode);
	END
	
	CREATE TABLE VariantsItems(
	VariantID BIGINT  primary key,
	Title varchar(255) null,
	ItemID  varchar(50) null,
	Price decimal null,
	ComparePrice bigint null,
	constraint fk_VariantItemID foreign KEY (ItemId) references ItemBarCode (ItemBarCode)
	)
	EXEC('BULK INSERT dbo.VariantsItems
	FROM ''' + 
	@CSV_PATH + 
	''' WITH (FIRSTROW = 2,ROWTERMINATOR=''\n'', FIELDTERMINATOR = '','', MAXERRORS = 0);')
	ALTER TABLE VariantsItems
	drop column Title,Price,ComparePrice
	END;`
	SQL_DB.MustExec(createProcStmt)
	if len(csvPath) > 0 {
		SQL_DB.MustExec("EXEC BootVarToItemTable @CSV_PATH='" + csvPath + "'")
	}
}

func buildSQLServerURL(host, database string, port int, user, password string, windowsAuth bool) string {
	query := url.Values{}
	query.Add("database", database)

	if windowsAuth {
		// Windows Authentication (integrated security)
		query.Add("trusted_connection", "yes")
		return fmt.Sprintf("sqlserver://%s:%d?%s", host, port, query.Encode())
	}

	// SQL Authentication with username + password
	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", user, url.QueryEscape(password), host, port, database)

	return connString
}
func Init_kv_db() {
	flatTransform := func(s string) []string { return []string{} }
	KV_DB = diskv.New(diskv.Options{
		BasePath:     "../portal_DB",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
}
