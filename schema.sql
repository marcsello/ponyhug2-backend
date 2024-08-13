CREATE TABLE card_base
(
    id        SERIAL PRIMARY KEY,
    key       VARCHAR(9),

    name      VARCHAR(64)    NOT NULL,
    story     TEXT           NULL,
    source    VARCHAR(255)   NULL,
    image_url VARCHAR(255)   NOT NULL,

    place     INTEGER UNIQUE NOT NULL -- known as "order" in PonyHug1
);


CREATE TABLE card_copy
(
    id         SERIAL PRIMARY KEY,
    player_id  SERIAL             NOT NULL,
    base_id    SERIAL             NOT NULL,

    timestamp  TIMESTAMP          NOT NULL DEFAULT now(),
    wear_level SMALLINT           NOT NULL DEFAULT 0, -- 0 indicates an "original copy"
    key        VARCHAR(10) UNIQUE NOT NULL,

    UNIQUE (player_id, base_id),                      -- <- include wear level here to allow having multiple cards of the same base

    CONSTRAINT card_copy_base
        FOREIGN KEY (base_id)
            REFERENCES card_base (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE,

    CONSTRAINT card_player
        FOREIGN KEY (player_id)
            REFERENCES player (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE
);



CREATE TABLE player
(
    id         SERIAL PRIMARY KEY,
    name       varchar(60) UNIQUE NOT NULL,
    registered TIMESTAMP          NOT NULL DEFAULT now(),

    is_admin   BOOLEAN            NOT NULL
);
