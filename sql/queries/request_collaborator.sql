-- name: GetRequestCollaboratorsByAdminID :many
SELECT * FROM request_collaborators
WHERE admin_id = $1 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateRequestColloboratorsAdmin :exec
UPDATE request_collaborators SET admin_id = $2
WHERE request_id = $1;


-- name: UpdateRequestColloboratorsValidators :exec
UPDATE request_collaborators SET validators = $2
WHERE request_id = $1;

-- name: GetCountRequestCollaboratorsByAdminID :one
SELECT COUNT(*) FROM request_collaborators
WHERE admin_id = $1;