package main

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestRandom_float(t *testing.T) {
	rate := 0.1
	delay := 1000.0

	for i := 0; i < 10000; i++ {
		if rand.Float64() <= rate/100 {
			duration := time.Duration(float64(time.Millisecond) * delay)
			log.Printf("[INFO] Will delay for %s.", duration.String())
		}
	}
}

func TestMain(m *testing.M) {

}
