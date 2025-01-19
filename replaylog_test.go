package sphinx

import (
	"testing"
)

// TestMemoryReplayLogStorageAndRetrieval tests the basic storage and retrieval
// functionality of MemoryReplayLog.
func TestMemoryReplayLogStorageAndRetrieval(t *testing.T) {
	rl := NewMemoryReplayLog()
	rl.Start()
	defer rl.Stop()

	var hashPrefix HashPrefix
	hashPrefix[0] = 1
	var cltv1 uint32 = 1

	// Lookup an unknown sphinx packet.
	_, err := rl.Get(&hashPrefix)
	if err == nil {
		t.Fatalf("Expected ErrLogEntryNotFound")
	}
	if err != ErrLogEntryNotFound {
		t.Fatalf("Unexpected error on Get: %v", err)
	}

	// Log a new sphinx packet.
	err = rl.Put(&hashPrefix, cltv1)
	if err != nil {
		t.Fatalf("Unexpected error on Put: %v", err)
	}

	// Attempt to replay the sphinx packet.
	err = rl.Put(&hashPrefix, cltv1)
	if err == nil {
		t.Fatalf("Expected ErrReplayedPacket")
	}
	if err != ErrReplayedPacket {
		t.Fatalf("Unexpected error on Put: %v", err)
	}

	// Retrieve the logged sphinx packet.
	cltv, err := rl.Get(&hashPrefix)
	if err != nil {
		t.Fatalf("Unexpected error on Get: %v", err)
	}
	if cltv != cltv1 {
		t.Fatalf("Get returned wrong value: expected %v, got %v", cltv1, cltv)
	}

	// Delete the sphinx packet from the log.
	err = rl.Delete(&hashPrefix)
	if err != nil {
		t.Fatalf("Unexpected error on Delete: %v", err)
	}

	// Attempt to retrieve the deleted sphinx packet.
	_, err = rl.Get(&hashPrefix)
	if err == nil {
		t.Fatalf("Expected ErrLogEntryNotFound")
	}
	if err != ErrLogEntryNotFound {
		t.Fatalf("Unexpected error on Get: %v", err)
	}

	// Reinsert the sphinx packet with a new value.
	var cltv2 uint32 = 2
	err = rl.Put(&hashPrefix, cltv2)
	if err != nil {
		t.Fatalf("Unexpected error on Put: %v", err)
	}

	// Retrieve the updated sphinx packet.
	cltv, err = rl.Get(&hashPrefix)
	if err != nil {
		t.Fatalf("Unexpected error on Get: %v", err)
	}
	if cltv != cltv2 {
		t.Fatalf("Get returned wrong value: expected %v, got %v", cltv2, cltv)
	}
}

// TestMemoryReplayLogPutBatch tests batch insertion and replay handling in
// MemoryReplayLog.
func TestMemoryReplayLogPutBatch(t *testing.T) {
	rl := NewMemoryReplayLog()
	rl.Start()
	defer rl.Stop()

	var hashPrefix1, hashPrefix2 HashPrefix
	hashPrefix1[0] = 1
	hashPrefix2[0] = 2

	// Create a batch with a duplicated packet.
	batch1 := NewBatch([]byte{1})
	if err := batch1.Put(1, &hashPrefix1, 1); err != nil {
		t.Fatalf("Unexpected error adding entry to batch: %v", err)
	}
	if err := batch1.Put(1, &hashPrefix1, 1); err != nil {
		t.Fatalf("Unexpected error adding duplicate entry to batch: %v", err)
	}

	replays, err := rl.PutBatch(batch1)
	if replays.Size() != 1 || !replays.Contains(1) {
		t.Fatalf("Unexpected replay set after adding batch 1: %v", err)
	}

	// Create a batch with one replayed packet and one valid packet.
	batch2 := NewBatch([]byte{2})
	if err := batch2.Put(1, &hashPrefix1, 1); err != nil {
		t.Fatalf("Unexpected error adding entry to batch: %v", err)
	}
	if err := batch2.Put(2, &hashPrefix2, 2); err != nil {
		t.Fatalf("Unexpected error adding entry to batch: %v", err)
	}

	replays, err = rl.PutBatch(batch2)
	if replays.Size() != 1 || !replays.Contains(1) {
		t.Fatalf("Unexpected replay set after adding batch 2: %v", err)
	}

	// Reprocess batch 2, ensuring idempotency.
	replays, err = rl.PutBatch(batch2)
	if replays.Size() != 1 || !replays.Contains(1) {
		t.Fatalf("Unexpected replay set after reprocessing batch 2: %v", err)
	}
}
