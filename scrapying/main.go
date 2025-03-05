package main

import (
	"fmt"
	"log"
	"net/url"
    "strings"
    "time"
    "math/rand"

	"github.com/gocolly/colly/v2"
)

// 隨機延遲 1~3 秒，模擬人類行為
func randomDelay() {
	delay := time.Duration(rand.Intn(2)+1) * time.Second
	log.Println("隨機延遲:", delay)
	time.Sleep(delay)
}

func FindBookURL(bookName string) (string, error) {
	c := colly.NewCollector()

	// 設定 User-Agent 和 Referer
	c.OnRequest(func(r *colly.Request) {
		log.Println("正在訪問:", r.URL.String()) // 紀錄訪問的網址
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		r.Headers.Set("Referer", "https://www.bookwalker.com.tw/")
	})

	// 轉換搜尋字串為 URL 格式
	query := url.QueryEscape(bookName)
	searchURL := fmt.Sprintf("https://www.bookwalker.com.tw/search?w=%s&series_display=1", query)
	log.Println("搜尋 URL:", searchURL)

	var bookURL string

	// 解析搜尋結果列表，尋找第一本書的超連結
	c.OnHTML(".bwbookitem a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		// title := e.Text
		// log.Println("找到鏈接:", href, "標題:", title)
        log.Println("找到鏈接:", href)

		if bookURL == "" { // 只抓取第一本書的網址
			bookURL = href
			log.Println("選擇的書籍網址:", bookURL)
		}
	})

	randomDelay() // 訪問前隨機延遲

	// 開始爬取
	err := c.Visit(searchURL)
	if err != nil {
		log.Println("Error visiting page:", err)
		return "", err
	}

	// 檢查是否成功取得書籍網址
	if bookURL == "" {
		log.Println("未找到符合的書籍")
		return "", fmt.Errorf("未找到書籍: %s", bookName)
	}

	// 返回書籍的完整網址
	finalURL := "https://www.bookwalker.com.tw" + bookURL
	log.Println("最終書籍網址:", finalURL)
	return finalURL, nil
}

func FindBookDetails(seriesURL string, targetVolume string) (string, error) {
	c := colly.NewCollector()

	// 設定 User-Agent 和 Referer
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		r.Headers.Set("Referer", "https://www.bookwalker.com.tw/")
	})

	var bookURL string

	// 抓取該系列頁面上的所有書籍
    c.OnHTML(".listbox_bwmain2 a", func(e *colly.HTMLElement) {
        bookTitle := strings.TrimSpace(e.DOM.Find("h4.bookname").Text()) // 抓取書名
        href := e.Attr("href")                                         // 抓取超連結
    
        log.Println("找到鏈接:", href, "標題:", bookTitle)
    
        // 檢查書名是否包含目標卷數 (targetVolume)，例如 "(6)"
        if strings.Contains(bookTitle, "("+targetVolume+")") {
            bookURL = href
            log.Println("找到符合的書籍:", bookTitle, "網址:", bookURL)
        }
    })

    randomDelay() // 訪問前隨機延遲
	// 開始抓取該系列的頁面
	err := c.Visit(seriesURL)
	if err != nil {
		log.Println("Error visiting page:", err)
		return "", err
	}

	if bookURL == "" {
		return "", fmt.Errorf("未找到符合卷數 (%s) 的書籍", targetVolume)
	}

	// 返回完整書籍詳細頁面 URL
	return "https://www.bookwalker.com.tw" + bookURL, nil
}

func main() {
	seriesURL, err := FindBookURL("結緣甘神神社")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("#找到的書籍網址:", seriesURL)

	// 查詢該系列的指定書籍
	url, err := FindBookDetails(seriesURL, "3")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("找到的書籍詳細頁面網址:", url)
}
