-- name: GetItems :many
SELECT * FROM items;

-- name: CreateItem :one
INSERT INTO items (
  name, capacity
) VALUES (
  @name, @capacity
) RETURNING *;

-- name: UpdateItem :one
UPDATE items
SET name = @name, capacity = @capacity
WHERE id = @id
RETURNING *;

-- name: DeleteItem :exec
DELETE FROM items
WHERE id = @id;