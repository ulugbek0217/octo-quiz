-- name: AddTestSetToClass :exec
INSERT INTO class_test_sets (
    class_id, test_set_id
) VALUES (
    $1, $2
) ON CONFLICT DO NOTHING;

-- name: DeleteTestSetFromClass :exec
DELETE FROM class_test_sets
WHERE test_set_id = $1;
