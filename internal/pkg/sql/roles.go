package sql

import "github.com/jmoiron/sqlx"

type BalancerNodeRole interface {
	String() string
}

type BalancerRoleFallback interface {
	Fallback() BalancerNodeRole
}

type BalancerStrategyEventType int

const (
	BalancerStrategyEventAddNode BalancerStrategyEventType = iota
	BalancerStrategyEventRemoveNode
	BalancerStrategyEventHealthcheck
)

type BalancerStrategyEvent struct {
	Event    BalancerStrategyEventType
	Database *sqlx.DB
	Role     BalancerNodeRole
	Done     chan struct{}
}

// BalancerStrategy defines the behavior of the balancer
type BalancerStrategy interface {
	// Next returns next available Database node by BalancerNodeRole
	Next(BalancerNodeRole) *sqlx.DB
	// Node returns node by role and index
	Node(BalancerNodeRole, int) *sqlx.DB
	// Nodes returns all nodes
	Nodes() []BalancerNode
	// Run listens update event from chan and async checks available nodes
	Run(<-chan BalancerStrategyEvent)
}

// BalancerNode returns node
type BalancerNode interface {
	// Index returns sequence number of the node
	Index() int
	// Role returns role of the node
	Role() BalancerNodeRole
	// Database returns database object of the node
	Database() *sqlx.DB
	// Available returns true if the last database checks is ok
	Available() bool
}

type role int

func (r role) String() string {
	switch r {
	case 0:
		return "read"
	case 1:
		return "write"
	default:
		return "unknown"
	}
}

var (
	// Read is the node role for slave replica
	Read BalancerNodeRole = role(0)
	// Write is the node role for master replica
	Write BalancerNodeRole = role(1)
)

func (r role) Fallback() BalancerNodeRole {
	if r == Read {
		return Write
	}
	return nil
}
