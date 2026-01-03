#!/bin/bash

echo "🔍 API Path Doubling Fix Verification"
echo "======================================"
echo ""

# Check 1: Verify no hardcoded /api/v1 paths
echo "1. Checking for hardcoded /api/v1 paths..."
if grep -r "/api/v1" src/ 2>/dev/null; then
  echo "   ❌ FAILED: Found hardcoded /api/v1 paths"
  exit 1
else
  echo "   ✅ PASSED: No hardcoded /api/v1 paths"
fi
echo ""

# Check 2: Verify axios baseURL uses window.location.origin
echo "2. Checking axios baseURL configuration..."
if grep -q "window.location.origin" src/api/axios.js; then
  echo "   ✅ PASSED: Axios uses runtime origin resolution"
else
  echo "   ❌ FAILED: Axios doesn't use window.location.origin"
  exit 1
fi
echo ""

# Check 3: Verify build succeeds
echo "3. Testing production build..."
if npm run build >/dev/null 2>&1; then
  echo "   ✅ PASSED: Production build successful"
else
  echo "   ❌ FAILED: Production build failed"
  exit 1
fi
echo ""

# Check 4: Verify runtime resolution in built assets
echo "4. Checking built assets for runtime resolution..."
if grep -q "window.location.origin" dist/assets/*.js; then
  echo "   ✅ PASSED: Runtime resolution found in build"
else
  echo "   ❌ FAILED: Runtime resolution not in build"
  exit 1
fi
echo ""

# Check 5: Verify Vite config still has correct base
echo "5. Checking Vite base configuration..."
if grep -q "base: '/admin/'" vite.config.js; then
  echo "   ✅ PASSED: Vite base path correct"
else
  echo "   ❌ FAILED: Vite base path incorrect"
  exit 1
fi
echo ""

echo "======================================"
echo "✅ All verification checks passed!"
echo ""
echo "Manual Testing Required:"
echo "1. Start server: ./build/gomailserver run"
echo "2. Visit: http://localhost:8980/admin/"
echo "3. Check Network tab: requests should show /api/v1/... (not /api/api/v1/...)"
echo "4. Verify: Login succeeds with 200 status"
