# Phase 2: Webmail Migration Plan

**Status**: Ready to Execute
**Prerequisites**: Phase 1 Complete ✅
**Complexity**: Medium (webmail already Vue 3 + Vite)

## Discovery Summary

**Good News**: The webmail has already been migrated from Nuxt 3 to Vue 3 + Vite!

### Current Webmail Structure
```
web/webmail/src/
├── pages/
│   ├── index.vue           # Landing/inbox
│   ├── login.vue           # Login page
│   ├── mail/
│   │   ├── [mailboxId].vue # Mailbox view (inbox, sent, etc.)
│   │   ├── compose.vue      # Email composer
│   │   └── [mailboxId]/message/[messageId].vue # Message detail
│   └── settings/
│       └── pgp.vue          # PGP settings
├── stores/
│   ├── auth.ts             # Auth store (Pinia)
│   └── mail.ts             # Mail store (Pinia)
├── router/
│   └── index.ts            # Router config (Vue Router 4)
├── assets/                 # CSS and assets
├── App.vue                 # Root component
└── main.ts                 # Entry point
```

### Current Router Structure
```javascript
createWebHistory('/webmail/')  // Already configured for /webmail/ base!
routes: [
  { path: '/', component: index.vue },
  { path: '/login', component: login.vue },
  { path: '/mail', component: mailbox, children: [...] },
  { path: '/settings', children: [...] }
]
```

## Migration Strategy

### Option A: Copy Files (Recommended)
**Pros**:
- Simple and straightforward
- Preserves existing structure
- Easy to test incrementally
- Can compare with original if needed

**Cons**:
- Duplicates code temporarily
- Need to sync auth stores

### Option B: Shared Modules
**Pros**:
- No code duplication
- Single source of truth

**Cons**:
- Complex build configuration
- Harder to maintain separate dev environments
- Risk of breaking existing webmail

**Decision**: Use Option A for Phase 2, consolidate in Phase 3

## Migration Steps

### Step 1: Copy Webmail Files to Unified App
```bash
# Copy pages
cp -r web/webmail/src/pages/* web/unified/src/views/webmail/

# Copy stores (will need merging)
cp web/webmail/src/stores/mail.ts web/unified/src/stores/

# Copy assets
cp -r web/webmail/src/assets/* web/unified/src/assets/webmail/
```

### Step 2: Update Router Configuration
Add webmail routes to `web/unified/src/router/index.js`:

```javascript
// Webmail module routes
{
  path: '/webmail',
  component: () => import('@/components/layout/WebmailLayout.vue'),
  meta: { requiresAuth: true, module: 'webmail' },
  children: [
    {
      path: '',
      name: 'WebmailInbox',
      component: () => import('@/views/webmail/index.vue')
    },
    {
      path: 'mail/:mailboxId',
      name: 'WebmailMailbox',
      component: () => import('@/views/webmail/mail/[mailboxId].vue'),
      children: [
        {
          path: 'message/:messageId',
          name: 'WebmailMessage',
          component: () => import('@/views/webmail/mail/[mailboxId]/message/[messageId].vue')
        }
      ]
    },
    {
      path: 'compose',
      name: 'WebmailCompose',
      component: () => import('@/views/webmail/mail/compose.vue')
    },
    {
      path: 'settings/pgp',
      name: 'WebmailPGP',
      component: () => import('@/views/webmail/settings/pgp.vue')
    }
  ]
}
```

### Step 3: Create Webmail Layout
Option 1: Use same AppLayout with conditional styling
Option 2: Create separate WebmailLayout.vue

**Decision**: Create WebmailLayout.vue for email-specific UI

### Step 4: Merge Auth Stores
The webmail has its own auth.ts. Need to:
1. Compare with unified auth.js
2. Merge webmail-specific auth logic
3. Ensure token sharing works correctly
4. Test authentication across both modules

### Step 5: Update API Calls
Webmail uses axios - need to ensure it uses the unified axios instance:
- Import from `@/api/axios` instead of creating new instance
- Verify baseURL configuration works
- Test all API endpoints

### Step 6: Handle Dependencies
Webmail specific dependencies:
- `@tiptap/*` - Rich text editor for email composer
- Already in webmail package.json
- Add to unified package.json

### Step 7: Update Navigation
Add webmail link to AppLayout navigation:
```javascript
{ name: 'Webmail', path: '/webmail', icon: '...' }
```

## Testing Checklist

### Authentication
- [ ] Login from webmail
- [ ] Token sharing with admin
- [ ] Logout from webmail
- [ ] Auto-redirect if not authenticated

### Email Viewing
- [ ] View inbox
- [ ] View sent folder
- [ ] View drafts folder
- [ ] View trash folder
- [ ] View custom folders
- [ ] Open individual email
- [ ] View email attachments

### Email Composition
- [ ] Compose new email
- [ ] Reply to email
- [ ] Reply all
- [ ] Forward email
- [ ] Add recipients (To, CC, BCC)
- [ ] Add subject
- [ ] Rich text editing (Tiptap)
- [ ] Add attachments
- [ ] Save as draft
- [ ] Send email

### Email Management
- [ ] Mark as read/unread
- [ ] Delete email
- [ ] Move to folder
- [ ] Search emails
- [ ] Filter emails
- [ ] Sort emails

### Settings
- [ ] PGP key management
- [ ] Import PGP key
- [ ] Export PGP key
- [ ] Generate PGP key

### Navigation
- [ ] Switch between admin and webmail
- [ ] Webmail navigation works
- [ ] Return to inbox
- [ ] Sidebar navigation

## Dependencies to Add

Add to `web/unified/package.json`:
```json
{
  "dependencies": {
    "@tiptap/extension-placeholder": "^3.14.0",
    "@tiptap/pm": "^3.14.0",
    "@tiptap/starter-kit": "^3.14.0",
    "@tiptap/vue-3": "^3.14.0",
    "@vueuse/core": "^14.1.0"
  }
}
```

## Potential Issues

### Issue 1: Dynamic Route Params
Webmail uses Nuxt-style dynamic routes: `[mailboxId].vue`

**Solution**: Rename to standard Vue Router format or use params in component

### Issue 2: Different Auth Logic
Webmail might have email-specific auth requirements

**Solution**: Extend unified auth store with webmail-specific methods

### Issue 3: Asset Loading
Webmail assets might not load correctly in unified app

**Solution**: Update asset imports to use `@/assets/webmail/`

### Issue 4: TypeScript
Webmail uses TypeScript (.ts files)

**Solution**:
- Configure unified app to support TypeScript
- Or convert to JavaScript
- Recommend: Add TypeScript support to unified app

## Timeline Estimate

**Optimistic**: 4-6 hours (files already Vue 3)
**Realistic**: 8-10 hours (integration and testing)
**Pessimistic**: 12-16 hours (unforeseen issues)

## Success Criteria

✅ All webmail pages accessible at `/webmail/*`
✅ Authentication shared between admin and webmail
✅ Email viewing works
✅ Email composition and sending works
✅ Navigation between admin and webmail works
✅ No console errors
✅ All API calls use correct paths
✅ Build process successful
✅ Production deployment works

## Rollback Plan

If Phase 2 fails:
1. Webmail remains at `web/webmail/` (unchanged)
2. Can be served separately if needed
3. Unified app continues with admin-only (Phase 1)

## Next Steps After Phase 2

1. Comprehensive testing of webmail functionality
2. Performance testing with large mailboxes
3. Security audit of email handling
4. Proceed to Phase 3: Portal Module
