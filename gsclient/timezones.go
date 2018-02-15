package gsclient

import (
	"fmt"
)

// UserTime contains information about users and their time
type UserTime struct {
	UserName    string
	CurrentTime string
	Offset      string
}

// GetTimeZones returns the list with users and their corresponding offset and currentTime
func (s *Sheet) GetTimeZones() ([]UserTime, error) {
	const userColumn = "A"
	const offsetColumn = "C"
	users, err := s.service.Spreadsheets.Values.Get(s.id, fmt.Sprintf("%s!%s%d:%s%d", tzSheet, userColumn, minRowN, offsetColumn, maxRowN)).Do()
	if err != nil {
		return nil, err
	}

	if len(users.Values) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	values := getDataSubset(users.Values)
	var result []UserTime
	for _, v := range values {
		if len(v) == 0 {
			// empty row, skip
			continue
		}
		entry := UserTime{}
		switch len(v) {
		case 3:
			entry.Offset = fmt.Sprintf("%s", v[2])
			fallthrough
		case 2:
			entry.CurrentTime = fmt.Sprintf("%s", v[1])
			fallthrough
		case 1:
			entry.UserName = fmt.Sprintf("%s", v[0])
		}
		result = append(result, entry)
	}

	return result, nil
}
