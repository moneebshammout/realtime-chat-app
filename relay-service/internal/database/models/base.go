package models

import (
	"reflect"
	"time"

	"relay-service/internal/database"
	"relay-service/pkg/utils"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3/qb"
	"github.com/scylladb/gocqlx/v3/table"
)

var logger = utils.GetLogger()

type model struct {
	dbClient  *database.DBClient
	model     *table.Table
	modelType reflect.Type
}

func (m *model) newInstance() interface{} {
	return reflect.New(m.modelType).Interface()
}

func (m *model) newSlice() interface{} {
	sliceType := reflect.SliceOf(m.modelType)
	return reflect.New(sliceType).Interface()
}

func (m *model) GetList(where qb.M) (any, error) {
	whereClause, binder := buildWhereClauses(where)
	query := m.dbClient.Session.
		Query(qb.Select(m.model.Name()).
			Columns("*").
			Where(whereClause...).
			ToCql()).BindMap(binder)

	result := m.newSlice()
	if err := query.SelectRelease(result); err != nil {
		logger.Errorf("Failed to get list: %v", err)
		return nil, err
	}

	if reflect.ValueOf(result).Elem().Len() == 0 {
		return []any{}, nil
	}

	return result, nil
}

func (m *model) Create(data map[string]interface{}) (map[string]interface{}, error) {
	data["id"] = gocql.TimeUUID()
	if data["created_at"] == nil {
		data["created_at"] = time.Now().UnixNano() / int64(time.Millisecond)
	}
	
	query := m.dbClient.Session.Query(m.model.Insert()).BindMap(data)
	if err := query.ExecRelease(); err != nil {
		logger.Errorf("Failed to create: %v", err)
		return nil, err
	}

	return data, nil
}

func (m *model) Delete(where qb.M) error {
	whereClause, binder := buildWhereClauses(where)
	query := m.dbClient.Session.
		Query(qb.Delete(m.model.Name()).
			Where(whereClause...).
			ToCql()).BindMap(binder)

	if err := query.ExecRelease(); err != nil {
		logger.Errorf("Failed to delete: %v", err)
		return err
	}

	return nil
}
