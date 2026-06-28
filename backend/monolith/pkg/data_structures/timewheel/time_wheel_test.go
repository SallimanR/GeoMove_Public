package timewheel

import (
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

type TestData struct {
	ID uint32
}

var (
	output = zerolog.ConsoleWriter{Out: os.Stdout}
	logger = zerolog.New(output).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Logger()
)

func TestAddSingle(t *testing.T) {}

func TestAddBatchWithConstDiff(t *testing.T) {
	// testCases := []struct {
	// 	name       string
	// 	wheelTick       time.Duration
	// 	wheelSize  int
	// 	taskNumber int
	// 	wantError  bool
	// }{
	// 	{wheelTick: 10 * time.Millisecond, wheelSize: 100, taskNumber: 10, wantError: false},
	// }
	//
	// for i := 0; i < len(testCases); i++ {
	// 	tc := testCases[i]
	// 	t.Run(tc.name, func(t *testing.T) {
	// 	})
	// }
	//
	wheelTick := 10 * time.Millisecond
	wheelSize := 1000
	wheel := NewTimeWheel[TestData](wheelTick, wheelSize)

	base := int64(0)
	tasksNumber := 100
	padding := 200
	tick := int64(wheelTick * time.Duration(padding))
	for i := int64(0); i < int64(tasksNumber); i++ {
		task := TestData{
			ID: uint32(i),
		}
		expires := base + i*tick
		wheelSlot := wheel.Add(task, expires)
		logger.Debug().Int("added to slot", wheelSlot).Send()
	}

	wheelTasks := wheel.GetTasksCount()
	logger.Debug().Any("wheelTasks", wheelTasks).Send()
	if tasksNumber != wheelTasks {
		t.Fatalf("Tasks count is not equal to task number")
	}
	wheelNonEmptySlots := wheel.GetNonEmptySlotsCount()
	logger.Debug().Any("wheelSlots", wheelNonEmptySlots).Send()
	expectedNumber := wheelSize / padding
	log.Println("expected: ", expectedNumber)
	if wheelNonEmptySlots != expectedNumber {
		t.Fatalf("Filled slot count is not equal to expected")
	}

	for tick := 0; tick <= wheelSize; tick += 1 {
		due := wheel.Advance()
		if due != nil {
			logger.Debug().Int("tick", tick).Any("due", due).Send()
		}
	}
}

func TestUpdate(t *testing.T) {}

func TestRemove(t *testing.T) {}

func TestAdvanceSingle(t *testing.T) {
	wheelTick := 10 * time.Millisecond
	wheelSize := 100
	wheel := NewTimeWheel[TestData](wheelTick, wheelSize)

	base := int64(0)
	tasksNumber := wheelSize
	padding := int64(20)
	tick := int64(wheelTick * time.Duration(padding))
	for i := 0; i < tasksNumber; i++ {
		task := TestData{
			ID: uint32(i),
		}
		expires := base + int64(i)*tick
		wheelSlot := wheel.Add(task, expires)
		logger.Debug().Int("added to slot", wheelSlot).Send()
	}

	taskCount := wheel.GetTasksCount()
	logger.Debug().Int("task count", taskCount).Send()
	wheelSlots := wheel.GetSlots()
	logger.Debug().Any("wheelSlots before test", wheelSlots).Send()

	t.Run("single advance", func(t *testing.T) {
		wheelSingle := wheel
		firstSlot := wheelSingle.slots[0]
		wheelSingle.Advance()
		if reflect.DeepEqual(firstSlot, wheelSingle.slots[0]) {
			t.Fatalf("First slot hasn't been removed")
		}
		wheelSlots = wheelSingle.GetSlots()
		logger.Debug().Any("wheelSingle slots", wheelSlots).Send()
	})
	wheelSlots = wheel.GetSlots()
	logger.Debug().Any("wheelSlots after test", wheelSlots).Send()
}

func TestAdvanceInSlots(t *testing.T) {
	wheelTick := 100 * time.Millisecond
	wheelSize := 200
	wheel := NewTimeWheel[TestData](wheelTick, wheelSize)

	base := int64(0)
	slotsNumber := int64(100)
	taskNumber := 20
	slotPadding := int64(2)
	for i := int64(0); i < slotsNumber; i++ {
		task := TestData{
			ID: uint32(i),
		}
		expires := base + i*int64(wheelTick)*slotPadding
		for range taskNumber {
			_ = wheel.Add(task, expires)
			// wheelSlot := wheel.Add(task, expires)
			// logger.Debug().Int("added to slot", wheelSlot).Send()
		}

	}

	// periodTicks := 200
	periodTicks := wheelSize
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	currentTick := 0
	for range ticker.C {
		due := wheel.Advance()
		if len(due) > 0 {
			logger.Debug().Any("due:", due).Send()

			// logger.Debug().Int("currentTick", currentTick).Send()
			nextExp := base + int64(currentTick+periodTicks)*int64(wheelTick)
			logger.Debug().Int64("nextExp", nextExp).Send()
			for _, connSlot := range due {
				_ = wheel.Add(connSlot, nextExp)
				// wheelSlot := wheel.Add(connSlot, nextExp)
				// logger.Debug().Int("wheelSlot", wheelSlot).Send()
			}
		}
		currentTick++
	}
}

func TestAdvanceInTime(t *testing.T) {
	wheelTick := 100 * time.Millisecond
	wheelSize := 200
	wheel := NewTimeWheel[TestData](wheelTick, wheelSize)

	base := int64(0)
	slotsNumber := int64(100)
	taskNumber := 20
	slotPadding := int64(2)
	for i := int64(0); i < slotsNumber; i++ {
		task := TestData{
			ID: uint32(i),
		}
		expires := base + i*int64(wheelTick)*slotPadding
		for range taskNumber {
			_ = wheel.Add(task, expires)
			// wheelSlot := wheel.Add(task, expires)
			// logger.Debug().Int("added to slot", wheelSlot).Send()
		}

	}

	periodTime := 20000 * time.Millisecond
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	currentTick := 0
	for range ticker.C {
		due := wheel.Advance()
		if len(due) > 0 {
			logger.Debug().Any("due:", due).Send()

			// logger.Debug().Int("currentTick", currentTick).Send()
			nextExp := time.Now().Add(periodTime).UnixNano()
			logger.Debug().Int64("nextExp", nextExp).Send()
			for _, connSlot := range due {
				_ = wheel.Add(connSlot, nextExp)
				// wheelSlot := wheel.Add(connSlot, nextExp)
				// logger.Debug().Int("wheelSlot", wheelSlot).Send()
			}
		}
		currentTick++
	}
}

func BenchmarkAdd(b *testing.B) {
	wheelTick := 10 * time.Millisecond
	wheelSize := 100000
	wheel := NewTimeWheel[TestData](wheelTick, wheelSize)

	base := int64(0)
	tasksNumber := int64(1000000)
	tick := wheelTick.Milliseconds()
	// tasksArray := make([]int, tasksNumber)
	for b.Loop() {
		for i := int64(0); i < tasksNumber; i++ {
			task := TestData{
				ID: uint32(i),
			}
			expires := base + i*tick
			_ = wheel.Add(task, expires)
			// wheelSlot := wheel.Add(task, expires)
			// tasksArray[i] = wheelSlot
		}
	}

	// log.Println(wheel.slots)
	// log.Println("test finished")
}

func BenchmarkAdvance(b *testing.B) {
	wheelTick := 10 * time.Millisecond
	wheelSize := 1000
	wheel := NewTimeWheel[TestData](wheelTick, wheelSize)

	base := int64(0)
	tasksNumber := int64(100)
	tick := wheelTick.Milliseconds()
	// log.Println("tick: ", tick)
	// log.Println("wheelTick: ", int64(wheelTick))
	for i := int64(0); i < tasksNumber*10; i += 10 {
		task := TestData{
			ID: uint32(i),
		}
		// expires := base + i*tick
		expires := base + i*int64(wheelTick)
		_ = wheel.Add(task, expires)
	}
	log.Println(wheel.slots)
	for b.Loop() {
		for i := int64(0); i < tasksNumber; i++ {
			due := wheel.Advance()
			if len(due) > 0 {
				logger.Debug().Any("due:", due).Send()

				nextExp := base + i*tick
				// logger.Debug().Int64("nextExp", nextExp).Send()
				for _, connSlot := range due {
					// wheelSlot := wheel.Add(connSlot, nextExp)
					// logger.Debug().Int("wheelSlot", wheelSlot).Send()
					_ = wheel.Add(connSlot, nextExp)
				}
			}
		}
	}

	// log.Println(wheel.slots)
	// log.Println("test finished")
}
