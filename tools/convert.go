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
	Logo      string    `json:"logo"`
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

	if len(rows) < 1 {
		log.Fatal("工作表為空")
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
		slug := getVal(row, headers, "slug")
		if len(row) == 0 || slug == "" || strings.HasPrefix(slug, "demo") {
			continue
		}

		// Helper to normalize filenames (convert full-width to half-width)
		normalizeFilename := func(s string) string {
			r := strings.NewReplacer(
				"０", "0", "１", "1", "２", "2", "３", "3", "４", "4",
				"５", "5", "６", "6", "７", "7", "８", "8", "９", "9",
				"＿", "_", "－", "-", "．", ".", "　", " ",
			)
			return r.Replace(s)
		}

		// Helper to resolve asset paths
		resolveAssetPath := func(raw string) string {
			raw = strings.TrimSpace(raw)
			if raw == "" {
				return ""
			}

			// Normalize full-width characters
			raw = normalizeFilename(raw)

			// If it's already a full URL or starts from root, keep it as is
			if strings.HasPrefix(raw, "http") || strings.HasPrefix(raw, "//") || strings.HasPrefix(raw, "/") || strings.HasPrefix(raw, "data:") {
				return raw
			}

			// Remove any legacy relative markers first
			cleaned := raw
			for strings.HasPrefix(cleaned, "../") {
				cleaned = strings.TrimPrefix(cleaned, "../")
			}

			// If it ALREADY contains public/assets (even after cleaning), return it directly
			if strings.HasPrefix(cleaned, "public/assets") {
				return cleaned
			}

			// Prepend the standard path
			return fmt.Sprintf("public/assets/%s/%s", slug, cleaned)
		}

		brand := Brand{
			Slug:      slug,
			Name:      getVal(row, headers, "品牌名稱"),
			Logo:      resolveAssetPath(getVal(row, headers, "logo")),
			Tagline:   getVal(row, headers, "標語"),
			Origin:    getVal(row, headers, "產地"),
			Story:     getVal(row, headers, "品牌故事"),
			HeroImage: resolveAssetPath(getVal(row, headers, "hero image", "主圖圖片url", "主圖 url")),
			QRSlug:    getVal(row, headers, "qrslug", "qr slug"),
		}

		// Hero Words
		words := getVal(row, headers, "hero 巨字標語", "hero標語", "hero巨字標語（以 | 分隔）")
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
			imgUrl := resolveAssetPath(getVal(row, headers, fmt.Sprintf("產品 %d 圖片url", k), fmt.Sprintf("產品%d圖片url", k), fmt.Sprintf("產品 %d 圖片", k), fmt.Sprintf("產品%d圖片", k)))
			brand.Products = append(brand.Products, Product{
				ID:      fmt.Sprintf("%s-%03d", brand.Slug, k),
				Name:    pName,
				Summary: getVal(row, headers, fmt.Sprintf("產品 %d 摘要", k), fmt.Sprintf("產品%d摘要", k)),
				Price:   price,
				Spec:    getVal(row, headers, fmt.Sprintf("產品 %d 規格", k), fmt.Sprintf("產品%d規格", k)),
				Image:   imgUrl,
				Link:    getVal(row, headers, fmt.Sprintf("產品 %d 連結", k), fmt.Sprintf("產品%d連結", k)),
			})
		}

		// Fallback HeroImage if empty
		if brand.HeroImage == "" && len(brand.Products) > 0 {
			brand.HeroImage = brand.Products[0].Image
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

	// Generate Directories for each brand
	fmt.Println("正在生成品牌專屬資料夾...")
	templatePath := "brand.html"
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		log.Fatalf("無法讀取模板檔 %s: %v", templatePath, err)
	}

	for _, brand := range brands {
		if brand.Slug == "" {
			continue
		}
		brandDir := fmt.Sprintf("brands/%s", brand.Slug)
		err := os.MkdirAll(brandDir, 0755)
		if err != nil {
			log.Printf("無法建立資料夾 %s: %v", brandDir, err)
			continue
		}

		// Adjust paths in the template to work from the subdirectory
		html := string(templateContent)
		html = strings.ReplaceAll(html, `href="style.css"`, `href="../../style.css"`)
		html = strings.ReplaceAll(html, `src="app.js"`, `src="../../app.js"`)
		html = strings.ReplaceAll(html, `href="index.html"`, `href="../../index.html"`)

		// Basic SSG Injection
		html = strings.ReplaceAll(html, `<h1 class="hero-brand"></h1>`, fmt.Sprintf(`<h1 class="hero-brand">%s</h1>`, brand.Name))
		html = strings.ReplaceAll(html, `<p class="hero-tagline"></p>`, fmt.Sprintf(`<p class="hero-tagline">%s</p>`, brand.Tagline))
		html = strings.ReplaceAll(html, `<p class="hero-story"></p>`, fmt.Sprintf(`<p class="hero-story">%s</p>`, brand.Story))
		html = strings.ReplaceAll(html, `<p class="story-text"></p>`, fmt.Sprintf(`<p class="story-text">%s</p>`, brand.Story))
		html = strings.ReplaceAll(html, `<title>品牌頁｜示範</title>`, fmt.Sprintf(`<title>%s</title>`, brand.Name))

		// Hero Background SSG Injection
		if brand.HeroImage != "" {
			heroImgURL := "../../" + brand.HeroImage
			html = strings.ReplaceAll(html, `<div class="hero-bg" aria-hidden="true"></div>`, fmt.Sprintf(`<div class="hero-bg" aria-hidden="true" style="--hero-bg: url('%s');"></div>`, heroImgURL))
		}

		// Logo SSG Injection
		if brand.Logo != "" {
			logoURL := "../../" + brand.Logo
			html = strings.ReplaceAll(html, `id="brand-logo" src="" alt="" style="display:none;`, fmt.Sprintf(`id="brand-logo" src="%s" alt="%s" style="display:block;`, logoURL, brand.Name))
			html = strings.ReplaceAll(html, `id="brand-name-text">Pingtung Friends</span>`, fmt.Sprintf(`id="brand-name-text" style="display:none;">%s</span>`, brand.Name))
		} else if brand.Name != "" {
			html = strings.ReplaceAll(html, `id="brand-name-text">Pingtung Friends</span>`, fmt.Sprintf(`id="brand-name-text">%s</span>`, brand.Name))
		}

		// HeroWords SSG Injection
		if len(brand.HeroWords) > 0 {
			var wordsHTML strings.Builder
			for _, w := range brand.HeroWords {
				wordsHTML.WriteString(fmt.Sprintf("<span>%s</span>", w))
			}
			html = strings.ReplaceAll(html, `<div class="hero-big-text"></div>`, fmt.Sprintf(`<div class="hero-big-text">%s</div>`, wordsHTML.String()))
		}

		// Products SSG Injection
		if len(brand.Products) > 0 {
			var prodHTML strings.Builder
			for _, p := range brand.Products {
				img := "../../" + p.Image
				prodHTML.WriteString(fmt.Sprintf(`
					<div class="card">
						<img src="%s" alt="%s" loading="lazy" />
						<div class="card-body">
							<h3>%s</h3>
							<p class="summary">%s</p>
							<p class="price">$%d · %s</p>
							<a class="btn-primary" href="%s" target="_blank" rel="noopener">立即購買</a>
						</div>
					</div>`, img, p.Name, p.Name, p.Summary, p.Price, p.Spec, p.Link))
			}
			html = strings.ReplaceAll(html, `<div class="product-grid"></div>`, fmt.Sprintf(`<div class="product-grid">%s</div>`, prodHTML.String()))
		}

		// Videos SSG Injection
		if len(brand.Videos) > 0 {
			var videoHTML strings.Builder
			for _, v := range brand.Videos {
				videoHTML.WriteString(fmt.Sprintf(`
					<div class="video-item">
						<div>%s</div>
						<a href="%s" target="_blank" rel="noopener">立即播放</a>
					</div>`, v.Title, v.URL))
			}
			html = strings.ReplaceAll(html, `<div class="video-list"></div>`, fmt.Sprintf(`<div class="video-list">%s</div>`, videoHTML.String()))
		}

		indexPath := fmt.Sprintf("%s/index.html", brandDir)
		err = os.WriteFile(indexPath, []byte(html), 0644)
		if err != nil {
			log.Printf("無法寫入品牌頁 %s: %v", indexPath, err)
		} else {
			fmt.Printf("  -> 已生成: %s\n", indexPath)
		}
	}

	fmt.Printf("成功轉換 %d 個品牌資料至 %s 並生成專屬頁面\n", len(brands), outputFile)
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
