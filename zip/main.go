package main

import (
	"archive/zip"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"mygo/database"
)

func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return true
	default:
		return false
	}
}

var numRegex = regexp.MustCompile(`\d+`)

func extractNumber(s string) int {
	matches := numRegex.FindString(s)
	if matches == "" {
		return 0
	}
	num, _ := strconv.Atoi(matches)
	return num
}

func calculateFileSHA(filePath string) string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("無法打開檔案: %v", err)
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		log.Fatalf("無法計算 SHA: %v", err)
	}

	return hex.EncodeToString(hash.Sum(nil))
}

func processZip(zipPath string) []database.ImageData {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Fatalf("無法打開 ZIP 檔案: %v", err)
	}
	defer reader.Close()

	var images []database.ImageData

	for index, file := range reader.File {
		if !file.FileInfo().IsDir() && isImageFile(file.Name) {
			images = append(images, database.ImageData{Filename: file.Name, ImgIndex: index})
		}
	}

	sortedImages := make([]database.ImageData, len(images))
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

func checkCacheExpiration(zipSHA string) bool {
	// Implement logic to check if the cache is expired based on the SHA of the ZIP file
	// Return true if cache is expired, otherwise false
	return true // Placeholder
}

func generateThumbnails(images []database.ImageData) {
	// Implement logic to generate thumbnails
}

func main() {
	zipPath := "example.zip"
	dbPath := "images.db"

	images := processZip(zipPath)

	db := database.InitDB(dbPath)
	defer db.Close()

	// Calculate SHA-256 hash of the ZIP file
	zipSHA := calculateFileSHA(zipPath)

	// Check if a ZipFile with the same SHA-256 already exists
	existingZipFileID, err := database.CheckZipFileExists(db, zipSHA)
	// if err != nil {
	// 	log.Fatalf("無法查詢 ZipFile: %v", err)
	// }
	if existingZipFileID != 0 {
		log.Printf("ZipFile with SHA-256 %s already exists with ID %d. Skipping insertion.", zipSHA, existingZipFileID)
		displayData(db, zipSHA, filepath.Base(zipPath))
		return
	}

	// Create ZipFile struct
	zipFile := database.ZipFile{
		FileName:  filepath.Base(zipPath),
		SHA:       zipSHA,
		Timestamp: time.Now().Unix(),
	}

	// 轉成通用切片存入DB
	data := []database.Storable{zipFile}
	database.StoreData(db, data)

	// Get the ID of the inserted ZipFile
	var zipFileID int
	err = db.QueryRow("SELECT id FROM zip_files WHERE filename = ? AND sha = ?", zipFile.FileName, zipFile.SHA).Scan(&zipFileID)
	if err != nil {
		log.Fatalf("無法獲取 ZipFile ID: %v", err)
	}

	// Store ImageData entries
	for _, img := range images {
		img.ZipID = zipFileID
		database.StoreData(db, []database.Storable{img})
	}

	displayData(db, zipFile.SHA, zipFile.FileName)
}

func displayData(db *sql.DB, sha string, filename string) {
	// Read and display data from the database
	rows, err := db.Query("SELECT filename, sha, timestamp FROM zip_files WHERE sha = ? AND filename = ?", sha, filename)
	if err != nil {
		log.Fatalf("無法讀取資料: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var fileName, sha string
		var timestamp int64
		if err := rows.Scan(&fileName, &sha, &timestamp); err != nil {
			log.Fatalf("無法解析資料: %v", err)
		}
		fmt.Printf("FileName: %s\n", fileName)
		fmt.Printf("SHA-256: %s\n", sha)
		fmt.Printf("Timestamp: %d\n", timestamp)
	}

	// Read and display ImageData from the database
	imgRows, err := db.Query("SELECT filename, img_index, page FROM image_data WHERE zip_id = (SELECT id FROM zip_files WHERE sha = ? AND filename = ?)", sha, filename)
	if err != nil {
		log.Fatalf("無法讀取圖片資料: %v", err)
	}
	defer imgRows.Close()

	fmt.Println("圖片索引對應表:")
	fmt.Println("-------------------------------------------------")
	fmt.Println("| Filename                         | ImgIndex | Page |")
	fmt.Println("-------------------------------------------------")
	for imgRows.Next() {
		var filename string
		var imgIndex, page int
		if err := imgRows.Scan(&filename, &imgIndex, &page); err != nil {
			log.Fatalf("無法解析圖片資料: %v", err)
		}
		fmt.Printf("| %-32s | %-8d | %-4d |\n", filename, imgIndex, page)
	}
	fmt.Println("-------------------------------------------------")
}
