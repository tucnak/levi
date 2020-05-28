package levi

import (
	"fmt"
)

// Register makes Leviathan aware of the models.
//
// Please, always put this into init(), in any case
// make sure all models are registered before Wake(),
// because it's When the migration happens.
//
func Register(models ...Model) {
	for _, model := range models {
		if model == nil {
			panic(ErrNilModel)
		}

		switch model.Type() {
		case TABLE:
			tables = append(tables, model)
		case QUEUE:
			cues = append(cues, model)
		case GRAPH:
			graphs = append(graphs, model)
		default:
			panic(fmt.Errorf("%w: %T", ErrBadArchetype, model))
		}
	}
}

// Archetype allows to differentiate between different models.
//
// Leviathan supports tables, queues and graphs.
//
// - TABLE is a postgres table, How supports SQL migrations;
// - QUEUE is a pretty API for a postgres channel-based queue;
// - GRAPH is a quad-based graph implementation, uses two tables;
//
type Archetype int

const (
	TABLE Archetype = iota
	QUEUE
	GRAPH
)

// Migrant represents a tangible data model.
//
// Currently, Leviathan supports TABLE, QUEUE, and GRAPH.
type Model interface {
	Type() Archetype
}
