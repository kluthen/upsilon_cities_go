package tools

import (
	"math/rand"
	"time"
)

//IntRange range of int ;)
type IntRange struct {
	Min int
	Max int
}

//Min returns min of two int
func Min(lhs, rhs int) int {
	if lhs < rhs {
		return lhs
	}
	return rhs
}

//Max returns max of two int
func Max(lhs, rhs int) int {
	if lhs > rhs {
		return lhs
	}
	return rhs
}

//In value in range
func In(value, left, right int) bool {
	return value > left && value < right
}

//InRange value in range
func InRange(value int, rg IntRange) bool {
	return value > rg.Min && value < rg.Max
}

//InEq value in range inclusive
func InEq(value, left, right int) bool {
	return value >= left && value <= right
}

//InEqRange value in range
func InEqRange(value int, rg IntRange) bool {
	return value >= rg.Min && value <= rg.Max
}

//Abs absolute value in int.
func Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

//Roll Random in intrange.
func (ir IntRange) Roll() int {
	return rand.Intn(ir.Max-ir.Min) + ir.Min
}

//CycleLength duration of a cycle
var CycleLength time.Duration

//InitCycle initialize cycle lenght.
func InitCycle() {
	CycleLength, _ = time.ParseDuration("10s")
}

//RoundTime rounds up time up to cycle.
func RoundTime(base time.Time) time.Time {
	return base.Round(CycleLength)
}

//AddCycles tell what time it will be in cycles cycles.
func AddCycles(base time.Time, cycles int) time.Time {
	return base.Round(CycleLength).Add(time.Duration(cycles) * CycleLength)
}

//MinTime return lesser time of the two
func MinTime(lhs time.Time, rhs time.Time) time.Time {
	if lhs.Before(rhs) {
		return lhs
	}
	return rhs
}

//MaxTime return lesser time of the two
func MaxTime(lhs time.Time, rhs time.Time) time.Time {
	if lhs.After(rhs) {
		return lhs
	}
	return rhs
}
