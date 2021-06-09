package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
)

func distance(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(float64((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)))
}

func calculateTempForOnePoint(sensorValues map[int]map[int]float64, x, y int) string {
	if sensorValues[x][y] > 0 {
		return fmt.Sprintf("%.2f", sensorValues[x][y])
	}

	var d, sum, result float64
	for xi, row := range sensorValues {
		for yi, v := range row {
			if v <= 0 {
				continue
			}
			d = distance(xi, yi, x, y)
			sum += 1 / d
		}
	}

	for xi, row := range sensorValues {
		for yi, v := range row {
			if v <= 0 {
				continue
			}
			d = distance(xi, yi, x, y)
			result += v / (sum * d)
		}
	}

	return fmt.Sprintf("%.2f", result)
}

func interpolate(sensorValues map[int]map[int]float64, x_max, y_max int) [][]string {
	heatmap := make([][]string, x_max)
	for x := 0; x < x_max; x++ {
		heatmap[x] = make([]string, y_max)
		for y := 0; y < y_max; y++ {
			heatmap[x][y] = calculateTempForOnePoint(sensorValues, x, y)
		}
	}

	return heatmap
}

type sensorValue struct {
	x, y  int
	value float64
}

func main() {
	inputs := []sensorValue{
		{
			x:     5,
			y:     0,
			value: 85,
		},
		{
			x:     2,
			y:     1,
			value: 70,
		},
		{
			x:     5,
			y:     4,
			value: 80,
		},
		{
			x:     0,
			y:     4,
			value: 60,
		},
		{
			x:     0,
			y:     8,
			value: 75,
		},
	}
	sensorValues := make(map[int]map[int]float64)
	for _, v := range inputs {
		if _, ok := sensorValues[v.x]; !ok {
			sensorValues[v.x] = make(map[int]float64)
		}
		sensorValues[v.x][v.y] = v.value
	}

	results := interpolate(sensorValues, 10, 10)
	file, err := os.Create("result.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, row := range results {
		err := writer.Write(row)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
}
