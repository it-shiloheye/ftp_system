

-- name: UploadStoreStepOnePeerDir :one
INSERT INTO peer_dirs(peer_id,dir_path)
VALUES ($1, $2)
ON CONFLICT (uniq_dirs) DO NOTHING
RETURNING   
    id;

-- name: UploadStoreStepTwoUploadFile :one
INSERT INTO file_data(
    mod_time,
    file_size,
    file_data_b
) VALUES ($1, $2, $3)
ON CONFLICT (file_hash) DO NOTHING
RETURNING 
    id, file_hash;

-- name: UploadStoreStepThreeUpdateMetadata :one
INSERT INTO file_metadata(
    file_path, -- relative
    file_type,
    file_state, 
    file_data_id, -- file_data(id)
    file_mode,
    mod_time,
    dir_id
) VALUES ($1,$2,'store',$3,$4, $5, $6)
ON CONFLICT (file_path)
DO UPDATE
    SET 
        file_data_id = $4,
        mod_time = $6,
        file_state = 'store'
RETURNING 
    id;


-- name: UpdateFileTrackerMarkStored :exec
INSERT INTO file_tracker (
    peer_id,
    file_meta_id,
    current_hash_id,
    file_state
) VALUES ($1, $2, $3, 'stored');

-- name: DownloadStoreBulk :many
SELECT * from file_metadata
JOIN file_tracker on file_tracker.file_meta_id = file_metadata.id
JOIN file_data on file_tracker.current_hash_id = file_data.id
WHERE 
	current_hash_id NOT IN (
		SELECT current_hash_id from file_metadata 
		where file_state = 'stored'
		AND peer_id = $1
	);
