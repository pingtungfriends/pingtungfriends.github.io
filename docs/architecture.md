# 系統架構說明 (System Architecture)

本專案「屏東農特產品牌官網」採用 **混合式渲染 (Hybrid Rendering)** 架構，結合靜態頁面生成 (SSG) 與客戶端動態載入 (CSR)，兼顧 SEO 效能與開發彈性。

## 1. 核心流程 (Core Workflow)

1.  **資料源 (Data Source)**: 管理者維護 `docs/brands.xlsx` Excel 檔案。
2.  **轉換與生成 (Convert & Generate)**: 執行 `go run tools/convert.go`：
    *   讀取 Excel 內容並驗證欄位。
    *   將內容轉換為 `data/brands.json`（供前端動態載入使用）。
    *   基於 `brand.html` 模板，為每個品牌產生獨立的目錄及 `index.html`。
3.  **渲染機制 (Rendering)**:
    *   **SSG (Server-Side Generation)**: 品牌名稱、故事、產品、影片、背景圖等靜態內容由 Go 程式直接寫入各品牌 HTML，確保即時開啟與 SEO 抓取。
    *   **CSR (Client-Side Rendering)**: `app.js` 負責處理複雜互動或二次渲染（如 FAQ、Logo 切換邏輯等）。

## 2. 目錄結構 (Folder Structure)

*   `docs/`: 存放原始文件與資料原 (`brands.xlsx`)。
*   `tools/`: 存放資料處理程式 (`convert.go`)。
*   `data/`: 存放生成後的 JSON 資料 (`brands.json`)。
*   `public/assets/`: 存放所有靜態資源（圖片、商標），按品牌 `slug` 分類。
*   `brands/`: 生成後的各品牌網站目錄。
*   `index.html`: 品牌地圖首頁。
*   `brand.html`: 品牌頁面基礎模板。
*   `style.css`: 全站統一樣式表。
*   `app.js`: 前端互動邏輯。

## 3. 技術棧 (Tech Stack)

*   **Frontend**: Vanilla HTML/JS, Modern CSS (Variable/Grid/Flex), CSS Backdrop Filter.
*   **Backend Tooling**: Go 1.20+ (using `excelize` library).
*   **Design**: Premium Botanical Aesthetic (Playfair Display font, Glassmorphism).
