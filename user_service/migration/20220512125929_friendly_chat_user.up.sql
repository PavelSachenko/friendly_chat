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
