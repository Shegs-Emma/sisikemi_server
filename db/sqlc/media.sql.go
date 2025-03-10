// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: media.sql

package db

import (
	"context"
)

const createMedia = `-- name: CreateMedia :one
INSERT INTO media (
    media_ref, url, aws_id
) VALUES (
    $1, $2, $3
)
RETURNING id, media_ref, url, aws_id, created_at
`

type CreateMediaParams struct {
	MediaRef string `json:"media_ref"`
	Url      string `json:"url"`
	AwsID    string `json:"aws_id"`
}

func (q *Queries) CreateMedia(ctx context.Context, arg CreateMediaParams) (Medium, error) {
	row := q.db.QueryRow(ctx, createMedia, arg.MediaRef, arg.Url, arg.AwsID)
	var i Medium
	err := row.Scan(
		&i.ID,
		&i.MediaRef,
		&i.Url,
		&i.AwsID,
		&i.CreatedAt,
	)
	return i, err
}

const getMedia = `-- name: GetMedia :one
SELECT id, media_ref, url, aws_id, created_at FROM media
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetMedia(ctx context.Context, id int64) (Medium, error) {
	row := q.db.QueryRow(ctx, getMedia, id)
	var i Medium
	err := row.Scan(
		&i.ID,
		&i.MediaRef,
		&i.Url,
		&i.AwsID,
		&i.CreatedAt,
	)
	return i, err
}

const getMediaByRef = `-- name: GetMediaByRef :one
SELECT id, media_ref, url, aws_id, created_at FROM media
WHERE media_ref = $1 LIMIT 1
`

func (q *Queries) GetMediaByRef(ctx context.Context, mediaRef string) (Medium, error) {
	row := q.db.QueryRow(ctx, getMediaByRef, mediaRef)
	var i Medium
	err := row.Scan(
		&i.ID,
		&i.MediaRef,
		&i.Url,
		&i.AwsID,
		&i.CreatedAt,
	)
	return i, err
}

const listMedia = `-- name: ListMedia :many
SELECT id, media_ref, url, aws_id, created_at FROM media
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListMediaParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListMedia(ctx context.Context, arg ListMediaParams) ([]Medium, error) {
	rows, err := q.db.Query(ctx, listMedia, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Medium{}
	for rows.Next() {
		var i Medium
		if err := rows.Scan(
			&i.ID,
			&i.MediaRef,
			&i.Url,
			&i.AwsID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
