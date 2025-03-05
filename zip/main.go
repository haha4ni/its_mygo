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

func listImagesInZip(zipPath string) map[int]string {
	// 打開 ZIP 檔案
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Fatalf("無法打開 ZIP 檔案: %v", err)
	}
	defer reader.Close()

	var imageFiles []string

	// 遍歷 ZIP 內的所有檔案
	for _, file := range reader.File {
		if !file.FileInfo().IsDir() && isImageFile(file.Name) {
			imageFiles = append(imageFiles, file.Name)
		}
	}

	// 進行自然排序
	sort.Slice(imageFiles, func(i, j int) bool {
		return extractNumber(imageFiles[i]) < extractNumber(imageFiles[j])
	})

	// 建立索引對應表
	imageIndexMap := make(map[int]string)
	for i, name := range imageFiles {
		imageIndexMap[i] = name
	}

	return imageIndexMap
}

func main() {
	zipPath := "example.zip" // 替換為你的 ZIP 檔案路徑
	imageIndexMap := listImagesInZip(zipPath)

	fmt.Println("圖片索引對應表:")
	for index, name := range imageIndexMap {
		fmt.Printf("%d -> %s\n", index, name)
	}

	// 測試查找某個索引位置的檔案
	targetIndex := 10 // 目標索引
	if name, exists := imageIndexMap[targetIndex]; exists {
		fmt.Printf("索引 %d 對應的檔案: %s\n", targetIndex, name)
	} else {
		fmt.Printf("索引 %d 超出範圍\n", targetIndex)
	}
}
