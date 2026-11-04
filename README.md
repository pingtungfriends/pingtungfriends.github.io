# Pingtung Friends｜屏東農特產品牌官網

介紹屏東在地社區與農特產品牌的全靜態網站。品牌資料以 Excel 維護，透過 Go 工具一鍵產生所有頁面。

## 快速開始

```bash
# 1. 編輯品牌資料
open docs/brands.xlsx

# 2. 重新產生網站（需 Go 1.20+）
go run tools/convert.go

# 3. 預覽：直接開啟 index.html，或啟動本地伺服器
python3 -m http.server 8000
```

## 專案結構

| 路徑 | 說明 |
| --- | --- |
| `docs/brands.xlsx` | 品牌資料來源 |
| `tools/convert.go` | Excel → JSON / 靜態頁面轉換工具 |
| `data/brands.json` | 產生的資料備份 |
| `brands/<slug>/` | 產生的各品牌頁面 |
| `index.html` | 產生的首頁（模板：`index_template.html`） |
| `brand.html` | 品牌頁模板 |
| `public/assets/<slug>/` | 各品牌圖片素材 |
| `style.css` / `app.js` | 全站樣式與前端腳本 |

## 文件

- [系統架構](docs/architecture.md)
- [品牌資料管理與 Excel 欄位說明](docs/brands.md)
