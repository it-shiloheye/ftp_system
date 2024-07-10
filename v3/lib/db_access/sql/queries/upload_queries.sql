

-- name: UploadFilesStepOneUploadData :one
INSERT INTO file_data(
    mod_time,
    file_size,
    file_data_b
) VALUES ($1, $2, $3)
ON CONFLICT (file_hash) DO NOTHING
RETURNING 
    id, file_hash;


-- name: UploadFilesStepTwoUploadMetadata :one
INSERT INTO file_metadata(
    file_path, -- relative
    file_type,
    file_state, 
    file_data_id, -- file_data(id)
    file_mode,
    mod_time
) VALUES ($1,$2,'current',$3,$4, $5)
ON CONFLICT (file_path)
DO UPDATE
    SET 
        file_data_id = $4,
        mod_time = $6
RETURNING 
    id;



-- name: UpdateFileTrackerMarkUploaded :exec
INSERT INTO file_tracker (
    peer_id,
    file_meta_id,
    current_hash_id,
    file_state
) VALUES ($1, $2, $3, 'uploaded');


-- name: CheckChangesStepOne :exec
SELECT DISTINCT file_metadata.file_path, file_metadata.file_data_id, file_data.file_hash,file_metadata.mod_time FROM file_metadata
JOIN file_data ON file_metadata.file_data_id = file_data.id
ORDER BY
    file_metadata.file_data_id DESC,
    file_metadata.mod_time DESC;