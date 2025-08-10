package segment

import (
	"bytes"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/ryan961/jtt"
)

// helper to create deterministic but distinct bodies per index
func makeBody(index uint16) []byte {
	return []byte("part_" + string(rune('A'+index)))
}

func TestPool_Cache_ConcurrentSameKey(t *testing.T) {
	pool := NewPool()

	phone := "13800138000"
	msgID := jtt.MsgID(0x0001)
	total := uint16(8)

	// Prepare shuffled indices
	idxs := make([]uint16, 0, total)
	for i := uint16(0); i < total; i++ {
		idxs = append(idxs, i+1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(idxs), func(i, j int) { idxs[i], idxs[j] = idxs[j], idxs[i] })

	var wg sync.WaitGroup
	wg.Add(len(idxs))

	completedCh := make(chan []byte, 4)

	for _, idx := range idxs {
		idx := idx // capture
		go func() {
			defer wg.Done()
			header := createTestHeader(phone, msgID, total, idx)
			body := makeBody(idx)
			isCompleted, buffers := pool.Cache(header, body)
			if isCompleted && buffers != nil {
				completedCh <- buffers
			}
		}()
	}

	wg.Wait()
	close(completedCh)

	// Expect exactly one completion
	var completions int
	var result []byte
	for b := range completedCh {
		completions++
		result = b
	}
	if completions != 1 {
		t.Fatalf("expected exactly 1 completion, got %d", completions)
	}

	// Build expected merged body by index order
	indices := make([]int, 0, total)
	for i := 1; i <= int(total); i++ {
		indices = append(indices, i)
	}
	sort.Ints(indices)
	expected := make([]byte, 0)
	for _, i := range indices {
		expected = append(expected, makeBody(uint16(i))...)
	}

	if !bytes.Equal(result, expected) {
		t.Fatalf("merged body mismatch: expected %v, got %v", expected, result)
	}
}

func TestPool_CacheWithTimeout_ConcurrentSameKey(t *testing.T) {
	pool := NewPool()

	phone := "13800138001"
	msgID := jtt.MsgID(0x0001)
	total := uint16(5)
	ttl := 2 * time.Second

	idxs := make([]uint16, 0, total)
	for i := uint16(0); i < total; i++ {
		idxs = append(idxs, i+1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(idxs), func(i, j int) { idxs[i], idxs[j] = idxs[j], idxs[i] })

	var wg sync.WaitGroup
	wg.Add(len(idxs))

	completedCh := make(chan []byte, 4)

	for _, idx := range idxs {
		idx := idx // capture
		go func() {
			defer wg.Done()
			header := createTestHeader(phone, msgID, total, idx)
			body := makeBody(idx)
			isCompleted, buffers := pool.CacheWithTimeout(header, body, ttl)
			if isCompleted && buffers != nil {
				completedCh <- buffers
			}
		}()
	}

	wg.Wait()
	close(completedCh)

	var completions int
	var result []byte
	for b := range completedCh {
		completions++
		result = b
	}
	if completions != 1 {
		t.Fatalf("expected exactly 1 completion, got %d", completions)
	}

	expected := make([]byte, 0)
	for i := 1; i <= int(total); i++ {
		expected = append(expected, makeBody(uint16(i))...)
	}
	if !bytes.Equal(result, expected) {
		t.Fatalf("merged body mismatch: expected %v, got %v", expected, result)
	}
}
