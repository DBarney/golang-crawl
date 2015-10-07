package queue

import (
	"testing"
)

func TestFifoQueue(test *testing.T) {
	queue := NewFifo()

	if queue.Length() != 0 {
		test.Fail()
	}

	values := []string{
		"testing1",
		"testing2",
		"testing3",
	}

	for idx, value := range values {
		queue.Push(value)

		// idx is 0 based
		if queue.Length() != idx+1 {
			test.Error("queue did not increase in size")
			test.Fail()
		}
	}

	for idx, value := range values {
		item, err := queue.Pop()

		if err != nil {
			test.Error("queue should not have returned an error")
			test.Fail()
		}

		// idx is 0 based
		if queue.Length() != len(values)-(idx+1) {
			test.Error("queue did not shrink in size")
			test.Fail()
		}

		if item != value {
			test.Error("item was out retrieved out of order")
			test.Fail()
		}
	}

	if _, err := queue.Pop(); err == nil {
		test.Error("queue should have been empty")
		test.Fail()
	}

}
