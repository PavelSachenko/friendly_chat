CREATE TABLE chats
(
    id            SERIAL       NOT NULL UNIQUE,
    created_at    timestamp             DEFAULT now(),
    updated_at    timestamp             DEFAULT now(),
    type          SMALLINT DEFAULT  1,
    description  varchar(500) NOT NULL DEFAULT '',
    PRIMARY KEY (id)
);

CREATE TABLE labels_chats
(
    id            SERIAL       NOT NULL UNIQUE,
    title    varchar(500) NOT NULL DEFAULT '',
    user_id  INTEGER NOT NULL,
    chat_id  INTEGER NOT NULL,

    UNIQUE (user_id, chat_id)
);

CREATE TABLE users_chats
(
    id            SERIAL       NOT NULL UNIQUE,
    user_id  INTEGER NOT NULL,
    chat_id  INTEGER NOT NULL,
    is_owner BOOLEAN DEFAULT FALSE,

    UNIQUE (user_id, chat_id)
);

CREATE TABLE messages_chats
(
    id            SERIAL       NOT NULL UNIQUE,
    chat_id  INTEGER NOT NULL,
    sender_user_id INTEGER NOT NULL,
    body TEXT NOT NULL ,
    created_at    timestamp             DEFAULT now(),
    updated_at    timestamp             DEFAULT now()
);

CREATE TABLE users_unread_messages
(
    id            SERIAL       NOT NULL UNIQUE,
    chat_id  INTEGER NOT NULL,
    message_id  INTEGER NOT NULL,
    user_id  INTEGER NOT NULL,

    UNIQUE (user_id, chat_id, message_id)
);

CREATE TABLE chats_messages_attachments
(
    id            SERIAL       NOT NULL UNIQUE,
    chat_id  INTEGER NOT NULL,
    attachment varchar(500) NOT NULL
);

alter table labels_chats
    add constraint fk_chat_id
        foreign key (chat_id)
            REFERENCES chats (id)
            ON DELETE CASCADE;

alter table users_chats
    add constraint fk_chat_id
        foreign key (chat_id)
            REFERENCES chats (id)
            ON DELETE CASCADE;

alter table messages_chats
    add constraint fk_chat_id
        foreign key (chat_id)
            REFERENCES chats (id)
            ON DELETE CASCADE;


alter table users_unread_messages
    add constraint fk_chat_id
        foreign key (chat_id)
            REFERENCES chats (id)
            ON DELETE CASCADE,
    add constraint fk_message_id
        foreign key (message_id)
            REFERENCES messages_chats (id)
            ON DELETE CASCADE;

alter table chats_messages_attachments
    add constraint fk_chat_id
        foreign key (chat_id)
            REFERENCES chats (id)
            ON DELETE CASCADE;