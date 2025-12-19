# 系統架構說明 (System Architecture)

本專案「屏東農特產品牌官網」採用 **全靜態生成 (Full Static Site Generation, SSG)** 架構，確保最佳的載入效能與最高的相容性。

## 1. 核心流程 (Core Workflow)

1.  **資料源 (Data Source)**: 管理者維護 `docs/brands.xlsx` Excel 檔案。
2.  **轉換與生成 (Convert & Generate)**: 執行 `go run tools/convert.go`：
    *   讀取 Excel 內容並驗證欄位。
    *   讀取 Excel 內容並驗證欄位。
    *   將內容轉換為 JSON 資料（做為資料備份）。
    *   基於 `brand.html` 模板，為每個品牌產生獨立的靜態頁面 (`brands/*/index.html`)。
    *   基於 `index_template.html` 模板，生成全靜態的 `index.html` 首頁。
3.  **渲染機制 (Rendering)**:
    *   **Full SSG (Full Static Site Generation)**: 所有內容（包括品牌列表、故事、產品、FAQ）皆在建置時期直接寫入 HTML。
    *   **Zero CSR Data Fetching**: 前端不再進行任何 API 請求，開啟速度極快且無 CORS 限制。

## 2. 目錄結構 (Folder Structure)

*   `docs/`: 存放原始文件與資料原 (`brands.xlsx`)。
*   `tools/`: 存放資料處理程式 (`convert.go`)。
*   `data/`: 存放生成後的 JSON 資料 (`brands.json`)。
*   `public/assets/`: 存放所有靜態資源（圖片、商標），按品牌 `slug` 分類。
*   `brands/`: 生成後的各品牌 Landing Page 目錄。
*   `index.html`: 生成後的品牌地圖首頁。
*   `index_template.html`: 首頁生成模板。
*   `brand.html`: 品牌頁面基礎模板。
*   `style.css`: 全站統一樣式表。
*   `app.js`: 前端互動邏輯。

## 4. 開發預覽 (Local Development)

由於網站已全面靜態化，您可以直接透過 **「雙擊檔案 (Double-click)」** 的方式開啟 `index.html` 進行預覽，無需啟動本地伺服器。

當然，您也可以使用標準的本地伺服器進行開發：
```bash
python3 -m http.server 8000
```

## 5. 技術棧 (Tech Stack)

*   **Frontend**: Vanilla HTML/JS, Modern CSS (Variable/Grid/Flex), CSS Backdrop Filter.
*   **Backend Tooling**: Go 1.20+ (using `excelize` library).
*   **Design**: Premium Botanical Aesthetic (Playfair Display font, Glassmorphism).
