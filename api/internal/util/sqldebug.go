package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var sqlPlaceholderRe = regexp.MustCompile(`\$(\d+)`)

// debugInterpolateSQL substitutes $1, $2, ... placeholders in query with their
// bound values so the result can be pasted straight into psql. It's a
// best-effort debug helper only - not safe for building queries to execute
// against the database (no protection against injection).
func debugInterpolateSQL(query string, args []any) string {
	return sqlPlaceholderRe.ReplaceAllStringFunc(query, func(m string) string {
		idx, err := strconv.Atoi(m[1:])
		if err != nil || idx < 1 || idx > len(args) {
			return m
		}
		return sqlLiteral(args[idx-1])
	})
}

func sqlLiteral(v any) string {
	switch val := v.(type) {
	case nil:
		return "NULL"
	case string:
		return "'" + strings.ReplaceAll(val, "'", "''") + "'"
	case []byte:
		return "'" + strings.ReplaceAll(string(val), "'", "''") + "'"
	case time.Time:
		return "'" + val.Format(time.RFC3339Nano) + "'"
	case bool:
		if val {
			return "TRUE"
		}
		return "FALSE"
	default:
		return fmt.Sprintf("%v", val)
	}
}
