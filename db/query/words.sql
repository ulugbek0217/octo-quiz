-- name: InsertWords :one
INSERT INTO words (
    test_set_id, english_word, uzbek_word
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: DeleteWords :exec
DELETE FROM words
WHERE words_id = $1;