package pipeline

import (
	"errors"
	"testing"
)

func TestSinglePipeline(t *testing.T) {
	count := 0
	source := make(chan interface{}, 1)
	handle := func(in interface{}) (interface{}, error) {
		count++
		switch {
		case count == 5:
			return "", errors.New("this is a mistake")
		default:
			return in, nil
		}
	}
	pipe := NewPipeline(source, 1, handle)
	for i := 0; i < 10; i++ {
		source <- string(i)
		select {
		case res := <-pipe.Output():
			if res.(string) != string(i) {
				t.Logf("wrong thing was received %v", res)
				t.Fail()
			}
		case <-pipe.Err():
			if i != 4 {
				t.Log("error was sent at the wrong time")
				t.Fail()
			}
		}

		// i is always 1 less then count
		if count != i+1 {
			t.Log("wrong number of messages were sent")
			t.Fail()
		}
	}

	close(source)
	_, open := <-pipe.Output()
	if open {
		t.Logf("the destination was not closed correctly")
		t.Fail()
	}
}
