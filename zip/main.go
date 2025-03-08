package main

import (
	"archive/zip"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

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

func processZip(zipPath string) []database.ImageEntry {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Fatalf("無法打開 ZIP 檔案: %v", err)
	}
	defer reader.Close()

	var images []database.ImageEntry

	for page, file := range reader.File {
		if !file.FileInfo().IsDir() && isImageFile(file.Name) {
			images = append(images, database.ImageEntry{Filename: file.Name, ImgIndex: page})
		}
	}

	sortedImages := make([]database.ImageEntry, len(images))
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

func main() {
	zipPath := "example.zip"
	dbPath := "images.db"

	images := processZip(zipPath)

	db := database.InitDB(dbPath)
	defer db.Close()

	// 轉換成 Storable 介面類型
	data := make([]database.Storable, len(images))
	for i, img := range images {
		data[i] = img
	}

	database.StoreData(db, data)

	fmt.Println("圖片索引對應表:")
	for _, img := range images {
		fmt.Printf("Filename: %s, ImgIndex: %d, Page: %d\n", img.Filename, img.ImgIndex, img.Page)
	}
}
