CREATE TABLE IF NOT EXISTS card_base
(
    id     SMALLSERIAL PRIMARY KEY,
    key    VARCHAR(9)      NULL,    -- these cards can only be obtained by some special actions

    name   VARCHAR(64)     NOT NULL,
    source VARCHAR(255)    NULL,

    place  SMALLINT UNIQUE NOT NULL -- known as "order" in PonyHug1
);

CREATE TABLE IF NOT EXISTS card_wear_img
(
    base_id    SMALLINT     NOT NULL,
    wear_level SMALLINT     NOT NULL,
    image_url  VARCHAR(255) NOT NULL,
    PRIMARY KEY (base_id, wear_level),

    CONSTRAINT card_wear_base
        FOREIGN KEY (base_id)
            REFERENCES card_base (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS player
(
    id         SERIAL PRIMARY KEY,
    name       varchar(64) UNIQUE NOT NULL,
    registered TIMESTAMP          NOT NULL DEFAULT now(),

    is_admin   BOOLEAN            NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS card_copy
(
    id                 SERIAL PRIMARY KEY,
    player_id          INTEGER            NOT NULL,
    base_id            SMALLINT           NOT NULL,

    copied_from_player INTEGER            NULL,
    timestamp          TIMESTAMP          NOT NULL DEFAULT now(),
    wear_level         SMALLINT           NOT NULL DEFAULT 0, -- 0 indicates an "original copy"
    key                VARCHAR(10) UNIQUE NOT NULL,

    UNIQUE (player_id, base_id),                              -- <- include wear level here to allow having multiple cards of the same base

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

CREATE OR REPLACE FUNCTION random_string(length integer) RETURNS TEXT AS
$$
DECLARE
    chars  text[]  := '{2,3,4,5,6,7,8,9,A,B,C,D,E,F,G,H,I,J,K,L,M,N,P,Q,R,S,T,U,V,W,X,Y,Z}';
    result text    := '';
    i      integer := 0;
BEGIN
    if length < 0 then
        raise exception 'Given length cannot be less than 0';
    end if;
    for i in 1..length
        loop
            result := result || chars[1 + random() * (array_length(chars, 1) - 1)];
        end loop;
    return result;
END;
$$ language plpgsql;