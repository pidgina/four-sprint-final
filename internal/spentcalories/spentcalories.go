package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

var (
	SliceNotValid   = errors.New("переданные данные не соответствуют, необходима отправка в верном формате")
	ErrorParceConv  = errors.New("ошибка конвертации, проверьте корректность введенных данных, ожидается запись в привычном формате")
	ErrorParceTime  = errors.New("ошибка преобразования переданных данных, проверьте корректность введенных данных")
	StepZeroNegativ = errors.New("неверное число шагов")
	InvalidWeHe     = errors.New("некорректный вес или рост")
	InvalidTime     = errors.New("некорректное число времени")
	ErrorViem       = errors.New("неизвестный тип тренировки")
)

func parseTraining(data string) (int, string, time.Duration, error) {
	//"3456,Ходьба,3h00m"
	// 0.      1.    2
	slice := strings.Split(data, ",")

	if len(slice) != 3 {
		return 0, "", 0, SliceNotValid
	}

	stepInt, err := strconv.Atoi(slice[0]) // количество шагов
	if err != nil {
		return 0, "", 0, ErrorParceConv
	}

	if stepInt <= 0 {
		return 0, "", 0, StepZeroNegativ
	}

	d, err := time.ParseDuration(slice[2]) // продолжительность (в duration)
	if err != nil {
		return 0, "", 0, ErrorParceTime
	}

	if d <= 0 {
		return 0, "", 0, InvalidTime
	}

	return stepInt, slice[1], d, nil

}

func distance(steps int, height float64) float64 {

	lenStep := height * stepLengthCoefficient
	distance := float64(steps) * lenStep

	return distance / mInKm

}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration.Hours() <= 0 {
		return 0
	}

	dist := distance(steps, height) // дистанция

	d := duration.Hours() // переводим продолжительность наносекунд в часы

	return dist / d
	//return distance(steps, height) / duration.Hours()

}

func TrainingInfo(data string, weight, height float64) (string, error) {
	//"3456,Ходьба,3h00m"
	// 0.      1.    2
	step, view, times, err := parseTraining(data) // шаги, вид тренировки, длительность
	if err != nil {
		log.Println(err)
		return "", err
	}

	distanceTraining := distance(step, height)      // дистанция за тренировку
	speedTraining := meanSpeed(step, height, times) // средняя скорость тренировки

	ccalTrainingRun, err := RunningSpentCalories(step, weight, height, times) // сожжено при беге
	if err != nil {
		log.Println(err)
		return "", err
	}

	ccalTrainingWalk, err := WalkingSpentCalories(step, weight, height, times) // сожжено при ходьбе
	if err != nil {
		log.Println(err)
		return "", err
	}

	timeses := times.Hours()
	switch view {
	case "Ходьба":
		return fmt.Sprintf(
			"Тип тренировки: %s\n"+
				"Длительность: %.2f ч.\n"+
				"Дистанция: %.2f км.\n"+
				"Скорость: %.2f км/ч\n"+
				"Сожгли калорий: %.2f\n",
			view,
			timeses,
			distanceTraining,
			speedTraining,
			ccalTrainingWalk,
		), err
	case "Бег":
		return fmt.Sprintf(
			"Тип тренировки: %s\n"+
				"Длительность: %.2f ч.\n"+
				"Дистанция: %.2f км.\n"+
				"Скорость: %.2f км/ч\n"+
				"Сожгли калорий: %.2f\n",
			view,
			timeses,
			distanceTraining,
			speedTraining,
			ccalTrainingRun,
		), err
	default:
		return "", ErrorViem
	}

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if duration.Hours() <= 0 {
		return 0, InvalidTime
	}

	if steps <= 0 {
		return 0, StepZeroNegativ
	}

	if weight <= 0 || height <= 0 {
		return 0, InvalidWeHe
	}

	averageSpeed := meanSpeed(steps, height, duration)
	dMin := duration.Minutes()
	ccal := (weight * averageSpeed * dMin) / minInH
	return ccal, nil

	//return (weight * meanSpeed(steps, height, duration*time.Hour) * float64(duration*time.Minute) / minInH), nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if duration.Hours() <= 0 {
		return 0, InvalidTime
	}

	if steps <= 0 {
		return 0, StepZeroNegativ
	}

	if weight <= 0 || height <= 0 {
		return 0, InvalidWeHe
	}

	averageSpeed := meanSpeed(steps, height, duration)
	dMin := duration.Minutes()
	ccal := (weight * averageSpeed * dMin) / minInH

	return ccal * walkingCaloriesCoefficient, nil

	//return (weight * meanSpeed(steps, height, duration*time.Hour) * float64(duration*time.Minute) / minInH) / walkingCaloriesCoefficient, nil

}
