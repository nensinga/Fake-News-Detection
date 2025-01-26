package sphinx

import (
	"encoding/binary"
	"io"
)

// ReplaySet is a data structure used to efficiently track replays by sequence
// numbers during batch processing. It supports operations like membership
// queries, merging, and serialization.
type ReplaySet struct {
	replays map[uint16]struct{}
}

// NewReplaySet initializes an empty ReplaySet.
func NewReplaySet() *ReplaySet {
	return &ReplaySet{
		replays: make(map[uint16]struct{}),
	}
}

// Size returns the number of elements in the ReplaySet.
func (rs *ReplaySet) Size() int {
	return len(rs.replays)
}

// Add inserts a sequence number into the ReplaySet.
func (rs *ReplaySet) Add(idx uint16) {
	rs.replays[idx] = struct{}{}
}

// Contains checks if a sequence number exists in the ReplaySet.
func (rs *ReplaySet) Contains(idx uint16) bool {
	_, exists := rs.replays[idx]
	return exists
}

// Merge combines another ReplaySet into the current ReplaySet.
func (rs *ReplaySet) Merge(other *ReplaySet) {
	for seqNum := range other.replays {
		rs.Add(seqNum)
	}
}

// Encode serializes the ReplaySet into an `io.Writer` for storage or transmission.
// The replay set can be reconstructed using `Decode`.
func (rs *ReplaySet) Encode(w io.Writer) error {
	for seqNum := range rs.replays {
		if err := binary.Write(w, binary.BigEndian, seqNum); err != nil {
			return err
		}
	}
	return nil
}

// Decode reconstructs a ReplaySet from an `io.Reader`.
// Expects the input to be a valid sequence of `uint16` values.
func (rs *ReplaySet) Decode(r io.Reader) error {
	for {
		var seqNum uint16

		// Attempt to read the next sequence proposal
		err := binary.Read(r, binary.BigEndian, &seqNum)
		if err == io.EOF {
			// Successfully reached the end of the input.
			return nil
		} else if err != nil {
			// Handle unexpected read errors.
			return err
		}

		// Add the decoded sequence number to the ReplaySet.
		rs.Add(seqNum)
	}
}
