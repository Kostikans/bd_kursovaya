package sql

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

var _ BalancerStrategy = &RoundRobin{}

const defaultCheckInterval = 10 * time.Second

func NewBalancerStrategy() *RoundRobin {
	rr := &RoundRobin{
		checkInterval: defaultCheckInterval,
		indexes:       make(map[BalancerNodeRole]int),
		nodes:         make(map[BalancerNodeRole][]*node),
	}
	return rr
}

type RoundRobin struct {
	sync.Mutex

	indexes       map[BalancerNodeRole]int
	nodes         map[BalancerNodeRole][]*node
	checkInterval time.Duration
}

func (b *RoundRobin) Node(role BalancerNodeRole, index int) *sqlx.DB {
	b.Lock()
	defer b.Unlock()

	nodes := b.nodes[role]
	if index >= len(nodes) || nodes[index] == nil {
		return nil
	}

	return nodes[index].db
}

func (b *RoundRobin) Next(role BalancerNodeRole) *sqlx.DB {
	b.Lock()
	defer b.Unlock()

	var rr map[BalancerNodeRole]struct{}

	for {
		db, err := b.next(role)
		if err == nil {
			return db
		}

		r, ok := role.(BalancerRoleFallback)
		if !ok || r.Fallback() == nil {
			return nil
		}

		if rr == nil {
			rr = map[BalancerNodeRole]struct{}{}
		}
		rr[role] = struct{}{}

		role = r.Fallback()
		_, processed := rr[role]
		if processed {
			return nil
		}
	}
}

func (b *RoundRobin) next(role BalancerNodeRole) (*sqlx.DB, error) {
	if role == nil {
		return nil, errors.New("unknown role")
	}

	i, ok := b.indexes[role]

	if !ok {
		return nil, fmt.Errorf("unknown role %d", role)
	}
	nodes := len(b.nodes[role])
	if nodes == 0 {
		return nil, fmt.Errorf("no nodes, role %d", role)
	}

	for c := 0; c < nodes; c++ {
		n := b.nodes[role][(i+c)%nodes]
		b.indexes[role]++
		if n.ok() {
			return n.db, nil
		}

	}
	return nil, fmt.Errorf("no nodes, role %d", role)
}

func (b *RoundRobin) Run(events <-chan BalancerStrategyEvent) {
	done := make(chan struct{})
	go func() {
		for event := range events {
			switch event.Event {
			case BalancerStrategyEventAddNode:
				b.Lock()
				b.nodes[event.Role] = append(b.nodes[event.Role], &node{
					db:      event.Database,
					lastErr: nil,
				})
				b.indexes[event.Role] = 0
				b.Unlock()
			case BalancerStrategyEventRemoveNode:
				b.Lock()
				for i, n := range b.nodes[event.Role] {
					if n.db == event.Database {
						b.removeNode(event.Role, i)
						break
					}
				}
				b.Unlock()
			}
			close(event.Done)
		}
		close(done)
	}()
}

func (b *RoundRobin) removeNode(r BalancerNodeRole, i int) {
	if len(b.nodes[r]) <= i {
		return
	}

	b.nodes[r] = append(b.nodes[r][:i], b.nodes[r][i+1:]...)
}

func (b *RoundRobin) Nodes() (result []BalancerNode) {
	b.Lock()
	defer b.Unlock()

	for role, nodes := range b.nodes {
		for index, n := range nodes {
			if n == nil {
				continue
			}
			result = append(result, &node{
				index:   index,
				role:    role,
				db:      n.db,
				lastErr: n.lastErr,
			})
		}
	}
	return
}

type node struct {
	db      *sqlx.DB
	lastErr error
	index   int
	role    BalancerNodeRole
}

func (n *node) Index() int {
	return n.index
}

func (n *node) ok() bool {
	return n.lastErr == nil
}

func (n *node) Role() BalancerNodeRole {
	return n.role
}

func (n *node) Database() *sqlx.DB {
	return n.db
}

func (n *node) Available() bool {
	return n.ok()
}
