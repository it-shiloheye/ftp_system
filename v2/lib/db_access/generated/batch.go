// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: batch.go

package db_access

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrBatchAlreadyClosed = errors.New("batch already closed")
)

const addDirForPeer = `-- name: AddDirForPeer :batchone
INSERT INTO peer_dirs(
    peer_id,
    dir_path
) VALUES ($1, $2)
RETURNING id, peer_id, creation_time, dir_path
`

type AddDirForPeerBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

type AddDirForPeerParams struct {
	PeerID  pgtype.UUID `json:"peer_id"`
	DirPath string      `json:"dir_path"`
}

// AddDirForPeer
//
//	INSERT INTO peer_dirs(
//	    peer_id,
//	    dir_path
//	) VALUES ($1, $2)
//	RETURNING id, peer_id, creation_time, dir_path
func (q *Queries) AddDirForPeer(ctx context.Context, db DBTX, arg []*AddDirForPeerParams) *AddDirForPeerBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.PeerID,
			a.DirPath,
		}
		batch.Queue(addDirForPeer, vals...)
	}
	br := db.SendBatch(ctx, batch)
	return &AddDirForPeerBatchResults{br, len(arg), false}
}

func (b *AddDirForPeerBatchResults) QueryRow(f func(int, *PeerDir, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		var i PeerDir
		if b.closed {
			if f != nil {
				f(t, nil, ErrBatchAlreadyClosed)
			}
			continue
		}
		row := b.br.QueryRow()
		err := row.Scan(
			&i.ID,
			&i.PeerID,
			&i.CreationTime,
			&i.DirPath,
		)
		if f != nil {
			f(t, &i, err)
		}
	}
}

func (b *AddDirForPeerBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}

const downloadFileData = `-- name: DownloadFileData :batchmany
SELECT 
    file_data.file_data_id, file_hash, prev_file_hash, file_status, modification_date, file_data.creation_time, file_data, id, peer_id, dir_id, file_metadata.file_data_id, file_name, file_path, file_type, creation_day, file_metadata.creation_time
FROM 
    file_data 
    INNER JOIN file_metadata USING (file_data_id)
WHERE 
    file_metadata.peer_id = $1
    AND file_data.file_data IS NOT NULL
    AND file_data.modification_date >= $2
ORDER BY
    modification_date DESC
`

type DownloadFileDataBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

type DownloadFileDataParams struct {
	PeerID           pgtype.UUID      `json:"peer_id"`
	ModificationDate pgtype.Timestamp `json:"modification_date"`
}

type DownloadFileDataRow struct {
	FileDataID       int32              `json:"file_data_id"`
	FileHash         *string            `json:"file_hash"`
	PrevFileHash     *int32             `json:"prev_file_hash"`
	FileStatus       NullFileStatusType `json:"file_status"`
	ModificationDate pgtype.Timestamp   `json:"modification_date"`
	CreationTime     pgtype.Timestamptz `json:"creation_time"`
	FileData         []byte             `json:"file_data"`
	ID               int32              `json:"id"`
	PeerID           pgtype.UUID        `json:"peer_id"`
	DirID            *int32             `json:"dir_id"`
	FileDataID_2     *int32             `json:"file_data_id_2"`
	FileName         string             `json:"file_name"`
	FilePath         string             `json:"file_path"`
	FileType         string             `json:"file_type"`
	CreationDay      pgtype.Date        `json:"creation_day"`
	CreationTime_2   pgtype.Timestamptz `json:"creation_time_2"`
}

// DownloadFileData
//
//	SELECT
//	    file_data.file_data_id, file_hash, prev_file_hash, file_status, modification_date, file_data.creation_time, file_data, id, peer_id, dir_id, file_metadata.file_data_id, file_name, file_path, file_type, creation_day, file_metadata.creation_time
//	FROM
//	    file_data
//	    INNER JOIN file_metadata USING (file_data_id)
//	WHERE
//	    file_metadata.peer_id = $1
//	    AND file_data.file_data IS NOT NULL
//	    AND file_data.modification_date >= $2
//	ORDER BY
//	    modification_date DESC
func (q *Queries) DownloadFileData(ctx context.Context, db DBTX, arg []*DownloadFileDataParams) *DownloadFileDataBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.PeerID,
			a.ModificationDate,
		}
		batch.Queue(downloadFileData, vals...)
	}
	br := db.SendBatch(ctx, batch)
	return &DownloadFileDataBatchResults{br, len(arg), false}
}

func (b *DownloadFileDataBatchResults) Query(f func(int, []*DownloadFileDataRow, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		items := []*DownloadFileDataRow{}
		if b.closed {
			if f != nil {
				f(t, items, ErrBatchAlreadyClosed)
			}
			continue
		}
		err := func() error {
			rows, err := b.br.Query()
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var i DownloadFileDataRow
				if err := rows.Scan(
					&i.FileDataID,
					&i.FileHash,
					&i.PrevFileHash,
					&i.FileStatus,
					&i.ModificationDate,
					&i.CreationTime,
					&i.FileData,
					&i.ID,
					&i.PeerID,
					&i.DirID,
					&i.FileDataID_2,
					&i.FileName,
					&i.FilePath,
					&i.FileType,
					&i.CreationDay,
					&i.CreationTime_2,
				); err != nil {
					return err
				}
				items = append(items, &i)
			}
			return rows.Err()
		}()
		if f != nil {
			f(t, items, err)
		}
	}
}

func (b *DownloadFileDataBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}

const uploadFile = `-- name: UploadFile :batchone
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
    file_data_id, file_hash
`

type UploadFileBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

type UploadFileParams struct {
	FileStatus       NullFileStatusType `json:"file_status"`
	FileData         []byte             `json:"file_data"`
	ModificationDate pgtype.Timestamp   `json:"modification_date"`
}

type UploadFileRow struct {
	FileDataID int32   `json:"file_data_id"`
	FileHash   *string `json:"file_hash"`
}

// UploadFile
//
//	INSERT INTO file_data(
//	    file_status,
//	    file_data,
//	    modification_date
//	) VALUES ($1, $2,$3)
//	ON CONFLICT (file_hash)
//	DO UPDATE
//	    SET
//	        file_data = $2
//	RETURNING
//	    file_data_id, file_hash
func (q *Queries) UploadFile(ctx context.Context, db DBTX, arg []*UploadFileParams) *UploadFileBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.FileStatus,
			a.FileData,
			a.ModificationDate,
		}
		batch.Queue(uploadFile, vals...)
	}
	br := db.SendBatch(ctx, batch)
	return &UploadFileBatchResults{br, len(arg), false}
}

func (b *UploadFileBatchResults) QueryRow(f func(int, *UploadFileRow, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		var i UploadFileRow
		if b.closed {
			if f != nil {
				f(t, nil, ErrBatchAlreadyClosed)
			}
			continue
		}
		row := b.br.QueryRow()
		err := row.Scan(&i.FileDataID, &i.FileHash)
		if f != nil {
			f(t, &i, err)
		}
	}
}

func (b *UploadFileBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}

const uploadMetadata = `-- name: UploadMetadata :batchone
INSERT INTO file_metadata(
    peer_id,
    dir_id,
    file_data_id,
    file_name,
    file_path,
    file_type
) VALUES($1, $2, $3, $4, $5, $6)
RETURNING id
`

type UploadMetadataBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

type UploadMetadataParams struct {
	PeerID     pgtype.UUID `json:"peer_id"`
	DirID      *int32      `json:"dir_id"`
	FileDataID *int32      `json:"file_data_id"`
	FileName   string      `json:"file_name"`
	FilePath   string      `json:"file_path"`
	FileType   string      `json:"file_type"`
}

// UploadMetadata
//
//	INSERT INTO file_metadata(
//	    peer_id,
//	    dir_id,
//	    file_data_id,
//	    file_name,
//	    file_path,
//	    file_type
//	) VALUES($1, $2, $3, $4, $5, $6)
//	RETURNING id
func (q *Queries) UploadMetadata(ctx context.Context, db DBTX, arg []*UploadMetadataParams) *UploadMetadataBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.PeerID,
			a.DirID,
			a.FileDataID,
			a.FileName,
			a.FilePath,
			a.FileType,
		}
		batch.Queue(uploadMetadata, vals...)
	}
	br := db.SendBatch(ctx, batch)
	return &UploadMetadataBatchResults{br, len(arg), false}
}

func (b *UploadMetadataBatchResults) QueryRow(f func(int, int32, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		var id int32
		if b.closed {
			if f != nil {
				f(t, id, ErrBatchAlreadyClosed)
			}
			continue
		}
		row := b.br.QueryRow()
		err := row.Scan(&id)
		if f != nil {
			f(t, id, err)
		}
	}
}

func (b *UploadMetadataBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}
