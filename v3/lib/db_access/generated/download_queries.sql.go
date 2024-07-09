// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: download_queries.sql

package db_access

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const downloadFileBulkGetLatest = `-- name: DownloadFileBulkGetLatest :many
SELECT DISTINCT file_path, 
dir_id, file_hash, file_data.mod_time, file_size, file_data_b,
file_metadata.id AS file_metadata_id,
 file_type, file_state, file_data_id, file_mode
FROM file_metadata 
JOIN file_data ON file_data.id = file_metadata.file_data_id
WHERE file_state = 'current'
`

type DownloadFileBulkGetLatestRow struct {
	FilePath       string             `json:"file_path"`
	DirID          *int32             `json:"dir_id"`
	FileHash       *string            `json:"file_hash"`
	ModTime        pgtype.Timestamptz `json:"mod_time"`
	FileSize       int32              `json:"file_size"`
	FileDataB      []byte             `json:"file_data_b"`
	FileMetadataID int32              `json:"file_metadata_id"`
	FileType       string             `json:"file_type"`
	FileState      string             `json:"file_state"`
	FileDataID     int32              `json:"file_data_id"`
	FileMode       int32              `json:"file_mode"`
}

// DownloadFileBulkGetLatest
//
//	SELECT DISTINCT file_path,
//	dir_id, file_hash, file_data.mod_time, file_size, file_data_b,
//	file_metadata.id AS file_metadata_id,
//	 file_type, file_state, file_data_id, file_mode
//	FROM file_metadata
//	JOIN file_data ON file_data.id = file_metadata.file_data_id
//	WHERE file_state = 'current'
func (q *Queries) DownloadFileBulkGetLatest(ctx context.Context, db DBTX) ([]*DownloadFileBulkGetLatestRow, error) {
	rows, err := db.Query(ctx, downloadFileBulkGetLatest)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*DownloadFileBulkGetLatestRow{}
	for rows.Next() {
		var i DownloadFileBulkGetLatestRow
		if err := rows.Scan(
			&i.FilePath,
			&i.DirID,
			&i.FileHash,
			&i.ModTime,
			&i.FileSize,
			&i.FileDataB,
			&i.FileMetadataID,
			&i.FileType,
			&i.FileState,
			&i.FileDataID,
			&i.FileMode,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const downloadFileDataWithFilePath = `-- name: DownloadFileDataWithFilePath :one
SELECT file_data.id, file_hash, file_data.mod_time, file_size, file_data_b, file_data.creation_time, 
file_metadata.id AS file_metadata_id,
dir_id, file_path, file_type, file_state, file_data_id, file_mode, file_metadata.mod_time, file_metadata.creation_time
from file_data 
JOIN file_metadata ON file_metadata.file_data_id = file_data.id
where file_metadata.file_path = $1
`

type DownloadFileDataWithFilePathRow struct {
	ID             int32              `json:"id"`
	FileHash       *string            `json:"file_hash"`
	ModTime        pgtype.Timestamptz `json:"mod_time"`
	FileSize       int32              `json:"file_size"`
	FileDataB      []byte             `json:"file_data_b"`
	CreationTime   pgtype.Timestamptz `json:"creation_time"`
	FileMetadataID int32              `json:"file_metadata_id"`
	DirID          *int32             `json:"dir_id"`
	FilePath       string             `json:"file_path"`
	FileType       string             `json:"file_type"`
	FileState      string             `json:"file_state"`
	FileDataID     int32              `json:"file_data_id"`
	FileMode       int32              `json:"file_mode"`
	ModTime_2      pgtype.Timestamptz `json:"mod_time_2"`
	CreationTime_2 pgtype.Timestamptz `json:"creation_time_2"`
}

// DownloadFileDataWithFilePath
//
//	SELECT file_data.id, file_hash, file_data.mod_time, file_size, file_data_b, file_data.creation_time,
//	file_metadata.id AS file_metadata_id,
//	dir_id, file_path, file_type, file_state, file_data_id, file_mode, file_metadata.mod_time, file_metadata.creation_time
//	from file_data
//	JOIN file_metadata ON file_metadata.file_data_id = file_data.id
//	where file_metadata.file_path = $1
func (q *Queries) DownloadFileDataWithFilePath(ctx context.Context, db DBTX, filePath string) (*DownloadFileDataWithFilePathRow, error) {
	row := db.QueryRow(ctx, downloadFileDataWithFilePath, filePath)
	var i DownloadFileDataWithFilePathRow
	err := row.Scan(
		&i.ID,
		&i.FileHash,
		&i.ModTime,
		&i.FileSize,
		&i.FileDataB,
		&i.CreationTime,
		&i.FileMetadataID,
		&i.DirID,
		&i.FilePath,
		&i.FileType,
		&i.FileState,
		&i.FileDataID,
		&i.FileMode,
		&i.ModTime_2,
		&i.CreationTime_2,
	)
	return &i, err
}

const downloadFileStepOneGetLatestData = `-- name: DownloadFileStepOneGetLatestData :one
SELECT file_data.id, file_hash, file_data.mod_time, file_size, file_data_b, file_data.creation_time, 
file_metadata.id AS file_metadata_id,
dir_id, file_path, file_type, file_state, file_data_id, file_mode, file_metadata.mod_time, file_metadata.creation_time
from file_data 
JOIN file_metadata ON file_metadata.file_data_id = file_data.id
WHERE file_data.file_hash = $1
`

type DownloadFileStepOneGetLatestDataRow struct {
	ID             int32              `json:"id"`
	FileHash       *string            `json:"file_hash"`
	ModTime        pgtype.Timestamptz `json:"mod_time"`
	FileSize       int32              `json:"file_size"`
	FileDataB      []byte             `json:"file_data_b"`
	CreationTime   pgtype.Timestamptz `json:"creation_time"`
	FileMetadataID int32              `json:"file_metadata_id"`
	DirID          *int32             `json:"dir_id"`
	FilePath       string             `json:"file_path"`
	FileType       string             `json:"file_type"`
	FileState      string             `json:"file_state"`
	FileDataID     int32              `json:"file_data_id"`
	FileMode       int32              `json:"file_mode"`
	ModTime_2      pgtype.Timestamptz `json:"mod_time_2"`
	CreationTime_2 pgtype.Timestamptz `json:"creation_time_2"`
}

// DownloadFileStepOneGetLatestData
//
//	SELECT file_data.id, file_hash, file_data.mod_time, file_size, file_data_b, file_data.creation_time,
//	file_metadata.id AS file_metadata_id,
//	dir_id, file_path, file_type, file_state, file_data_id, file_mode, file_metadata.mod_time, file_metadata.creation_time
//	from file_data
//	JOIN file_metadata ON file_metadata.file_data_id = file_data.id
//	WHERE file_data.file_hash = $1
func (q *Queries) DownloadFileStepOneGetLatestData(ctx context.Context, db DBTX, fileHash *string) (*DownloadFileStepOneGetLatestDataRow, error) {
	row := db.QueryRow(ctx, downloadFileStepOneGetLatestData, fileHash)
	var i DownloadFileStepOneGetLatestDataRow
	err := row.Scan(
		&i.ID,
		&i.FileHash,
		&i.ModTime,
		&i.FileSize,
		&i.FileDataB,
		&i.CreationTime,
		&i.FileMetadataID,
		&i.DirID,
		&i.FilePath,
		&i.FileType,
		&i.FileState,
		&i.FileDataID,
		&i.FileMode,
		&i.ModTime_2,
		&i.CreationTime_2,
	)
	return &i, err
}

const updateFileTrackerMarkDownloaded = `-- name: UpdateFileTrackerMarkDownloaded :exec
INSERT INTO file_tracker (
    peer_id,
    file_meta_id,
    current_hash_id,
    file_state
) VALUES ($1, $2, $3, 'downloaded')
`

type UpdateFileTrackerMarkDownloadedParams struct {
	PeerID        uuid.UUID `json:"peer_id"`
	FileMetaID    int32     `json:"file_meta_id"`
	CurrentHashID int32     `json:"current_hash_id"`
}

// UpdateFileTrackerMarkDownloaded
//
//	INSERT INTO file_tracker (
//	    peer_id,
//	    file_meta_id,
//	    current_hash_id,
//	    file_state
//	) VALUES ($1, $2, $3, 'downloaded')
func (q *Queries) UpdateFileTrackerMarkDownloaded(ctx context.Context, db DBTX, arg *UpdateFileTrackerMarkDownloadedParams) error {
	_, err := db.Exec(ctx, updateFileTrackerMarkDownloaded, arg.PeerID, arg.FileMetaID, arg.CurrentHashID)
	return err
}