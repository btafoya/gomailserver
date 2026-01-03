# API Path Doubling Fix - Testing Checklist

## Quick Test (5 minutes)

### ✅ Production Mode Test
```bash
# 1. Build and start server
cd /home/btafoya/projects/gomailserver
make build
./build/gomailserver run

# 2. Open browser
# Navigate to: http://localhost:8980/admin/

# 3. Open DevTools (F12) → Network tab

# 4. Attempt login with test credentials

# 5. Check Network tab:
# ✅ Request URL should be: http://localhost:8980/api/v1/auth/login
# ❌ NOT: http://localhost:8980/api/api/v1/auth/login

# 6. Check Response:
# ✅ Status: 200 OK
# ❌ NOT: 401 Unauthorized

# 7. Verify login succeeds and redirects to dashboard
```

### ✅ Development Mode Test
```bash
# 1. Start Vite dev server
cd web/admin
npm run dev

# 2. Open browser
# Navigate to: http://localhost:5173

# 3. Open DevTools (F12) → Network tab

# 4. Attempt login

# 5. Check Network tab:
# ✅ Proxy should forward to localhost:8980
# ✅ Request path: /api/v1/auth/login

# 6. Verify login works
```

## Expected Results

### ✅ What Should Happen
- Request URLs: `/api/v1/...` (single `/api/` prefix)
- Response status: 200 OK
- Login succeeds
- Dashboard loads
- No console errors

### ❌ What Should NOT Happen
- Double paths: `/api/api/v1/...`
- 401 Unauthorized errors
- Network request failures
- Console errors about CORS or paths

## Automated Verification

Run automated checks:
```bash
cd web/admin
./verify-fix.sh
```

Expected output:
```
✅ All verification checks passed!
```

## Additional Endpoints to Test (Optional)

After successful login, test these endpoints:

### Dashboard
- Navigate to Dashboard
- Check Network tab: `/api/v1/stats/dashboard`
- ✅ Should be 200 OK

### Queue Management
- Navigate to Queue
- Check Network tab: `/api/v1/queue`
- ✅ Should be 200 OK

### Settings
- Navigate to Settings
- Check Network tabs:
  - `/api/v1/settings/server`
  - `/api/v1/settings/security`
  - `/api/v1/settings/tls`
- ✅ All should be 200 OK

### Domains
- Navigate to Domains
- Check Network tab: `/api/v1/domains`
- ✅ Should be 200 OK

## Troubleshooting

### If you see `/api/api/v1/...` paths:
1. Clear browser cache (Ctrl+Shift+Del)
2. Rebuild: `npm run build`
3. Restart server
4. Hard refresh browser (Ctrl+F5)

### If you see 401 errors:
1. Check that credentials are correct
2. Verify admin user exists (create with `create-admin` command)
3. Check server logs for authentication errors

### If build fails:
1. Clean install: `rm -rf node_modules package-lock.json && npm install`
2. Try build again: `npm run build`

## Success Criteria

✅ **Fix is successful when**:
- All API requests show single `/api/` prefix
- Login returns 200 OK status
- No 401 Unauthorized errors
- Dashboard loads successfully
- All API endpoints respond correctly
- Works in both dev and production modes

## Contact

If issues persist after testing:
1. Review `ISSUE-API-PATH-DOUBLING-RESOLUTION.md` for technical details
2. Check `FIX-SUMMARY.md` for quick reference
3. Examine server logs for additional context
4. Verify Vite config and axios config unchanged

---

**Fix Applied**: 2026-01-02
**Files Changed**: 12 files
**Status**: ✅ Automated checks passed, manual testing pending
