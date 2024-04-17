-- name: CreateRequest :one
INSERT INTO requests (
  id,
  requester_id,
  full_name,
  age,
  location_lat,
  location_lon,
  status,
  photo,
  msgs
) VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8,$9
) RETURNING *;

-- name: GetRequestByID :one
SELECT * FROM requests
WHERE id = $1;

-- name: GetRequestByRequesterID :many
SELECT * FROM requests
WHERE requester_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetOrphanRequests :many
SELECT * FROM requests
WHERE status = 'admin_pending' 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateRequestStatus :exec
UPDATE requests SET status = $2
WHERE id = $1;

-- name: UpdateRequestQuestionAnswers :exec
UPDATE requests SET question_answers = $2
WHERE id = $1;

-- name: GetRequestQuestionAnswers :one
SELECT question_answers FROM requests
WHERE id = $1;
