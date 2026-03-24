-- name: GetNextURLID :one

SELECT nextval('public.urls_id_seq');
-- name: CreateUrl :exec

INSERT INTO URLS ( id, code, original_url, counter, updated_at) 
VALUES ($1,$2,$3,$4 , NOW());

-- name: GetOriginalUrl :one

SELECT
    original_url
FROM URLS
WHERE code = $1;


-- name: GetStats :one

SELECT 
    original_url,
    counter,
    created_at
FROM URLS
WHERE code = $1;


-- name: UpdateCounter :exec

UPDATE URLS
SET counter = counter + $2
WHERE code = $1;