package levi

import (
	"crypto/rand"
	"io"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/graph/path"
	"github.com/cayleygraph/cayley/schema"
	"github.com/cayleygraph/cayley/voc/rdf"
	"github.com/cayleygraph/quad"
)

type Graph struct {
	graph.Handle
	sch *schema.Config
}

// ID generates a random quad.IRI for a given datatype,
// or constructs an ID.
//
// Assuming "tas:" is the designated namespace.
//
//   ID(tas.Knol)           	  -> "tas:Knol/u4i0ddazo2hcczjr"
//   ID(tas.Message, "custom_id") -> "tas:Message/custom_id"
//
func ID(datatype string, assumed ...string) quad.IRI {
	if len(assumed) == 1 {
		return quad.IRI(datatype + "/" + assumed[0])
	}

	const (
		charset = `abcdefghijklmnopqrstuvwxyz0123456789`
		space   = byte(len(charset))
		length  = 16
	)

	id := make([]byte, length)
	n, err := io.ReadFull(rand.Reader, id)
	if n != len(id) || err != nil {
		panic(err)
	}

	for i := 0; i < length; i++ {
		id[i] = charset[id[i]%space]
	}

	return quad.IRI(datatype + "/" + string(id))
}

func (g *Graph) Schema() *schema.Config {
	return g.sch
}

func (g *Graph) V() *path.Path {
	return path.NewPath(g)
}

func (g *Graph) Select(iri string) *path.Path {
	return path.StartPath(g, quad.IRI(iri)).In(quad.IRI(rdf.Type))
}

func (g *Graph) Load(p *path.Path, dst interface{}) error {
	return g.sch.LoadPathTo(nil, g, dst, p)
}

func (g *Graph) Get(dst interface{}, depth int, ids ...quad.Value) error {
	return g.sch.LoadToDepth(nil, g, dst, depth, ids...)
}

func (g *Graph) Put(it interface{}) (quad.Value, error) {
	tx := graph.NewTransaction()

	result, err := g.PutInto(tx, it)
	if err != nil {
		return nil, err
	}

	err = g.ApplyTransaction(tx)
	if err != nil {
		return nil, err
	}

	return result[0], nil
}

func (g *Graph) PutInto(tx *graph.Transaction, values ...interface{}) ([]quad.Value, error) {
	qw := graph.NewTxWriter(tx, graph.Add)

	result := make([]quad.Value, len(values))

	for i, value := range values {
		id, err := g.sch.WriteAsQuads(qw, value)
		if err != nil {
			return nil, err
		}

		result[i] = id
	}

	return result, nil
}

func (g *Graph) Delete(it interface{}) error {
	tx := graph.NewTransaction()

	if err := g.DeleteFrom(tx, it); err != nil {
		return err
	}

	err := g.ApplyTransaction(tx)
	return err
}

func (g *Graph) DeleteFrom(tx *graph.Transaction, values ...interface{}) error {
	qw := graph.NewTxWriter(tx, graph.Delete)

	for _, value := range values {
		_, err := g.sch.WriteAsQuads(qw, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func newSchema() *schema.Config {
	sch := schema.NewConfig()
	sch.GenerateID = func(unknown interface{}) quad.Value {
		return quad.IRI(ID("_"))
	}

	return sch
}

type nologger struct{}

func (nologger) V(int) bool                                  { return false }
func (nologger) Infof(string, ...interface{})                {}
func (nologger) Warningf(format string, args ...interface{}) {}
func (nologger) Errorf(format string, args ...interface{})   {}
func (nologger) Fatalf(format string, args ...interface{})   {}
func (nologger) SetV(int)                                    {}
