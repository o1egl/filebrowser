package common

import (
	"context"
	"strconv"

	"github.com/doug-martin/goqu/v9"

	"github.com/filebrowser/filebrowser/domain"
	"github.com/filebrowser/filebrowser/storage/sql/conv"
	"github.com/filebrowser/filebrowser/storage/sql/model"
)

const (
	uploadsTable = "uploads"
)

func (s *Store) GetUpload(ctx context.Context, id int64) (*domain.Upload, error) {
	var uploadModel model.Upload
	found, err := s.db.Select().
		From(uploadsTable).
		Where(goqu.C("id").Eq(id)).
		Executor().
		ScanStructContext(ctx, &uploadModel)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, domain.NewNotFoundError(err, domain.ResourceUpload, "id", strconv.FormatInt(id, 10))
	}

	return conv.UploadToDomain(&uploadModel), nil
}

func (s *Store) CreateUpload(ctx context.Context, upload *domain.Upload) (id int64, err error) {
	res, err := s.db.Insert(uploadsTable).
		Rows(conv.UploadFromDomain(upload)).
		Executor().
		ExecContext(ctx)
	if err != nil {
		return 0, err
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	upload.ID = id

	return id, nil
}

func (s *Store) DeleteUpload(ctx context.Context, id int64) error {
	_, err := s.db.Delete(uploadsTable).Where(goqu.C("id").Eq(id)).Executor().ExecContext(ctx)
	return err
}
