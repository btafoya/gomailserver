package database

const migrationV4Up = `
-- Add role column to users table for admin/user distinction
ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'user' CHECK(role IN ('admin', 'user'));

-- Create index on role for faster admin queries
CREATE INDEX idx_users_role ON users(role);
`

const migrationV4Down = `
-- Remove role column index
DROP INDEX IF EXISTS idx_users_role;

-- Note: SQLite doesn't support ALTER TABLE DROP COLUMN directly
-- In production, you would need to recreate the table without the role column
-- For development, this is a reminder that downgrades are not fully supported
`
