package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

const (
	TrainingTypeWalking = "Ходьба"
	TrainingTypeRunning = "Бег"
)

var (
	ErrInvalidFormat   = errors.New("invalid training format: expected '<steps>,<type>,<duration>'")
	ErrInvalidSteps    = errors.New("steps must be positive")
	ErrInvalidDuration = errors.New("duration must be positive")
	ErrUnknownType     = errors.New("неизвестный тип тренировки")
)

type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	personaldata.Personal
}

func (t *Training) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		return ErrInvalidFormat
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("parse steps: %w", err)
	}
	if steps <= 0 {
		return ErrInvalidSteps
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return fmt.Errorf("parse duration: %w", err)
	}
	if duration <= 0 {
		return ErrInvalidDuration
	}

	t.Steps = steps
	t.TrainingType = parts[1]
	t.Duration = duration
	return nil
}

func (t Training) ActionInfo() (string, error) {
	calorieCalc, ok := t.calorieCalculator()
	if !ok {
		return "", ErrUnknownType
	}

	calories, err := calorieCalc()
	if err != nil {
		return "", fmt.Errorf("calculate training calories: %w", err)
	}

	distance := spentenergy.Distance(t.Steps, t.Height)
	speed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, t.Duration.Hours(), distance, speed, calories,
	), nil
}

func (t Training) calorieCalculator() (func() (float64, error), bool) {
	switch t.TrainingType {
	case TrainingTypeWalking:
		return func() (float64, error) {
			return spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
		}, true
	case TrainingTypeRunning:
		return func() (float64, error) {
			return spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
		}, true
	}
	return nil, false
}
