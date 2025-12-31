# ISSUE004: Phase 3 Web Interfaces - Admin & Portal

## Status: In Progress
## Priority: High (MVP Critical)
## Phase: 3 - Web Interfaces
## Started: 2025-12-30

## Summary

Implementing Phase 3 - Web Interfaces for gomailserver, including REST API foundation, Admin Web UI, and User Self-Service Portal. This phase delivers the management and user-facing interfaces critical for the MVP.

## Phase 3 Requirements (from PR.md)

### 3.1 REST API Foundation [MVP]
- JSON-based REST API
- Domain management (CRUD)
- User management (CRUD)
- Alias management (CRUD)
- Quota management
- Statistics and monitoring
- Log retrieval
- Queue management
- Authentication (API keys, JWT tokens)
- Rate limiting
- OpenAPI/Swagger documentation

### 3.2 Admin Web Interface [MVP]
- Modern, responsive web UI (Vue.js 3 + Tailwind CSS)
- Domain management (CRUD)
- User management (CRUD)
- Alias management (CRUD)
- Quota management and visualization
- Real-time statistics dashboard
- Log viewer with filtering
- Queue management interface
- Security settings (DKIM, SPF, DMARC per domain)
- TLS certificate status and management
- Backup/restore interface
- System health monitoring
- Role-based access control (admin/read-only)

### 3.3 User Self-Service Portal [MVP]
- Password change
- 2FA setup (TOTP)
- Alias management (create/delete own aliases)
- Quota usage display
- Forwarding rules
- Sieve filter management (visual editor)
- PGP key management
- Spam quarantine review
- Session management
- Email signature settings

### 3.4 System Features [MVP]
- Let's Encrypt ACME integration (Cloudflare DNS)
- Initial setup wizard
- Dashboard and statistics
- System monitoring interface

## Implementation Plan

### Task Breakdown

#### Epic 1: REST API Foundation
**Story 1.1**: Core API Infrastructure
- [x] Create API package structure
- [ ] Implement JWT authentication middleware
- [ ] Implement API key authentication
- [ ] Create rate limiting middleware
- [ ] Implement error handling middleware
- [ ] Create response helpers
- [ ] Set up CORS configuration

**Story 1.2**: Domain Management API
- [ ] POST /api/v1/domains - Create domain
- [ ] GET /api/v1/domains - List domains
- [ ] GET /api/v1/domains/:id - Get domain details
- [ ] PUT /api/v1/domains/:id - Update domain
- [ ] DELETE /api/v1/domains/:id - Delete domain
- [ ] POST /api/v1/domains/:id/dkim - Generate DKIM keys

**Story 1.3**: User Management API
- [ ] POST /api/v1/users - Create user
- [ ] GET /api/v1/users - List users
- [ ] GET /api/v1/users/:id - Get user details
- [ ] PUT /api/v1/users/:id - Update user
- [ ] DELETE /api/v1/users/:id - Delete user
- [ ] POST /api/v1/users/:id/password - Reset password

**Story 1.4**: Alias Management API
- [ ] POST /api/v1/aliases - Create alias
- [ ] GET /api/v1/aliases - List aliases
- [ ] GET /api/v1/aliases/:id - Get alias
- [ ] DELETE /api/v1/aliases/:id - Delete alias

**Story 1.5**: Statistics & Monitoring API
- [ ] GET /api/v1/stats/dashboard - Dashboard statistics
- [ ] GET /api/v1/stats/domains/:id - Domain statistics
- [ ] GET /api/v1/stats/users/:id - User statistics
- [ ] GET /api/v1/logs - Retrieve logs with filtering
- [ ] GET /api/v1/health - System health check

**Story 1.6**: Queue Management API
- [ ] GET /api/v1/queue - List queued messages
- [ ] GET /api/v1/queue/:id - Get queue item
- [ ] POST /api/v1/queue/:id/retry - Retry message
- [ ] DELETE /api/v1/queue/:id - Remove from queue

#### Epic 2: Admin Web UI (Vue.js 3 + Tailwind)
**Story 2.1**: UI Foundation
- [ ] Initialize Vue.js 3 project with Vite
- [ ] Install and configure Tailwind CSS
- [ ] Install shadcn-vue component library
- [ ] Create layout components (sidebar, header, footer)
- [ ] Implement routing (Vue Router)
- [ ] Create authentication views (login, logout)

**Story 2.2**: Domain Management Views
- [ ] Domain list view with search/filter
- [ ] Domain create form
- [ ] Domain edit form with tabs (general, security, quotas)
- [ ] DKIM key generation interface
- [ ] DNS record display/copy
- [ ] Domain deletion with confirmation

**Story 2.3**: User Management Views
- [ ] User list view with pagination
- [ ] User create form
- [ ] User edit form
- [ ] Password reset interface
- [ ] Quota visualization
- [ ] User status toggle (active/disabled)

**Story 2.4**: Alias Management Views
- [ ] Alias list view
- [ ] Alias create form
- [ ] Alias edit form
- [ ] Bulk alias operations

**Story 2.5**: Dashboard View
- [ ] Real-time statistics cards (domains, users, messages, queue)
- [ ] Charts (message volume, storage usage)
- [ ] Recent activity feed
- [ ] System health indicators

**Story 2.6**: System Monitoring
- [ ] Log viewer with filtering
- [ ] Queue management interface
- [ ] TLS certificate status
- [ ] Service health dashboard
- [ ] Backup/restore interface

#### Epic 3: User Self-Service Portal
**Story 3.1**: Portal Foundation
- [ ] Create separate Vue.js 3 app for portal
- [ ] User authentication
- [ ] Portal layout components

**Story 3.2**: Account Management
- [ ] Profile view/edit
- [ ] Password change form
- [ ] 2FA setup (TOTP with QR code)
- [ ] Session management

**Story 3.3**: Alias Management
- [ ] User alias list
- [ ] Create personal alias
- [ ] Delete personal alias

**Story 3.4**: Sieve Filter Editor
- [ ] Visual rule builder
- [ ] Sieve script preview
- [ ] Test filter interface
- [ ] Filter activation

**Story 3.5**: Additional Features
- [ ] Quota usage display with chart
- [ ] Spam quarantine viewer
- [ ] PGP key upload/management
- [ ] Email signature editor
- [ ] Forwarding rules

#### Epic 4: Let's Encrypt ACME Integration
**Story 4.1**: ACME Client
- [ ] Integrate lego/certmagic ACME library
- [ ] Implement Cloudflare DNS provider
- [ ] Certificate request flow
- [ ] Certificate renewal worker
- [ ] Certificate storage in SQLite

**Story 4.2**: Admin Interface
- [ ] Certificate status view
- [ ] Manual certificate request
- [ ] Certificate renewal trigger
- [ ] Domain verification status

#### Epic 5: Initial Setup Wizard
**Story 5.1**: Wizard Flow
- [ ] Welcome screen
- [ ] First domain configuration
- [ ] First admin user creation
- [ ] DKIM key generation
- [ ] DNS record display
- [ ] ACME certificate setup
- [ ] Test email configuration

## Technology Stack

### Backend (Go)
- **HTTP Router**: `chi` or `gin`
- **Authentication**: `golang-jwt/jwt`
- **ACME Client**: `go-acme/lego` or `caddyserver/certmagic`
- **WebSocket**: For real-time updates
- **OpenAPI**: `swaggo/swag`

### Frontend (Admin UI)
- **Framework**: Vue.js 3 (Composition API)
- **Build Tool**: Vite
- **UI Library**: shadcn-vue + Tailwind CSS
- **State Management**: Pinia
- **HTTP Client**: Axios
- **Charts**: Chart.js or Apache ECharts
- **Forms**: Vee-Validate + Yup
- **Icons**: Lucide Icons

### Frontend (User Portal)
- Same stack as Admin UI (separate app)

## Database Schema Updates

### New Tables for Phase 3

#### API Keys Table
```sql
CREATE TABLE api_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key_hash TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    user_id INTEGER,
    domain_id INTEGER,
    permissions TEXT NOT NULL, -- JSON array
    rate_limit_per_hour INTEGER DEFAULT 1000,
    last_used_at DATETIME,
    expires_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
);
```

#### TLS Certificates Table
```sql
CREATE TABLE tls_certificates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL UNIQUE,
    certificate_pem TEXT NOT NULL,
    private_key_pem TEXT NOT NULL,
    issuer TEXT NOT NULL,
    not_before DATETIME NOT NULL,
    not_after DATETIME NOT NULL,
    acme_account_url TEXT,
    acme_order_url TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

#### Setup Wizard State Table
```sql
CREATE TABLE setup_wizard_state (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    step TEXT NOT NULL,
    completed INTEGER DEFAULT 0,
    data TEXT, -- JSON
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## Acceptance Criteria

### REST API
- [ ] All CRUD operations functional for domains, users, aliases
- [ ] JWT authentication working
- [ ] API key authentication working
- [ ] Rate limiting enforced
- [ ] OpenAPI documentation generated
- [ ] All endpoints returning proper HTTP status codes
- [ ] Error responses properly formatted

### Admin UI
- [ ] Responsive design (mobile, tablet, desktop)
- [ ] All CRUD forms working
- [ ] Real-time dashboard updates
- [ ] Log viewer with filtering
- [ ] Queue management functional
- [ ] DKIM key generation working
- [ ] TLS certificate display

### User Portal
- [ ] Password change working
- [ ] 2FA setup functional with QR code
- [ ] Alias management working
- [ ] Sieve filter editor functional
- [ ] Quota visualization accurate
- [ ] Spam quarantine viewer working

### ACME Integration
- [ ] Automatic certificate issuance via Cloudflare DNS
- [ ] Certificate renewal working
- [ ] Certificate storage in SQLite
- [ ] Multi-domain certificate support

### Setup Wizard
- [ ] First-run detection
- [ ] Wizard guides through domain setup
- [ ] Admin user created successfully
- [ ] DKIM keys generated
- [ ] DNS records displayed
- [ ] Test email sent successfully

## Testing Strategy

### Unit Tests
- [ ] API handler tests
- [ ] Middleware tests
- [ ] Service layer tests
- [ ] ACME client tests

### Integration Tests
- [ ] Full API flow tests
- [ ] Authentication flow tests
- [ ] CRUD operation tests
- [ ] ACME certificate issuance test

### E2E Tests (Playwright)
- [ ] Admin UI login flow
- [ ] Domain creation flow
- [ ] User creation flow
- [ ] Setup wizard flow
- [ ] User portal flows

## Security Considerations

- **API Authentication**: JWT tokens with short expiry (1 hour) and refresh tokens
- **API Authorization**: Role-based access control (admin, user)
- **CORS**: Restrict to configured origins
- **Rate Limiting**: Per IP and per API key
- **Input Validation**: All API inputs validated
- **SQL Injection Prevention**: Parameterized queries only
- **XSS Prevention**: Output encoding in UI
- **CSRF Protection**: CSRF tokens for state-changing operations

## Performance Targets

- API response time < 100ms for simple queries
- API response time < 500ms for complex queries
- Dashboard real-time updates via WebSocket
- UI initial load < 2 seconds
- UI interactions < 200ms

## Dependencies

### Go Dependencies
```
github.com/go-chi/chi/v5
github.com/golang-jwt/jwt/v5
github.com/go-acme/lego/v4
github.com/swaggo/swag
github.com/gorilla/websocket
```

### Frontend Dependencies
```json
{
  "vue": "^3.5.13",
  "vite": "^6.0.3",
  "tailwindcss": "^4.1.5",
  "shadcn-vue": "latest",
  "pinia": "^2.3.1",
  "vue-router": "^4.5.0",
  "axios": "^1.7.9",
  "chart.js": "^4.4.7",
  "vee-validate": "^4.14.11",
  "yup": "^1.6.1",
  "lucide-vue-next": "^0.468.0"
}
```

## Timeline Estimate

- **Epic 1** (REST API): 1.5 weeks
- **Epic 2** (Admin UI): 1.5 weeks
- **Epic 3** (User Portal): 1 week
- **Epic 4** (ACME): 0.5 weeks
- **Epic 5** (Setup Wizard): 0.5 weeks

**Total**: 5 weeks (matches PR.md estimate of 4-5 weeks)

## Current Status

- [x] Project structure created
- [x] Phase 3 issue file created
- [ ] REST API foundation
- [ ] Admin Web UI
- [ ] User Portal
- [ ] ACME integration
- [ ] Setup wizard

## Next Immediate Steps

1. Initialize Vue.js projects for admin and portal
2. Set up Go HTTP router and middleware
3. Implement JWT authentication
4. Create domain management API endpoints
5. Build admin UI domain management views

## Notes

Following autonomous work mode per CLAUDE.md:
- Proceeding with full implementation
- Using Vue.js 3 with shadcn-vue for modern, accessible UI
- Tailwind CSS for responsive design
- SQLite-first architecture maintained
- No "Generated with Claude Code" in commits
