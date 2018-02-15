package gsclient

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/makpoc/hades-sheet/models"
)

const tzSheet = "Timezones"

// GetTimeZones returns the list with users and their corresponding offset and currentTime
func (s *Sheet) GetTimeZones() ([]models.UserTime, error) {
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
	var result []models.UserTime
	for _, v := range values {
		if len(v) == 0 {
			// empty row, skip
			continue
		}

		result = append(result, buildUserTime(v))
	}

	return result, nil
}

// buildUserTime builds UserTime from sheet cell values
func buildUserTime(v []interface{}) models.UserTime {
	entry := models.UserTime{}
	if len(v) >= 1 {
		entry.UserName = fmt.Sprintf("%s", v[0])
	}
	if len(v) == 3 {
		entry.Offset = fmt.Sprintf("%s", v[2])
	}

	currTime, err := getCurrentTime(entry.Offset)
	if err != nil {
		log.Printf("Failed to calculate current time for user %s: %v", entry.UserName, err)
	}
	entry.CurrentTime = currTime

	return entry
}

// getCurrentTime calculates the time based on given offset
func getCurrentTime(offset string) (time.Time, error) {
	// test if it's a number at all
	_, err := strconv.ParseFloat(offset, 64)
	if err != nil {
		return time.Time{}, err
	}

	var offsetH, offsetM = "0", "0"
	if strings.Contains(offset, ".") {
		// float offset in form "h.m"
		parts := strings.Split(offset, ".")
		if len(parts) != 2 {
			return time.Time{}, fmt.Errorf("invalid offset format")
		}
		offsetH, offsetM = parts[0], parts[1]
	} else {
		offsetH = offset
	}

	durationH, err := strconv.Atoi(offsetH)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse hours offset: %v", err)
	}
	durationM, err := strconv.Atoi(offsetM)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse minutes offset: %v", err)
	}

	now := time.Now().UTC()
	userTime := now.Add(time.Hour*time.Duration(durationH) + time.Minute*time.Duration(durationM))

	return userTime, nil
}
