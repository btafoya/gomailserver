# Phase 3: Web Interfaces (Weeks 8-12)

**Status**: Not Started
**Priority**: MVP - Required
**Estimated Duration**: 4-5 weeks
**Dependencies**: Phase 1 (Core Mail), Phase 2 (Security)

---

## Overview

Build the REST API foundation, Admin Web UI for server management, User Self-Service Portal, Let's Encrypt integration for automatic TLS certificates, and a web-based setup wizard for initial configuration.

---

## 3.1 REST API Foundation [MVP]

| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| API-001 | Set up Echo web framework | [ ] | F-002 | `labstack/echo/v4` |
| API-002 | JWT authentication middleware | [ ] | API-001 | `golang-jwt/jwt/v5` |
| API-003 | API key authentication | [ ] | API-002 |
| API-004 | Request rate limiting middleware | [ ] | API-001 |
| API-005 | CORS configuration | [ ] | API-001 |
| API-006 | OpenAPI/Swagger documentation | [ ] | API-001 |
| API-007 | Request validation middleware | [ ] | API-001 |

### API-001: Echo Setup

```go
// internal/api/server.go
package api

import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

type Server struct {
    echo   *echo.Echo
    config *config.API
}

func NewServer(cfg *config.API) *Server {
    e := echo.New()

    // Global middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.RequestID())
    e.Use(middleware.Gzip())

    return &Server{
        echo:   e,
        config: cfg,
    }
}

func (s *Server) SetupRoutes(
    authMiddleware echo.MiddlewareFunc,
    domainHandler *handlers.DomainHandler,
    userHandler *handlers.UserHandler,
    // ... other handlers
) {
    // Health check
    s.echo.GET("/health", func(c echo.Context) error {
        return c.JSON(200, map[string]string{"status": "ok"})
    })

    // API v1
    v1 := s.echo.Group("/api/v1")

    // Public routes
    v1.POST("/auth/login", authHandler.Login)
    v1.POST("/auth/refresh", authHandler.RefreshToken)

    // Protected routes
    admin := v1.Group("/admin", authMiddleware, adminOnlyMiddleware)
    admin.GET("/domains", domainHandler.List)
    admin.POST("/domains", domainHandler.Create)
    admin.GET("/domains/:id", domainHandler.Get)
    admin.PUT("/domains/:id", domainHandler.Update)
    admin.DELETE("/domains/:id", domainHandler.Delete)

    admin.GET("/users", userHandler.List)
    admin.POST("/users", userHandler.Create)
    // ... etc
}

func (s *Server) Start() error {
    return s.echo.Start(fmt.Sprintf(":%d", s.config.Port))
}
```

### API-002: JWT Authentication

```go
// internal/api/middleware/auth.go
package middleware

import (
    "github.com/golang-jwt/jwt/v5"
    "github.com/labstack/echo/v4"
)

type Claims struct {
    UserID   int64  `json:"user_id"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

func JWTAuth(secret string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            authHeader := c.Request().Header.Get("Authorization")
            if authHeader == "" {
                return echo.NewHTTPError(401, "Missing authorization header")
            }

            tokenString := strings.TrimPrefix(authHeader, "Bearer ")

            token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
                return []byte(secret), nil
            })

            if err != nil || !token.Valid {
                return echo.NewHTTPError(401, "Invalid token")
            }

            claims := token.Claims.(*Claims)
            c.Set("user_id", claims.UserID)
            c.Set("email", claims.Email)
            c.Set("role", claims.Role)

            return next(c)
        }
    }
}

func GenerateToken(secret string, user *domain.User, expiry time.Duration) (string, error) {
    claims := &Claims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

**Acceptance Criteria**:
- [ ] JWT tokens expire after 24 hours (configurable)
- [ ] HS256 signing algorithm (symmetric key)
- [ ] Bearer token format: "Authorization: Bearer <token>"
- [ ] Token includes user_id, email, role claims
- [ ] Expired tokens rejected with 401 Unauthorized
- [ ] Invalid signature rejected with 401 Unauthorized
- [ ] Missing Authorization header rejected with 401
- [ ] Refresh token flow with 7-day expiry for long-lived sessions
- [ ] Secret key minimum 256 bits (32 bytes)
- [ ] Token validation on every protected endpoint

**Structured Logging (slog)**:
- [ ] **INFO**: JWT token generated (user_id, email, role, expires_at, issued_at, token_id, request_id)
- [ ] **INFO**: JWT authentication successful (user_id, email, role, endpoint, method, ip_address, session_id, request_id)
- [ ] **WARN**: JWT token expired (user_id, email, expires_at, endpoint, ip_address, age_hours, request_id)
- [ ] **ERROR**: JWT invalid signature (attempted_user_id, endpoint, ip_address, error="tampered_token", request_id, trace_id)
- [ ] **ERROR**: JWT missing authorization header (endpoint, method, ip_address, request_id)
- [ ] **INFO**: JWT token refreshed (user_id, email, old_token_id, new_token_id, new_expires_at, request_id)
- [ ] **DEBUG**: JWT validation started (endpoint, method, token_id, ip_address, request_id)
- [ ] **FATAL**: JWT secret key weak (key_length_bits, minimum_required=256, security_risk="high")
- [ ] **Fields**: user_id, email, role, token_id, endpoint, method, ip_address, session_id, request_id, trace_id, expires_at, issued_at

**Given/When/Then Scenarios**:
```
Given user "alice@example.com" with ID 123 and role "user"
When token is generated with 24-hour expiry
Then JWT contains claims: user_id=123, email="alice@example.com", role="user"
And exp claim is set to now + 24 hours
And iat claim is set to current timestamp
And token is signed with HS256 algorithm

Given valid JWT token for user 123
When API request is made to /api/v1/messages with Authorization header
Then token is parsed and validated
And user_id=123 is set in request context
And request proceeds to handler
And protected resource is accessed

Given JWT token issued 25 hours ago (expired)
When API request is made with expired token
Then token validation fails
And 401 Unauthorized is returned
And error message is "Token has expired"
And request is rejected without reaching handler

Given JWT token with tampered claims (user_id changed)
When API request is made with tampered token
Then signature validation fails
And 401 Unauthorized is returned
And error message is "Invalid token signature"
And security event is logged

Given no Authorization header
When API request is made to protected endpoint /api/v1/settings
Then middleware returns 401 Unauthorized
And error message is "Missing authorization header"
And request is rejected immediately

Given user successfully authenticated
When token refresh is requested before expiry
Then new token is generated with fresh 24-hour expiry
And old token is invalidated (if using token rotation)
And new token contains same user claims with updated iat/exp
```

### API-003: API Key Authentication

```go
// internal/api/middleware/apikey.go
package middleware

func APIKeyAuth(keyService *service.APIKeyService) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            apiKey := c.Request().Header.Get("X-API-Key")
            if apiKey == "" {
                return next(c) // Fall through to JWT
            }

            key, err := keyService.Validate(apiKey)
            if err != nil {
                return echo.NewHTTPError(401, "Invalid API key")
            }

            c.Set("user_id", key.UserID)
            c.Set("api_key", true)
            c.Set("permissions", key.Permissions)

            return next(c)
        }
    }
}
```

---

## 3.2 Admin API Endpoints [MVP]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| AA-001 | Domain CRUD endpoints | [ ] | API-002, U-003 |
| AA-002 | User CRUD endpoints | [ ] | API-002, U-001 |
| AA-003 | Alias CRUD endpoints | [ ] | API-002, U-004 |
| AA-004 | Quota management endpoints | [ ] | API-002, U-006 |
| AA-005 | Statistics endpoints | [ ] | API-002 |
| AA-006 | Log retrieval endpoints | [ ] | API-002, F-010 |
| AA-007 | Queue management endpoints | [ ] | API-002, Q-002 |
| AA-008 | DKIM key management endpoints | [ ] | API-002, DK-007 |
| AA-009 | System health endpoints | [ ] | API-002 |
| AA-010 | Backup/restore endpoints | [ ] | API-002 |

### AA-001: Domain Handler

```go
// internal/api/handlers/domain_handler.go
package handlers

type DomainHandler struct {
    domainService *service.DomainService
}

type CreateDomainRequest struct {
    Name           string `json:"name" validate:"required,fqdn"`
    MaxUsers       int    `json:"max_users"`
    MaxMailboxSize int64  `json:"max_mailbox_size"`
    DefaultQuota   int64  `json:"default_quota"`
}

func (h *DomainHandler) Create(c echo.Context) error {
    var req CreateDomainRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(400, "Invalid request body")
    }

    if err := c.Validate(&req); err != nil {
        return echo.NewHTTPError(400, err.Error())
    }

    domain := &domain.Domain{
        Name:           req.Name,
        MaxUsers:       req.MaxUsers,
        MaxMailboxSize: req.MaxMailboxSize,
        DefaultQuota:   req.DefaultQuota,
        Status:         "active",
    }

    if err := h.domainService.Create(domain); err != nil {
        return echo.NewHTTPError(500, err.Error())
    }

    return c.JSON(201, domain)
}

func (h *DomainHandler) List(c echo.Context) error {
    offset, _ := strconv.Atoi(c.QueryParam("offset"))
    limit, _ := strconv.Atoi(c.QueryParam("limit"))
    if limit == 0 {
        limit = 20
    }

    domains, total, err := h.domainService.List(offset, limit)
    if err != nil {
        return echo.NewHTTPError(500, err.Error())
    }

    return c.JSON(200, map[string]interface{}{
        "data":   domains,
        "total":  total,
        "offset": offset,
        "limit":  limit,
    })
}
```

### AA-005: Statistics Endpoint

```go
// internal/api/handlers/stats_handler.go
package handlers

type StatsHandler struct {
    statsService *service.StatsService
}

type SystemStats struct {
    TotalDomains      int64 `json:"total_domains"`
    TotalUsers        int64 `json:"total_users"`
    TotalMessages     int64 `json:"total_messages"`
    StorageUsed       int64 `json:"storage_used"`
    QueueSize         int64 `json:"queue_size"`
    ActiveConnections int64 `json:"active_connections"`
    MessagesToday     int64 `json:"messages_today"`
    SpamBlocked       int64 `json:"spam_blocked"`
    VirusBlocked      int64 `json:"virus_blocked"`
}

func (h *StatsHandler) GetSystemStats(c echo.Context) error {
    stats, err := h.statsService.GetSystemStats()
    if err != nil {
        return echo.NewHTTPError(500, err.Error())
    }
    return c.JSON(200, stats)
}

func (h *StatsHandler) GetDomainStats(c echo.Context) error {
    domainID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    stats, err := h.statsService.GetDomainStats(domainID)
    if err != nil {
        return echo.NewHTTPError(500, err.Error())
    }
    return c.JSON(200, stats)
}
```

---

## 3.3 Admin Web UI [MVP]

| ID | Task | Status | Dependencies | Framework |
|----|------|--------|--------------|-----------|
| AUI-001 | Set up Vue.js 3 project with Vite | [ ] | - | Vue 3 + Vite |
| AUI-002 | Admin authentication flow | [ ] | AUI-001, API-002 |
| AUI-003 | Domain management UI | [ ] | AUI-002, AA-001 |
| AUI-004 | User management UI | [ ] | AUI-002, AA-002 |
| AUI-005 | Alias management UI | [ ] | AUI-002, AA-003 |
| AUI-006 | Quota visualization | [ ] | AUI-002, AA-004 |
| AUI-007 | Real-time statistics dashboard | [ ] | AUI-002, AA-005 |
| AUI-008 | Log viewer with filtering | [ ] | AUI-002, AA-006 |
| AUI-009 | Queue management interface | [ ] | AUI-002, AA-007 |
| AUI-010 | DKIM/SPF/DMARC settings per domain | [ ] | AUI-002, AA-008 |
| AUI-011 | TLS certificate status | [ ] | AUI-002, LE-004 |
| AUI-012 | System health monitoring | [ ] | AUI-002, AA-009 |
| AUI-013 | Role-based access control | [ ] | AUI-002 |

### AUI-001: Vue Project Structure

```
web/admin/
├── src/
│   ├── main.ts
│   ├── App.vue
│   ├── router/
│   │   └── index.ts
│   ├── stores/
│   │   ├── auth.ts
│   │   ├── domains.ts
│   │   └── users.ts
│   ├── api/
│   │   ├── client.ts
│   │   ├── domains.ts
│   │   └── users.ts
│   ├── views/
│   │   ├── LoginView.vue
│   │   ├── DashboardView.vue
│   │   ├── DomainsView.vue
│   │   ├── UsersView.vue
│   │   ├── AliasesView.vue
│   │   ├── QueueView.vue
│   │   ├── LogsView.vue
│   │   └── SettingsView.vue
│   ├── components/
│   │   ├── common/
│   │   │   ├── DataTable.vue
│   │   │   ├── Modal.vue
│   │   │   ├── Pagination.vue
│   │   │   └── Toast.vue
│   │   ├── dashboard/
│   │   │   ├── StatsCard.vue
│   │   │   └── ActivityChart.vue
│   │   ├── domains/
│   │   │   ├── DomainForm.vue
│   │   │   └── DKIMSettings.vue
│   │   └── users/
│   │       └── UserForm.vue
│   └── layouts/
│       ├── AdminLayout.vue
│       └── AuthLayout.vue
├── package.json
├── vite.config.ts
└── tsconfig.json
```

### AUI-002: Authentication Store

```typescript
// web/admin/src/stores/auth.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, logout, refreshToken } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<User | null>(null)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function doLogin(email: string, password: string) {
    const response = await login(email, password)
    token.value = response.token
    user.value = response.user
    localStorage.setItem('token', response.token)
  }

  async function doLogout() {
    await logout()
    token.value = null
    user.value = null
    localStorage.removeItem('token')
  }

  return { token, user, isAuthenticated, isAdmin, doLogin, doLogout }
})
```

### AUI-007: Dashboard Component

```vue
<!-- web/admin/src/views/DashboardView.vue -->
<template>
  <div class="dashboard">
    <h1>Dashboard</h1>

    <div class="stats-grid">
      <StatsCard
        title="Total Domains"
        :value="stats.totalDomains"
        icon="globe"
      />
      <StatsCard
        title="Total Users"
        :value="stats.totalUsers"
        icon="users"
      />
      <StatsCard
        title="Messages Today"
        :value="stats.messagesToday"
        icon="mail"
      />
      <StatsCard
        title="Queue Size"
        :value="stats.queueSize"
        icon="inbox"
        :alert="stats.queueSize > 100"
      />
      <StatsCard
        title="Storage Used"
        :value="formatBytes(stats.storageUsed)"
        icon="database"
      />
      <StatsCard
        title="Spam Blocked"
        :value="stats.spamBlocked"
        icon="shield"
      />
    </div>

    <div class="charts-row">
      <ActivityChart :data="activityData" />
      <QueueChart :data="queueData" />
    </div>

    <div class="recent-activity">
      <h2>Recent Activity</h2>
      <ActivityLog :entries="recentLogs" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getSystemStats, getActivityData } from '@/api/stats'

const stats = ref({})
const activityData = ref([])
const recentLogs = ref([])

onMounted(async () => {
  stats.value = await getSystemStats()
  activityData.value = await getActivityData()
})
</script>
```

---

## 3.4 User Self-Service Portal [MVP]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| USP-001 | User authentication API | [ ] | API-002 |
| USP-002 | Password change endpoint | [ ] | USP-001 |
| USP-003 | 2FA setup endpoint | [ ] | USP-001, AU-001 |
| USP-004 | User alias management | [ ] | USP-001 |
| USP-005 | Quota usage display | [ ] | USP-001, U-006 |
| USP-006 | Forwarding rules API | [ ] | USP-001 |
| USP-007 | Session management | [ ] | USP-001, F-022 |

### USP-002: Password Change

```go
// internal/api/handlers/user_self_handler.go
package handlers

type ChangePasswordRequest struct {
    CurrentPassword string `json:"current_password" validate:"required"`
    NewPassword     string `json:"new_password" validate:"required,min=8"`
}

func (h *UserSelfHandler) ChangePassword(c echo.Context) error {
    userID := c.Get("user_id").(int64)

    var req ChangePasswordRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(400, "Invalid request")
    }

    if err := h.userService.ChangePassword(userID, req.CurrentPassword, req.NewPassword); err != nil {
        if err == service.ErrInvalidPassword {
            return echo.NewHTTPError(400, "Current password is incorrect")
        }
        return echo.NewHTTPError(500, err.Error())
    }

    return c.JSON(200, map[string]string{"message": "Password changed successfully"})
}
```

### USP-003: 2FA Setup

```go
// internal/api/handlers/user_self_handler.go

type TOTPSetupResponse struct {
    Secret     string `json:"secret"`
    QRCodeURL  string `json:"qr_code_url"`
}

type TOTPVerifyRequest struct {
    Code string `json:"code" validate:"required,len=6"`
}

func (h *UserSelfHandler) SetupTOTP(c echo.Context) error {
    userID := c.Get("user_id").(int64)
    user, _ := h.userService.GetByID(userID)

    setup, err := h.totpService.GenerateSecret(user.Email)
    if err != nil {
        return echo.NewHTTPError(500, err.Error())
    }

    // Store pending secret
    h.userService.SetPendingTOTPSecret(userID, setup.Secret)

    return c.JSON(200, TOTPSetupResponse{
        Secret:    setup.Secret,
        QRCodeURL: setup.URL,
    })
}

func (h *UserSelfHandler) VerifyTOTP(c echo.Context) error {
    userID := c.Get("user_id").(int64)

    var req TOTPVerifyRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(400, "Invalid request")
    }

    if err := h.userService.EnableTOTP(userID, req.Code); err != nil {
        return echo.NewHTTPError(400, "Invalid code")
    }

    return c.JSON(200, map[string]string{"message": "2FA enabled successfully"})
}
```

---

## 3.5 User Portal UI [MVP]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| UP-001 | User portal Vue.js project | [ ] | AUI-001 |
| UP-002 | Password change UI | [ ] | UP-001, USP-002 |
| UP-003 | 2FA setup wizard | [ ] | UP-001, USP-003 |
| UP-004 | Alias management UI | [ ] | UP-001, USP-004 |
| UP-005 | Quota display widget | [ ] | UP-001, USP-005 |
| UP-006 | Forwarding rules editor | [ ] | UP-001, USP-006 |
| UP-007 | Spam quarantine viewer | [ ] | UP-001, AS-004 |

### UP-003: 2FA Setup Wizard

```vue
<!-- web/portal/src/components/TwoFactorSetup.vue -->
<template>
  <div class="two-factor-setup">
    <div v-if="step === 1" class="step-intro">
      <h2>Enable Two-Factor Authentication</h2>
      <p>Protect your account with an additional layer of security.</p>
      <button @click="initSetup" class="btn-primary">Get Started</button>
    </div>

    <div v-if="step === 2" class="step-scan">
      <h2>Scan QR Code</h2>
      <p>Scan this QR code with your authenticator app:</p>
      <div class="qr-code">
        <img :src="qrCodeUrl" alt="2FA QR Code" />
      </div>
      <p class="manual-entry">
        Or enter this code manually: <code>{{ secret }}</code>
      </p>
      <button @click="step = 3" class="btn-primary">Next</button>
    </div>

    <div v-if="step === 3" class="step-verify">
      <h2>Verify Setup</h2>
      <p>Enter the 6-digit code from your authenticator app:</p>
      <input
        v-model="verifyCode"
        type="text"
        maxlength="6"
        placeholder="000000"
        class="code-input"
      />
      <button @click="verifySetup" class="btn-primary">Verify & Enable</button>
      <p v-if="error" class="error">{{ error }}</p>
    </div>

    <div v-if="step === 4" class="step-complete">
      <h2>Two-Factor Authentication Enabled!</h2>
      <p>Your account is now protected with 2FA.</p>
      <div class="backup-codes">
        <h3>Backup Codes</h3>
        <p>Save these codes in a safe place:</p>
        <ul>
          <li v-for="code in backupCodes" :key="code">{{ code }}</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { setupTOTP, verifyTOTP } from '@/api/auth'

const step = ref(1)
const secret = ref('')
const qrCodeUrl = ref('')
const verifyCode = ref('')
const backupCodes = ref<string[]>([])
const error = ref('')

async function initSetup() {
  const response = await setupTOTP()
  secret.value = response.secret
  qrCodeUrl.value = response.qr_code_url
  step.value = 2
}

async function verifySetup() {
  try {
    const response = await verifyTOTP(verifyCode.value)
    backupCodes.value = response.backup_codes
    step.value = 4
  } catch (e) {
    error.value = 'Invalid code. Please try again.'
  }
}
</script>
```

---

## 3.6 Let's Encrypt Integration [MVP]

| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| LE-001 | ACME client integration | [ ] | F-002 | `go-acme/lego/v4` |
| LE-002 | Cloudflare DNS challenge | [ ] | LE-001 |
| LE-003 | Automatic certificate renewal | [ ] | LE-002 |
| LE-004 | Certificate storage and loading | [ ] | LE-001, T-001 |
| LE-005 | Per-domain certificate support | [ ] | LE-004 |

### LE-001: ACME Client

```go
// internal/tls/acme/client.go
package acme

import (
    "github.com/go-acme/lego/v4/certcrypto"
    "github.com/go-acme/lego/v4/certificate"
    "github.com/go-acme/lego/v4/lego"
    "github.com/go-acme/lego/v4/registration"
)

type ACMEClient struct {
    client     *lego.Client
    user       *ACMEUser
    certStore  *CertificateStore
}

type ACMEUser struct {
    Email        string
    Registration *registration.Resource
    key          crypto.PrivateKey
}

func (u *ACMEUser) GetEmail() string                        { return u.Email }
func (u *ACMEUser) GetRegistration() *registration.Resource { return u.Registration }
func (u *ACMEUser) GetPrivateKey() crypto.PrivateKey        { return u.key }

func NewACMEClient(email string, certStore *CertificateStore) (*ACMEClient, error) {
    // Generate user private key
    privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

    user := &ACMEUser{
        Email: email,
        key:   privateKey,
    }

    config := lego.NewConfig(user)
    config.CADirURL = lego.LEDirectoryProduction
    config.Certificate.KeyType = certcrypto.RSA2048

    client, err := lego.NewClient(config)
    if err != nil {
        return nil, err
    }

    // Register
    reg, err := client.Registration.Register(registration.RegisterOptions{
        TermsOfServiceAgreed: true,
    })
    if err != nil {
        return nil, err
    }
    user.Registration = reg

    return &ACMEClient{
        client:    client,
        user:      user,
        certStore: certStore,
    }, nil
}
```

### LE-002: Cloudflare DNS Challenge

```go
// internal/tls/acme/cloudflare.go
package acme

import (
    "github.com/go-acme/lego/v4/providers/dns/cloudflare"
)

func (c *ACMEClient) SetupCloudflareDNS(apiToken string) error {
    config := cloudflare.NewDefaultConfig()
    config.AuthToken = apiToken

    provider, err := cloudflare.NewDNSProviderConfig(config)
    if err != nil {
        return err
    }

    return c.client.Challenge.SetDNS01Provider(provider)
}

func (c *ACMEClient) ObtainCertificate(domains []string) (*tls.Certificate, error) {
    request := certificate.ObtainRequest{
        Domains: domains,
        Bundle:  true,
    }

    certificates, err := c.client.Certificate.Obtain(request)
    if err != nil {
        return nil, err
    }

    // Store certificate
    c.certStore.Store(domains[0], certificates)

    // Parse into tls.Certificate
    cert, err := tls.X509KeyPair(certificates.Certificate, certificates.PrivateKey)
    if err != nil {
        return nil, err
    }

    return &cert, nil
}
```

### LE-003: Auto-Renewal

```go
// internal/tls/acme/renewer.go
package acme

type Renewer struct {
    client    *ACMEClient
    certStore *CertificateStore
    logger    logger.Logger
}

func (r *Renewer) Start(ctx context.Context) {
    ticker := time.NewTicker(24 * time.Hour) // Check daily

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            r.checkAndRenew()
        }
    }
}

func (r *Renewer) checkAndRenew() {
    certs := r.certStore.ListExpiringSoon(30 * 24 * time.Hour) // 30 days

    for _, cert := range certs {
        r.logger.Info("Renewing certificate", zap.String("domain", cert.Domain))

        renewed, err := r.client.RenewCertificate(cert)
        if err != nil {
            r.logger.Error("Failed to renew certificate",
                zap.String("domain", cert.Domain),
                zap.Error(err))
            continue
        }

        r.certStore.Update(cert.Domain, renewed)
        r.logger.Info("Certificate renewed", zap.String("domain", cert.Domain))
    }
}
```

---

## 3.7 Setup Wizard [MVP]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| SW-001 | First-run detection | [ ] | F-020 |
| SW-002 | First domain configuration | [ ] | SW-001, U-003 |
| SW-003 | First admin user creation | [ ] | SW-002, U-001 |
| SW-004 | TLS certificate setup flow | [ ] | SW-002, LE-001 |
| SW-005 | DKIM key generation UI | [ ] | SW-002, DK-001 |
| SW-006 | DNS record suggestions | [ ] | SW-002 |
| SW-007 | Pre-flight checks (ports, services) | [ ] | SW-001 |

### SW-001: First-Run Detection

```go
// internal/setup/wizard.go
package setup

type Wizard struct {
    db       *database.DB
    config   *config.Config
}

func (w *Wizard) IsFirstRun() bool {
    var count int
    w.db.QueryRow("SELECT COUNT(*) FROM domains").Scan(&count)
    return count == 0
}

func (w *Wizard) GetSetupStatus() *SetupStatus {
    return &SetupStatus{
        IsFirstRun:        w.IsFirstRun(),
        HasDomain:         w.hasDomain(),
        HasAdminUser:      w.hasAdminUser(),
        HasTLSCertificate: w.hasTLSCertificate(),
        HasDKIMKey:        w.hasDKIMKey(),
    }
}
```

### SW-006: DNS Record Suggestions

```go
// internal/setup/dns_helper.go
package setup

type DNSRecord struct {
    Type    string `json:"type"`
    Name    string `json:"name"`
    Value   string `json:"value"`
    Purpose string `json:"purpose"`
}

func (w *Wizard) GetDNSSuggestions(domain string) []DNSRecord {
    hostname, _ := os.Hostname()
    publicIP := getPublicIP()

    records := []DNSRecord{
        {
            Type:    "MX",
            Name:    domain,
            Value:   fmt.Sprintf("10 %s.", hostname),
            Purpose: "Mail exchange record - directs email to your server",
        },
        {
            Type:    "A",
            Name:    hostname,
            Value:   publicIP,
            Purpose: "Points your mail server hostname to this server",
        },
        {
            Type:    "TXT",
            Name:    domain,
            Value:   fmt.Sprintf("v=spf1 mx a:%s -all", hostname),
            Purpose: "SPF record - authorizes your server to send email",
        },
        {
            Type:    "TXT",
            Name:    "_dmarc." + domain,
            Value:   "v=DMARC1; p=quarantine; rua=mailto:postmaster@" + domain,
            Purpose: "DMARC policy - protects against email spoofing",
        },
    }

    // Add DKIM record if key exists
    if dkimKey := w.getDKIMKey(domain); dkimKey != nil {
        records = append(records, DNSRecord{
            Type:    "TXT",
            Name:    dkimKey.Selector + "._domainkey." + domain,
            Value:   dkimKey.DNSRecord(),
            Purpose: "DKIM public key - verifies email authenticity",
        })
    }

    return records
}
```

### SW-007: Pre-flight Checks

```go
// internal/setup/preflight.go
package setup

type PreflightCheck struct {
    Name    string `json:"name"`
    Status  string `json:"status"` // pass, fail, warn
    Message string `json:"message"`
}

func (w *Wizard) RunPreflightChecks() []PreflightCheck {
    checks := []PreflightCheck{}

    // Check port 25
    checks = append(checks, w.checkPort(25, "SMTP (port 25)"))

    // Check port 587
    checks = append(checks, w.checkPort(587, "SMTP Submission (port 587)"))

    // Check port 465
    checks = append(checks, w.checkPort(465, "SMTPS (port 465)"))

    // Check port 143
    checks = append(checks, w.checkPort(143, "IMAP (port 143)"))

    // Check port 993
    checks = append(checks, w.checkPort(993, "IMAPS (port 993)"))

    // Check ClamAV
    checks = append(checks, w.checkClamAV())

    // Check SpamAssassin
    checks = append(checks, w.checkSpamAssassin())

    // Check DNS resolution
    checks = append(checks, w.checkDNS())

    return checks
}

func (w *Wizard) checkPort(port int, name string) PreflightCheck {
    ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        return PreflightCheck{
            Name:    name,
            Status:  "fail",
            Message: fmt.Sprintf("Port %d is not available: %s", port, err),
        }
    }
    ln.Close()
    return PreflightCheck{
        Name:    name,
        Status:  "pass",
        Message: fmt.Sprintf("Port %d is available", port),
    }
}
```

---

## Acceptance Criteria

### API
- [ ] All endpoints return correct JSON responses
- [ ] JWT authentication works
- [ ] API key authentication works
- [ ] Rate limiting prevents abuse
- [ ] CORS configured correctly
- [ ] OpenAPI documentation generated

### Admin UI
- [ ] Login/logout works
- [ ] Domain CRUD operations work
- [ ] User CRUD operations work
- [ ] Dashboard shows real-time stats
- [ ] Log viewer displays and filters logs
- [ ] Queue management interface works
- [ ] DKIM/SPF/DMARC settings configurable

### User Portal
- [ ] Password change works
- [ ] 2FA setup wizard works
- [ ] Alias management works
- [ ] Quota display accurate
- [ ] Spam quarantine viewable

### Let's Encrypt
- [ ] Certificates obtained automatically
- [ ] Cloudflare DNS challenge works
- [ ] Auto-renewal runs on schedule
- [ ] Certificates loaded by SMTP/IMAP

### Setup Wizard
- [ ] First-run detection works
- [ ] Domain creation works
- [ ] Admin user creation works
- [ ] TLS setup flow works
- [ ] DNS suggestions accurate
- [ ] Pre-flight checks identify issues

---

## Go Dependencies for Phase 3

```go
// Additional go.mod entries
require (
    github.com/labstack/echo/v4 v4.11.4
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/go-playground/validator/v10 v10.18.0
    github.com/go-acme/lego/v4 v4.15.0
    github.com/swaggo/echo-swagger v1.4.1
)
```

---

## Frontend Dependencies

```json
// web/admin/package.json
{
  "dependencies": {
    "vue": "^3.4.0",
    "vue-router": "^4.2.0",
    "pinia": "^2.1.0",
    "axios": "^1.6.0",
    "@vueuse/core": "^10.7.0"
  },
  "devDependencies": {
    "vite": "^5.0.0",
    "typescript": "^5.3.0",
    "@vitejs/plugin-vue": "^5.0.0",
    "tailwindcss": "^3.4.0"
  }
}
```

---

## Next Phase

After completing Phase 3, proceed to [TASKS4.md](TASKS4.md) - CalDAV/CardDAV.

**Note**: Phases 1-3 constitute the MVP milestone.
