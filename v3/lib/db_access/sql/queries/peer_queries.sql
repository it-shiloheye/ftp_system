
-- name: CreateClient :one
INSERT INTO peers_table(ip_address,peer_name,pem)
VALUES ($1,$2,$3)
RETURNING 
    peer_id, peer_role;

-- name: ConnectClient :many
SELECT * FROM peers_table
WHERE
    peer_id = $1
LIMIT 1;

-- name: GetPEMs :many
SELECT pem FROM peers_table;



-- name: GetFilesList :many
SELECT DISTINCT file_path, dir_id,  file_type, file_state, file_data.id, 
file_metadata.id as file_metadata_id,file_mode, 
file_metadata.mod_time,  file_data.file_hash, file_data.file_size FROM file_metadata
JOIN file_data ON file_metadata.file_data_id = file_data.id
WHERE file_state = 'current'
ORDER BY file_metadata.mod_time DESC;

-- name: UpdateFileTracker :exec
INSERT INTO file_tracker (
    peer_id,
    file_meta_id,
    current_hash_id,
    file_state
) VALUES ($1, $2, $3, $4);

-- name: DBTotalStorageSize :one 
SELECT SUM(file_size) FROM file_data;

-- name: DBCurrentStorageSize :one
SELECT SUM(file_size) FROM file_metadata 
JOIN file_data ON file_metadata.file_data_id = file_data.id
WHERE file_metadata.file_state = 'current';




