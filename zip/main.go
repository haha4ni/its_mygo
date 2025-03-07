package main

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

// 判斷是否為圖片檔案
func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return true
	default:
		return false
	}
}

// 提取檔名中的數字部分，用於自然排序
var numRegex = regexp.MustCompile(`\d+`)

func extractNumber(s string) int {
	matches := numRegex.FindString(s)
	if matches == "" {
		return 0
	}
	num, _ := strconv.Atoi(matches)
	return num
}

type ImageEntry struct {
	Filename string
	ImgIndex int
	Page     int
}

func processZip(zipPath string) []ImageEntry {
	// 打開 ZIP 檔案
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Fatalf("無法打開 ZIP 檔案: %v", err)
	}
	defer reader.Close()

	var images []ImageEntry

	// 遍歷 ZIP 內的所有檔案，記錄原始順序
	for page, file := range reader.File {
		if !file.FileInfo().IsDir() && isImageFile(file.Name) {
			images = append(images, ImageEntry{Filename: file.Name, ImgIndex: page})
		}
	}

	// 依據數字部分排序，設定 Page
	sortedImages := make([]ImageEntry, len(images))
	copy(sortedImages, images)
	sort.Slice(sortedImages, func(i, j int) bool {
		return extractNumber(sortedImages[i].Filename) < extractNumber(sortedImages[j].Filename)
	})

	for page, entry := range sortedImages {
		for i := range images {
			if images[i].Filename == entry.Filename {
				images[i].Page = page
				break
			}
		}
	}

	return images
}

func storeInDatabase(dbPath string, images []ImageEntry) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("無法打開資料庫: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT,
		img_index INTEGER,
		page INTEGER
	)`) 
	if err != nil {
		log.Fatalf("建立資料表失敗: %v", err)
	}

	stmt, err := db.Prepare("INSERT INTO images (filename, img_index, page) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatalf("準備 SQL 語句失敗: %v", err)
	}
	defer stmt.Close()

	for _, img := range images {
		_, err := stmt.Exec(img.Filename, img.ImgIndex, img.Page)
		if err != nil {
			log.Printf("插入數據失敗: %v", err)
		}
	}
}

func main() {
	zipPath := "example.zip" // 替換為你的 ZIP 檔案路徑
	dbPath := "images.db"     // SQLite 資料庫路徑

	images := processZip(zipPath)
	storeInDatabase(dbPath, images)

	fmt.Println("圖片索引對應表:")
	for _, img := range images {
		fmt.Printf("Filename: %s, ImgIndex: %d, Page: %d\n", img.Filename, img.ImgIndex, img.Page)
	}
}
