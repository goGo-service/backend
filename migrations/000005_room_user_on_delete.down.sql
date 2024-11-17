ALTER TABLE rooms_users
    DROP CONSTRAINT fk_user,
    DROP CONSTRAINT fk_room,
    DROP CONSTRAINT fk_role;

ALTER TABLE rooms_users
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    ADD CONSTRAINT fk_room FOREIGN KEY (room_id) REFERENCES rooms(id),
    ADD CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles(id);
