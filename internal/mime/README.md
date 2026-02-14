# mime

## 1. Package 概述（Overview）

### Package 名稱
`mime` - 一個用於檔案副檔名與 MIME 類型轉換的 Go 套件

### 簡介 / 目的
`mime` 套件提供了一個簡單且完整的方式來處理檔案副檔名與 MIME 類型之間的轉換。本套件支援廣泛的檔案類型，包括文字、圖片、音訊、視訊、文件、字型等多種格式，適合用於檔案上傳、下載、內容類型判斷等場景。

### 主要功能摘要
- **副檔名轉 MIME 類型**：將檔案副檔名轉換為對應的 MIME 類型
- **檔案名稱轉 MIME 類型**：從完整檔案名稱中提取副檔名並轉換為 MIME 類型
- **MIME 類型轉副檔名**：將 MIME 類型轉換為對應的檔案副檔名
- **自動正規化**：自動處理副檔名的格式（移除或添加點號、轉為小寫）
- **廣泛支援**：支援超過 200 種常見檔案格式

### 適合使用情境
- 檔案上傳時的內容類型驗證
- HTTP 回應的 Content-Type 設定
- 檔案下載時的檔案名稱與類型處理
- 檔案管理系統的類型識別
- API 端點的檔案類型檢查

## 2. 快速開始（Quick Start）

### 最小可用範例

基本轉換功能：

```go
package main

import (
	"fmt"
	"your-module/backend/pkg/mime"
)

func main() {
	convertor := mime.NewConvertor()

	// 副檔名轉 MIME 類型
	mimeType := convertor.FileExtensionToMime(".jpg")
	fmt.Println(mimeType) // 輸出: image/jpeg

	// 檔案名稱轉 MIME 類型
	mimeType2 := convertor.FileNameToMime("document.pdf")
	fmt.Println(mimeType2) // 輸出: application/pdf

	// MIME 類型轉副檔名
	ext := convertor.MimeToFileExtension("image/png")
	fmt.Println(ext) // 輸出: .png
}
```

### 簡單範例程式碼

處理檔案上傳：

```go
package main

import (
	"fmt"
	"your-module/backend/pkg/mime"
)

func main() {
	convertor := mime.NewConvertor()

	fileName := "photo.jpg"
	mimeType := convertor.FileNameToMime(fileName)
	
	if mimeType == "" {
		fmt.Println("Unsupported file type")
		return
	}

	fmt.Printf("File: %s\n", fileName)
	fmt.Printf("MIME Type: %s\n", mimeType)
	
	// 驗證是否為圖片類型
	if len(mimeType) >= 6 && mimeType[:6] == "image/" {
		fmt.Println("This is an image file")
	}
}
```

### 預期輸出或行為
- `FileExtensionToMime()` 和 `FileNameToMime()` 會回傳對應的 MIME 類型字串（如 "image/jpeg"）
- `MimeToFileExtension()` 會回傳對應的副檔名字串（如 ".jpg"）
- 如果找不到對應的轉換，會回傳空字串
- 副檔名會自動正規化（移除或添加點號、轉為小寫）

## 3. 使用範例（Examples / Use Cases）

### 常見使用情境

#### 1. 基本副檔名轉 MIME 類型
```go
convertor := mime.NewConvertor()
mimeType := convertor.FileExtensionToMime(".pdf")
// 輸出: application/pdf
```

#### 2. 從檔案名稱取得 MIME 類型
```go
convertor := mime.NewConvertor()
mimeType := convertor.FileNameToMime("document.pdf")
// 輸出: application/pdf
```

#### 3. MIME 類型轉副檔名
```go
convertor := mime.NewConvertor()
ext := convertor.MimeToFileExtension("application/json")
// 輸出: .json
```

#### 4. 處理不同格式的副檔名輸入
```go
convertor := mime.NewConvertor()

// 以下都會正確處理
convertor.FileExtensionToMime(".jpg")   // 有前綴點號
convertor.FileExtensionToMime("jpg")    // 無前綴點號
convertor.FileExtensionToMime(".JPG")   // 大寫
convertor.FileExtensionToMime("JPG")    // 大寫且無前綴點號
// 全部都會回傳: image/jpeg
```

#### 5. 檢查不支援的檔案類型
```go
convertor := mime.NewConvertor()
mimeType := convertor.FileExtensionToMime(".unknown")
if mimeType == "" {
	fmt.Println("Unsupported file type")
}
```

### 進階使用方式

#### 在 HTTP 處理器中使用
```go
package main

import (
	"net/http"
	"your-module/backend/pkg/mime"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	convertor := mime.NewConvertor()
	
	// 從檔案名稱取得 MIME 類型
	expectedMimeType := convertor.FileNameToMime(header.Filename)
	if expectedMimeType == "" {
		http.Error(w, "Unsupported file type", http.StatusBadRequest)
		return
	}

	// 驗證上傳的檔案類型
	detectedMimeType := header.Header.Get("Content-Type")
	if detectedMimeType != expectedMimeType {
		http.Error(w, "File type mismatch", http.StatusBadRequest)
		return
	}

	// 處理檔案上傳...
	w.WriteHeader(http.StatusOK)
}
```

#### 在檔案下載服務中使用
```go
type FileService struct {
	convertor mime.Convertor
}

func NewFileService() *FileService {
	return &FileService{
		convertor: mime.NewConvertor(),
	}
}

func (s *FileService) DownloadFile(w http.ResponseWriter, fileName string, fileData []byte) {
	// 根據檔案名稱設定 Content-Type
	mimeType := s.convertor.FileNameToMime(fileName)
	if mimeType == "" {
		mimeType = "application/octet-stream" // 預設類型
	}

	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Write(fileData)
}
```

#### 檔案類型驗證中介軟體
```go
func ValidateFileType(allowedTypes []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			file, header, err := r.FormFile("file")
			if err != nil {
				http.Error(w, "No file provided", http.StatusBadRequest)
				return
			}
			file.Close()

			convertor := mime.NewConvertor()
			mimeType := convertor.FileNameToMime(header.Filename)

			allowed := false
			for _, allowedType := range allowedTypes {
				if mimeType == allowedType {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, "File type not allowed", http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// 使用範例
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploadHandler)
	
	// 只允許圖片類型
	handler := ValidateFileType([]string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
	})(mux)
	
	http.ListenAndServe(":8080", handler)
}
```

#### 根據 MIME 類型決定處理方式
```go
func ProcessFile(fileName string, fileData []byte) error {
	convertor := mime.NewConvertor()
	mimeType := convertor.FileNameToMime(fileName)

	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return processImage(fileData)
	case strings.HasPrefix(mimeType, "video/"):
		return processVideo(fileData)
	case strings.HasPrefix(mimeType, "audio/"):
		return processAudio(fileData)
	case mimeType == "application/pdf":
		return processPDF(fileData)
	default:
		return processGeneric(fileData)
	}
}
```

### 實務案例

#### 案例 1：檔案上傳 API
```go
type UploadResponse struct {
	FileName string `json:"fileName"`
	MimeType string `json:"mimeType"`
	Size     int64  `json:"size"`
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	convertor := mime.NewConvertor()
	mimeType := convertor.FileNameToMime(header.Filename)

	if mimeType == "" {
		http.Error(w, "Unsupported file type", http.StatusBadRequest)
		return
	}

	// 儲存檔案...
	fileInfo, _ := file.Stat()
	
	response := UploadResponse{
		FileName: header.Filename,
		MimeType: mimeType,
		Size:     fileInfo.Size(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

#### 案例 2：檔案下載 API
```go
func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name required", http.StatusBadRequest)
		return
	}

	// 從資料庫或儲存系統取得檔案
	fileData, err := getFileData(fileName)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	convertor := mime.NewConvertor()
	mimeType := convertor.FileNameToMime(fileName)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Write(fileData)
}
```

#### 案例 3：檔案類型過濾
```go
type FileTypeFilter struct {
	convertor mime.Convertor
	allowed   map[string]bool
}

func NewFileTypeFilter(allowedTypes []string) *FileTypeFilter {
	allowed := make(map[string]bool)
	for _, t := range allowedTypes {
		allowed[t] = true
	}
	return &FileTypeFilter{
		convertor: mime.NewConvertor(),
		allowed:   allowed,
	}
}

func (f *FileTypeFilter) IsAllowed(fileName string) bool {
	mimeType := f.convertor.FileNameToMime(fileName)
	return f.allowed[mimeType]
}

// 使用範例
func main() {
	filter := NewFileTypeFilter([]string{
		"image/jpeg",
		"image/png",
		"application/pdf",
	})

	if filter.IsAllowed("photo.jpg") {
		fmt.Println("File is allowed")
	}
}
```

#### 案例 4：批次處理檔案
```go
func ProcessFiles(fileNames []string) map[string]string {
	convertor := mime.NewConvertor()
	result := make(map[string]string)

	for _, fileName := range fileNames {
		mimeType := convertor.FileNameToMime(fileName)
		if mimeType != "" {
			result[fileName] = mimeType
		}
	}

	return result
}
```

## 4. 限制（Limitations）

### 已知限制
1. **單一對應關係**：一個 MIME 類型對應一個副檔名，但實際上一個 MIME 類型可能對應多個副檔名（例如 "image/jpeg" 對應 ".jpg" 和 ".jpeg"）。本套件只回傳其中一個預設的副檔名。
2. **不支援的格式**：如果輸入的副檔名或 MIME 類型不在支援列表中，會回傳空字串。需要呼叫端自行處理這種情況。
3. **大小寫敏感性**：雖然副檔名會自動轉為小寫，但 MIME 類型的大小寫必須完全匹配（例如 "Image/JPEG" 不會匹配 "image/jpeg"）。
4. **無參數驗證**：不會驗證 MIME 類型字串的格式是否正確，只進行簡單的字串匹配。

### 不支援的情況
1. **自訂 MIME 類型**：不支援新增或修改 MIME 類型與副檔名的對應關係。對應關係是固定的。
2. **MIME 類型參數**：不支援帶有參數的 MIME 類型（例如 "text/html; charset=utf-8"）。需要先移除參數部分再進行轉換。
3. **多副檔名檔案**：不支援處理多個副檔名的檔案（例如 "file.tar.gz"），只會處理最後一個副檔名。
4. **動態 MIME 類型檢測**：不支援根據檔案內容動態檢測 MIME 類型，只根據副檔名進行轉換。

### 建議
- **錯誤處理**：始終檢查轉換結果是否為空字串，並提供適當的錯誤處理或預設值。
- **MIME 類型正規化**：在進行 MIME 類型轉換前，建議先移除參數部分（例如 "text/html; charset=utf-8" → "text/html"）。
- **預設值**：對於不支援的檔案類型，建議使用 "application/octet-stream" 作為預設的 MIME 類型。
- **驗證流程**：在檔案上傳場景中，建議同時驗證檔案名稱的 MIME 類型和實際檔案的 Content-Type，以確保安全性。
- **效能考量**：`Convertor` 實例是無狀態的，可以安全地在多個 goroutine 之間共享，或為每個請求建立新實例。

### 支援的檔案類型類別
本套件支援以下主要類別的檔案類型：
- **文字檔案**：txt, html, css, csv, xml, markdown, yaml 等
- **圖片檔案**：jpeg, png, gif, bmp, webp, svg, tiff 等
- **音訊檔案**：mp3, wav, ogg, m4a, flac, aac 等
- **視訊檔案**：mp4, webm, avi, mov, mkv, wmv 等
- **文件檔案**：pdf, doc, docx, xls, xlsx, ppt, pptx 等
- **壓縮檔案**：zip, gz, tar, 7z, rar 等
- **字型檔案**：ttf, otf, woff, woff2 等
- **程式碼檔案**：js, ts, py, java, c, cpp 等
- **其他**：模型檔案、化學檔案、訊息檔案等


