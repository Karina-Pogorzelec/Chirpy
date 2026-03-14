-- name: GetRefreshToken :one
select * from refresh_tokens 
join users on users.id = refresh_tokens.user_id
where token = $1
    AND expires_at > NOW()
    and revoked_at is null;