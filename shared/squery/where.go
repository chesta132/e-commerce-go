package squery

import (
	"fmt"
	"reflect"
	"strings"
)

type Where struct {
	Name  string // Field to filter | REQUIRED
	Value any    // Value of Name 	| REQUIRED
	Ind   string // Where indent 		| e.g. ("OR", "NOT") 											| Default: "AND"
	Op    string // Where operator  | e.g. ("=", ">", "<", "LIKE", "IN", etc) | Default: "="
}

func BuildWhere(where []Where) (q string, v []any) {
	query := ""
	value := []any{}
	for i, w := range where {
		indent := ""
		operator := "="

		if i != 0 {
			if w.Ind != "" {
				indent = fmt.Sprintf(" %v ", w.Ind)
			} else {
				indent = " AND "
			}
		}

		if w.Op != "" {
			operator = w.Op
		}

		if strings.ToUpper(w.Op) == "IN" {
			rv := reflect.ValueOf(w.Value)
			if rv.Kind() == reflect.Slice {
				placeholders := strings.Repeat("?,", rv.Len())
				placeholders = placeholders[:len(placeholders)-1]
				query += fmt.Sprintf("%v%v IN (%v)", indent, w.Name, placeholders)
				for j := 0; j < rv.Len(); j++ {
					value = append(value, rv.Index(j).Interface())
				}
			}
			continue
		}

		query = query + fmt.Sprintf("%v%v %v ?", indent, w.Name, operator)
		value = append(value, w.Value)
	}
	return query, value
}
