package gsclient

import (
	"fmt"
	//	sheets "google.golang.org/api/sheets/v4"
)

const wsFleetSheet = "WS Fleet"
const tzSheet = "Timezones"
const minRowN = 1
const maxRowN = 999

// GetUsers returns all usernames from the WS Fleet sheet
func (s *Sheet) GetUsers() ([]string, error) {
	const userColumn = "B"

	users, err := s.service.Spreadsheets.Values.Get(s.id, fmt.Sprintf("%s!%s%d:%s%d", wsFleetSheet, userColumn, minRowN, userColumn, maxRowN)).Do()
	if err != nil {
		return nil, err
	}

	if len(users.Values) == 0 {
		return []string{}, nil
	}

	var result []string
	values := getDataSubset(users.Values)
	for _, u := range values {
		result = append(result, fmt.Sprintf("%s", u[0]))
	}
	return result, nil
}
