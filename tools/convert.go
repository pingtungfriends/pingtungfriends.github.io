package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Brand struct {
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	Tagline   string    `json:"tagline"`
	HeroWords []string  `json:"heroWords"`
	Origin    string    `json:"origin"`
	Story     string    `json:"story"`
	HeroImage string    `json:"heroImage"`
	CTA       CTA       `json:"cta"`
	Badges    []string  `json:"badges"`
	Products  []Product `json:"products"`
	Videos    []Video   `json:"videos"`
	FAQ       []FAQ     `json:"faq"`
	Social    Social    `json:"social"`
	QRSlug    string    `json:"qrSlug"`
	SEO       SEO       `json:"seo"`
}

type CTA struct {
	Label string `json:"label"`
	Link  string `json:"link"`
}

type Product struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Summary string `json:"summary"`
	Price   int    `json:"price"`
	Spec    string `json:"spec"`
	Image   string `json:"image"`
	Link    string `json:"link"`
}

type Video struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type FAQ struct {
	Q string `json:"q"`
	A string `json:"a"`
}

type Social struct {
	Line     string `json:"line"`
	Facebook string `json:"facebook"`
}

type SEO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	xlsxFile := "docs/brands.xlsx"
	outputFile := "data/brands.json"

	f, err := excelize.OpenFile(xlsxFile)
	if err != nil {
		log.Fatalf("無法讀取 XLSX 檔案: %v", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		log.Fatal("XLSX 檔案中沒有工作表")
	}
	sheetName := sheets[0]

	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("無法讀取工作表內容: %v", err)
	}

	if len(rows) < 2 {
		log.Fatal("工作表中沒有足夠的資料列（需包含標題列與至少一列資料）")
	}

	headers := make(map[string]int)
	fmt.Println("正在讀取標題列...")
	for i, cell := range rows[0] {
		name := strings.ToLower(strings.TrimSpace(cell))
		if name != "" {
			fmt.Printf("找到標題: [%s] (index: %d)\n", name, i)
			headers[name] = i
		}
	}

	var brands []Brand

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 || getVal(row, headers, "slug") == "" {
			continue
		}

		brand := Brand{
			Slug:      getVal(row, headers, "slug"),
			Name:      getVal(row, headers, "品牌名稱"),
			Tagline:   getVal(row, headers, "標語"),
			Origin:    getVal(row, headers, "產地"),
			Story:     getVal(row, headers, "品牌故事"),
			HeroImage: getVal(row, headers, "主圖圖片url", "主圖 url"),
			QRSlug:    getVal(row, headers, "qrslug", "qr slug"),
		}

		// Hero Words
		words := getVal(row, headers, "hero 巨字標語", "hero標語")
		if words != "" {
			brand.HeroWords = strings.Split(words, "|")
		}

		// CTA
		brand.CTA = CTA{
			Label: getVal(row, headers, "cta 文字", "cta文字"),
			Link:  getVal(row, headers, "cta 連結", "cta連結"),
		}

		// Badges
		for k := 1; k <= 3; k++ {
			b := getVal(row, headers, fmt.Sprintf("徽章 %d", k), fmt.Sprintf("徽章%d", k))
			if b != "" {
				brand.Badges = append(brand.Badges, b)
			}
		}

		// Products
		for k := 1; k <= 10; k++ {
			pName := getVal(row, headers, fmt.Sprintf("產品 %d 名稱", k), fmt.Sprintf("產品%d名稱", k))
			if pName == "" {
				continue
			}
			priceStr := getVal(row, headers, fmt.Sprintf("產品 %d 價格", k), fmt.Sprintf("產品%d價格", k))
			price, _ := strconv.Atoi(priceStr)
			brand.Products = append(brand.Products, Product{
				ID:      fmt.Sprintf("%s-%03d", brand.Slug, k),
				Name:    pName,
				Summary: getVal(row, headers, fmt.Sprintf("產品 %d 摘要", k), fmt.Sprintf("產品%d摘要", k)),
				Price:   price,
				Spec:    getVal(row, headers, fmt.Sprintf("產品 %d 規格", k), fmt.Sprintf("產品%d規格", k)),
				Image:   getVal(row, headers, fmt.Sprintf("產品 %d 圖片url", k), fmt.Sprintf("產品%d圖片url", k), fmt.Sprintf("產品 %d 圖片", k), fmt.Sprintf("產品%d圖片", k)),
				Link:    getVal(row, headers, fmt.Sprintf("產品 %d 連結", k), fmt.Sprintf("產品%d連結", k)),
			})
		}

		// Videos
		for k := 1; k <= 3; k++ {
			vTitle := getVal(row, headers, fmt.Sprintf("影片 %d 標題", k), fmt.Sprintf("影片%d標題", k))
			if vTitle != "" {
				brand.Videos = append(brand.Videos, Video{
					Title: vTitle,
					URL:   getVal(row, headers, fmt.Sprintf("影片 %d 網址", k), fmt.Sprintf("影片%d網址", k)),
				})
			}
		}

		// Social
		brand.Social = Social{
			Line:     getVal(row, headers, "line 連結", "line連結"),
			Facebook: getVal(row, headers, "fb 連結", "fb連結", "facebook連結"),
		}

		// SEO
		brand.SEO = SEO{
			Title:       getVal(row, headers, "seo 標題", "seo標題"),
			Description: getVal(row, headers, "seo 描述", "seo描述"),
		}

		brands = append(brands, brand)
	}

	result := struct {
		Brands []Brand `json:"brands"`
	}{
		Brands: brands,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("無法轉換為 JSON: %v", err)
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("無法寫入檔案: %v", err)
	}

	fmt.Printf("成功轉換 %d 個品牌資料至 %s\n", len(brands), outputFile)
}

func getVal(row []string, headers map[string]int, names ...string) string {
	for _, name := range names {
		idx, ok := headers[strings.ToLower(name)]
		if !ok || idx >= len(row) {
			continue
		}
		val := row[idx]
		if trimmed := strings.TrimSpace(val); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
