package validate

import (
	"fmt"
	"net/mail"
	"strings"
	"time"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
)

func Email(email string) error {
	if email == "" {
		return cerr.NewBadRequestError("email is required")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return cerr.NewBadRequestError("invalid email format")
	}
	return nil
}

func PasswordMinLength(password string) error {
	if len(password) < 8 {
		return cerr.NewBadRequestError("password must be at least 8 characters")
	}
	return nil
}

func NonEmpty(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return cerr.NewBadRequestError(fmt.Sprintf("%s is required", field))
	}
	return nil
}

func MaxLength(value string, max int, field string) error {
	if len(value) > max {
		return cerr.NewBadRequestError(fmt.Sprintf("%s must be at most %d characters", field, max))
	}
	return nil
}

func Positive(value float64, field string) error {
	if value <= 0 {
		return cerr.NewBadRequestError(fmt.Sprintf("%s must be greater than 0", field))
	}
	return nil
}

func NonNegative(value float64, field string) error {
	if value < 0 {
		return cerr.NewBadRequestError(fmt.Sprintf("%s must be non-negative", field))
	}
	return nil
}

func InRange(value, min, max float64, field string) error {
	if value < min || value > max {
		return cerr.NewBadRequestError(fmt.Sprintf("%s must be between %.0f and %.0f", field, min, max))
	}
	return nil
}

func NotFuture(t time.Time, field string) error {
	if t.After(time.Now()) {
		return cerr.NewBadRequestError(fmt.Sprintf("%s cannot be in the future", field))
	}
	return nil
}

func PositivePtr(value *float64, field string) error {
	if value != nil {
		return Positive(*value, field)
	}
	return nil
}

func NonNegativePtr(value *float64, field string) error {
	if value != nil {
		return NonNegative(*value, field)
	}
	return nil
}

func InRangePtr(value *float64, min, max float64, field string) error {
	if value != nil {
		return InRange(*value, min, max, field)
	}
	return nil
}

func NonNegativeIntPtr(value *int, field string) error {
	if value != nil && *value < 0 {
		return cerr.NewBadRequestError(fmt.Sprintf("%s must be non-negative", field))
	}
	return nil
}

func OneOf(value string, allowed []string, field string) error {
	for _, a := range allowed {
		if value == a {
			return nil
		}
	}
	return cerr.NewBadRequestError(fmt.Sprintf("%s must be one of: %s", field, strings.Join(allowed, ", ")))
}

func OneOfPtr(value *string, allowed []string, field string) error {
	if value != nil {
		return OneOf(*value, allowed, field)
	}
	return nil
}
