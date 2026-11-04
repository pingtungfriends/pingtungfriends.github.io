package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"sort"
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
			path := fmt.Sprintf("public/assets/%s/%s", slug, cleaned)
			// Ensure forward slashes for web usage, even on Windows
			return strings.ReplaceAll(path, "\\", "/")
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

	// Keep the source JSON in workbook order; only presentation follows the north–south route.
	sort.SliceStable(brands, func(i, j int) bool { return routeRank(brands[i].Slug) < routeRank(brands[j].Slug) })
	mapSVG, err := townshipMap(brands)
	if err != nil {
		log.Fatal(err)
	}
	funcs := template.FuncMap{
		"splitProduct": func(id string) []string { return strings.Split(id, "-") },
		"short":        shortName, "town": townName, "asset": assetPath,
		"videos": videoViews, "paragraphs": storyParagraphs,
	}
	pages, err := template.New("site").Funcs(funcs).ParseFiles("index_template.html", "brand.html")
	if err != nil {
		log.Fatal(err)
	}
	tiles := []Product{}
	for _, choice := range []struct{ slug, file string }{
		{"jiadong", "B_8_1"}, {"xinlong", "B_3_1"}, {"chishan", "B_5_1"},
		{"fenggang", "B_2_2"}, {"chengde", "B_4_2"}, {"wandan", "B_6_2"}, {"ailiao", "B_1_1"},
	} {
		for _, b := range brands {
			if b.Slug == choice.slug {
				for _, p := range b.Products {
					if strings.Contains(p.Image, choice.file) {
						tiles = append(tiles, p)
					}
				}
			}
		}
	}
	writePage := func(path, name string, data PageData) {
		var rendered strings.Builder
		if err := pages.ExecuteTemplate(&rendered, name, data); err != nil {
			log.Fatal(err)
		}
		lines := strings.Split(rendered.String(), "\n")
		for i := range lines {
			lines[i] = strings.TrimRight(lines[i], " \t")
		}
		if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644); err != nil {
			log.Fatal(err)
		}
		fmt.Println("已生成:", path)
	}

	writePage("index.html", "index", PageData{Brands: brands, Map: mapSVG, Tiles: tiles})
	for i, b := range brands {
		dir := "brands/" + b.Slug
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal(err)
		}
		writePage(dir+"/index.html", "brand", PageData{Brand: b, Next: brands[(i+1)%len(brands)], Prefix: "../../"})
	}
	fmt.Printf("成功轉換 %d 個品牌；JSON 保持原始工作表順序。\n", len(brands))
}

type PageData struct {
	Brands      []Brand
	Brand, Next Brand
	Map         template.HTML
	Tiles       []Product
	Prefix      string
}

func routeRank(slug string) int {
	for i, s := range []string{"pingtung", "ailiao", "wandan", "chengde", "chishan", "jiadong", "xinlong", "fenggang"} {
		if slug == s {
			return i
		}
	}
	return 999
}

func shortName(b Brand) string {
	if strings.Contains(b.Name, "社區發展協會") {
		name := strings.TrimSuffix(b.Name, "發展協會")
		if i := strings.LastIndexAny(name, "鄉市鎮"); i >= 0 {
			name = name[i+len("鄉"):]
		}
		return name
	}
	// 屏東縣農會 is the county's name, rather than a township prefix.
	if b.Name == "屏東縣農會" {
		return b.Name
	}
	return strings.TrimPrefix(b.Name, "屏東縣")
}

func townName(b Brand) string {
	if b.Origin == "屏東市" {
		return b.Origin
	}
	return strings.TrimPrefix(strings.TrimPrefix(b.Origin, "屏東縣"), "屏東")
}

func assetPath(prefix, path string) string {
	if path == "" || strings.HasPrefix(path, "/") || strings.HasPrefix(path, "http") || strings.HasPrefix(path, "data:") {
		return path
	}
	return prefix + path
}

func storyParagraphs(story string) []string {
	var result []string
	for _, p := range strings.Split(story, "\n") {
		if p = strings.TrimSpace(p); p != "" {
			result = append(result, p)
		}
	}
	return result
}

type VideoView struct {
	ID, Title, URL string
	Shorts         bool
}

var youtubeID = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

func videoViews(videos []Video) []VideoView {
	var result []VideoView
	for _, v := range videos {
		u, err := url.Parse(v.URL)
		if err != nil {
			continue
		}
		id := ""
		shorts := false
		switch strings.ToLower(u.Hostname()) {
		case "youtu.be":
			id = strings.Trim(u.Path, "/")
		case "youtube.com", "www.youtube.com", "m.youtube.com":
			parts := strings.Split(strings.Trim(u.Path, "/"), "/")
			if len(parts) == 2 && (parts[0] == "shorts" || parts[0] == "embed") {
				id = parts[1]
				shorts = parts[0] == "shorts"
			} else {
				id = u.Query().Get("v")
			}
		}
		if !youtubeID.MatchString(id) {
			continue
		}
		title := strings.Map(func(r rune) rune {
			if r >= 0x1F000 && r <= 0x1FAFF || r == 0xFE0F || r == 0x200D {
				return -1
			}
			return r
		}, v.Title)
		result = append(result, VideoView{ID: id, Title: title, URL: v.URL, Shorts: shorts})
	}
	return result
}

// Read only path geometry and titles: election colours and original strokes never enter the page.
func townshipMap(brands []Brand) (template.HTML, error) {
	file, err := os.Open("public/assets/common/pingtung-map.svg")
	if err != nil {
		return "", err
	}
	defer file.Close()
	dec := xml.NewDecoder(file)
	number := regexp.MustCompile(`-?\d+(?:\.\d+)?`)
	var svg, labels strings.Builder
	// Label centres are spaced for the narrowest map as well as desktop.
	labelPositions := map[string][2]int{
		"屏東市": {335, 605}, "內埔鄉": {655, 650}, "萬丹鄉": {325, 900},
		"萬巒鄉": {670, 950}, "佳冬鄉": {442, 1308}, "枋寮鄉": {694, 1478}, "枋山鄉": {894, 1917},
	}
	svg.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1440 2777" role="group" aria-label="屏東縣鄉鎮品牌地圖"><title>屏東縣鄉鎮品牌地圖</title>`)
	for {
		token, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "path" {
			continue
		}
		var path struct {
			D     string `xml:"d,attr"`
			Title string `xml:"title"`
		}
		if err := dec.DecodeElement(&path, &start); err != nil {
			return "", err
		}
		fields := strings.Fields(path.Title)
		if len(fields) == 0 {
			return "", fmt.Errorf("地圖 path 缺少鄉鎮名稱")
		}
		town := fields[0]
		d := number.ReplaceAllStringFunc(path.D, func(n string) string { f, _ := strconv.ParseFloat(n, 64); return strconv.FormatFloat(f, 'f', 1, 64) })
		var slugs, names []string
		for _, b := range brands {
			if strings.Contains(b.Origin, town) {
				slugs = append(slugs, b.Slug)
				names = append(names, shortName(b))
			}
		}
		if len(slugs) > 0 {
			fmt.Fprintf(&svg, `<a class="town has-brand" href="#brand-%s" data-town="%s" data-slugs="%s" aria-label="%s：%s">`, template.HTMLEscapeString(slugs[0]), template.HTMLEscapeString(town), strings.Join(slugs, " "), template.HTMLEscapeString(town), template.HTMLEscapeString(strings.Join(names, "、")))
		}
		fmt.Fprintf(&svg, `<path class="town-shape" data-town="%s" d="%s"><title>%s</title></path>`, template.HTMLEscapeString(town), d, template.HTMLEscapeString(town))
		if len(slugs) > 0 {
			svg.WriteString(`</a>`)
			if point, ok := labelPositions[town]; ok {
				fmt.Fprintf(&labels, `<text class="town-label" data-town="%s" data-slugs="%s" x="%d" y="%d" text-anchor="middle" dominant-baseline="central">%s</text>`, template.HTMLEscapeString(town), strings.Join(slugs, " "), point[0], point[1], template.HTMLEscapeString(town))
			}
		}
	}
	svg.WriteString(`<g class="map-labels" aria-hidden="true">` + labels.String() + `</g></svg>`)
	return template.HTML(svg.String()), nil
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
