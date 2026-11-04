# 系統架構說明 (System Architecture)

本專案「屏東好朋友｜屏東農特產品牌官網」採用 **全靜態生成 (Full Static Site Generation, SSG)** 搭配 **漸進增強 (Progressive Enhancement)** 架構。所有內容在建置階段即完整寫入 HTML，不依賴客戶端 API 請求；在未啟用 JavaScript 的環境下仍具備完整的閱讀與導覽能力。

---

## 1. 生成流程 (Build & Generation Workflow)

核心建置工具為 `tools/convert.go`（執行 `go run tools/convert.go`）：

1. **讀取與正規化資料 (`docs/brands.xlsx`)**：
   - 使用 `excelize/v2` 讀取第一張工作表，解析品牌基本資料、標語、產地 (`origin`)、故事、Hero 巨字標語（以 `|` 分隔）、CTA、徽章（最多 3 個）、產品（最多 10 項，自動編號 `<slug>-001` 起；價格支援 `300`、`1,200`、`NT$300`、`300元` 等格式，無法解析時於終端機印出警告）、影片（最多 3 支，標題與有效 YouTube 網址皆為必填）、社群與 SEO 欄位（目前 `docs/brands.xlsx` 預設僅建立產品 1～2、影片 1 與徽章 1～2 欄位，可依命名規則自行擴充標題列）。
   - 略過空白列、`slug` 為空或以 `demo` 開頭的測試列（`slug` 支援連字號 `-`）。
   - **資源路徑正規化**：自動將檔名中的全形數字與符號（`０-９`、`＿`、`－`、`．`、全形空白）轉為半形；相對檔名自動補全為 `public/assets/<slug>/<檔名>`。若未指定 `HeroImage`，自動以第一筆已收錄產品的圖片作為備援（沒有已收錄產品時不提供備援）。
2. **輸出資料備份 (`data/brands.json`)**：
   - **保持 Excel 工作表的原始列順序**輸出為縮排 JSON，供資料備份與檢視。
3. **展示排序與模板渲染 (`html/template`)**：
   - 將記憶體中的品牌切片依**由北到南路線**重新穩定排序 (`sort.SliceStable`)。
   - 解析 `public/assets/common/pingtung-map.svg` 產生內嵌互動地圖 SVG (`townshipMap`)。
   - 挑選首頁 Hero 拼貼牆的 7 件指定代表產品 (`Tiles`)。
   - 載入 `index_template.html` 與 `brand.html` 共同編譯為單一 `html/template` 模板集，去除每行尾端多餘空白後，輸出根目錄 `index.html` 與各品牌頁 `brands/<slug>/index.html`（傳入 `Prefix: "../../"` 處理相對資源路徑，並注入循環排列的下一站品牌 `.Next`）。

---

## 2. 排序、命名推導與影片解析規則

`tools/convert.go` 定義了供模板與地圖共用的轉換函式：

### 由北到南排序 (`routeRank`)
網站展示順序（地圖對應、首頁品牌列表、品牌頁「下一站」導覽）皆固定依地理位置由北到南排列，與 Excel 原始順序脫鉤：
1. `pingtung`（屏東市｜屏東縣農會）
2. `ailiao`（內埔鄉｜隘寮社區）
3. `wandan`（萬丹鄉｜萬丹鄉農會）
4. `chengde`（萬巒鄉｜成德社區）
5. `chishan`（萬巒鄉｜赤山社區）
6. `jiadong`（佳冬鄉｜佳冬鄉農會）
7. `xinlong`（枋寮鄉｜新龍社區）
8. `fenggang`（枋山鄉｜楓港社區）

未列入上述清單的 `slug` 權重設為 `999`，排在最後並維持原相對順序。最後一站（`fenggang`）的「下一站」會循環回到第一站（`pingtung`）。

### 品牌短名推導 (`shortName`)
模板函式 `short` 用於大標題、導覽與地圖標籤：
- **社區發展協會**（名稱含 `社區發展協會`）：先移除字尾 `發展協會`，再尋找最後一個 `鄉`、`市` 或 `鎮` 字元並裁切其前綴（例如：`屏東縣內埔鄉隘寮社區發展協會` → `隘寮社區`）。
- **屏東縣農會**：因「屏東縣」為縣級農會本名而非鄉鎮前綴，特例保留全名 `屏東縣農會`。
- **其他品牌**（如鄉鎮農會）：移除開頭的 `屏東縣` 前綴（例如：`屏東縣萬丹鄉農會` → `萬丹鄉農會`）。

### 鄉鎮名推導 (`townName`)
模板函式 `town` 由 `Origin`（產地）欄位推導簡短鄉鎮名：
- 若 `Origin` 為 `屏東市`，保留 `屏東市`。
- 其他情況依序移除開頭的 `屏東縣` 與 `屏東` 前綴（例如：`屏東縣萬巒鄉` → `萬巒鄉`）。

### 故事段落與 YouTube 影片解析 (`storyParagraphs` / `videoViews`)
- **故事分段 (`paragraphs`)**：將 `Story` 依換行符 `\n` 切分，去除前後空白並忽略空行，輸出為多個 `<p>` 段落。
- **影片解析 (`videos`)**：
  - 影片標題與網址皆須填寫（標題留空時於資料讀取階段即整筆略過）。
  - 支援三類網址格式：`youtu.be/<id>`、`youtube.com/watch?v=<id>`（含 `www.`、`m.`），以及 `youtube.com/shorts/<id>` 或 `youtube.com/embed/<id>`。
  - 驗證擷取出的影片 ID 必須符合 11 碼英數與 `_-` 正規表示式 (`^[A-Za-z0-9_-]{11}$`)，不符者自動略過。
  - 若路徑為 `/shorts/<id>`，標記 `Shorts: true`（模板加上 `.is-short` 以 `9:16` 直式比例呈現）。
  - 自動濾除影片標題中的 Emoji 與變體選擇符（Unicode `U+1F000`–`U+1FAFF`、`U+FE0F`、`U+200D`）。

---

## 3. 模板分工 (`index_template.html` vs `brand.html`)

兩個模板檔案由 `template.ParseFiles("index_template.html", "brand.html")` 一併載入並共用子模板：

### `index_template.html`
- **`{{define "index"}}`（首頁主體）**：
  - **Hero 區塊 (`.home-hero`)**：左側主視覺文案與右側 `specimen-mosaic` 7 件產品拼貼牆（固定挑選 `jiadong/B_8_1`、`xinlong/B_3_1`、`chishan/B_5_1`、`fenggang/B_2_2`、`chengde/B_4_2`、`wandan/B_6_2`、`ailiao/B_1_1`，點擊可跳轉至首頁對應品牌錨點 `#brand-<slug>`）。
  - **品牌地圖與列表 (`#map`)**：左側粘性定位 (`sticky`) 內嵌屏東鄉鎮 SVG 地圖與即時狀態說明 (`.map-caption`)；右側依由北到南順序渲染 `.brand-row`（含條件式渲染的 Logo、短名、鄉鎮名、標語、產品與價格規格列表，以及進入品牌頁連結；若 `Logo` 留空則首頁與品牌頁皆不顯示標誌）。
  - **三種風土 (`#places`)**：內容與分頁結構**直接寫死於 `index_template.html`**（不由 Excel 生成），分為「大武山腳」（隘寮、成德、赤山）、「屏東平原」（萬丹、屏東縣農會）、「南方海岸」（佳冬、新龍、楓港）三個頁籤面板（新增品牌時需手動加入對應面板）。
  - **關於計畫 (`#about`)**：信件體裁的計畫介紹與署名，**直接寫死於 `index_template.html`**（連同 `<meta name="description">` 與 `.home-hero` 皆有寫死的品牌與產品數量敘述）。
- **`{{define "fonts"}}`（共用字型宣告）**：載入 Google Fonts 的 `Work Sans`（300/400/500）與 `Noto Sans TC`（300/400/500）。
- **`{{define "footer"}}`（共用全幅照片頁尾）**：透過 `{{.Prefix}}` 同時相容首頁 (`""`) 與品牌子頁 (`"../../"`)，包含頁尾導覽、大型 `屏東好朋友` 字標，以及屏東縣鄉鎮地圖的 CC BY-SA 4.0 授權出處連結。

### `brand.html`
- **`{{define "brand"}}`（品牌獨立頁）**：
  - 採用 **Split-stage（左圖右文）** 版面：左側 `.brand-visual` 在桌面版為 `100svh` 粘性全幅主圖與左下角 `HeroWords` 巨字標語；右側 `.brand-content` 為可捲動的品牌敘事區。
  - **條件式區塊輸出**：
    - 品牌標誌 (`Logo`)、產品購買連結 (`Product.Link`) 有值時才渲染。
    - **產品 (`Products`)**、**影片 (`Videos`)**、**常見問題 (`FAQ`)**、**行動呼籲按鈕 (`CTA.Link`)** 皆透過 `{{if ...}}` 判斷，僅在該品牌有對應資料時才輸出 HTML 區塊（若 `CTA.Link` 存在但 `CTA.Label` 為空，按鈕文字預設顯示「聯絡品牌」）。
    - *註：目前 `convert.go` 保留 `Brand.FAQ` 結構與模板渲染邏輯，若資料源未填寫則自動隱藏。*
  - **下一站導覽 (`.next-stop`)**：頁面底部顯示由北到南路線的下一個鄉鎮與品牌短名，引導使用者連續探索全縣品牌。

---

## 4. 互動式屏東鄉鎮地圖機制

1. **原始底圖與授權**：
   - 原始未修改的 SVG 存放於 `public/assets/common/pingtung-map.svg`（源自 Wikimedia Commons Ythlev 製作之 2018 屏東縣長選舉地圖，採用 CC BY-SA 4.0 授權，詳見 `docs/credits.md`）。
2. **建置期 SVG 淨化與內嵌 (`townshipMap`)**：
   - `convert.go` 使用 `encoding/xml` 串流解析 `<path>` 元素與其 `<title>` 的第一個詞作為鄉鎮名稱 (`town`)。
   - **移除選舉配色與原描邊**：只讀取幾何路徑 `d` 與 `<title>`，原有的填充色與描邊屬性完全不進入頁面；同時將路徑座標數值統一四捨五入至小數點後一位 (`%.1f`) 以縮減 HTML 體積。
   - **產地自動對應**：比對所有品牌的 `Origin` 是否包含該鄉鎮名稱 (`strings.Contains(b.Origin, town)`)。同一鄉鎮可對應多個品牌（如萬巒鄉同時包含 `chengde` 與 `chishan`）。
   - **語意化包裝與標籤定位**：
     - 有品牌的鄉鎮會以 `<a class="town has-brand" href="#brand-<首個slug>" data-town="<鄉鎮>" data-slugs="<以空白分隔的slug清單>" aria-label="<鄉鎮>：<品牌短名頓號清單>">` 包覆 `<path>`。
     - 在 `<g class="map-labels" aria-hidden="true">` 中為有品牌的鄉鎮輸出 `<text class="town-label">` 鄉鎮文字標籤：新鄉鎮預設會由路徑座標自動計算中心點置中，僅在需要微調視覺位置時才透過 `labelPositions`（座標系對應 `viewBox="0 0 1440 2777"`）指定覆寫座標。

---

## 5. 前端互動與視覺設計 (`style.css` & `app.js`)

前端不依賴任何第三方 JS 框架，所有互動皆以無障礙與漸進增強為前提實作（腳本執行後會在 `<html>` 加上 `.js` class）：

1. **地圖與品牌列表雙向同步**：
   - **捲動偵測**：透過 `IntersectionObserver`（上下內縮 20% 視窗高度：`rootMargin: '-20% 0px -20% 0px'`）搭配 `requestAnimationFrame` 捲動監聽，計算目前最接近視窗垂直中央的品牌列 (`.brand-row`)，自動將對應鄉鎮 `<a class="town has-brand">` 與 `<text class="town-label">` 設為 `.is-active`（並設定 `aria-current="location"`），同時更新地圖下方的 `.map-caption`（顯示「品牌短名 · 鄉鎮名」）。
   - **點擊跳轉與鍵盤接續**：點擊地圖上的品牌鄉鎮時，平滑捲動（若開啟減少動態則瞬間跳轉）至該鄉鎮第一個品牌列，更新網址 Hash，並將焦點 (`focus`) 移至目標品牌 `<article>`，方便鍵盤使用者直接繼續往下閱讀。
   - **SVG 標籤字級自適應**：使用 `ResizeObserver` 監聽地圖 SVG 的 `getScreenCTM().a` 縮放比例，動態計算 `--map-label-size`（`13.5 / scale` px），確保桌面與行動裝置上的鄉鎮文字保持清晰一致的視覺大小。
2. **「三種風土」無障礙頁籤 (`role="tablist"`)**：
   - 無 JS 時三個風土面板垂直展開皆可閱讀；啟用 JS 後（配合 `.js .places-section` Grid 疊合版面）轉為互動頁籤，支援滑鼠點擊與鍵盤 `ArrowLeft`、`ArrowRight`、`Home`、`End` 切換焦點與 `aria-selected` / `hidden` 狀態。
3. **品牌頁點擊載入 YouTube (`[data-video]`)**：
   - 初始僅載入輕量的 YouTube 預覽縮圖 (`i.ytimg.com/vi/<ID>/hqdefault.jpg`) 與指向原影片網址的「播放影片」連結（無 JS 時可直接另開新分頁觀看）。
   - 點擊「播放影片」時攔截預設行為，動態替換為 `https://www.youtube-nocookie.com/embed/<ID>?autoplay=1` 的 `<iframe>` 並移交焦點，兼顧隱私與頁面初始載入效能。
4. **動態揭示與減敏支援 (`prefers-reduced-motion`)**：
   - 當使用者未要求減少動態時，透過 `IntersectionObserver` 為 `[data-reveal]` 元素加上交錯延遲淡入上移效果，並為 `[data-shutter]` 圖片容器加上橫向百葉揭幕 (`clip-path`) 動畫。
   - 若系統設定 `prefers-reduced-motion: reduce`，自動停用捲動平滑過渡與揭示動畫。
5. **視覺系統 (`style.css`)**：
   - 採用溫暖石色與植物綠調色盤（`--bg: #eeecea`、`--surface: #e4e1dd`、`--ink: #333231`、`--accent: #5c6b54`、`--map-brand: #b5a793`）。
   - 產品圖與標本牆使用 `mix-blend-mode: multiply` 融入石色底圖。

---

## 6. 目錄結構 (Folder Structure)

* `docs/brands.xlsx`：品牌與產品主資料源（Excel）。
* `docs/architecture.md`：本系統架構說明文件。
* `docs/brands.md`：Excel 欄位定義與素材命名指南。
* `docs/credits.md`：屏東縣鄉鎮地圖（CC BY-SA 4.0）與素材授權紀錄。
* `tools/convert.go`：建置與靜態網頁生成工具（Go + `excelize/v2` + `html/template`）。
* `data/brands.json`：建置產生的結構化資料備份（維持 Excel 原始順序）。
* `index_template.html`：首頁模板（含寫死的「三種風土」、「關於計畫」及共用 `fonts` / `footer` 子模板）。
* `brand.html`：品牌頁模板（含條件式區塊與「下一站」導覽）。
* `index.html`：生成後的靜態首頁。
* `brands/<slug>/index.html`：生成後的各品牌獨立靜態頁面。
* `public/assets/common/pingtung-map.svg`：保留未修改的屏東縣鄉鎮原始 SVG 底圖。
* `public/assets/<slug>/`：各品牌圖片與 Logo 靜態資源。
* `style.css`：全站統一樣式表（Split-stage 版面、響應式斷點 `1250px` / `1000px` / `740px`、無障礙與動態樣式）。
* `app.js`：前端漸進增強腳本（地圖同步、風土頁籤、點擊載入影片、捲動揭示動畫）。

---

## 7. 開發預覽與檢查方式 (Local Development & Verification)

### 重新生成與檢查
修改 `docs/brands.xlsx`、`index_template.html`、`brand.html` 或 `tools/convert.go` 後，請於專案根目錄執行：

```bash
# 1. 重新編譯並產生 data/brands.json、index.html 與 brands/*/index.html
go run tools/convert.go

# 2. 執行 Go 靜態分析與檢查格式（gofmt -l 應無輸出）
go vet ./tools/...
gofmt -l tools/convert.go

# 3. 檢視生成結果與變更差異
git status
git diff
```

### 本地預覽
由於全站為無 API 依賴的純靜態 HTML 與相對路徑資源，可直接以瀏覽器開啟 `index.html`，或啟動本地 HTTP 伺服器預覽：

```bash
python3 -m http.server 8000
# 開啟 http://localhost:8000
```
