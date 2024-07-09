package base

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetUUID(str string) (uuid_ pgtype.UUID, err error) {
	err = uuid_.Scan(str)
	if err != nil {
		err = fmt.Errorf("invalid uuid: %s", PrettyPrintValue(err))

	}
	return
}
