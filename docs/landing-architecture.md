# 多品牌一頁式網站架構方案

## 目標
- 快速為多個在地農產品牌生成自然風格的 Landing Page，支援 QRCode 導流與 CTA 轉換。
- 可重複套用的模板，新增品牌只需填資料與素材，不需改程式。
- 內建成效追蹤（掃碼、CTA 點擊、影片播放）與周報輸出。

## 頁面模板骨架（依序）
- Hero：品牌名稱、標語、關鍵賣點、主 CTA（立即購買/聯絡）、自然風景/田野背景。
- 品牌故事：來源地、產地特色、加工亮點、認證徽章。
- 產品精選：至少 10 款，卡片含縮圖、短賣點、價格/規格、CTA（購買/加入清單）。
- 產地/合作社區：地圖或定位，強調在地性；可放數據亮點或里程碑。
- 影片區：至少 3 支（品牌介紹/機器操作/產品故事）內嵌，播放事件送出追蹤。
- 評語/媒體露出：社群或媒體摘錄。
- FAQ：常見疑問、配送/保存方式。
- 底部 CTA：固定浮動按鈕（行動優先），次 CTA 為追蹤社群或加入 Line。

## 資料模型（範例 JSON/YAML）
```json
{
  "slug": "green-mango",
  "brandName": "綠意芒果工坊",
  "tagline": "屏東日照，成就香甜芒果乾",
  "heroWords": ["FLOW", "WITH", "IT"],
  "origin": "屏東內埔",
  "story": "以低溫烘焙保留果香，與在地小農共製。",
  "heroImage": "/assets/green-mango/hero.jpg",
  "gallery": ["/assets/green-mango/hero2.jpg", "/assets/green-mango/field.jpg"],
  "badges": ["產銷履歷", "無添加糖"],
  "products": [
    {
      "id": "gm-001",
      "name": "日曬芒果乾",
      "summary": "果肉厚實、低溫烘焙保留纖維",
      "price": 220,
      "spec": "120g/包",
      "image": "/assets/green-mango/prod-dry.jpg",
      "ctaLabel": "立即購買",
      "ctaLink": "https://shop.example.com/gm-001"
    }
  ],
  "videos": [
    {"title": "品牌故事", "embedUrl": "https://www.youtube.com/embed/..."}
  ],
  "contact": {"ctaLabel": "批發/合作", "ctaLink": "https://line.me/..."},
  "qrSlug": "gm",
  "seo": {
    "title": "綠意芒果工坊｜屏東在地芒果乾",
    "description": "日照香甜、低溫烘焙，與在地小農共製的芒果乾。"
  }
}
```

## 資料蒐集與轉換（XLSX 表單）
- 範例檔：`docs/brands.xlsx`（sheet 名稱 `brands`），欄位含 slug、品牌名稱、標語、Hero 巨字標語（例：FLOW|WITH|IT）、產地、品牌故事、主圖 URL、CTA 文字/連結、產品 1/2（名稱/摘要/價格/規格/圖片/連結）、影片標題/網址、徽章 1/2、SEO 標題/描述。
- 交付流程：你填好 XLSX 後，使用 Go 轉換程式讀取試算表，輸出 `data/brands.json`，並同步下載圖片到 `public/assets/<slug>/`。
- 轉換工具可以做欄位驗證（必填 CTA、產品至少 10 項等）並自動補上 UTM。

## 檔案與模板結構 (Vanilla JS + HTML)
- `data/brands.json`：匯總所有品牌資料的 JSON 檔案，由 Go 工具生成。
- `public/assets/<slug>/`：品牌圖片、影片縮圖、QR 圖檔。
- `index.html`：品牌列表與單品牌 Landing Page 渲染入口。
- `brand.html`：專門處理單一品牌的頁面（透過 slug 參數）。
- `app.js`：核心渲染邏輯，使用 Vanilla JS 確保 GitHub Pages 相容性。

## 生成流程（新增品牌）
1) 收集素材：hero 風景照、產品照 >=10、徽章/認證、短影片 3+、CTA 連結。  
2) 更新 `docs/brands.xlsx`。
3) 執行 `go run tools/convert.go` 更新 `data/brands.json`。
4) 將素材放置於指定路徑（或由腳本自動下載）。
5) 驗收：檢查手機端 CTA 浮動、影片播放事件、產品卡 CTA 正常。

## QR 與導流
- 為每品牌產生 QR 圖檔（SVG/PNG），檔名 `qr-<slug>.svg`；可在後台或 CI 以 `qr-image` 等工具生成。
- 路由 `/qr/<qrSlug>` 302 到 `/brands/<slug>`，便於短域名印刷。
- 建議 UTM：`?utm_source=qr&utm_medium=offline&utm_campaign=<slug>`.

## 追蹤與周報
- 事件：`view_brand`, `cta_click`, `product_click`, `video_play`, `qr_scan`。  
- 欄位：`brand`, `productId`, `position`(hero/footer/products)、`device`、`utm_*`。  
- 資料管道：前端送 API -> 存入 DB (如 Supabase / Postgres) 或 CSV，後端排程周報。  
- 周報內容：掃碼次數、UV、CTA 點擊率、影片完成率、熱門產品 Top 5。  

## UI/UX 重點
- 自然風格：田園/作物紋理背景、暖色系點綴；避免單色平鋪，可用柔和漸層或紙質紋路。  
- 手機優先：Hero CTA 與浮動底部 CTA 保持可見；產品卡左右滑；影片縮圖點擊播放。  
- 速度：圖片 WebP、自適應尺寸；影片使用平台嵌入或懶載。  
- 無障礙：語意標籤、替代文字、對比度合規。  

## 未來擴充
- 多語系（繁中/英），資料檔支援 `locales`。  
- 搜尋/篩選品牌列表頁（地區、品類）。  
- 行銷模組：限定檔期 Banner、優惠碼欄位、A/B 測試不同 CTA 文案。  

## 範例靜態頁
- 位置：根目錄 `index.html`（品牌列表）與 `brand.html?slug=<品牌 slug>`（單品牌頁）。  
- 資料：`data/brands.json`，直接讀取並渲染；可透過前述 XLSX 轉換腳本自動生成。  
- 風格：自然背景、CTA 浮動，展示 Hero、產品列表、影片、FAQ。  
- 使用方式：啟動本機伺服器後瀏覽 `http://localhost:8000/index.html`，點品牌卡片即可查看示範一頁式。  
- 海浪 Hero 風格示範：`hero-concept.html` 搭配 `hero-concept.css`，採用 Noto Sans TC、背景圖覆疊與字體背景裁切（background-clip）呈現文字與圖片融合效果，可作為第一屏視覺參考。  
- 首頁整合海浪風格：`index.html` 已加入海浪 Hero，動態讀取品牌的 heroWords（巨字標語）、主圖、標語、CTA。  
- Hero 對比自動判斷：頁面會讀取 heroImage，計算亮度，若偏亮則加強覆蓋濾鏡與文字陰影（保持白字），提升可讀性；另在 Hero 文字背後加入 liquid glass 效果的玻璃卡片，凸顯文字。  
- 長頁結構：`index.html` 以單品牌長卷頁呈現（Hero + Story + 精選產品 <=3 + 護照 CTA + 據點），背景疊加植物裝飾圖與紙質感底色。  

## 開發路線圖 (Current Roadmap)
- [x] **Phase 1: 基礎自動化 [DONE]**
  - 實作 Go 轉換程式 `tools/convert.go`，成功連結 `brands.xlsx` 與 `brands.json`。
  - 使用 `excelize` 套件進行資料讀取。
- [ ] **Phase 2: 自然風格視覺強化 [進行中]**
  - 使用 AI 生成高品質農地與工坊景觀圖 (Done)。
  - 引入紙張紋理 (Natural Paper Texture) 與植物插圖裝飾。
  - 實作 Hero 區塊的漸層覆蓋與高對比文字設計。
- [ ] **Phase 3: CTA 與成效優化**
  - 實作行動端常駐 CTA 按鈕。
  - 增加產品卡片的 Hover/Tap 效果，強調購買連結。
- [ ] **Phase 4: 生成與發布優化**
  - 整合單一品牌頁面範本渲染邏輯。
  - 驗證 GitHub Pages 相容性。
