package db_helpers

import "strings"

func CheckNoRowsInResultSet(err error) bool {
	return strings.Contains(err.Error(), "no rows in result set")
}
