package spentenergy

import (
	"errors"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе.
)

var (
	ErrInvalidSteps    = errors.New("steps must be positive")
	ErrInvalidWeight   = errors.New("weight must be positive")
	ErrInvalidHeight   = errors.New("height must be positive")
	ErrInvalidDuration = errors.New("duration must be positive")
)

func Distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0
	}
	stepLength := height * stepLengthCoefficient
	return float64(steps) * stepLength / mInKm
}

func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 || height <= 0 || duration <= 0 {
		return 0
	}
	hours := duration.Hours()
	if hours <= 0 {
		return 0
	}
	return Distance(steps, height) / hours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateInput(steps, weight, height, duration); err != nil {
		return 0, err
	}
	speed := MeanSpeed(steps, height, duration)
	return calories(weight, speed, duration), nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateInput(steps, weight, height, duration); err != nil {
		return 0, err
	}
	speed := MeanSpeed(steps, height, duration)
	return calories(weight, speed, duration) * walkingCaloriesCoefficient, nil
}

func validateInput(steps int, weight, height float64, duration time.Duration) error {
	switch {
	case steps <= 0:
		return ErrInvalidSteps
	case weight <= 0:
		return ErrInvalidWeight
	case height <= 0:
		return ErrInvalidHeight
	case duration <= 0:
		return ErrInvalidDuration
	}
	return nil
}

func calories(weight, speed float64, duration time.Duration) float64 {
	minutes := duration.Minutes()
	return weight * speed * minutes / minInH
}
