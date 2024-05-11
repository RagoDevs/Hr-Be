-- name: CreateAnnoucemnt :exec
INSERT INTO announcement (
    description,
    announcement_date,
    created_by
)
VALUES ($1, $2, $3);


-- name: GetAnnouncement :one
SELECT *
FROM  announcement
WHERE id = $1;


-- name: GetAnnouncements :many
SELECT 
    a.id, 
    a.description, 
    a.announcement_date, 
    e.name AS created_by,
    a.created_at 

FROM  announcement a  
JOIN employee e ON e.id = a.created_by
ORDER BY 
   a.created_at DESC;


-- name: UpdateAnnouncement :exec
UPDATE announcement
SET description = $1,
    announcement_date = $2
WHERE id = $3;


-- name: DeleteAnnouncement :exec
DELETE FROM announcement WHERE id = $1;