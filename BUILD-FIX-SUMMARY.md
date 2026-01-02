# Build Fix Summary

**Date**: 2026-01-02
**Status**: ✅ Resolved

## Issue
The webmail UI build was failing due to incompatible imports from the previous Nuxt 3 framework. After migrating to a vanilla Vite + Vue 3 setup, several files still referenced Nuxt-specific composables and APIs.

## Root Cause
The webmail codebase was converted from Nuxt 3 to Vite, but several Vue component files retained Nuxt-specific patterns:

1. **Missing composable**: `useAuth()` was imported from `~/composables/useAuth` which didn't exist
2. **Nuxt auto-imports**: Files used `useRoute()`, `computed()`, `ref()`, etc. without importing them
3. **Nuxt metadata**: Files used `definePageMeta()` which is Nuxt-specific

## Files Fixed

### 1. `/web/webmail/src/pages/settings/pgp.vue`
**Problem**: Imported non-existent `~/composables/useAuth`
**Solution**:
- Replaced with direct Pinia store import: `import { useAuthStore } from '../../stores/auth'`
- Changed from `const { user } = useAuth()` to `const authStore = useAuthStore(); const user = ref(authStore.user)`

### 2. `/web/webmail/src/pages/login.vue`
**Problem**: Used auto-imported `useAuth()` composable
**Solution**:
- Added explicit imports for Vue and router functionality
- Replaced `const { login } = useAuth()` with direct store usage: `const authStore = useAuthStore()`
- Updated login handler to use `authStore.login()` and proper router navigation

### 3. `/web/webmail/src/pages/index.vue`
**Problem**: Used `useAuth()` and Nuxt's `navigateTo()`
**Solution**:
- Added explicit imports for Vue and router
- Replaced `checkAuth()` with `authStore.initializeAuth()` and `authStore.isAuthenticated`
- Replaced `navigateTo()` with Vue Router's `router.push()`

### 4. `/web/webmail/src/pages/mail/compose.vue`
**Problem**: Used `definePageMeta()` and auto-imported Vue functions
**Solution**:
- Removed `definePageMeta()` call
- Added explicit imports for `computed` and `useRoute`

### 5. `/web/webmail/src/pages/mail/[mailboxId].vue`
**Problem**: Used `definePageMeta()` and auto-imported Vue functions
**Solution**:
- Removed `definePageMeta()` call
- Added explicit imports for `computed` and `useRoute`

### 6. `/web/webmail/src/pages/mail/[mailboxId]/message/[messageId].vue`
**Problem**: Used `definePageMeta()` and auto-imported Vue functions
**Solution**:
- Removed `definePageMeta()` call
- Added explicit imports for `computed` and `useRoute`

## Build Results

### Before Fix
```
error during build:
[vite]: Rollup failed to resolve import "~/composables/useAuth" from
"/home/btafoya/projects/gomailserver/web/webmail/src/pages/settings/pgp.vue"
```

### After Fix
```
✓ built in 944ms
Webmail UI build complete
Building gomailserver...
Build complete: ./build/gomailserver
```

## Verification

1. ✅ Admin UI builds successfully
2. ✅ Webmail UI builds successfully
3. ✅ Go binary compiles without errors
4. ✅ Binary executes and shows version: `gomailserver version dev`

## Pattern Applied

For future Vue 3 + Vite development in this project:

1. **Always use explicit imports**: Import `ref`, `computed`, `onMounted`, etc. from `vue`
2. **Use Pinia stores directly**: Import and use `useAuthStore()`, `useMailStore()` instead of composables
3. **Use Vue Router directly**: Import `useRoute()` and `useRouter()` from `vue-router`
4. **No Nuxt APIs**: Avoid `definePageMeta()`, `navigateTo()`, auto-imports, and other Nuxt-specific features

## Related Documentation
- WEBUI-FIXES.md - Previous UI fixes applied
- web/webmail/package.json - Current Vite configuration
- web/webmail/src/stores/auth.ts - Pinia auth store implementation
