-- Migration: CalDAV and CardDAV support
-- Description: Add tables for calendars, events, addressbooks, and contacts

-- Calendars table
CREATE TABLE IF NOT EXISTS calendars (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    display_name TEXT,
    color TEXT,
    description TEXT,
    timezone TEXT DEFAULT 'UTC',
    sync_token TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, name)
);

CREATE INDEX IF NOT EXISTS idx_calendars_user_id ON calendars(user_id);
CREATE INDEX IF NOT EXISTS idx_calendars_sync_token ON calendars(sync_token);

-- Events table
CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    calendar_id INTEGER NOT NULL,
    uid TEXT NOT NULL,
    summary TEXT,
    description TEXT,
    location TEXT,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    all_day INTEGER DEFAULT 0,
    rrule TEXT,
    attendees TEXT,
    organizer TEXT,
    status TEXT DEFAULT 'CONFIRMED',
    sequence INTEGER DEFAULT 0,
    etag TEXT NOT NULL,
    ical_data TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (calendar_id) REFERENCES calendars(id) ON DELETE CASCADE,
    UNIQUE(calendar_id, uid)
);

CREATE INDEX IF NOT EXISTS idx_events_calendar_id ON events(calendar_id);
CREATE INDEX IF NOT EXISTS idx_events_uid ON events(uid);
CREATE INDEX IF NOT EXISTS idx_events_etag ON events(etag);
CREATE INDEX IF NOT EXISTS idx_events_time_range ON events(calendar_id, start_time, end_time);

-- Addressbooks table
CREATE TABLE IF NOT EXISTS addressbooks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    display_name TEXT,
    description TEXT,
    sync_token TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, name)
);

CREATE INDEX IF NOT EXISTS idx_addressbooks_user_id ON addressbooks(user_id);
CREATE INDEX IF NOT EXISTS idx_addressbooks_sync_token ON addressbooks(sync_token);

-- Contacts table
CREATE TABLE IF NOT EXISTS contacts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    addressbook_id INTEGER NOT NULL,
    uid TEXT NOT NULL,
    fn TEXT NOT NULL,
    email TEXT,
    tel TEXT,
    etag TEXT NOT NULL,
    vcard_data TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (addressbook_id) REFERENCES addressbooks(id) ON DELETE CASCADE,
    UNIQUE(addressbook_id, uid)
);

CREATE INDEX IF NOT EXISTS idx_contacts_addressbook_id ON contacts(addressbook_id);
CREATE INDEX IF NOT EXISTS idx_contacts_uid ON contacts(uid);
CREATE INDEX IF NOT EXISTS idx_contacts_etag ON contacts(etag);
CREATE INDEX IF NOT EXISTS idx_contacts_fn ON contacts(fn);
CREATE INDEX IF NOT EXISTS idx_contacts_email ON contacts(email);
