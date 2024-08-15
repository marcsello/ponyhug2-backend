-- name: CreateCardBase :one
INSERT INTO card_base (key, name, source, place)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: AssignCardImageToWearLevel :one
INSERT INTO card_wear_img (base_id, wear_level, image_url)
VALUES ($1, $2, $3)
ON CONFLICT (base_id, wear_level) DO UPDATE SET image_url = EXCLUDED.image_url
RETURNING *;


-- name: GetCardBases :many
SELECT *
FROM card_base
         LEFT JOIN card_wear_img ON card_base.id = card_wear_img.base_id;


-- name: GetPlayerCards :many
SELECT card_copy.id, card_copy.key, card_base.place, card_base.name, card_base.source, card_copy.wear_level, card_wear_img.image_url, player.name
FROM card_copy
         JOIN card_base ON card_copy.base_id = card_base.id
         JOIN card_wear_img
              ON card_copy.base_id = card_wear_img.base_id AND card_wear_img.wear_level = (
                  SELECT MAX(cwi.wear_level)
                  FROM card_wear_img cwi
                  WHERE cwi.base_id = card_copy.base_id
                    AND cwi.wear_level <= card_copy.wear_level
              )
         LEFT JOIN player ON player.id = card_copy.copied_from_player
WHERE player_id = $1;

-- name: GetCardCopy :one
SELECT card_copy.id, card_copy.key, card_base.place, card_base.name, card_base.source, card_copy.wear_level, card_wear_img.image_url, player.name
FROM card_copy
         JOIN card_base ON card_copy.base_id = card_base.id
         JOIN card_wear_img
              ON card_copy.base_id = card_wear_img.base_id AND card_wear_img.wear_level <= card_copy.wear_level
         LEFT JOIN player ON player.id = card_copy.copied_from_player
WHERE card_copy.id = $1
ORDER BY card_wear_img.wear_level DESC
LIMIT 1;


-- name: MakeFirstCopy :one
INSERT INTO card_copy (player_id, base_id, key)
VALUES ($1,
        (SELECT id FROM card_base WHERE card_base.key = $2 AND card_base.key IS NOT NULL),
        random_string(10))
ON CONFLICT (player_id, base_id) DO UPDATE SET wear_level = 0, copied_from_player = NULL
RETURNING *;


-- name: MakeSubsequentCopy :one
WITH t AS (SELECT src.base_id, src.player_id, (src.wear_level + 1) AS new_wear_level
           FROM card_copy AS src
           WHERE src.key = $2)
INSERT
INTO card_copy (player_id, base_id, copied_from_player, wear_level, key)
VALUES ($1,
        (SELECT base_id FROM t),
        (SELECT t.player_id FROM t),
        (SELECT new_wear_level FROM t),
        random_string(10)) -- hopefully this does not result in a conflict
ON CONFLICT (player_id, base_id) DO UPDATE SET wear_level = EXCLUDED.wear_level
WHERE card_copy.wear_level > EXCLUDED.wear_level
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

-- name: DemotePlayer :exec
UPDATE player
SET is_admin = false
WHERE id = $1;