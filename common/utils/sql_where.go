package utils

import (
	"errors"
	"fmt"
	"strings"
)

func GetWhereCondition(where string, args []interface{}) (string, error) {

	needArgsNum := strings.Count(where, "?")
	if needArgsNum > len(args) {
		return "", errors.New(fmt.Sprintf("need %d args, receive %d", needArgsNum, len(args)))
	} else if needArgsNum == 0 {
		return "", nil
	}

	whereFormatStr := strings.Replace(where, "?", "%v", -1)

	whereStr := fmt.Sprintf(whereFormatStr, args...)

	return "where " + whereStr, nil
}
