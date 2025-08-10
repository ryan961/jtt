package segment

import (
	"bytes"
	"testing"
	"time"

	"github.com/ryan961/jtt"
)

// createTestHeader creates a test message header
func createTestHeader(phoneNumber string, msgID jtt.MsgID, total, index uint16) *jtt.MsgHeader {
	return &jtt.MsgHeader{
		MsgID:        msgID,
		PhoneNumber:  phoneNumber,
		SerialNumber: 1,
		Property: &jtt.Property{
			BodyLength:   100,
			Encryption:   0,
			Segmentation: 1,
			VersionSign:  jtt.Version2013,
		},
		SegmentInfo: &jtt.SegmentInfo{
			Total: total,
			Index: index,
		},
	}
}

// createTestBody creates a test message body
func createTestBody(index uint16, content string) []byte {
	return []byte(content + "_part" + string(rune('0'+index)))
}

// TestNewPool tests the creation of cache pool and options configuration
func TestNewPool(t *testing.T) {
	// Test default configuration
	pool := NewPool()
	if pool == nil {
		t.Fatal("Failed to create pool with default options")
	}

	// Test custom capacity
	customCapacity := 1024 * 1024 // 1MB
	pool = NewPool(WithCapacity(customCapacity))
	if pool.capacity != customCapacity {
		t.Errorf("Expected capacity %d, got %d", customCapacity, pool.capacity)
	}

	// Test custom initial capacity
	customInitialCapacity := 500
	pool = NewPool(WithInitialCapacity(customInitialCapacity))
	if pool.initialCapacity != customInitialCapacity {
		t.Errorf("Expected initialCapacity %d, got %d", customInitialCapacity, pool.initialCapacity)
	}

	// Test custom TTL
	customTTL := 60 * time.Second
	pool = NewPool(WithVariableTTL(customTTL))
	if pool.variableTTL != customTTL {
		t.Errorf("Expected variableTTL %v, got %v", customTTL, pool.variableTTL)
	}

	// Test combined options
	pool = NewPool(
		WithCapacity(customCapacity),
		WithInitialCapacity(customInitialCapacity),
		WithVariableTTL(customTTL),
	)
	if pool.capacity != customCapacity {
		t.Errorf("Expected capacity %d, got %d", customCapacity, pool.capacity)
	}
	if pool.initialCapacity != customInitialCapacity {
		t.Errorf("Expected initialCapacity %d, got %d", customInitialCapacity, pool.initialCapacity)
	}
	if pool.variableTTL != customTTL {
		t.Errorf("Expected variableTTL %v, got %v", customTTL, pool.variableTTL)
	}
}

// TestPool_Cache tests basic caching functionality
func TestPool_Cache(t *testing.T) {
	pool := NewPool()

	// Create test data
	phoneNumber := "13800138000"
	msgID := jtt.MsgID(0x0200)

	// Test single segment (incomplete)
	header := createTestHeader(phoneNumber, msgID, 3, 1)
	body := createTestBody(1, "test")

	isCompleted, buffers := pool.Cache(header, body)
	if isCompleted {
		t.Error("Expected incomplete segment, but got completed")
	}
	if buffers != nil {
		t.Error("Expected nil buffers for incomplete segment")
	}
}

// TestPool_Cache_Complete tests the merging of a complete message
func TestPool_Cache_Complete(t *testing.T) {
	pool := NewPool()

	// Create test data
	phoneNumber := "13800138000"
	msgID := jtt.MsgID(0x0200)

	// Total number of segments
	total := uint16(3)

	// Create and cache the first segment
	header1 := createTestHeader(phoneNumber, msgID, total, 1)
	body1 := createTestBody(1, "test")
	isCompleted, _ := pool.Cache(header1, body1)
	if isCompleted {
		t.Error("Expected incomplete segment after first part")
	}

	// Create and cache the second segment
	header2 := createTestHeader(phoneNumber, msgID, total, 2)
	body2 := createTestBody(2, "test")
	isCompleted, _ = pool.Cache(header2, body2)
	if isCompleted {
		t.Error("Expected incomplete segment after second part")
	}

	// Create and cache the third segment (the last one)
	header3 := createTestHeader(phoneNumber, msgID, total, 3)
	body3 := createTestBody(3, "test")
	isCompleted, buffers := pool.Cache(header3, body3)

	// Now it should be complete
	if !isCompleted {
		t.Error("Expected completed segment after all parts")
	}
	if buffers == nil {
		t.Fatal("Expected non-nil buffers for complete segment")
	}

	// Verify the merged content
	expected := append(append(body1, body2...), body3...)
	if !bytes.Equal(buffers, expected) {
		t.Errorf("Expected merged body %v, got %v", expected, buffers)
	}
}

// TestPool_CacheWithTimeout tests caching functionality with timeout
func TestPool_CacheWithTimeout(t *testing.T) {
	pool := NewPool()

	// Create test data
	phoneNumber := "13800138000"
	msgID := jtt.MsgID(0x0300)
	timeout := 2 * time.Second

	// Total number of segments
	total := uint16(2)

	// Create and cache the first segment with timeout
	header1 := createTestHeader(phoneNumber, msgID, total, 1)
	body1 := createTestBody(1, "timeout")
	isCompleted, buffers := pool.CacheWithTimeout(header1, body1, timeout)
	if isCompleted {
		t.Error("Expected incomplete segment after first part")
	}

	// Create and cache the second segment with timeout
	header2 := createTestHeader(phoneNumber, msgID, total, 2)
	body2 := createTestBody(2, "timeout")
	isCompleted, _ = pool.CacheWithTimeout(header2, body2, timeout)

	// Now it should be complete
	if !isCompleted {
		t.Error("Expected completed segment after all parts")
	}
	if buffers == nil {
		t.Fatal("Expected non-nil buffers for complete segment")
	}

	// Verify the merged content
	expected := append(body1, body2...)
	if !bytes.Equal(buffers, expected) {
		t.Errorf("Expected merged body %v, got %v", expected, buffers)
	}
}

// TestPool_Cache_NilHeader tests handling of nil header
func TestPool_Cache_NilHeader(t *testing.T) {
	pool := NewPool()

	// Test nil header
	isCompleted, buffers := pool.Cache(nil, []byte("test"))
	if isCompleted {
		t.Error("Expected incomplete segment with nil header")
	}
	if buffers != nil {
		t.Error("Expected nil buffers with nil header")
	}

	// Test nil header with timeout
	isCompleted, buffers = pool.CacheWithTimeout(nil, []byte("test"), 1*time.Second)
	if isCompleted {
		t.Error("Expected incomplete segment with nil header in CacheWithTimeout")
	}
	if buffers != nil {
		t.Error("Expected nil buffers with nil header in CacheWithTimeout")
	}
}

// TestPool_Cache_NoSegmentInfo tests handling of no segment info
func TestPool_Cache_NoSegmentInfo(t *testing.T) {
	pool := NewPool()

	// Create a header without segment info
	header := &jtt.MsgHeader{
		MsgID:        jtt.MsgID(0x0200),
		PhoneNumber:  "13800138000",
		SerialNumber: 1,
		Property: &jtt.Property{
			BodyLength:   100,
			Encryption:   0,
			Segmentation: 0, // No segmentation
			VersionSign:  jtt.Version2013,
		},
		SegmentInfo: nil, // No segment info
	}

	// Test no segment info
	isCompleted, buffers := pool.Cache(header, []byte("test"))
	if isCompleted {
		t.Error("Expected incomplete segment with no segment info")
	}
	if buffers != nil {
		t.Error("Expected nil buffers with no segment info")
	}

	// Test no segment info with timeout
	isCompleted, buffers = pool.CacheWithTimeout(header, []byte("test"), 1*time.Second)
	if isCompleted {
		t.Error("Expected incomplete segment with no segment info in CacheWithTimeout")
	}
	if buffers != nil {
		t.Error("Expected nil buffers with no segment info in CacheWithTimeout")
	}
}
