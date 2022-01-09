package sql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type key struct{}

var pinKey key

type Balancer struct {
	events    chan BalancerStrategyEvent
	strategy  BalancerStrategy
	writeRole BalancerNodeRole
	readRole  BalancerNodeRole
}

func New() *Balancer {
	b := &Balancer{
		strategy:  NewBalancerStrategy(),
		events:    make(chan BalancerStrategyEvent),
		writeRole: Write,
		readRole:  Read,
	}
	b.strategy.Run(b.events)

	return b
}

func (b *Balancer) AddNode(role BalancerNodeRole, db *sqlx.DB) *Balancer {
	event := BalancerStrategyEvent{
		Event:    BalancerStrategyEventAddNode,
		Role:     role,
		Database: db,
		Done:     make(chan struct{}),
	}
	b.events <- event
	<-event.Done
	return b
}

func (b *Balancer) RemoveNode(role BalancerNodeRole, db *sqlx.DB) *Balancer {
	event := BalancerStrategyEvent{
		Event:    BalancerStrategyEventRemoveNode,
		Role:     role,
		Database: db,
		Done:     make(chan struct{}),
	}
	b.events <- event
	<-event.Done
	return b
}

func (b *Balancer) Next(ctx context.Context, role BalancerNodeRole) *sqlx.DB {
	if v, ok := ctx.Value(pinKey).(*sqlx.DB); ok && v != nil {
		return v
	}

	return b.strategy.Next(role)
}

func (b *Balancer) Node(role BalancerNodeRole, index int) *sqlx.DB {
	return b.strategy.Node(role, index)
}

func (b *Balancer) Ping(ctx context.Context) error {
	var errs []string
	for _, n := range b.strategy.Nodes() {
		var err error
		if n.Available() {
			err = n.Database().Ping()
		} else {
			err = fmt.Errorf("node is not available")
		}
		if err == nil {
			return nil
		}
		errs = append(errs, fmt.Sprintf("got err during db.Ping() from %#q node %d: %s", n.Role(), n.Index(), err))
	}
	if len(errs) > 0 {
		return fmt.Errorf("Balancer.Ping(): %s", strings.Join(errs, "; "))
	}
	return nil
}

func (b *Balancer) Close() error {
	close(b.events)

	var errs []string
	for _, n := range b.strategy.Nodes() {
		err := n.Database().Close()
		if err == nil {
			return nil
		}
		errs = append(errs, fmt.Sprintf("got err during db.Close() from %#q node %d: %s", n.Role(), n.Index(), err))
	}
	if len(errs) > 0 {
		return fmt.Errorf("Balancer.Close(): %s", strings.Join(errs, "; "))
	}
	return nil
}

func (b *Balancer) Read(ctx context.Context) *sqlx.DB {
	return b.Next(ctx, b.readRole)
}

func (b *Balancer) Write(ctx context.Context) *sqlx.DB {
	return b.Next(ctx, b.writeRole)
}

func (b *Balancer) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return b.Write(ctx).ExecContext(ctx, query, args...)
}

func (b *Balancer) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return b.Read(ctx).QueryContext(ctx, query, args...)
}

func (b *Balancer) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return b.Read(ctx).QueryRowContext(ctx, query, args...)
}

func (b *Balancer) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return b.Write(ctx).PrepareContext(ctx, query)
}

func (b *Balancer) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return b.Read(ctx).GetContext(ctx, dest, query, args...)
}

func (b *Balancer) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return b.Read(ctx).SelectContext(ctx, dest, query, args...)
}

func (b *Balancer) Begin(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return b.Write(ctx).BeginTx(ctx, opts)
}

func (b *Balancer) Pin(ctx context.Context, db *sqlx.DB) context.Context {
	return context.WithValue(ctx, pinKey, db)
}

func (b *Balancer) Unpin(ctx context.Context) context.Context {
	return b.Pin(ctx, nil)
}
