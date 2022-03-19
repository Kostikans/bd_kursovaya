package updater

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
)

const (
	checkTime = 30
)

type ResharderRepo interface {
	GetPartitions(ctx context.Context, partitionType model.PartitionType) (res []model.Partition, err error)
}

type Resharder struct {
	repo ResharderRepo
}

func NewResharder(repo ResharderRepo) *Resharder {
	return &Resharder{
		repo: repo,
	}
}

func (r *Resharder) StartResharder(ctx context.Context) {
	go func(ctx context.Context) {
		randomizedRunTime := time.Duration(rand.Int63n(int64(checkTime*time.Second))) + checkTime*time.Second // nolint
		ticker := time.NewTicker(randomizedRunTime)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				log.Printf("start partitions checking")
				err := r.ObservePartitions(ctx)
				if err != nil {
					log.Println("can't dump background job metrics")
				}
				randomizedRunTime = time.Duration(rand.Int63n(int64(checkTime*time.Second))) + checkTime*time.Second // nolint
				ticker.Reset(randomizedRunTime)
			}
		}
	}(ctx)
}

func (r *Resharder) ObservePartitions(ctx context.Context) (err error) {
	partitionTypes := []model.PartitionType{
		model.PartitionTypeComments,
		model.PartitionTypeVote,
	}
	for _, partitionType := range partitionTypes {
		partitionType := partitionType
		go func() {
			err := r.CheckPartition(ctx, partitionType)
			if err != nil {
				log.Printf("error while check partition: %s", err.Error())
			}
		}()
	}

	return
}

func (r *Resharder) CheckPartition(ctx context.Context, partitionType model.PartitionType) (err error) {
	partitions, err := r.repo.GetPartitions(ctx, partitionType)
	if err != nil {
		return err
	}
	var count uint64
	for _, partition := range partitions {
		count += partition.Count
	}
	return
}
