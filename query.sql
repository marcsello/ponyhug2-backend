-- name: GetCardBase :one
SELECT *
FROM card_base
WHERE id = $1
LIMIT 1;

-- name: GetPlayerCards :many
SELECT sqlc.embed(card_copy), sqlc.embed(card_base)
FROM card_copy
         JOIN card_base ON card_copy.base_id = card_base.id
WHERE player_id = $1;


-- name: MakeFirstCopy :one
INSERT INTO card_copy (player_id, base_id, key)
VALUES ($1,
        (SELECT id FROM card_base WHERE card_base.key = $2),
        random_string(10))
RETURNING *;

-- name: MakeSubsequentCopy :one
INSERT INTO card_copy (player_id, base_id, wear_level, key)
WITH
    t AS (
    SELECT card_base.id as base_id, card_copy.wear_level+1 as new_wear_level
    FROM card_copy JOIN card_base ON card_copy.base_id = card_base.id
    WHERE card_copy.key = $2
    )
VALUES ($1,
        t.base_id,
        t.new_wear_level,
        random_string(10)) -- hopefully this does not result in a conflict
RETURNING *;


-- Player stuffs

-- name: GetPlayers :many
SELECT *
FROM player;

-- name: GetPlayer :one
SELECT *
FROM player
WHERE id = $1
LIMIT 1;

-- name: CreatePlayer :one
INSERT INTO player (name)
VALUES ($1)
RETURNING *;

-- name: PromotePlayer :exec
UPDATE player
SET is_admin = true
WHERE id = $1;