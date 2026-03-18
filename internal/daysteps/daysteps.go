package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	slice := strings.Split(data, ",")
	if len(slice) != 2 {
		return 0, 0, spentcalories.SliceNotValid
	}

	stepInt, err := strconv.Atoi(slice[0]) // количество шагов
	if err != nil {
		return 0, 0, spentcalories.ErrorParceConv
	}

	if stepInt <= 0 {
		return 0, 0, spentcalories.StepZeroNegativ
	}

	d, err := time.ParseDuration(slice[1]) // продолжительность в duration
	if err != nil {
		return 0, 0, spentcalories.ErrorParceTime
	}

	if d <= 0 {
		return 0, 0, spentcalories.InvalidTime
	}

	return stepInt, d, nil
}

func DayActionInfo(data string, weight, height float64) string {
	step, times, err := parsePackage(data) // step = количество шагов, times = продолжительность
	if err != nil {
		log.Println(err)
		return ""
	}

	distantionMetr := float64(step) * stepLength // дистанция в метрах
	distantionKm := distantionMetr / mInKm       // дистанция в км

	ccalTraining, err := spentcalories.WalkingSpentCalories(step, weight, height, times) // считаем ккал
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.\n",
		step,
		distantionKm,
		ccalTraining,
	)
}
