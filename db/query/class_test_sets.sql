-- name: AddTestSetToClass :exec
INSERT INTO class_test_sets (
    class_id, test_set_id
) VALUES (
    $1, $2
) ON CONFLICT DO NOTHING;

-- -- name: GetTestSetsByClassID :many
-- SELECT * FROM class_test_sets
-- WHERE class_id = $1;

-- name: ListTestSetsByClassID :many
SELECT
  ts.test_set_id,
  ts.test_set_name,
  ts.time_limit
FROM
  test_sets ts
INNER JOIN
  class_test_sets cts ON cts.test_set_id = ts.test_set_id
WHERE
  cts.class_id = $1
LIMIT $2
OFFSET $3;

-- name: DeleteTestSetFromClass :exec
DELETE FROM class_test_sets
WHERE test_set_id = $1;
