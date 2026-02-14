# jwt

## 1. Package 概述（Overview）

### Package 名稱
`jwt` - 一個用於 JSON Web Token (JWT) 編碼、解碼與管理的 Go 套件

### 簡介 / 目的
`jwt` 套件提供了一個完整且易用的方式來處理 JWT token 的生成、編碼、解碼與驗證。本套件支援多種簽名方法（無簽名、HMAC、RSA），並提供靈活的 claims 建構器，適合用於身份驗證、授權、API 安全等場景。

### 主要功能摘要
- **多種簽名方法**：支援無簽名（none）、HMAC（HS256）、RSA（RS256）三種簽名方法
- **Claims 建構器**：使用 builder 模式輕鬆建立 JWT claims，支援標準與自訂欄位
- **RSA 金鑰生成**：提供 RSA 金鑰對生成功能
- **類型安全的 Claims 存取**：提供多種類型轉換方法（String, Int, Float64, Bool 等）
- **標準 JWT Claims**：完整支援 JWT 標準 claims（iss, sub, aud, exp, nbf, iat, jti）

### 適合使用情境
- 使用者身份驗證與授權
- API 端點的 token 驗證
- 微服務間的安全通訊
- Session 管理
- 單一登入（SSO）系統
- 任何需要安全 token 傳遞的應用場景

## 2. 快速開始（Quick Start）

### 最小可用範例

使用 Secret Encoder 建立和驗證 JWT：

```go
package main

import (
	"fmt"
	"time"
	"your-module/backend/pkg/jwt"
)

func main() {
	// 建立 Secret Encoder
	encoder := jwt.NewSecretEncoder("my-secret-key")

	// 建立 claims
	claims := jwt.NewClaimsBuilder().
		Subject("user123").
		Issuer("my-app").
		ExpirationTime(time.Now().Add(1 * time.Hour).Unix()).
		IssuedAt(time.Now().Unix()).
		Build()

	// 編碼 token
	token, err := encoder.Encode(claims)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Token: %s\n", token)

	// 解碼並驗證 token
	decodedClaims, err := encoder.Decode(token)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Subject: %s\n", decodedClaims.Subject())
}
```

### 簡單範例程式碼

完整的登入與驗證流程：

```go
package main

import (
	"fmt"
	"time"
	"your-module/backend/pkg/jwt"
)

func main() {
	secret := "my-secret-key"
	encoder := jwt.NewSecretEncoder(secret)

	// 使用者登入：產生 token
	claims := jwt.NewClaimsBuilder().
		Subject("user123").
		Issuer("my-app").
		Audience("api").
		ExpirationTime(time.Now().Add(24 * time.Hour).Unix()).
		IssuedAt(time.Now().Unix()).
		ID("unique-token-id").
		Build()

	token, err := encoder.Encode(claims)
	if err != nil {
		panic(err)
	}

	// 驗證 token
	decodedClaims, err := encoder.Decode(token)
	if err != nil {
		fmt.Println("Invalid token")
		return
	}

	// 檢查是否過期
	if time.Now().Unix() > decodedClaims.ExpirationTime() {
		fmt.Println("Token expired")
		return
	}

	fmt.Printf("User: %s\n", decodedClaims.Subject())
	fmt.Printf("Issuer: %s\n", decodedClaims.Issuer())
}
```

### 預期輸出或行為
- `Encode()` 方法會回傳一個 JWT token 字串
- `Decode()` 方法會回傳 `Claims` 介面實例，可用於存取 claims 值
- 如果 token 無效、過期或簽名錯誤，`Decode()` 會回傳錯誤
- Claims 提供類型安全的方法來存取各種類型的值

## 3. 使用範例（Examples / Use Cases）

### 常見使用情境

#### 1. 使用 Secret Encoder（最常用）
```go
encoder := jwt.NewSecretEncoder("my-secret-key")

claims := jwt.NewClaimsBuilder().
	Subject("user123").
	ExpirationTime(time.Now().Add(1 * time.Hour).Unix()).
	Build()

token, err := encoder.Encode(claims)
// 驗證
decodedClaims, err := encoder.Decode(token)
```

#### 2. 使用 RSA Encoder（更高安全性）
```go
// 生成 RSA 金鑰對
generator := jwt.NewRsaGenerator()
privateKey, publicKey, err := generator.GenerateKeyPair()

// 使用私鑰編碼
encoder := jwt.NewRsaEncoder(privateKey)
claims := jwt.NewClaimsBuilder().
	Subject("user123").
	Build()
token, err := encoder.Encode(claims)

// 使用公鑰解碼（驗證）
decodedClaims, err := encoder.PublicDecode(token)
```

#### 3. 建立包含自訂欄位的 Claims
```go
claims := jwt.NewClaimsBuilder().
	Subject("user123").
	Issuer("my-app").
	Key("role").Set("admin").
	Key("permissions").Set([]string{"read", "write"}).
	Build()

// 存取自訂欄位
role := claims.String("role")
permissions := claims.StringSlice("permissions")
```

#### 4. 使用 Unverified Encoder（開發/測試用）
```go
encoder := jwt.NewUnverifiedEncoder()

claims := jwt.NewClaimsBuilder().
	Subject("user123").
	Build()

token, err := encoder.Encode(claims)
// 注意：此方法不進行簽名驗證，僅用於開發測試
```

#### 5. 存取標準 Claims
```go
decodedClaims, _ := encoder.Decode(token)

issuer := decodedClaims.Issuer()        // "iss"
subject := decodedClaims.Subject()       // "sub"
audience := decodedClaims.Audience()    // "aud"
exp := decodedClaims.ExpirationTime()    // "exp"
nbf := decodedClaims.NotBefore()        // "nbf"
iat := decodedClaims.IssuedAt()         // "iat"
jti := decodedClaims.ID()               // "jti"
```

### 進階使用方式

#### 在 HTTP 中介軟體中使用
```go
func JWTAuthMiddleware(secret string) func(http.Handler) http.Handler {
	encoder := jwt.NewSecretEncoder(secret)
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			// 移除 "Bearer " 前綴
			token := strings.TrimPrefix(authHeader, "Bearer ")
			
			claims, err := encoder.Decode(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// 檢查過期時間
			if time.Now().Unix() > claims.ExpirationTime() {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}

			// 將 claims 存入 context
			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
```

#### 在服務中使用 JWT
```go
type AuthService struct {
	encoder jwt.SecretEncoder
}

func NewAuthService(secret string) *AuthService {
	return &AuthService{
		encoder: jwt.NewSecretEncoder(secret),
	}
}

func (s *AuthService) GenerateToken(userID string, role string) (string, error) {
	claims := jwt.NewClaimsBuilder().
		Subject(userID).
		Issuer("my-app").
		ExpirationTime(time.Now().Add(24 * time.Hour).Unix()).
		IssuedAt(time.Now().Unix()).
		Key("role").Set(role).
		Build()

	return s.encoder.Encode(claims)
}

func (s *AuthService) ValidateToken(token string) (jwt.Claims, error) {
	claims, err := s.encoder.Decode(token)
	if err != nil {
		return nil, err
	}

	// 檢查過期
	if time.Now().Unix() > claims.ExpirationTime() {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}
```

#### 使用 RSA 進行微服務間通訊
```go
type ServiceAuth struct {
	privateEncoder jwt.RsaEncoder
	publicEncoder  jwt.RsaEncoder
}

func NewServiceAuth(privateKey, publicKey string) *ServiceAuth {
	return &ServiceAuth{
		privateEncoder: jwt.NewRsaEncoder(privateKey),
		publicEncoder:  jwt.NewRsaEncoder(publicKey),
	}
}

func (s *ServiceAuth) GenerateServiceToken(serviceName string) (string, error) {
	claims := jwt.NewClaimsBuilder().
		Subject(serviceName).
		Issuer("service-registry").
		ExpirationTime(time.Now().Add(1 * time.Hour).Unix()).
		Build()

	return s.privateEncoder.Encode(claims)
}

func (s *ServiceAuth) ValidateServiceToken(token string) (jwt.Claims, error) {
	return s.publicEncoder.PublicDecode(token)
}
```

#### 動態更新 Claims
```go
func RefreshToken(oldToken string, encoder jwt.SecretEncoder) (string, error) {
	// 解碼舊 token
	oldClaims, err := encoder.Decode(oldToken)
	if err != nil {
		return "", err
	}

	// 使用舊 claims 建立新 claims
	newClaims := oldClaims.Builder().
		ExpirationTime(time.Now().Add(24 * time.Hour).Unix()).
		IssuedAt(time.Now().Unix()).
		Build()

	return encoder.Encode(newClaims)
}
```

### 實務案例

#### 案例 1：使用者登入 API
```go
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 驗證使用者（省略實際驗證邏輯）
	userID := "user123"
	
	encoder := jwt.NewSecretEncoder(os.Getenv("JWT_SECRET"))
	claims := jwt.NewClaimsBuilder().
		Subject(userID).
		Issuer("my-app").
		ExpirationTime(time.Now().Add(24 * time.Hour).Unix()).
		IssuedAt(time.Now().Unix()).
		Key("email").Set(req.Email).
		Build()

	token, err := encoder.Encode(claims)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

#### 案例 2：受保護的 API 端點
```go
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt.Claims)
	
	userID := claims.Subject()
	role := claims.String("role")

	response := map[string]interface{}{
		"userID": userID,
		"role":   role,
		"message": "This is a protected resource",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 使用方式
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/protected", ProtectedHandler)
	
	handler := JWTAuthMiddleware(os.Getenv("JWT_SECRET"))(mux)
	http.ListenAndServe(":8080", handler)
}
```

#### 案例 3：Token 刷新機制
```go
type RefreshRequest struct {
	Token string `json:"token"`
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	encoder := jwt.NewSecretEncoder(os.Getenv("JWT_SECRET"))
	
	// 解碼舊 token（允許稍微過期）
	oldClaims, err := encoder.Decode(req.Token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// 檢查是否過期太久（例如超過 7 天）
	expTime := oldClaims.ExpirationTime()
	if time.Now().Unix()-expTime > 7*24*3600 {
		http.Error(w, "Token expired too long ago", http.StatusUnauthorized)
		return
	}

	// 產生新 token
	newClaims := oldClaims.Builder().
		ExpirationTime(time.Now().Add(24 * time.Hour).Unix()).
		IssuedAt(time.Now().Unix()).
		Build()

	newToken, err := encoder.Encode(newClaims)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{Token: newToken}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

#### 案例 4：多租戶系統的 JWT
```go
func GenerateTenantToken(userID, tenantID string) (string, error) {
	encoder := jwt.NewSecretEncoder(os.Getenv("JWT_SECRET"))
	
	claims := jwt.NewClaimsBuilder().
		Subject(userID).
		Issuer("my-app").
		Audience(tenantID).
		ExpirationTime(time.Now().Add(24 * time.Hour).Unix()).
		IssuedAt(time.Now().Unix()).
		Key("tenant_id").Set(tenantID).
		Key("permissions").Set([]string{"read", "write"}).
		Build()

	return encoder.Encode(claims)
}

func ValidateTenantAccess(claims jwt.Claims, tenantID string) bool {
	// 檢查 audience 或自訂欄位
	return claims.Audience() == tenantID || 
		   claims.String("tenant_id") == tenantID
}
```

#### 案例 5：使用 RSA 進行跨服務驗證
```go
// 服務 A：產生 token
func ServiceAHandler(w http.ResponseWriter, r *http.Request) {
	privateKey := os.Getenv("SERVICE_A_PRIVATE_KEY")
	encoder := jwt.NewRsaEncoder(privateKey)
	
	claims := jwt.NewClaimsBuilder().
		Subject("service-a").
		Issuer("service-a").
		ExpirationTime(time.Now().Add(1 * time.Hour).Unix()).
		Build()

	token, _ := encoder.Encode(claims)
	
	// 呼叫服務 B 時帶上此 token
	// ...
}

// 服務 B：驗證 token
func ServiceBHandler(w http.ResponseWriter, r *http.Request) {
	publicKey := os.Getenv("SERVICE_A_PUBLIC_KEY")
	encoder := jwt.NewRsaEncoder(publicKey)
	
	token := r.Header.Get("X-Service-Token")
	claims, err := encoder.PublicDecode(token)
	if err != nil {
		http.Error(w, "Invalid service token", http.StatusUnauthorized)
		return
	}

	// 驗證服務身份
	if claims.Subject() != "service-a" {
		http.Error(w, "Unauthorized service", http.StatusUnauthorized)
		return
	}

	// 處理請求...
}
```

## 4. 限制（Limitations）

### 已知限制
1. **Secret Encoder 安全性**：使用 HS256 時，secret 必須足夠長且隨機。建議使用至少 32 字元的隨機字串。
2. **RSA 金鑰格式**：RSA encoder 支援 PKIX 和 PKCS1 格式的公鑰，但私鑰必須是 PKCS1 格式。
3. **Claims 類型轉換**：當 JSON 解碼時，數字可能被解析為 `float64`。Claims 的類型轉換方法會自動處理這種情況，但可能會有精度損失。
4. **Token 大小**：JWT token 的大小會隨著 claims 的增加而增長。過大的 token 可能導致 HTTP header 大小限制問題。
5. **時鐘同步**：exp、nbf、iat 等時間相關 claims 依賴系統時鐘。如果系統時鐘不同步，可能導致驗證問題。

### 不支援的情況
1. **其他簽名演算法**：目前只支援 none、HS256 和 RS256，不支援其他演算法（如 ES256、PS256 等）。
2. **Token 加密**：本套件只處理 JWT 的簽名，不提供 JWE（JSON Web Encryption）功能。
3. **Token 黑名單**：不提供內建的 token 撤銷或黑名單機制，需要在使用端自行實作。
4. **自動過期檢查**：`Decode()` 方法不會自動檢查過期時間，需要在使用端自行檢查 `ExpirationTime()`。
5. **Claims 驗證**：不提供內建的 claims 驗證邏輯（如 issuer、audience 驗證），需要在使用端自行實作。

### 建議
- **Secret 管理**：
  - 使用環境變數或密鑰管理服務儲存 secret
  - 定期輪換 secret（需要實作 token 遷移機制）
  - 不同環境使用不同的 secret
- **Token 過期時間**：
  - 存取 token：建議 15 分鐘到 1 小時
  - 刷新 token：建議 7 到 30 天
  - 根據安全需求調整過期時間
- **RSA 金鑰**：
  - 使用至少 2048 位元的 RSA 金鑰（本套件預設使用 4096 位元）
  - 妥善保管私鑰，不要提交到版本控制系統
  - 定期輪換金鑰對
- **Claims 設計**：
  - 不要在 claims 中儲存敏感資訊（JWT 是 base64 編碼，不是加密）
  - 限制 claims 的大小，避免 token 過大
  - 使用標準 claims（iss, sub, aud, exp 等）提高相容性
- **安全性最佳實踐**：
  - 始終使用 HTTPS 傳輸 JWT token
  - 實作 token 刷新機制而非延長過期時間
  - 在伺服器端驗證所有 claims，不要信任客戶端
  - 實作適當的錯誤處理，避免洩露敏感資訊
- **效能考量**：
  - `Encoder` 和 `Decoder` 實例可以安全地在多個 goroutine 之間共享
  - 對於高並發場景，考慮使用連線池或快取機制
  - RSA 簽名比 HMAC 慢，根據效能需求選擇適當的演算法

