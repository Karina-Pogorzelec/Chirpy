-- name: DeleteChirp :exec
DELETE FROM chirps
where user_id = $1 and id = $2;