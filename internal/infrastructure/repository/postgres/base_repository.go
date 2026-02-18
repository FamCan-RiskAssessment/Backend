package postgres

import (
	"fmt"
	"strings"

	"github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func applyQueryOptions(query *gorm.DB, options *postgres.QueryOptions) *gorm.DB {
	if options == nil {
		return query
	}

	if options.Search != nil && options.Search.Query != "" && len(options.Search.Columns) > 0 {
		var conditions []string
		var args []interface{}
		searchPattern := "%" + options.Search.Query + "%"
		for _, col := range options.Search.Columns {
			conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", col))
			args = append(args, searchPattern)
		}
		query = query.Where(strings.Join(conditions, " OR "), args...)
	}

	if options.Sorting != nil {
		if options.Sorting.Asc {
			query = query.Order(clause.OrderByColumn{
				Column: clause.Column{Name: options.Sorting.Column},
			})
		} else {
			query = query.Order(clause.OrderByColumn{
				Column: clause.Column{Name: options.Sorting.Column},
				Desc:   true,
			})
		}
	}

	if options.Pagination != nil {
		query = query.Offset(options.Pagination.Offset).Limit(options.Pagination.Limit)
	}

	return query
}

func applySearchOnly(query *gorm.DB, options *postgres.QueryOptions) *gorm.DB {
	if options == nil {
		return query
	}

	if options.Search != nil && options.Search.Query != "" && len(options.Search.Columns) > 0 {
		var conditions []string
		var args []interface{}
		searchPattern := "%" + options.Search.Query + "%"
		for _, col := range options.Search.Columns {
			conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", col))
			args = append(args, searchPattern)
		}
		query = query.Where(strings.Join(conditions, " OR "), args...)
	}

	return query
}
