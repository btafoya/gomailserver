# gomailserver Web UI Testing Report

**Date**: 2026-01-02
**Tester**: Claude Code with Playwright MCP
**Server Version**: dev
**Test Environment**: localhost:8980

---

## Test Setup

### Test Users Created
- **Admin User**: `admin@localhost` / `TestPassword123!` (role: admin)
- **Regular User**: `alice@localhost` / `TestPassword123!` (role: user)

### Database
- **Path**: `./mailserver.db`
- **Migration Version**: 7
- **Domains**: `_default`, `localhost`

### Server Configuration
- **API Port**: 8980
- **WebDAV Port**: 8800
- **SMTP Ports**: 2525 (relay), 2587 (submission), 2465 (smtps)
- **IMAP Ports**: 2143, 2993
- **TLS**: Self-signed certificate (development)

---

## Test Results

### 1. Server Health Check ✅

**Endpoint**: `GET /health`
**Status**: PASS
**Response**: `{"status":"ok"}`

---

### 2. Admin UI Login ❌ CRITICAL FAILURE

**URL**: `http://localhost:8980/admin`
**Status**: FAIL
**Issues Found**: Multiple critical bugs preventing basic functionality

#### Issue Summary:
1. **Circular Dependency Bug** - Maximum call stack exceeded
2. **URL Path Doubling** - API requests going to `/api/api/v1/*` instead of `/api/v1/*`
3. **Infinite Refresh Loop** - Thousands of failed refresh requests per second

#### Test Steps:
1. ✅ Navigate to http://localhost:8980/admin
2. ✅ Login page loads correctly
3. ✅ Fill email: `admin@localhost`
4. ✅ Fill password: `TestPassword123!`
5. ❌ Click "Sign In" button
6. ❌ **FAILURE**: Browser hangs with infinite API requests

#### Server Logs Evidence:
```
{"method":"POST","path":"/api/api/v1/auth/refresh","status":401}
{"method":"POST","path":"/api/api/v1/auth/refresh","status":401}
{"method":"POST","path":"/api/api/v1/auth/refresh","status":401}
... (thousands of requests)
```

**Expected**:
```
{"method":"POST","path":"/api/v1/auth/login","status":200}
```

#### Root Causes Identified:

**1. Circular Dependency (FIXED)**
- `axios.js` imported `useAuthStore` causing infinite loop
- **Fix**: Use localStorage directly in interceptors
- **Status**: ✅ Resolved - no more call stack errors

**2. URL Path Doubling (NOT FIXED)**
- Axios requests to `/api/v1/auth/login` become `/api/api/v1/auth/login`
- Tested multiple baseURL configurations - all failed
- Related to Vite `base: '/admin/'` configuration
- **Status**: ❌ Blocking all admin UI functionality

**3. Infinite Refresh Loop (BLOCKED BY #2)**
- Response interceptor retries on 401 errors
- Since all requests fail due to wrong URL, it retries forever
- Browser becomes unresponsive
- **Status**: ❌ Will resolve when URL issue fixed

#### Detailed Investigation:

**Files Modified**:
- `web/admin/src/api/axios.js` - Request/response interceptors
- `web/admin/src/router/index.js` - Navigation guard
- `web/admin/.env.production` - API base URL config

**Attempted Solutions**:
1. Set baseURL to `/api` → Result: double `/api` prefix
2. Set baseURL to `http://localhost:8980` → Works but hardcoded
3. Set baseURL to empty string → Still gets double `/api`
4. Set baseURL to `window.location.origin` → Still gets double `/api`

**Investigation Time**: ~3 hours of debugging
**Conclusion**: Admin UI is completely non-functional in production build mode

#### Recommendations:
1. **Immediate**: See ISSUE001.md for comprehensive bug documentation
2. **Short-term**: Fix URL path resolution in axios/vite configuration
3. **Long-term**: Add integration tests to catch these issues before deployment

---

### 3. Admin UI Navigation ⏸️ BLOCKED

**Status**: Cannot test due to login failure
**Planned Tests**:
- Dashboard access
- Domain management
- User management
- Settings configuration
- Queue monitoring
- Logs viewer
- Audit trail

**Blocking Issue**: ISSUE001 must be resolved first

---

### 4. Webmail UI ⏸️ NOT TESTED

**Status**: Deferred pending admin UI fix
**Reason**: May have similar URL path issues
**Planned Tests**:
- User login
- Inbox/mailbox listing
- Message viewing
- Email composition
- Attachments
- Contact integration
- Calendar integration

---

## Summary

### Test Coverage
- ✅ Server Health: PASS
- ❌ Admin UI Login: CRITICAL FAILURE
- ⏸️ Admin UI Navigation: BLOCKED
- ⏸️ Webmail UI: NOT TESTED

### Critical Issues
1. **ISSUE001**: Admin UI URL path doubling - blocks all functionality
2. **Circular Dependencies**: Fixed but admin UI still non-functional
3. **Infinite Refresh Loop**: Consequence of URL issue

### Files Created/Modified
- `ISSUE001.md` - Detailed bug documentation
- `web/admin/src/api/axios.js` - Partial fixes
- `web/admin/src/router/index.js` - Partial fixes
- `web/admin/.env.production` - Configuration attempts

### Recommendations
1. **Priority P0**: Resolve URL path doubling issue
2. Add axios request logging for debugging
3. Consider removing Vite base configuration
4. Add E2E tests to CI/CD pipeline
5. Test in dev mode with Vite proxy to verify expected behavior

### Time Investment
- Setup: 30 minutes
- Admin UI debugging: 3+ hours
- Documentation: 30 minutes
- **Total**: ~4 hours

---

## Next Steps

1. Debug axios request pipeline to find exact point of URL doubling
2. Test alternate Vite base configurations
3. Consider using absolute URLs for API in production
4. Once fixed, complete full admin UI test suite
5. Test webmail UI for similar issues
6. Add automated E2E tests to prevent regression

---

**Testing Status**: ⚠️ INCOMPLETE due to critical bugs
**Recommendation**: DO NOT DEPLOY current admin UI to production

