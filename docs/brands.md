# 品牌資料管理與轉換指南 (Brands Data Management)

本計畫所有品牌資訊皆透過 `docs/brands.xlsx` 進行管理。

## 1. Excel 欄位定義 (Field Definitions)

| 欄位名稱 | 說明 | 範例 |
| :--- | :--- | :--- |
| **Slug** | 品牌唯一識別碼（英文），會成為網址路徑 | `xinlong`, `ailiao` |
| **品牌名稱** | 顯示於頁面頂部與標題 | `新龍社區發展協會` |
| **Logo** | 品牌標誌檔名。存放於 `public/assets/{slug}/` 下 | `logo.png` |
| **Hero Image** | 首頁大背景圖檔名。存放於 `public/assets/{slug}/` 下 | `hero.jpg` |
| **標語** | 品牌副標題 | `食漁、樂漁、慢漁` |
| **品牌故事** | 品牌簡介文字 | `坐擁豐富海洋的環境...` |
| **Hero 巨字標語** | 顯示在毛玻璃裝飾框內的英文大字 (以 `|` 分隔) | `SEA | FOOD | SLOW` |
| **產品 [1-10] 名稱** | 產品名稱 | `龍膽石斑魚粥` |
| **產品 [1-10] 圖片url** | 產品圖檔名。存放於 `public/assets/{slug}/` 下 | `product_01.png` |
| **產品 [1-10] 摘要** | 產品短評 | `慢熬釋放膠原蛋白...` |
| **產品 [1-10] 價格** | 產品售價（整數） | `300` |
| **影片 [1-3] 標題** | 影片介紹文字 | `漁村體驗紀錄` |
| **影片 [1-3] 網址** | YouTube 或影片連結 | `https://youtu.be/...` |

## 2. 資產路徑慣例 (Asset Mapping)

轉換程式 `convert.go` 具備自動路徑補全功能：
1.  如果您在 Excel 填入 `H3.jpg` 且 Slug 為 `xinlong`。
2.  程式會自動對應到 `public/assets/xinlong/H3.jpg`。
3.  **注意**：檔名中的全形字元（如 `＿`、`１`）也會被自動轉成半形，以避免 404 錯誤。

## 3. 使用轉換工具 (Usage)

當您更新 `brands.xlsx` 後，必須執行以下指令來同步更新網站：

```bash
# 確保位於專案根目錄
go run tools/convert.go
```

**執行後的變化：**
*   `data/brands.json` 會根據最新資料重新生成。
*   `brands/` 目錄下的所有子頁面會根據 `brand.html` 模板重新產生，並注入最新的靜態內容。

## 4. 注意事項
*   **Slug 不可重複**：這是區分不同品牌的關鍵。
*   **檔名一致性**：請確保上傳到 `public/assets/{slug}/` 的圖檔名稱與 Excel 中填報的完全一致。
