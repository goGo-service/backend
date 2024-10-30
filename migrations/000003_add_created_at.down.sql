-- Удаление триггеров
DROP TRIGGER IF EXISTS set_timestamp_users ON users;
DROP TRIGGER IF EXISTS set_timestamp_rooms ON rooms;
DROP TRIGGER IF EXISTS set_timestamp_roles ON roles;
DROP TRIGGER IF EXISTS set_timestamp_rooms_users ON rooms_users;

-- Удаление функции обновления updated_at
DROP FUNCTION IF EXISTS update_updated_at_column;

-- Удаление колонок created_at и updated_at
ALTER TABLE users
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE rooms
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE roles
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE rooms_users
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;
