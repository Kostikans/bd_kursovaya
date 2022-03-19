package model

type Partition struct {
	ID                  uint64        `db:"id"`
	PartitionType       PartitionType `db:"partition_type"`
	PartitionStartIndex string        `db:"partition_start_index"`
	PartitionEndIndex   string        `db:"partition_end_index"`
	Count               uint64        `db:"elements_count"`
}

// PartitionType - partition type
type PartitionType int16

// PartitionType
const (
	PartitionTypeComments = 0
	PartitionTypeVote     = 1
)
