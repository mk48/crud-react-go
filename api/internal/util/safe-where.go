package util

import (
	"fmt"
	"regexp"
	"strings"
)

// allowedWhereKeywords are the only bare (unquoted) words permitted in a
// client-supplied WHERE fragment - everything react-querybuilder's SQL/
// parameterized formatters can emit for a comparison. Column/table names are
// never bare; they always come through quoted (see quotedIdentifierPattern),
// so any other bare word is either a stray SQL keyword (SELECT, UNION, DROP,
// ...) or an attempt to break out of the intended expression.
var allowedWhereKeywords = map[string]bool{
	"AND": true, "OR": true, "NOT": true,
	"IS": true, "NULL": true,
	"LIKE": true, "ILIKE": true, "IN": true, "BETWEEN": true,
	"TRUE": true, "FALSE": true,
}

var (
	quotedIdentifierPattern = regexp.MustCompile(`"([A-Za-z_][A-Za-z0-9_]*)"`)
	bareWordPattern         = regexp.MustCompile(`[A-Za-z_]+`)
	placeholderPattern      = regexp.MustCompile(`:p[0-9]+`)
	// What's left over once quoted identifiers, placeholders and allowed
	// keywords are stripped out should only be punctuation that glues
	// comparisons together - anything else (a stray literal like `1` in
	// `OR 1=1`, an operator we don't recognize, ...) means the fragment
	// carries more structure than a plain filter expression should.
	remainderAllowedPattern = regexp.MustCompile(`^[\s()=<>!,]*$`)
)

// ValidateWhereCondition checks a client-supplied, parameterized SQL WHERE
// fragment (as produced by react-querybuilder's `formatQuery` with the
// "postgresql" preset - see AdvancedQueryBuilder.tsx) before it's
// interpolated into a query string. Values are never inlined in this
// fragment - they arrive separately as bound `:pN` parameters - so this only
// needs to guard structure: which columns and keywords appear in it.
//
// allowedColumns should be every column of the table the query runs
// against (e.g. from util.TableColumns), keyed by column name.
func ValidateWhereCondition(where string, allowedColumns map[string]bool) error {
	if strings.TrimSpace(where) == "" {
		return fmt.Errorf("whereCondition must not be empty")
	}

	if strings.ContainsAny(where, ";") || strings.Contains(where, "--") ||
		strings.Contains(where, "/*") || strings.Contains(where, "*/") {
		return fmt.Errorf("whereCondition contains disallowed characters")
	}

	remainder := where
	for _, m := range quotedIdentifierPattern.FindAllStringSubmatch(where, -1) {
		if !allowedColumns[m[1]] {
			return fmt.Errorf("whereCondition references an unknown column: %s", m[1])
		}
	}
	remainder = quotedIdentifierPattern.ReplaceAllString(remainder, "")
	remainder = placeholderPattern.ReplaceAllString(remainder, "")

	for _, w := range bareWordPattern.FindAllString(remainder, -1) {
		if !allowedWhereKeywords[strings.ToUpper(w)] {
			return fmt.Errorf("whereCondition contains an unexpected keyword: %s", w)
		}
	}
	remainder = bareWordPattern.ReplaceAllString(remainder, "")

	if !remainderAllowedPattern.MatchString(remainder) {
		return fmt.Errorf("whereCondition has unexpected structure")
	}

	return nil
}

// QualifyWhereColumns rewrites every quoted column reference in an
// already-validated WHERE fragment (see ValidateWhereCondition) to be
// qualified with the given table alias, e.g. "deleted_at" -> k."deleted_at".
// Queries that join in other tables sharing column names with the primary
// one - every audited table joins in "user" for created_by/updated_by/
// deleted_by's emails, and "user" has its own created_at/updated_at/
// deleted_at - would otherwise make a bare column reference ambiguous.
func QualifyWhereColumns(where string, alias string) string {
	return quotedIdentifierPattern.ReplaceAllString(where, alias+`."$1"`)
}
