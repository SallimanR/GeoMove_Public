package timewheel

import (
	"sync"
	"time"
)

type position struct {
	slot int
	idx  int
}

type TimeWheel[SlotData comparable] struct {
	slots           [][]SlotData
	slotPositionMap map[SlotData]position
	Current         int
	mu              sync.Mutex

	Tick      time.Duration
	wheelSize int
}

func NewTimeWheel[SlotData comparable](tick time.Duration, wheelSize int) *TimeWheel[SlotData] {
	slots := make([][]SlotData, wheelSize)
	slotPositionMap := make(map[SlotData]position, wheelSize)
	return &TimeWheel[SlotData]{
		slots:           slots,
		slotPositionMap: slotPositionMap,
		Current:         0,
		Tick:            tick,
		wheelSize:       wheelSize,
	}
}

// Add schedules a connection at absolute time `expires` (UnixNano). Returns slot index.
func (tw *TimeWheel[SlotData]) Add(task SlotData, expires int64) int {
	slot := (expires / int64(tw.Tick)) % int64(tw.wheelSize)
	tw.mu.Lock()
	idx := len(tw.slots[slot])
	tw.slots[slot] = append(tw.slots[slot], task)
	tw.slotPositionMap[task] = position{slot: int(slot), idx: idx}
	tw.mu.Unlock()
	return int(slot)
}

// Update changes expiration time of and existing task.
// Equivalent to Remove followed by Add, but under sigle lock.
func (tw *TimeWheel[SlotData]) Update(task SlotData, oldExpires, newExpires int64) bool {
	tw.mu.Lock()
	pos, ok := tw.slotPositionMap[task]
	if !ok {
		tw.mu.Unlock()
		return false
	}
	slot := pos.slot
	idx := pos.idx

	lastIdx := len(tw.slots) - 1
	if lastIdx > 0 {
		lastTask := tw.slots[slot][lastIdx]
		tw.slots[slot][idx] = lastTask
		tw.slotPositionMap[lastTask] = position{slot: slot, idx: idx}
	}
	delete(tw.slotPositionMap, task)

	// Add new slot
	newSlot := (newExpires / int64(tw.Tick)) % int64(tw.wheelSize)
	newIdx := len(tw.slots[newSlot])
	tw.slots[newSlot] = append(tw.slots[newSlot], task)
	tw.slotPositionMap[task] = position{slot: int(newSlot), idx: newIdx}
	tw.mu.Unlock()
	return true
}

func (tw *TimeWheel[SlotData]) Remove(task SlotData) bool {
	tw.mu.Lock()
	pos, ok := tw.slotPositionMap[task]
	if !ok {
		tw.mu.Unlock()
		return false
	}
	slot := pos.slot
	idx := pos.idx

	// Swap and remove for efficient removal
	//
	// Swap with last
	lastIdx := len(tw.slots[slot]) - 1
	lastTask := tw.slots[slot][lastIdx]
	if idx != lastIdx {
		tw.slots[slot][idx] = lastTask
		tw.slotPositionMap[lastTask] = position{slot: slot, idx: idx}
	}
	// Remove last
	tw.slots[slot] = tw.slots[slot][:lastIdx]
	delete(tw.slotPositionMap, task)
	tw.mu.Unlock()
	return true
}

// Advance returns the list of due SlotData for the current tick and moves the hand to the next slot
func (tw *TimeWheel[SlotData]) Advance() []SlotData {
	tw.mu.Lock()
	due := tw.slots[tw.Current]
	tw.slots[tw.Current] = nil

	for _, task := range due {
		delete(tw.slotPositionMap, task)
	}

	tw.Current = (tw.Current + 1) % tw.wheelSize
	tw.mu.Unlock()
	return due
}

func (tw *TimeWheel[SlotData]) GetTasksCount() int {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	return len(tw.slotPositionMap)
}

func (tw *TimeWheel[SlotData]) GetNonEmptySlotsCount() int {
	tw.mu.Lock()
	count := 0
	for _, slot := range tw.slots {
		if len(slot) > 0 {
			count += 1
		}
	}
	tw.mu.Unlock()
	return count
}

func (tw *TimeWheel[SlotData]) GetSlots() [][]SlotData {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	return tw.slots
}
