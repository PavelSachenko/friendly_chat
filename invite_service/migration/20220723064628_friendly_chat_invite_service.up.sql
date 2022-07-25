CREATE TABLE users_friends
(
    id             SERIAL  NOT NULL UNIQUE,
    user_id        INTEGER NOT NULL,
    friend_user_id INTEGER NOT NULL,
    created_at     timestamp DEFAULT now(),
    PRIMARY KEY (id)
);

CREATE TABLE users_invites
(
    id               SERIAL  NOT NULL UNIQUE,
    sender_id        INTEGER NOT NULL,
    receiver_user_id INTEGER NOT NULL,
    description      varchar(500),
    created_at       timestamp DEFAULT now(),
    PRIMARY KEY (id)
);

alter table chats_messages_attachments
    add constraint fk_chat_id
        foreign key (chat_id)
            REFERENCES chats (id)
            ON DELETE CASCADE;