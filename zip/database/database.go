package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "modernc.org/sqlite"
)

// 定義通用存儲接口
type Storable interface {
	TableName() string
	Fields() map[string]any
}

// 初始化資料庫
func InitDB(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("無法打開資料庫: %v", err)
	}
	return db
}

// 創建表格（如果不存在）
func createTable(db *sql.DB, tableName string, fields map[string]any) {
	var fieldDefs []string
	for field := range fields {
		fieldDefs = append(fieldDefs, fmt.Sprintf("%s TEXT", field))
	}
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY AUTOINCREMENT, %s)", 
		tableName, strings.Join(fieldDefs, ", "))

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("建立表格 %s 失敗: %v", tableName, err)
	}
}

// 儲存通用數據
func StoreData(db *sql.DB, data []Storable) {
	if len(data) == 0 {
		return
	}

	tableName := data[0].TableName()
	fields := data[0].Fields()

	createTable(db, tableName, fields)

	// 準備 SQL 語句
	columns := []string{}
	values := []string{}
	for field := range fields {
		columns = append(columns, field)
		values = append(values, "?")
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("準備 SQL 失敗: %v", err)
	}
	defer stmt.Close()

	// 插入資料
	for _, record := range data {
		fieldValues := []any{}
		for _, value := range record.Fields() {
			fieldValues = append(fieldValues, value)
		}
		_, err := stmt.Exec(fieldValues...)
		if err != nil {
			log.Printf("插入數據失敗: %v", err)
		}
	}
}
