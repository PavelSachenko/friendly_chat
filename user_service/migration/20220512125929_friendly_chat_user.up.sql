CREATE TABLE users
(
    id            SERIAL       NOT NULL UNIQUE,
    created_at    timestamp             DEFAULT now(),
    updated_at    timestamp             DEFAULT now(),
    username          varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    description  varchar(500) NOT NULL DEFAULT '',
    avatar        varchar(500) NOT NULL DEFAULT '',

    PRIMARY KEY (id)
);

CREATE TABLE roles
(
    id          SERIAL       NOT NULL UNIQUE,
    title       varchar(20)  NOT NULL UNIQUE,
    description varchar(255) NOT NULL default '',

    PRIMARY KEY (id)
);

INSERT INTO roles (id,title)
VALUES (1,'admin'),
       (2,'client');

CREATE TABLE users_roles
(
    id      SERIAL  NOT NULL UNIQUE,
    user_id integer NOT NULL,
    role_id integer NOT NULL,
    PRIMARY KEY (id)
);

alter table users_roles
    add constraint fk_user
        foreign key (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE;

alter table users_roles
    add constraint fk_role
        foreign key (role_id)
            REFERENCES roles (id)
            ON DELETE SET NULL;

INSERT INTO users (id, username, password_hash)
VALUES
    (1, 'pavel', '4aa3366be51a719c1b2db3c207dabb2c1015f528'),
    (2, 'vlad', '4aa3366be51a719c1b2db3c207dabb2c1015f528'),
    (3, 'anton', '4aa3366be51a719c1b2db3c207dabb2c1015f528'),
    (4, 'test', '4aa3366be51a719c1b2db3c207dabb2c1015f528'),
    (5, 'masha', '4aa3366be51a719c1b2db3c207dabb2c1015f528'),
    (6, 'tihon', '4aa3366be51a719c1b2db3c207dabb2c1015f528');

INSERT INTO users_roles (user_id, role_id)
VALUES
(1, 2),
(2, 2),
(3, 2),
(4, 2),
(5, 2),
(6, 2);
