CREATE TABLE users
(
    id            serial            NOT NULL PRIMARY KEY,
    vk_id         bigint            NOT NULL UNIQUE,
    first_name    varchar(255)      NOT NULL,
    last_name     varchar(255)      NOT NULL,
    username      varchar(255)      NOT NULL,
    email         varchar(255)      NOT NULL UNIQUE
);

CREATE TABLE rooms
(
    id            serial            NOT NULL PRIMARY KEY,
    settings      json              NOT NULL
);

CREATE TABLE roles
(
    id            serial            NOT NULL PRIMARY KEY,
    role_name     varchar(255)      NOT NULL
);

CREATE TABLE rooms_users
(
    user_id       int               NOT NULL,
    room_id       int               NOT NULL,
    role_id       int               NOT NULL,
    PRIMARY KEY (user_id, room_id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_room FOREIGN KEY (room_id) REFERENCES rooms(id),
    CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles(id)
);
