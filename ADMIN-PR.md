# COMPREHENSIVE PLAN: Restore Reputation Management & Setup Wizard to Nuxt Admin WebUI

## 📋 EXECUTIVE SUMMARY

**What Was Lost During Unified Refactor (commit 9df2bcb4, Jan 6 2026):**

1. **9 Reputation Management Pages** - Complete UI for all reputation features
2. **4 Reputation Dashboard Components** - Reusable widgets for overview
3. **Setup Wizard Page** - First-run configuration tool
4. **Reputation API Composables** - API integration layer

**Backend Status:** ✅ **FULLY IMPLEMENTED**
- 38+ REST API endpoints exist across reputation_handler.go, reputation_phase5_handler.go
- Setup wizard API exists (4 endpoints in setup_handler.go)
- **Note:** Phase 6 endpoints are commented out due to compilation errors

**Current Unified UI Status:** ✅ **BASIC ADMIN ONLY**
- 12 admin pages (dashboard, domains, users, aliases, queue, settings)
- 1 admin component (Sidebar.vue)
- 6 basic API composables (domains, users, aliases, queue, auth)
- **NO reputation pages, components, or API integrations**

---

## 🎯 IMPLEMENTATION PLAN

### **PHASE 1: Foundation - API Composables** (Priority: HIGH)
*Create reusable API integration layer following existing patterns*

#### 1.1 Basic Reputation API (Phase 1-4)
**File:** `unified/composables/api/reputation.ts`

**Endpoints to implement:**
```typescript
// Domain audit
POST /api/v1/reputation/audit/{domain}

// Scores
GET /api/v1/reputation/scores (list all)
GET /api/v1/reputation/scores/{domain}

// Circuit breakers
GET /api/v1/reputation/circuit-breakers
GET /api/v1/reputation/circuit-breakers/{domain}/history

// Alerts
GET /api/v1/reputation/alerts
```

#### 1.2 Advanced Reputation API (Phase 5)
**File:** `unified/composables/api/reputation-phase5.ts`

**Endpoints to implement:**
```typescript
// DMARC Reports
GET /api/v1/reputation/dmarc/reports (list, paginated)
GET /api/v1/reputation/dmarc/reports/{id}
GET /api/v1/reputation/dmarc/stats/{domain}
GET /api/v1/reputation/dmarc/actions
POST /api/v1/reputation/dmarc/reports/{id}/export

// ARF Reports
GET /api/v1/reputation/arf/reports
GET /api/v1/reputation/arf/stats
POST /api/v1/reputation/arf/reports/{id}/process

// External Metrics
GET /api/v1/reputation/external/postmaster/{domain}
GET /api/v1/reputation/external/snds/{ip}
GET /api/v1/reputation/external/trends

// Provider Rate Limits
GET /api/v1/reputation/provider-limits
PUT /api/v1/reputation/provider-limits/{id}
POST /api/v1/reputation/provider-limits/init/{domain}
POST /api/v1/reputation/provider-limits/{id}/reset

// Custom Warmup
GET /api/v1/reputation/warmup/{domain}
POST /api/v1/reputation/warmup
PUT /api/v1/reputation/warmup/{id}
DELETE /api/v1/reputation/warmup/{id}
GET /api/v1/reputation/warmup/templates

// Predictions
GET /api/v1/reputation/predictions/latest
GET /api/v1/reputation/predictions/{domain}
POST /api/v1/reputation/predictions/generate/{domain}
GET /api/v1/reputation/predictions/{domain}/history

// Phase 5 Alerts
GET /api/v1/reputation/alerts/phase5
POST /api/v1/reputation/alerts/{id}/acknowledge
POST /api/v1/reputation/alerts/{id}/resolve
```

#### 1.3 Setup Wizard API
**File:** `unified/composables/api/setup.ts`

**Endpoints to implement:**
```typescript
// No authentication required
GET /api/v1/setup/status
GET /api/v1/setup/state
POST /api/v1/setup/admin (create first admin user)
POST /api/v1/setup/complete
```

#### 1.4 Phase 6 Features (Optional - backend commented out)
**File:** `unified/composables/api/reputation-phase6.ts`

**Endpoints (currently disabled in router.go):**
```typescript
// Operational Mail
GET /api/v1/reputation/operational-mail
POST /api/v1/reputation/operational-mail/{id}/read
DELETE /api/v1/reputation/operational-mail/{id}
POST /api/v1/reputation/operational-mail/{id}/spam
POST /api/v1/reputation/operational-mail/{id}/forward

// Deliverability Status
GET /api/v1/reputation/deliverability
GET /api/v1/reputation/deliverability/{domain}

// Enhanced Circuit Breaker Controls
GET /api/v1/reputation/circuit-breakers/active
POST /api/v1/reputation/circuit-breakers/{id}/resume
POST /api/v1/reputation/circuit-breakers/pause

// Enhanced Alerts
GET /api/v1/reputation/alerts/unread
POST /api/v1/reputation/alerts/{id}/read
```

---

### **PHASE 2: Setup Wizard** (Priority: HIGH)
*First-run configuration tool to guide system setup*

#### 2.1 Setup Pages
**File:** `unified/pages/admin/setup/index.vue`

**Structure:**
- Multi-step wizard with progress indicator
- 4 steps: System Configuration → Domain Setup → Admin User → Review
- Live DNS validation during domain setup
- Form validation at each step
- Summary and confirmation before completion

#### 2.2 Setup Components
**Directory:** `unified/components/admin/setup/`

**Components:**
- `Step1System.vue` - Server hostname, port, TLS settings
- `Step2Domain.vue` - Primary domain with SPF/DKIM/DMARC validation
- `Step3Admin.vue` - Admin user creation (email, name, password, 2FA)
- `Step4Review.vue` - Summary of all configuration, confirm button

#### 2.3 Setup Middleware
**File:** `unified/middleware/setup.ts`

**Logic:**
- Check `GET /api/v1/setup/status` on app load
- If `setup_complete: false`, redirect to `/admin/setup`
- If setup complete, redirect to `/admin/` when on `/admin/setup`

---

### **PHASE 3: Reputation Core Pages** (Priority: HIGH)
*Basic reputation management features (Phases 1-4)*

#### 3.1 Reputation Overview/Dashboard
**File:** `unified/pages/admin/reputation/index.vue`

**Features:**
- System-wide reputation statistics
- Active alerts count (with severity badges)
- Circuit breaker status summary
- Quick action buttons (run audit, view circuit breakers)
- Domain reputation score table (sorted by score)
- Links to detailed views

**Components Used:**
- `DeliverabilityCard.vue` (from Phase 6)
- `CircuitBreakersCard.vue` (from Phase 6)
- `RecentAlertsTimeline.vue` (from Phase 6)

#### 3.2 Circuit Breakers Management
**File:** `unified/pages/admin/reputation/circuit-breakers/index.vue`

**Features:**
- Table of all circuit breaker events
- Status indicators (active/paused/resumed)
- Filter by domain and status
- Manual resume button for active breakers
- View detailed history per domain
- Auto-resume countdown timers

#### 3.3 Warm-up Tracking
**File:** `unified/pages/admin/reputation/warmup/index.vue`

**Features:**
- List of domains in warm-up phase
- Current day of warm-up schedule (1-14)
- Daily volume limits and usage
- Progress bars for warm-up completion
- Create custom warm-up schedule button
- Edit/delete custom schedules

#### 3.4 Domain Audit Tool
**File:** `unified/pages/admin/reputation/audit/index.vue`

**Features:**
- Input field for domain to audit
- "Run Audit" button with loading state
- Detailed results display:
  - SPF status (pass/fail with record details)
  - DKIM status (pass/fail with selector)
  - DMARC status (pass/fail with policy)
  - rDNS/FCrDNS status
  - TLS certificate status
  - MTA-STS status
  - Postmaster/abuse mailbox check
- Overall deliverability score (0-100)
- List of issues with severity
- Export audit results (JSON/PDF)

---

### **PHASE 4: Reputation Advanced Pages** (Priority: MEDIUM)
*Advanced reputation features (Phase 5)*

#### 4.1 DMARC Reports Viewer
**File:** `unified/pages/admin/reputation/dmarc-reports/index.vue`

**Features:**
- Table of received DMARC reports
- Filters: Domain, date range, org name
- Pagination (20 reports per page)
- Expand row to see detailed breakdown:
  - SPF/DKIM alignment counts
  - Policy pass rates
  - Source IP distribution
- Statistics summary cards (alignment rates, failed messages)
- Export selected reports (XML/CSV)
- View individual report details

#### 4.2 External Metrics Dashboard
**File:** `unified/pages/admin/reputation/external-metrics/index.vue`

**Features:**
- **Gmail Postmaster Tools section:**
  - Domain reputation score (Good/Fair/Poor)
  - Spam rate percentage
  - Authentication rates (SPF/DKIM/DMARC)
  - Encryption rate
  - Trend chart (7-day history)

- **Microsoft SNDS section:**
  - IP reputation score
  - Spam trap hit count
  - Complaint rate
  - Filter level (GREEN/YELLOW/RED)

- Trend visualization (line charts for 7/30 days)
- Manual sync buttons (sync Gmail, sync SNDS)
- Alert indicators when metrics deteriorate

#### 4.3 Provider Rate Limits
**File:** `unified/pages/admin/reputation/provider-limits/index.vue`

**Features:**
- Table of provider-specific limits (Gmail, Outlook, Yahoo, etc.)
- Current usage vs. limit
- Edit limits button
- Reset usage counter button
- Initialize defaults for new domain button
- Percentage utilization bars (green < 50%, yellow 50-80%, red > 80%)

#### 4.4 Custom Warm-up Scheduler
**File:** `unified/pages/admin/reputation/warmup-scheduler/index.vue`

**Features:**
- List of custom warm-up schedules
- Template options:
  - Conservative (14-day, low volumes)
  - Moderate (10-day, medium volumes)
  - Aggressive (7-day, higher volumes)
- Create/edit schedule form:
  - Domain selection
  - Template selection or custom schedule
  - Day-by-day volume configuration
- Progress tracking for active schedules
- Delete schedule button

#### 4.5 AI Predictions
**File:** `unified/pages/admin/reputation/predictions/index.vue`

**Features:**
- Latest predictions for all domains
- Trend indicators:
  - Improving (green arrow up)
  - Stable (gray dash)
  - Declining (red arrow down)
- Confidence levels (high/medium/low)
- Prediction horizon (7/14/30 days)
- Domain-specific prediction details (click to expand):
  - Historical trend data
  - Factors influencing prediction
  - Recommended actions
- "Generate New Predictions" button
- Historical predictions table

---

### **PHASE 5: Reputation Dashboard Components** (Priority: MEDIUM)
*Reusable widgets for admin dashboard*

#### 5.1 Deliverability Card
**File:** `unified/components/admin/reputation/DeliverabilityCard.vue`

**Features:**
- Circular gauge or progress bar for overall reputation score (0-100)
- Color-coded: Green (>70), Yellow (50-70), Red (<50)
- Trend indicator (7-day change)
- Quick status badges for:
  - SPF (pass/fail)
  - DKIM (pass/fail)
  - DMARC (pass/fail)
  - rDNS (pass/fail)
- "Run Full Audit" button

#### 5.2 Circuit Breakers Card
**File:** `unified/components/admin/reputation/CircuitBreakersCard.vue`

**Features:**
- Count of active circuit breakers (badge with severity)
- List of active breakers:
  - Domain name
  - Trigger type (complaint/bounce/provider block)
  - Trigger value vs threshold
  - Paused timestamp
  - Auto-resume countdown
- Manual "Resume All" button
- Click row to view details

#### 5.3 Recent Alerts Timeline
**File:** `unified/components/admin/reputation/RecentAlertsTimeline.vue`

**Features:**
- Vertical timeline of last 10 alerts
- Timeline connectors (colored by severity)
- Alert cards showing:
  - Timestamp (relative: "2 hours ago")
  - Severity badge (critical/warning/info)
  - Alert type icon
  - Alert message (truncated)
- Filter by severity
- Click alert to view details

#### 5.4 Operational Mail Inbox
**File:** `unified/components/admin/reputation/OperationalMail.vue`

**Features:**
- List of postmaster@ and abuse@ messages
- Unread badge count
- Quick actions per message:
  - Mark as read
  - Mark as spam
  - Forward to admin email
  - Delete
- Message preview panel
- Filter by sender domain
- Auto-refresh every 30 seconds

#### 5.5 Score Gauge Component (Reusable)
**File:** `unified/components/admin/reputation/ScoreGauge.vue`

**Props:**
- score (0-100)
- size (small/large)
- showTrend (boolean)
- trendValue (number, percent change)

**Features:**
- Animated circular progress
- Color gradient (red→yellow→green)
- Center text with score
- Optional trend indicator (+X% or -X%)
- Tooltip on hover with details

---

### **PHASE 6: Integration & Navigation** (Priority: HIGH)
*Connect all pieces into unified admin interface*

#### 6.1 Update Admin Sidebar
**File:** `unified/components/admin/Sidebar.vue`

**Changes:**
- Add new navigation group "Reputation Management" between "Management" and "System"
- Navigation items:
  - Overview → /admin/reputation (home icon)
  - Circuit Breakers → /admin/reputation/circuit-breakers (alert icon)
  - Warm-up → /admin/reputation/warmup (trending-up icon)
  - Audit → /admin/reputation/audit (shield-check icon)
  - DMARC Reports → /admin/reputation/dmarc-reports (file-text icon)
  - External Metrics → /admin/reputation/external-metrics (bar-chart icon)
  - Provider Limits → /admin/reputation/provider-limits (sliders icon)
  - Warmup Scheduler → /admin/reputation/warmup-scheduler (calendar icon)
  - Predictions → /admin/reputation/predictions (brain icon)
- Collapsible group (same pattern as existing groups)

#### 6.2 Add Setup Redirect Middleware
**File:** `unified/middleware/setup-redirect.ts`

**Logic:**
```typescript
export default defineNuxtRouteMiddleware(async (to, from) => {
  // Skip if already going to setup
  if (to.path === '/admin/setup') return

  // Check setup status on first load or every page load
  const { setupComplete } = await $fetch('/api/v1/setup/status')

  if (!setupComplete) {
    return navigateTo('/admin/setup')
  }
})
```

**Apply to:** `unified/layouts/admin.vue` middleware configuration

#### 6.3 Update Admin Dashboard
**File:** `unified/pages/admin/index.vue`

**Changes:**
- Replace existing dashboard content with reputation overview
- Import and use Phase 6 components:
  - `<DeliverabilityCard />` - Top center
  - `<CircuitBreakersCard />` - Top right
  - `<RecentAlertsTimeline />` - Bottom section
- Keep existing quick links to other admin sections
- Add "Run Full Audit" button

---

## 🛠️ TECHNICAL IMPLEMENTATION GUIDELINES

### Nuxt 3 Conventions

**File-Based Routing:**
```
unified/pages/admin/reputation/index.vue         → /admin/reputation
unified/pages/admin/reputation/circuit-breakers/  → /admin/reputation/circuit-breakers
```

**Composition API:**
```vue
<script setup lang="ts">
// Auto-imports work automatically
const route = useRoute()
const router = useRouter()
const { data, pending, error, refresh } = await useFetch('/api/v1/reputation/scores')

// State management
const domainInput = ref('')
const auditResults = ref(null)
</script>
```

**Composables Pattern:**
```typescript
// unified/composables/api/reputation.ts
export const useReputationApi = () => {
  const API_BASE = useApiBase()

  const getScores = async () => {
    const response = await fetch(`${API_BASE}/reputation/scores`, {
      headers: getAuthHeaders()
    })
    if (!response.ok) throw new Error('Failed to fetch scores')
    return await response.json()
  }

  return { getScores, /* other methods */ }
}
```

### Shadcn UI Components

**Available Components:**
- Card, Button, Table, Badge, Alert, Dialog, Tabs
- Progress, Slider, Select, Input, Textarea
- Separator, Skeleton, LoadingSpinner

**Usage Example:**
```vue
<Card>
  <CardHeader>
    <CardTitle>Reputation Overview</CardTitle>
  </CardHeader>
  <CardContent>
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Domain</TableHead>
          <TableHead>Score</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="domain in domains" :key="domain.id">
          <TableCell>{{ domain.name }}</TableCell>
          <TableCell>{{ domain.score }}</TableCell>
        </TableRow>
      </TableBody>
    </Table>
  </CardContent>
</Card>
```

### Authentication & Authorization

**Setup Wizard (No Auth Required):**
- Direct API calls, no Authorization header
- Public route (no middleware)

**Admin Pages (JWT Required):**
```typescript
// Use existing auth composable
import { getAuthHeaders } from '~/composables/api/auth'

const response = await fetch(`${API_BASE}/reputation/scores`, {
  headers: getAuthHeaders()
})
```

**Token Storage:** `localStorage.getItem('token')`

### Error Handling

**Pattern:**
```vue
<script setup lang="ts">
const { data, pending, error } = await useReputationApi().getScores()

watch(error, (newError) => {
  if (newError) {
    showToast({
      title: 'Error',
      description: newError.message,
      variant: 'destructive'
    })
  }
})
</script>

<template>
  <Alert v-if="error" variant="destructive">
    {{ error.message }}
  </Alert>
</template>
```

### Loading States

**Pattern:**
```vue
<template>
  <div v-if="pending" class="flex justify-center">
    <LoadingSpinner />
  </div>
  <div v-else>
    <!-- Content -->
  </div>
</template>
```

### Real-Time Updates

**Polling Pattern:**
```vue
<script setup lang="ts">
let refreshInterval: any = null

onMounted(() => {
  refreshScores()
  refreshInterval = setInterval(refreshScores, 30000) // 30 seconds
})

onUnmounted(() => {
  if (refreshInterval) clearInterval(refreshInterval)
})
</script>
```

---

## 📁 FILE STRUCTURE SUMMARY

### New Composables (4 files)
```
unified/composables/api/
├── reputation.ts           # Phase 1-4 endpoints (250-300 lines)
├── reputation-phase5.ts   # Phase 5 endpoints (350-400 lines)
├── reputation-phase6.ts   # Phase 6 endpoints (150-200 lines)
└── setup.ts               # Setup wizard endpoints (80-100 lines)
```

### New Pages (13 files)
```
unified/pages/admin/
├── setup/
│   └── index.vue          # Setup wizard (300-400 lines)
└── reputation/
    ├── index.vue           # Overview/dashboard (200-250 lines)
    ├── circuit-breakers/
    │   └── index.vue       # Circuit breaker mgmt (250-300 lines)
    ├── warmup/
    │   └── index.vue       # Warm-up tracking (250-300 lines)
    ├── audit/
    │   └── index.vue       # Domain audit tool (300-350 lines)
    ├── dmarc-reports/
    │   └── index.vue       # DMARC reports viewer (400-500 lines)
    ├── external-metrics/
    │   └── index.vue       # External metrics dashboard (400-500 lines)
    ├── provider-limits/
    │   └── index.vue       # Provider limits (300-350 lines)
    ├── warmup-scheduler/
    │   └── index.vue       # Custom warmup schedules (350-400 lines)
    └── predictions/
        └── index.vue       # AI predictions (350-400 lines)
```

### New Components (5 files)
```
unified/components/admin/reputation/
├── DeliverabilityCard.vue          # Score gauge & DNS status (200-250 lines)
├── CircuitBreakersCard.vue        # Active breakers widget (250-300 lines)
├── RecentAlertsTimeline.vue       # Alert timeline (300-350 lines)
├── OperationalMail.vue            # Postmaster/abuse inbox (400-500 lines)
└── ScoreGauge.vue               # Reusable gauge (150-200 lines)
```

### Setup Components (4 files)
```
unified/components/admin/setup/
├── Step1System.vue              # Server configuration (150-200 lines)
├── Step2Domain.vue              # Domain setup w/ DNS validation (250-300 lines)
├── Step3Admin.vue              # Admin user creation (200-250 lines)
└── Step4Review.vue              # Summary & confirm (150-200 lines)
```

### Updated Files (2 files)
```
unified/components/admin/Sidebar.vue  # Add reputation navigation (20-30 lines)
unified/pages/admin/index.vue        # Replace with reputation dashboard (100-150 lines)
```

### New Middleware (1 file)
```
unified/middleware/setup-redirect.ts  # Setup status check (30-40 lines)
```

**Total New/Modified Files:** ~29 files
**Estimated Lines of Code:** ~5,000-6,000 lines

---

## 🚀 IMPLEMENTATION SEQUENCE

### Week 1: Foundation (API + Setup)
1. **Day 1-2:** Create all API composables (4 files)
2. **Day 3:** Build setup wizard page and components
3. **Day 4-5:** Test setup wizard end-to-end, add middleware

### Week 2: Core Reputation
4. **Day 1-2:** Build Overview and Circuit Breakers pages
5. **Day 3:** Build Warm-up and Audit pages
6. **Day 4-5:** Test Phase 1-4 pages, add to sidebar

### Week 3: Advanced Reputation
7. **Day 1-2:** Build DMARC Reports and External Metrics pages
8. **Day 3:** Build Provider Limits and Warmup Scheduler
9. **Day 4:** Build Predictions page
10. **Day 5:** Test Phase 5 pages integration

### Week 4: Polish & Dashboard
11. **Day 1-2:** Build Phase 6 components (4 cards)
12. **Day 3:** Update admin dashboard with components
13. **Day 4:** Update sidebar with reputation navigation
14. **Day 5:** Final testing and bug fixes

---

## ✅ ACCEPTANCE CRITERIA

### Phase 1: API Composables
- [ ] All composables follow existing patterns
- [ ] TypeScript interfaces exported for all API responses
- [ ] Error handling with user-friendly messages
- [ ] Authorization headers properly included
- [ ] Unit tests for each composable

### Phase 2: Setup Wizard
- [ ] Setup page accessible without authentication
- [ ] 4-step wizard with progress indicator
- [ ] Live DNS validation during domain setup
- [ ] Admin user creation works
- [ ] Setup completion marks server as configured
- [ ] Redirect to /admin/ after setup complete

### Phase 3: Core Pages
- [ ] Overview page shows system-wide reputation
- [ ] Circuit breakers page allows manual resume
- [ ] Warm-up page shows progress and schedules
- [ ] Audit page runs and displays detailed results
- [ ] All pages have loading states
- [ ] Error handling for failed API calls

### Phase 4: Advanced Pages
- [ ] DMARC reports display with filtering and pagination
- [ ] External metrics show Gmail and Microsoft data
- [ ] Provider limits table allows editing
- [ ] Warmup scheduler creates/edits/deletes schedules
- [ ] Predictions page shows AI forecasts

### Phase 5: Dashboard Components
- [ ] DeliverabilityCard shows score with color coding
- [ ] CircuitBreakersCard lists active breakers
- [ ] RecentAlertsTimeline shows last 10 alerts
- [ ] OperationalMail shows postmaster/abuse inbox
- [ ] ScoreGauge is reusable with proper props

### Phase 6: Integration
- [ ] Sidebar includes all reputation navigation
- [ ] Navigation groups collapse/expand properly
- [ ] Setup redirect works on unconfigured server
- [ ] Admin dashboard shows reputation cards
- [ ] All routes work without page reload
- [ ] Auto-refresh on data pages (30-60s polling)

---

## ⚠️ KNOWN RISKS & MITIGATIONS

### Risk 1: Phase 6 Backend Not Compiled
**Issue:** Phase 6 endpoints are commented out in router.go due to compilation errors
**Impact:** Operational mail, enhanced alerts, deliverability status endpoints unavailable
**Mitigation:**
- Implement Phase 6 features using available Phase 5 endpoints
- Document missing Phase 6 features in README
- Consider fixing Phase 6 backend in separate task

### Risk 2: Real-Time Updates
**Issue:** WebSocket not implemented, must use polling
**Impact:** Slightly delayed updates, increased server load
**Mitigation:**
- Use reasonable polling intervals (30-60 seconds)
- Implement manual refresh button
- Consider WebSocket implementation in future

### Risk 3: Complex State Management
**Issue:** Multiple pages need shared state (alerts, scores, circuit breakers)
**Impact:** Potential code duplication, inconsistent state
**Mitigation:**
- Create reputation store (stores/reputation.ts)
- Use Pinia for centralized state management
- Share state across components

### Risk 4: Large Page Bundle
**Issue:** Many new pages/components may increase bundle size
**Impact:** Slower initial page load
**Mitigation:**
- Use lazy loading for routes
- Code splitting by route
- Optimize imports (only what's needed)

---

## 📊 SUCCESS METRICS

**Completion Definition:**
- All 29 files implemented and tested
- All reputation pages accessible via sidebar navigation
- Setup wizard works end-to-end
- Admin dashboard displays reputation overview
- No build errors or TypeScript errors
- Manual testing passes for all features

**Performance Targets:**
- Page load time < 2 seconds (on 3G)
- API response time < 500ms (most endpoints)
- Polling interval: 30-60 seconds
- Bundle size increase < 200KB (gzipped)

**Quality Targets:**
- TypeScript strict mode enabled (no `any` types)
- ESLint passes with 0 errors
- All components responsive (mobile/tablet/desktop)
- Dark mode support (Nuxt handles automatically)
- WCAG AA accessibility compliance

---

## 🎉 CONCLUSION

This plan comprehensively restores all reputation management functionality lost during the unified refactor, including:

✅ **4 API composables** (1,000+ lines)
✅ **13 reputation pages** (3,500+ lines)
✅ **9 reputation components** (2,000+ lines)
✅ **Setup wizard** (1,000+ lines)
✅ **Navigation integration** (100+ lines)

**Total:** ~7,600 lines of code

The backend is complete with 38+ API endpoints ready for integration. This plan follows Nuxt 3 conventions, uses Shadcn UI components, and maintains consistency with existing admin UI patterns.

**Ready for implementation immediately.** 🚀
