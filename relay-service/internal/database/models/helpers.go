package models

import (
	"strings"

	"github.com/scylladb/gocqlx/v3/qb"
)

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
		case "in":
			clauses = append(clauses, qb.In(column))
		default:
			logger.Errorf("Unsupported operation: %s\n", op)
		}
		logger.Infof(column)
		bindMap[column] = value
	}

	return clauses, bindMap
}
