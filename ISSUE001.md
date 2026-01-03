# ISSUE001: Admin UI Critical Bugs - Circular Dependencies and URL Path Issues

**Date**: 2026-01-02
**Status**: ⚠️ CRITICAL - Admin UI non-functional
**Priority**: P0
**Assignee**: Development Team

---

## Summary

The admin UI has multiple critical bugs preventing basic authentication:
1. **Circular dependency** causing maximum call stack exceeded errors
2. **URL path doubling** causing API requests to `/api/api/v1/*` instead of `/api/v1/*`
3. **Infinite refresh loop** when authentication fails

These issues make the admin UI completely non-functional in production build mode.

---

## Bugs Identified

### 1. Circular Dependency (PARTIALLY FIXED)

**Root Cause**: `axios.js` imported `useAuthStore` at module level, and the auth store imported `api` from `axios.js`. When axios interceptors tried to access the store during initialization, it created an infinite loop.

**Files Affected**:
- `web/admin/src/api/axios.js`
- `web/admin/src/router/index.js`

**Fix Applied**:
- ✅ Changed axios request interceptor to read token directly from localStorage instead of calling `useAuthStore()`
- ✅ Changed router navigation guard to check localStorage instead of calling `useAuthStore()`
- ✅ Made response interceptor import `useAuthStore` dynamically only when needed (401 errors)

**Status**: Fixed - no more call stack exceeded errors

---

### 2. API URL Path Doubling (NOT FIXED)

**Symptom**: API requests to `/api/v1/auth/login` result in server seeing `/api/api/v1/auth/login`

**Server Logs Show**:
```
{"method":"POST","path":"/api/api/v1/auth/refresh","status":401}
```

**Expected**:
```
{"method":"POST","path":"/api/v1/auth/refresh","status":200}
```

**Investigation**:
1. Vite config has `base: '/admin/'` for asset paths
2. Admin UI is served from `/admin/*` by Go server
3. Axios `baseURL` has been tried as:
   - `/api` → double `/api` in requests
   - `http://localhost:8980` → works in dev but hardcoded
   - Empty string → still results in double `/api`
   - `window.location.origin` → still results in double `/api`

**Root Cause**: Unknown - the `/api` prefix is being added somewhere in the request pipeline even with empty baseURL. Possibly related to Vite's `base: '/admin/'` configuration affecting how relative URLs are resolved.

**Files Affected**:
- `web/admin/vite.config.js` (line 8: `base: '/admin/'`)
- `web/admin/src/api/axios.js` (baseURL configuration)
- `web/admin/.env.production` (VITE_API_BASE_URL)

**Status**: ❌ Not Fixed - blocking all API requests

---

### 3. Infinite Refresh Loop (CONSEQUENCE OF #2)

**Symptom**: Browser makes thousands of `/api/api/v1/auth/refresh` requests per second

**Root Cause**:
1. Initial login attempt fails due to wrong URL (`/api/api/v1/auth/login` → 404/401)
2. Axios response interceptor catches 401 error
3. Interceptor tries to refresh token at `/api/api/v1/auth/refresh`
4. Refresh also fails with 401 (wrong URL)
5. Interceptor retries infinitely

**Files Affected**:
- `web/admin/src/api/axios.js` (response interceptor, lines 28-64)

**Status**: ❌ Blocked by Bug #2 - will resolve once URL path issue is fixed

---

## Attempted Fixes

### Axios BaseURL Configurations Tried:

1. **`baseURL: '/api'`** (original .env.production)
   - Result: `/api` + `/api/v1/auth/login` = `/api/api/v1/auth/login` ❌

2. **`baseURL: 'http://localhost:8980'`** (test config)
   - Result: Works but hardcoded, breaks in different environments ⚠️

3. **`baseURL: ''`** (empty string)
   - Expected: Use document root
   - Result: Still `/api/api/v1/auth/login` ❌

4. **`baseURL: window.location.origin`**
   - Expected: Use current origin
   - Result: Still `/api/api/v1/auth/login` ❌

### Other Attempts:

- Removed Vite proxy configuration (no effect in production build)
- Checked for `<base>` tag in HTML (none found)
- Verified admin handler strips `/admin` prefix correctly (works)
- Checked for middleware adding prefixes (none found)

---

## Next Steps to Resolve

### Option 1: Remove Vite Base Configuration
- Change `vite.config.js` base from `/admin/` to `/`
- Update asset paths manually or with custom plugin
- Risk: May break static asset loading

### Option 2: Use Absolute URLs in Production
- Set `VITE_API_BASE_URL=http://localhost:8980` for local testing
- Use environment-specific URLs for deployment
- Risk: Requires configuration per environment

### Option 3: Add Custom Axios URL Resolver
- Create a custom function to resolve API URLs correctly
- Handle the `/admin/` base path explicitly
- Risk: Added complexity

### Option 4: Debug Axios Request Pipeline
- Add extensive logging to axios interceptors
- Check `config.url` and `config.baseURL` at each step
- Identify exact point where `/api` is duplicated
- Risk: Time-consuming debugging

---

## Recommended Solution

**Immediate**: Use Option 2 (absolute URLs) to unblock testing
**Long-term**: Implement Option 4 to find root cause, then apply proper fix

---

## Testing Evidence

### Server Logs (Incorrect Requests):
```bash
tail -f /tmp/gomailserver-test.log | grep auth
# Shows: POST /api/api/v1/auth/refresh (401)
# Expected: POST /api/v1/auth/login (200)
```

### Browser Behavior:
- Login button click triggers infinite refresh loop
- Browser DevTools shows thousands of failed requests
- Page becomes unresponsive due to request volume

---

## Related Files

### Modified During Investigation:
- `web/admin/src/api/axios.js`
- `web/admin/src/router/index.js`
- `web/admin/src/stores/auth.js`
- `web/admin/.env.production`
- `web/admin/vite.config.js`

### Needing Review:
- `internal/admin/handler.go` (admin UI serving logic)
- `internal/api/router.go` (API route definitions)

---

## Impact

**Severity**: P0 - Complete admin UI failure
**Users Affected**: All admin users
**Workaround**: Use API directly or wait for fix
**Business Impact**: Cannot manage mail server through web UI

---

## Additional Notes

- Bug only affects production builds (embedded in Go binary)
- Dev mode with Vite proxy works correctly
- Webmail UI may have similar issues (untested)
- Issue has consumed significant debugging time (3+ hours)

---

## Commit History

Partial fixes committed in: [pending commit]
- Circular dependency fixes
- Router guard improvements
- Various baseURL attempts

Full fix pending root cause identification.
