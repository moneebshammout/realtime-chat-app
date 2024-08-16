package models

import (
	"reflect"
	"strings"

	"relay-service/internal/database"
	"relay-service/pkg/utils"

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

func (m *model) GetList(where qb.M) (any,error) {
	whereClause,binder := buildWhereClauses(where)
	query := m.dbClient.Session.
	Query(qb.Select(m.model.Name()).
	Columns("*").
	Where(whereClause...).
	ToCql()).BindMap(binder)
	
	result := m.newSlice()
	if err := query.SelectRelease(result); err != nil {
		logger.Errorf("Failed to get list: %v", err)
		return nil,err
	}

	return result,nil
}

func buildWhereClauses(where qb.M) ([]qb.Cmp, qb.M) {
	var clauses []qb.Cmp
	bindMap := qb.M{}

	for key, value := range where {
		parts := strings.SplitN(key, " ", 2)
		if len(parts) != 2 {
			logger.Errorf("Invalid where clause: %s\n", key)
			continue
		}

		column, op := parts[0], parts[1]
		switch op {
		case "=":
			clauses = append(clauses, qb.Eq(column))
		case "!=":
			clauses = append(clauses, qb.Ne(column))
		case ">":
			clauses = append(clauses, qb.Gt(column))
		case ">=":
			clauses = append(clauses, qb.GtOrEq(column))
		case "<":
			clauses = append(clauses, qb.Lt(column))
		case "<=":
			clauses = append(clauses, qb.LtOrEq(column))
		default:
			logger.Errorf("Unsupported operation: %s\n", op)
		}
		
		bindMap[column] = value
	}

	return clauses, bindMap
}
