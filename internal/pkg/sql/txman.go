package sql

import (
	"context"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// TxManager - transaction manager
type TxManager struct {
	db *Balancer
}

// NewTxManager - return new transaction manager
func NewTxManager(db *Balancer) *TxManager {
	return &TxManager{
		db: db,
	}
}

// RunTX - run transaction
func (m *TxManager) RunTX(ctx context.Context, name string, op func(ctx context.Context) error) error {
	return m.runTx(ctx, name, op)
}

func (m *TxManager) runTx(ctx context.Context, name string, op func(ctx context.Context) error) error {
	var commitCalled bool
	ctx, tx, err := m.begin(ctx, name)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil && !commitCalled {
			tx.Rollback()
		}
	}()

	err = op(ctx)
	if err != nil {
		return err
	}

	commitCalled = true

	return tx.Commit()
}

// TX - transaction implementation
type TX struct {
	commitFn   func() error
	rollbackFn func()

	mtx  sync.Mutex
	done bool
}

// Commit - commit transaction if called in transaction context
func (tx *TX) Commit() error {
	tx.mtx.Lock()
	defer func() {
		tx.done = true
		tx.mtx.Unlock()
	}()

	if tx.commitFn == nil {
		return nil
	}
	return tx.commitFn()
}

// Rollback - rollback transaction if called in transaction context, safe for call any times
func (tx *TX) Rollback() {
	if tx.rollbackFn == nil {
		return
	}

	tx.mtx.Lock()
	defer func() {
		tx.done = true
		tx.mtx.Unlock()
	}()

	if tx.done {
		return
	}

	tx.rollbackFn()
}

type txKey struct{}

type txValue struct {
	name    string
	beginTs time.Time
	mtx     sync.RWMutex
	db      *sqlx.DB
}

var noopTX = &TX{}

func (m *TxManager) begin(ctx context.Context, name string) (ctxOut context.Context, txOut *TX, err error) {
	if m.isInTx(ctx) {
		return ctx, noopTX, nil
	}
	tx, err := m.db.Write(ctx).BeginTx(ctx, nil)
	if err != nil {
		return nil, noopTX, err
	}
	txval := txValue{
		name:    name,
		beginTs: time.Now(),
		db:      m.db.Write(ctx),
	}

	txHolder := &TX{
		commitFn: tx.Commit,
		rollbackFn: func() {
			err := tx.Rollback()
			if err != nil {
				logrus.Errorf("rollback %v failed, error: %v", name, err)
			}
		},
	}
	return context.WithValue(ctx, txKey{}, &txval), txHolder, nil
}

func (m *TxManager) isInTx(ctx context.Context) (ok bool) {
	tx, ok := m.getTx(ctx)
	if !ok {
		return
	}

	tx.mtx.RLock()
	ok = tx.db != nil
	tx.mtx.RUnlock()

	return
}

func (m *TxManager) getTx(ctx context.Context) (res *txValue, ok bool) {
	v := ctx.Value(txKey{})
	if v == nil {
		return
	}
	res, ok = ctx.Value(txKey{}).(*txValue)
	return
}
