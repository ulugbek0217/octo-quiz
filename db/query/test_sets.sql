-- name: CreateTestSet :one
INSERT INTO test_sets (
    test_set_name, creator_id, is_public, time_limit
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: DeleteTestSet :exec
DELETE FROM test_sets
WHERE test_set_id = $1;

-- name: MakeTestSetPublic :exec
UPDATE test_sets
SET is_public = TRUE
WHERE test_set_id = $1;
