-- Migration: CalDAV Tasks (VTODO) support
-- Description: Add table for CalDAV tasks/todos

-- Tasks table for VTODO components
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    calendar_id INTEGER NOT NULL,
    uid TEXT NOT NULL,
    summary TEXT,
    description TEXT,
    location TEXT,
    due DATETIME,
    start DATETIME,
    completed DATETIME,
    status TEXT DEFAULT 'NEEDS-ACTION',
    priority INTEGER DEFAULT 0,
    percent INTEGER DEFAULT 0,
    categories TEXT,
    organizer TEXT,
    attendees TEXT,
    sequence INTEGER DEFAULT 0,
    etag TEXT NOT NULL,
    ical_data TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (calendar_id) REFERENCES calendars(id) ON DELETE CASCADE,
    UNIQUE(calendar_id, uid)
);

CREATE INDEX IF NOT EXISTS idx_tasks_calendar_id ON tasks(calendar_id);
CREATE INDEX IF NOT EXISTS idx_tasks_uid ON tasks(uid);
CREATE INDEX IF NOT EXISTS idx_tasks_etag ON tasks(etag);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(calendar_id, status);
CREATE INDEX IF NOT EXISTS idx_tasks_due ON tasks(calendar_id, due);
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(calendar_id, priority);
