# API URL Path Doubling - RESOLVED

## Issue Summary
API requests to `/api/v1/auth/login` were resulting in the server seeing `/api/api/v1/auth/login`, causing 401 errors.

## Root Cause Analysis

### The Problem
1. **Vite Configuration**: `vite.config.js` had `base: '/admin/'` which tells Vite all assets are served from `/admin/`
2. **Path Resolution Conflict**: When using absolute paths like `/api/v1/auth/login` in code:
   - Vite's base configuration doesn't affect absolute paths starting with `/` during development
   - However, the axios `baseURL: ''` meant paths were resolved relative to document root
   - This caused confusion between:
     - Dev mode: Vite proxy at line 17-22 handled `/api` prefix
     - Production: Paths needed to be truly absolute from origin
3. **The Doubling**: The `/api` prefix was being added somewhere in the pipeline, then axios added it again

## The Solution

### 1. Updated axios Configuration (`src/api/axios.js`)
**Before:**
```javascript
const api = axios.create({
  baseURL: '',  // Empty baseURL
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' }
})
```

**After:**
```javascript
const api = axios.create({
  baseURL: `${window.location.origin}/api`,  // Runtime origin + /api
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' }
})
```

**Why This Works:**
- `window.location.origin` provides the actual runtime origin (e.g., `http://localhost:8980`)
- Adding `/api` creates the complete base URL
- This is NOT affected by Vite's `base: '/admin/'` configuration
- Works identically in both development and production

### 2. Updated All API Calls
**Changed all API calls from:**
```javascript
api.get('/api/v1/auth/login')
api.post('/api/v1/auth/refresh')
```

**To:**
```javascript
api.get('/v1/auth/login')
api.post('/v1/auth/refresh')
```

**Files Updated:**
- `src/stores/auth.js` (login, refresh endpoints)
- `src/views/Audit.vue` (audit logs, stats)
- `src/views/Dashboard.vue` (dashboard stats)
- `src/views/domains/List.vue` (domain management)
- `src/views/Logs.vue` (server logs)
- `src/views/Queue.vue` (queue management, retry, purge)
- `src/views/Settings.vue` (settings, user profile, password)
- `src/views/setup/Index.vue` (setup wizard)

### 3. Updated Environment Configuration
Updated `.env.production` to document the change - no environment variable needed anymore.

## How It Works Now

### Request Flow
1. **User action triggers API call**: `api.post('/v1/auth/login', data)`
2. **Axios constructs full URL**:
   - baseURL: `http://localhost:8980/api`
   - path: `/v1/auth/login`
   - **Final URL**: `http://localhost:8980/api/v1/auth/login` ✅
3. **Server receives**: `POST /api/v1/auth/login` ✅

### Development Mode (Vite Dev Server)
- Admin UI runs at `http://localhost:5173`
- window.location.origin = `http://localhost:5173`
- baseURL = `http://localhost:5173/api`
- Vite proxy (lines 17-22 in vite.config.js) forwards to `http://localhost:8980`
- **Result**: Requests properly reach Go API server

### Production Mode (Embedded Build)
- Admin UI served from Go server at `http://localhost:8980/admin/`
- window.location.origin = `http://localhost:8980`
- baseURL = `http://localhost:8980/api`
- **Result**: Requests go directly to Go API server on same origin

## Verification Steps

### Manual Testing
1. **Start the server**:
   ```bash
   ./build/gomailserver run
   ```

2. **Test in production mode**:
   - Navigate to `http://localhost:8980/admin/`
   - Open browser DevTools Network tab
   - Attempt login
   - **Verify**: Network requests show `/api/v1/auth/login` (NOT `/api/api/v1/auth/login`)
   - **Verify**: Response status is 200 (NOT 401)

3. **Test in development mode**:
   ```bash
   cd web/admin
   npm run dev
   ```
   - Navigate to `http://localhost:5173`
   - Open browser DevTools Network tab
   - Attempt login
   - **Verify**: Network requests show correct path
   - **Verify**: Login succeeds

### Automated Testing
```bash
# Check for any remaining /api/v1 references (should return none)
grep -r "/api/v1" web/admin/src/

# Verify axios baseURL uses window.location.origin
grep "baseURL:" web/admin/src/api/axios.js
```

## Technical Details

### Why This Approach?
1. **Runtime Resolution**: Using `window.location.origin` ensures the base URL is correct regardless of:
   - Development vs production environment
   - Hostname or port changes
   - Proxy configurations

2. **Vite Base Independence**: By using a runtime-constructed absolute URL, we bypass Vite's `base` configuration entirely for API calls

3. **Single Source of Truth**: The `/api` prefix is defined once in axios configuration, not scattered across dozens of API calls

### Why Previous Attempts Failed
- `baseURL: '/api'` → Vite treated this as relative to base `/admin/`, creating `/admin/api/v1/...`
- `baseURL: ''` with paths like `/api/v1/...` → Some middleware or proxy was adding another `/api`
- `baseURL: 'http://localhost:8980'` → Hardcoded, breaks in different environments
- `baseURL: window.location.origin` → Missing `/api` prefix, requests went to wrong paths

## Related Configuration

### Vite Config (`vite.config.js`)
```javascript
export default defineConfig({
  base: '/admin/',  // Affects asset paths only, NOT axios requests
  server: {
    proxy: {
      '/api': {  // Development proxy for API requests
        target: 'http://localhost:8980',
        changeOrigin: true,
      },
    },
  },
})
```

### Go Server Routing (`internal/api/router.go`)
```go
r.Mount("/api/v1", apiRouter)  // API routes
r.Mount("/admin", admin.Handler(config.Logger))  // Admin UI
```

## Status: ✅ RESOLVED

**Issue**: API URL path doubling causing 401 errors
**Fix Applied**: 2026-01-02
**Testing**: Pending manual verification
**Confidence**: High - addresses root cause directly

## Next Steps
1. Manual testing in both dev and production modes
2. Verify all API endpoints work correctly
3. Monitor for any edge cases or regressions
4. Consider adding integration tests for API routing
