# Web UI Fixes

This document details the fixes that were applied to the web UI to get it to build successfully.

## Admin UI (`web/admin`)

*   **Upgraded `tailwindcss` to v4:** The `tailwindcss` package was upgraded from v3 to v4 to fix a build issue.
*   **Installed `@tailwindcss/postcss`:** The `@tailwindcss/postcss` package was installed to support Tailwind CSS v4.
*   **Updated `postcss.config.js`:** The `postcss.config.js` file was updated to use `@tailwindcss/postcss`.
*   **Fixed `style.css`:** The `style.css` file was updated to be compatible with Tailwind CSS v4. This involved adding the `@config` directive and removing a global `*` style.

## Webmail UI (`web/webmail`)

*   **Downgraded `tailwindcss` to v3:** The `tailwindcss` package was downgraded from v4 to v3 to fix a build issue.
*   **Downgraded `@nuxtjs/tailwindcss`:** The `@nuxtjs/tailwindcss` package was downgraded to a version compatible with Tailwind CSS v3.
*   **Fixed `main.css`:** The `assets/css/main.css` file was updated to be compatible with Tailwind CSS v3. This involved adding the `@config` directive and removing a global `*` style.

## Go Backend

*   **Fixed `webmail_calendar.go`:**
    *   Corrected calls to `middleware.GetUserID`.
    *   Corrected a call to `eventService.GetEventsInRange`.
    *   Added a `Timezone` field to the `Event` struct in `internal/calendar/domain/event.go`.
    *   Updated the `CreateEvent` function in `internal/calendar/service/event.go` to handle the new `Timezone` field.
*   **Fixed `webmail_contacts.go`:**
    *   Corrected calls to `middleware.GetUserID`.
*   **Fixed `router.go`:**
    *   Corrected the types of the contact and calendar services in the `RouterConfig` struct.
*   **Fixed `server.go`:**
    *   Corrected the types of the contact and calendar services in the `NewServer` function.
*   **Fixed `migration_v7.go`:**
    *   Added `IF NOT EXISTS` to the `CREATE INDEX` statements to prevent errors if the indexes already exist.
*   **Fixed `migrations.go`:**
    *   Reverted the `splitSQL` function to its original implementation.

# WebUI Remaining Issues

   1. The server fails to start with a persistent "disk I/O error", even when using an in-memory database.
   2. I am unable to read the source code in the cmd/gomailserver directory, which prevents me from debugging the startup issue.