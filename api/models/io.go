package models

import (
	"database/sql"
	"leichil-backend/api/utils"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type CountResult struct {
	Count *int `json:"count"`
}

func GetListData[T any](
	filter []any, fields []any, model *T,
	count_query string, list_query string, offset int, limit int,
) (CountResult, []any, error) {
	var err error
	var count_value = 0
	var count = CountResult{
		Count: &count_value,
	}
	var item_count = []any{count}
	if len(count_query) > 0 {
		count_fields := []any{&count.Count}
		item_count, err = utils.RunQuery(count_query, count_fields, &count, filter)
		if err != nil {
			return item_count[0].(CountResult), nil, err
		}
	}
	query := list_query + " LIMIT ? OFFSET ?"
	filter = append(filter, limit, offset)
	result, err := utils.RunQuery(query, fields, model, filter)
	return item_count[0].(CountResult), result, err
}

func GetCount[T any](filter []any, model *T, query string) (CountResult, error) {
	var count_value = 0
	var count = CountResult{
		Count: &count_value,
	}
	var fields = []any{&count.Count}
	item_count, err := utils.RunQuery(query, fields, &count, filter)
	return item_count[0].(CountResult), err
}

func GetSingleData[T any](
	filter []any, fields []any, model *T, query string,
) ([]any, error) {
	return utils.RunQuery(query, fields, model, filter)
}

func SanitizeData[T int | float64 | float32 | string | time.Time | bool](data *T) T {
	var result T
	if data == nil {
		return result
	}
	return *data
}

func PrepareTransaction(insert_stmt **sql.Stmt, tx *sql.Tx, query, error_message string) error {
	stmt, err := tx.Prepare(query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(error_message, err.Error())
	}
	*insert_stmt = stmt
	return err
}

func RunTransaction(result *sql.Result, stmt *sql.Stmt, tx *sql.Tx, fields []any, message string) error {
	stmt_result, err := stmt.Exec(fields...)
	if err != nil {
		_ = tx.Rollback()
		log.Error(message, err.Error())
	}
	*result = stmt_result
	return err
}
