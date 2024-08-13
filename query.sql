-- name: GetCardBase :one
SELECT * FROM card_base
WHERE id = $1 LIMIT 1;

-- name: GetPlayer :one
SELECT * FROM player
WHERE id = $1 LIMIT 1;