package app

import "time"

// CalculateMissingDays implements the catch-up detection.
func CalculateMissingDays(existingDates map[string]bool, timezone string) ([]string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}

	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	var missing []string

	// Check last 14 days + today
	for i := 14; i >= 0; i-- {
		d := today.AddDate(0, 0, -i)
		dStr := d.Format("2006-01-02")

		if !existingDates[dStr] {
			missing = append(missing, dStr)
		}
	}

	return missing, nil
}
