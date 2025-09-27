package bot

import (
	"strconv"
	"strings"
	"time"

	"github.com/Alcamech/FitBoisBot/internal/constants"
	"github.com/Alcamech/FitBoisBot/internal/errors"
)

// validateDatePart validates that a date part is exactly 2 digits and within the given range.
func validateDatePart(part string, min, max int, name string) error {
	if len(part) != constants.DatePartLength {
		return errors.NewValidationError(name, part, "must be 2 digits")
	}

	value, err := strconv.Atoi(part)
	if err != nil {
		return errors.NewValidationError(name, part, "must be numeric")
	}

	if value < min || value > max {
		return errors.NewValidationError(name, part, "out of valid range")
	}

	return nil
}

// validateTimezone checks if a timezone string is valid.
func validateTimezone(timezone string) error {
	if timezone == "" {
		return errors.NewValidationError("timezone", timezone, "cannot be empty")
	}

	_, err := time.LoadLocation(timezone)
	if err != nil {
		return errors.NewValidationError("timezone", timezone, "invalid IANA timezone")
	}

	return nil
}

// isValidActivityName checks if an activity name is reasonable.
func isValidActivityName(activity string) bool {
	if len(activity) == 0 {
		return false
	}
	// Activity name should be reasonable length and contain printable characters
	return len(activity) <= constants.ActivityNameMaxLength && len(strings.TrimSpace(activity)) > 0
}