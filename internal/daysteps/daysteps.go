package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

var (
	ErrInvalidFormat   = errors.New("invalid day steps format: expected '<steps>,<duration>'")
	ErrInvalidSteps    = errors.New("steps must be positive")
	ErrInvalidDuration = errors.New("duration must be positive")
)

type DaySteps struct {
	Steps    int
	Duration time.Duration
	personaldata.Personal
}

func (ds *DaySteps) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 2 {
		return ErrInvalidFormat
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("parse steps: %w", err)
	}
	if steps <= 0 {
		return ErrInvalidSteps
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return fmt.Errorf("parse duration: %w", err)
	}
	if duration <= 0 {
		return ErrInvalidDuration
	}

	ds.Steps = steps
	ds.Duration = duration
	return nil
}

func (ds *DaySteps) ActionInfo() (string, error) {
	distance := spentenergy.Distance(ds.Steps, ds.Height)
	calories, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)
	if err != nil {
		return "", fmt.Errorf("calculate day steps calories: %w", err)
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps, distance, calories,
	), nil
}
