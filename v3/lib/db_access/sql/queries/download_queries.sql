

-- name: DownloadFileDataWithFilePath :one
SELECT file_data.id, file_hash, file_data.mod_time, file_size, file_data_b, file_data.creation_time, 
file_metadata.id AS file_metadata_id,
dir_id, file_path, file_type, file_state, file_data_id, file_mode, file_metadata.mod_time, file_metadata.creation_time
from file_data 
JOIN file_metadata ON file_metadata.file_data_id = file_data.id
where file_metadata.file_path = $1;

-- name: DownloadFileStepOneGetLatestData :one
SELECT file_data.id, file_hash, file_data.mod_time, file_size, file_data_b, file_data.creation_time, 
file_metadata.id AS file_metadata_id,
dir_id, file_path, file_type, file_state, file_data_id, file_mode, file_metadata.mod_time, file_metadata.creation_time
from file_data 
JOIN file_metadata ON file_metadata.file_data_id = file_data.id
WHERE file_data.file_hash = $1;


-- name: DownloadFileBulkGetLatest :many
SELECT DISTINCT file_path, 
dir_id, file_hash, file_data.mod_time, file_size, file_data_b,
file_metadata.id AS file_metadata_id,
 file_type, file_state, file_data_id, file_mode
FROM file_metadata 
JOIN file_data ON file_data.id = file_metadata.file_data_id
WHERE file_state = 'current';

-- name: UpdateFileTrackerMarkDownloaded :exec
INSERT INTO file_tracker (
    peer_id,
    file_meta_id,
    current_hash_id,
    file_state
) VALUES ($1, $2, $3, 'downloaded');
