# Pingtung Friends｜屏東農特產品牌官網

介紹屏東在地社區與農特產品牌的全靜態網站。品牌資料以 Excel 維護，透過 Go 工具一鍵產生所有頁面。

## 快速開始

```bash
# 1. 編輯品牌資料
open docs/brands.xlsx

# 2. 重新產生網站（需 Go 1.25+，版本依 go.mod 設定）
go run tools/convert.go

# 3. 預覽：直接開啟 index.html，或啟動本地伺服器
python3 -m http.server 8000
```

> 若為**新增品牌**或**新增圖片素材**，請先參閱 [`docs/brands.md`](docs/brands.md) 完成素材放置與路線設定。

## 專案結構

| 路徑 | 說明 |
| --- | --- |
| `docs/brands.xlsx` | 品牌資料來源 |
| `docs/architecture.md` | 系統架構與生成機制說明 |
| `docs/brands.md` | 品牌資料管理與 Excel 欄位說明 |
| `docs/credits.md` | 屏東縣鄉鎮地圖與素材授權說明 |
| `tools/convert.go` | Excel → JSON / 靜態頁面與互動地圖轉換工具 |
| `data/brands.json` | 產生的資料備份（維持原始工作表順序） |
| `brands/<slug>/` | 產生的各品牌頁面（依由北到南路線串聯） |
| `index.html` | 產生的首頁（模板：`index_template.html`，含共用字型 `fonts` 與頁尾 `footer` 子模板） |
| `brand.html` | 品牌頁模板 |
| `public/assets/common/pingtung-map.svg` | 屏東縣鄉鎮底圖原始 SVG（供 `convert.go` 內嵌轉換） |
| `public/assets/<slug>/` | 各品牌圖片與標誌素材 |
| `style.css` / `app.js` | 全站樣式與前端漸進增強腳本 |

## 文件

- [系統架構](docs/architecture.md)
- [品牌資料管理與 Excel 欄位說明](docs/brands.md)
- [素材授權](docs/credits.md)
