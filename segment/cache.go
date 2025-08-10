package segment

import (
	"strconv"
	"sync"
	"time"

	"github.com/maypok86/otter/v2"
	"github.com/ryan961/jtt"
)

var (
	segmentPool = &sync.Pool{
		New: func() any {
			return &jtt.Segment{Data: make(map[uint16][]byte)}
		},
	}
)

type Pool struct {
	pool *otter.Cache[string, *jtt.Segment]

	capacity        int
	initialCapacity int
	variableTTL     time.Duration
}

func NewPool(opts ...PoolOptions) *Pool {
	p := &Pool{
		capacity:        3 * 1024 * 1024 * 1024, // 3G
		initialCapacity: 1000,
		variableTTL:     300 * time.Second,
	}
	for _, opt := range opts {
		opt(p)
	}

	p.pool = otter.Must(&otter.Options[string, *jtt.Segment]{
		MaximumSize:      p.capacity,
		InitialCapacity:  p.initialCapacity,
		ExpiryCalculator: otter.ExpiryWriting[string, *jtt.Segment](p.variableTTL),
		OnAtomicDeletion: func(e otter.DeletionEvent[string, *jtt.Segment]) {
			if e.Cause == otter.CauseOverflow || e.Cause == otter.CauseExpiration {
				if e.Value != nil {
					e.Value.Reset()
					segmentPool.Put(e.Value)
				}
			}
		},
	})
	return p
}

// Cache adds a segment to the cache, returns true and the complete body if the segment is complete.
// Otherwise, returns false and nil.
func (s *Pool) Cache(header *jtt.MsgHeader, body []byte) (isCompleted bool, buffers []byte) {
	if header == nil || header.SegmentInfo == nil {
		return false, nil
	}
	key := header.PhoneNumber + ":" + strconv.FormatInt(int64(header.MsgID), 10)
	// Use atomic per-key compute to avoid global lock and ensure correctness under concurrency
	_, _ = s.pool.Compute(key, func(oldValue *jtt.Segment, found bool) (*jtt.Segment, otter.ComputeOp) {
		var segment *jtt.Segment
		if !found || oldValue == nil {
			segment = segmentPool.Get().(*jtt.Segment)
			segment.PhoneNumber = header.PhoneNumber
			segment.MsgID = header.MsgID
			segment.Total = header.SegmentInfo.Total
		} else {
			segment = oldValue
		}

		segment.Merge(header.SegmentInfo, body)
		if segment.IsComplete() {
			buffers = segment.GetBody()
			// recycle and remove from cache
			segment.Reset()
			segmentPool.Put(segment)
			return nil, otter.InvalidateOp
		}
		return segment, otter.WriteOp
	})
	// buffers != nil implies completion occurred during compute
	return buffers != nil, buffers
}

// CacheWithTimeout adds a segment to the cache with a specific timeout, returns true and the complete body
// if the segment is complete. Otherwise, returns false and nil.
func (s *Pool) CacheWithTimeout(header *jtt.MsgHeader, body []byte, timeout time.Duration) (isCompleted bool, buffers []byte) {
	if header == nil || header.SegmentInfo == nil {
		return false, nil
	}
	key := header.PhoneNumber + ":" + strconv.FormatInt(int64(header.MsgID), 10)
	_, ok := s.pool.Compute(key, func(oldValue *jtt.Segment, found bool) (*jtt.Segment, otter.ComputeOp) {
		var segment *jtt.Segment
		if !found || oldValue == nil {
			segment = segmentPool.Get().(*jtt.Segment)
			segment.PhoneNumber = header.PhoneNumber
			segment.MsgID = header.MsgID
			segment.Total = header.SegmentInfo.Total
		} else {
			segment = oldValue
		}

		segment.Merge(header.SegmentInfo, body)
		if segment.IsComplete() {
			buffers = segment.GetBody()
			// recycle and remove from cache
			segment.Reset()
			segmentPool.Put(segment)
			return nil, otter.InvalidateOp
		}
		return segment, otter.WriteOp
	})
	if ok && buffers == nil {
		// still incomplete, set TTL for the key
		s.pool.SetExpiresAfter(key, timeout)
	}
	// buffers != nil implies completion occurred during compute
	return buffers != nil, buffers
}
