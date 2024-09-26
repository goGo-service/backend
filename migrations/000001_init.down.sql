-- Удаление таблицы room_user (сначала удаляем таблицы с внешними ключами)
DROP TABLE IF EXISTS rooms_users;

-- Удаление таблицы User
DROP TABLE IF EXISTS users;

-- Удаление таблицы vk_id_user
DROP TABLE IF EXISTS vk_id_users;

-- Удаление таблицы room
DROP TABLE IF EXISTS rooms;

-- Удаление таблицы role
DROP TABLE IF EXISTS roles;