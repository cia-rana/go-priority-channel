package prioritychannel

import (
	"sort"
	"testing"
	"time"
)

func Test(t *testing.T) {
	testData := []Item{
		Item{0, "a"},
		Item{4, "e"},
		Item{1, "c"},
		Item{1, "b"},
		Item{3, "d"},
	}

	pc := New()
	go func() {
		for _, item := range testData {
			pc.In <- item
		}
	}()

	time.Sleep(100 * time.Millisecond)

	sort.Slice(testData, func(i, j int) bool {
		return testData[i].Priority > testData[j].Priority
	})

	for _, item := range testData {
		pcOut := <-pc.Out
		t.Log(pcOut, item)
		if pcOut != item {
			t.Fatal()
		}
	}
}
