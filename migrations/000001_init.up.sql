
CREATE TABLE vk_id_users
(
    user_id       bigserial         NOT NULL PRIMARY KEY,
    access_token  varchar(255)      NOT NULL UNIQUE
);

CREATE TABLE users
(
    id      serial               NOT NULL PRIMARY KEY,
    vk_id         bigserial               NOT NULL UNIQUE,
    first_name    varchar(255)      NOT NULL,
    last_name     varchar(255)      NOT NULL,
    username      varchar(255)      NOT NULL UNIQUE,
    email         varchar(255)      NOT NULL UNIQUE,
    CONSTRAINT fk_vk_id_users FOREIGN KEY (id) REFERENCES vk_id_users(user_id)
);


CREATE TABLE rooms
(
    id            serial               NOT NULL PRIMARY KEY,
    settings      json              NOT NULL
);

CREATE TABLE roles
(
    id       int               NOT NULL PRIMARY KEY,
    role_name          varchar(255)      NOT NULL
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