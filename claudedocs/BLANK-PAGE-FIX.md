# Blank Page Fix - Network Access Issue

## Problem
User reported blank page when accessing WebUI at `http://192.168.25.165:5173/admin/`

## Root Causes Identified

### 1. Orphaned Vite Process
**Issue**: An orphaned Vite process from a previous session was holding port 5173
**Symptom**: Current WebUI started on port 5174 instead of 5173
**Evidence**: `lsof -i :5173` showed PID 2763972 listening on port

### 2. Network Binding Configuration
**Issue**: Vite dev server was only listening on localhost (127.0.0.1)
**Symptom**: Vite logs showed "use --host to expose" message
**Impact**: WebUI not accessible via network IP addresses

## Solution Implemented

### 1. Killed Orphaned Process
```bash
kill 2763972
```
This freed port 5173 for the current WebUI instance.

### 2. Updated Vite Configuration
**File**: `web/unified/vite.config.js`

**Change**:
```javascript
server: {
  host: '0.0.0.0', // Listen on all network interfaces (NEW)
  port: 5173,
  proxy: {
    '/api': {
      target: 'http://localhost:8980',
      changeOrigin: true
    }
  }
}
```

**Effect**: Vite now listens on all network interfaces, making the WebUI accessible from:
- Local: `http://localhost:5173/admin/`
- Network: `http://192.168.25.165:5173/admin/`
- Network: `http://10.0.0.165:5173/admin/`
- All other network interfaces on the host

### 3. Updated Control Script Messages
Enhanced `scripts/gomailserver-control.sh` to inform users about network access:
- Start message now mentions network URL availability
- Status command directs users to check logs for all network URLs

## Verification

### Port Binding
```bash
$ lsof -i :5173 | grep LISTEN
node    2788204 btafoya   24u  IPv4 110462044      0t0  TCP *:5173 (LISTEN)
```
✅ Port 5173 listening on all interfaces (`*:5173` instead of `localhost:5173`)

### Network Access
```bash
$ curl -s http://192.168.25.165:5173/admin/ | head -15
<!doctype html>
<html lang="en">
  <head>
    <script type="module" src="/admin/@vite/client"></script>
    ...
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/admin/src/main.js"></script>
  </body>
</html>
```
✅ HTML served successfully via network IP

### Vite Startup Logs
```
VITE v6.4.1  ready in 157 ms

➜  Local:   http://localhost:5173/admin/
➜  Network: http://192.168.25.165:5173/admin/
➜  Network: http://10.0.0.165:5173/admin/
[... all network interfaces listed ...]
```
✅ All network URLs displayed in logs

## If Blank Page Persists

If you still see a blank page in your browser after these fixes, check the following:

### 1. Browser Console Errors
**Action**: Open browser DevTools (F12) → Console tab
**Look for**:
- JavaScript errors (red messages)
- Failed module loads
- 404 errors for `/admin/src/main.js` or other resources

**Common Issues**:
- Path resolution errors
- CORS errors
- Module import failures

### 2. Network Tab
**Action**: Open DevTools → Network tab → Reload page
**Check**:
- `/admin/` returns 200 (HTML page)
- `/admin/src/main.js` returns 200
- `/admin/@vite/client` returns 200
- No 404 errors for resources

### 3. Authentication Redirect
**Action**: Check if page redirects to `/admin/login`
**Reason**: Vue router has auth guard that redirects unauthenticated users
**Solution**: This is expected behavior - login page should appear, not blank page

**Router Auth Guard** (`web/unified/src/router/index.js:161-172`):
```javascript
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth !== false)

  if (requiresAuth && !token) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (to.name === 'Login' && token) {
    next({ name: 'AdminDashboard' })
  } else {
    next()
  }
})
```

### 4. Browser Cache
**Action**: Hard refresh or clear cache
**Methods**:
- Chrome/Firefox: Ctrl+Shift+R (Windows/Linux) or Cmd+Shift+R (Mac)
- Or: DevTools → Network tab → Check "Disable cache"

### 5. API Backend Connection
**Action**: Check if backend API is accessible
**Test**:
```bash
curl http://192.168.25.165:8980/api/v1/health
```
**Expected**: API response (not connection refused)

**Common Issue**: Vite proxy configured for localhost, but accessing via IP
**Solution**: Update `vite.config.js` proxy target to use IP address when accessing remotely

### 6. Firewall
**Action**: Verify port 5173 is accessible through firewall
**Test from another machine**:
```bash
curl http://192.168.25.165:5173/admin/
```
**If fails**: Check firewall rules on host machine

## Production Considerations

**Important**: The `host: '0.0.0.0'` configuration is for **development only**.

**Why**: In production:
- WebUI is built to static files (`web/unified/dist/`)
- Static files are served by the Go backend
- No Vite dev server runs
- Network access controlled by Go server configuration

**Production Deployment**:
```bash
# Build WebUI
cd web/unified
pnpm build

# Start production server (no WebUI dev server)
./scripts/gomailserver-control.sh start
```

## Network Security Considerations

### Development Mode Security
With `host: '0.0.0.0'`, the WebUI dev server is accessible from any network interface:
- ✅ Good: Easy development and testing from other devices
- ⚠️ Risk: Exposes dev server to local network
- ⚠️ Risk: Anyone on local network can access WebUI

### Recommendations

**For secure local development**:
1. Use firewall rules to restrict port 5173 access
2. Use VPN or SSH tunnel for remote access
3. Only enable network binding when needed

**For team development**:
```javascript
// vite.config.js
server: {
  host: process.env.VITE_HOST || 'localhost', // Default to localhost
  port: 5173,
  // ...
}
```
Then enable network access when needed:
```bash
VITE_HOST=0.0.0.0 pnpm dev
```

**For public networks** (coffee shops, conferences):
- Do NOT use `host: '0.0.0.0'`
- Use `host: 'localhost'` only
- Access via SSH tunnel if needed

## Troubleshooting Commands

### Check WebUI Status
```bash
./scripts/gomailserver-control.sh status
```

### View WebUI Logs
```bash
tail -f data/webui.log
```

### Check Port Binding
```bash
lsof -i :5173
```

### Kill Orphaned Processes
```bash
# Find Vite processes
ps aux | grep vite

# Kill specific PID
kill <PID>
```

### Test Network Access
```bash
# From same machine
curl http://localhost:5173/admin/

# From remote machine
curl http://192.168.25.165:5173/admin/
```

### Restart WebUI
```bash
./scripts/gomailserver-control.sh restart --dev
```

## Summary

**Problem**: Blank page at network IP address
**Causes**:
1. Orphaned Vite process on port 5173
2. Vite only listening on localhost

**Solution**:
1. Killed orphaned process to free port
2. Added `host: '0.0.0.0'` to Vite config
3. Restarted WebUI with new configuration

**Result**: WebUI now accessible from all network interfaces

**Files Modified**:
- `web/unified/vite.config.js` - Added network binding
- `scripts/gomailserver-control.sh` - Updated status messages

**Test**: http://192.168.25.165:5173/admin/ should now load
