

-- name: ConnectClient :one
SELECT *
FROM peers_table
WHERE peer_id = $1;

-- name: CreateClient :one
INSERT INTO peers_table(peer_id,peer_role,peer_name,pem, peer_config,ip_address)
VALUES (default, $1, $2, $3, $4, $5)
RETURNING *;

-- name: AddDirForPeer :batchone
INSERT INTO peer_dirs(
    peer_id,
    dir_path
) VALUES ($1, $2)
RETURNING *;

-- name: UploadFileData :batchone
INSERT INTO file_data(
    file_status,
    file_data,
    modification_date
) VALUES ($1, $2,$3)
ON CONFLICT (file_hash)
DO UPDATE 
    SET 
        file_data = $2
RETURNING 
    file_data_id, file_hash;

-- name: UploadMetadata :batchone
INSERT INTO file_metadata(
    peer_id,
    dir_id,
    file_data_id,
    file_name,
    file_path,
    file_type
) VALUES($1, $2, $3, $4, $5, $6)
RETURNING id;


-- name: DownloadFileData :batchmany
SELECT 
    *
FROM 
    file_data 
    INNER JOIN file_metadata USING (file_data_id)
WHERE 
    file_metadata.peer_id = $1
    AND file_data.file_data IS NOT NULL
    AND file_data.modification_date >= $2
ORDER BY
    modification_date DESC;

