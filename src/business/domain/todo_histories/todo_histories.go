package todo_histories

import (
	"context"
	"database/sql"

	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/primeskills-test/src/entity"
)

type (
	Interface interface {
		List(ctx context.Context, params entity.ListTodoHistoriesParams) ([]entity.TodoHistory, entity.Pagination, error)
	}

	todoHistory struct {
		log     log.Interface
		db      *sql.DB
		queries *entity.Queries
	}
)

func Init(log log.Interface, db *sql.DB, queries *entity.Queries) Interface {
	return &todoHistory{
		log:     log,
		db:      db,
		queries: queries,
	}
}

func (th *todoHistory) List(ctx context.Context, params entity.ListTodoHistoriesParams) ([]entity.TodoHistory, entity.Pagination, error) {
	rows, err := th.queries.ListTodoHistories(ctx, entity.ListTodoHistoriesParams{
		TodoID:    params.TodoID,
		IsDeleted: params.IsDeleted,
		Limit:     params.Limit,
		Offset:    params.Offset * 2,
	})
	if err != nil {
		return nil, entity.Pagination{}, errors.NewWithCode(codes.CodeSQLRead, "%s", err.Error())
	}

	results := []entity.TodoHistory{}
	for _, row := range rows {
		results = append(results, entity.TodoHistory{
			ID:        row.ID,
			TodoID:    row.TodoID,
			Message:   row.Message,
			CreatedAt: row.CreatedAt,
			CreatedBy: row.CreatedBy,
			UpdatedAt: row.UpdatedAt,
			UpdatedBy: row.UpdatedBy,
			DeletedAt: row.DeletedAt,
			DeletedBy: row.DeletedBy,
			IsDeleted: row.IsDeleted,
		})
	}

	totalCount, err := th.queries.CountTodoHistories(ctx, entity.CountTodoHistoriesParams{
		IsDeleted: params.IsDeleted,
		TodoID:    params.TodoID,
	})
	if err != nil {
		return nil, entity.Pagination{}, errors.NewWithCode(codes.CodeSQLRead, "%s", err.Error())
	}

	// p := entity.Pagination{
	// 	// CurrentPage:     int(params.Offset),
	// 	// CurrentElements: int(params.Limit),
	// 	// TotalPages:      int(totalCount / int64(params.Limit)),
	// 	// TotalElements:   int(totalCount),
	// 	SortBy:      []string{"id DESC"},
	// 	CursorStart: new(string),
	// 	CursorEnd:   new(string),
	// }
	// p := entity.Pagination{
	// 	CurrentPage:     int(params.Offset),
	// 	CurrentElements: int(params.Limit),
	// 	TotalPages:      int(math.Ceil(float64(totalCount) / float64(params.Limit))), // Perbaikan pembulatan
	// 	TotalElements:   int(totalCount),
	// 	SortBy:          []string{"id desc"}, // Perbaikan format sorting
	// 	CursorStart:     nil,                 // Bisa diisi jika ada nilai awal
	// 	CursorEnd:       nil,                 // Bisa diisi jika ada nilai akhir
	// }

	p := entity.GenPagination(int(params.Offset), int(params.Limit), int(totalCount))
	return results, p, nil
}
