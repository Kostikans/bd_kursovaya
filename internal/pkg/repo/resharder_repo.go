package repo

import (
	"context"

	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
	"github.com/kostikan/bd_kursovaya/internal/pkg/sql"
)

// ResharderRepo - resharder repo
type ResharderRepo struct {
	db *sql.Balancer
}

// NewResharderRepo - returns new postRepo
func NewResharderRepo(db *sql.Balancer) *ResharderRepo {
	return &ResharderRepo{
		db: db,
	}
}

func (r *ResharderRepo) GetPartitions(ctx context.Context, partitionType model.PartitionType) (res []model.Partition, err error) {
	query := `SELECT id, partition_type, partition_start_index, partition_end_index, elements_count 
FROM partitions WHERE partition_type = $1`

	err = r.db.Write(ctx).Select(&res, query, partitionType)
	return
}

func (r *ResharderRepo) AddCountTrigger(ctx context.Context, partitionType model.PartitionType) (err error) {
	query := `SELECT id, partition_type, partition_start_index, partition_end_index 
FROM partitions WHERE partition_type = $1`

	_, err = r.db.Write(ctx).Exec(query, partitionType)
	return
}
