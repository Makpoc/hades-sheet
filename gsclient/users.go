package gsclient

import (
	"fmt"
	"log"

	"github.com/makpoc/hades-sheet/models"
	//	sheets "google.golang.org/api/sheets/v4"
)

const wsFleetSheet = "WS Fleet"

// GetUsers returns all usernames from the WS Fleet sheet
func (s *Sheet) GetUsers() (models.Users, error) {
	const userColumn = "B"

	users, err := s.service.Spreadsheets.Values.Get(s.id, fmt.Sprintf("%s!%s%d:%s%d", wsFleetSheet, userColumn, minRowN, userColumn, maxRowN)).Do()
	if err != nil {
		return nil, err
	}

	if len(users.Values) == 0 {
		return models.Users{}, nil
	}

	var result models.Users
	values := getDataSubset(users.Values)
	for _, u := range values {
		usr, ok := u[0].(models.User)
		if !ok {
			log.Printf("Value not of type models.User: %v\n", u[0])
			continue
		}
		result = append(result, usr)
	}
	return result, nil
}
