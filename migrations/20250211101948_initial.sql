-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "user"
(
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(255)                        NOT NULL UNIQUE,
    password   VARCHAR(255)                        NOT NULL,
    coins      int       DEFAULT 1000              NOT NULL,
    is_active  BOOLEAN   DEFAULT TRUE              NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE user_token
(
    id        SERIAL PRIMARY KEY,
    jwt       text                       NOT NULL,
    user_id   int REFERENCES "user" (id) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE merch
(
    id    SERIAL PRIMARY KEY,
    name  VARCHAR(255) NOT NULL,
    price int          NOT NULL
);

INSERT INTO merch (name, price)
VALUES ('t-shirt', 80),
       ('cup', 20),
       ('book', 50),
       ('pen', 10),
       ('powerbank', 200),
       ('hoody', 300),
       ('umbrella', 200),
       ('socks', 10),
       ('wallet', 50),
       ('pink-hoody	', 500)
;

CREATE TABLE user_merch
(
    id       SERIAL PRIMARY KEY,
    user_id  INT REFERENCES "user" (id) NOT NULL,
    merch_id INT REFERENCES merch (id)  NOT NULL
);

CREATE TABLE coin_transfer
(
    id           serial primary key,
    user_id_from int references "user" (id) NOT NULL,
    user_id_to   int references "user" (id) NOT NULL,
    coins        int                        NOT NULL
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE coin_transfer;
DROP TABLE user_merch;
DROP TABLE merch;
DROP TABLE user_token;
DROP TABLE "user";
-- +goose StatementEnd
