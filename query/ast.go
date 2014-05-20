package query

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/msiebuhr/MetricBase/metrics"
	"github.com/msiebuhr/MetricBase/query/graphiteParser"
)

// astNode for handling part of the query
type astNode interface {
	String() string
	GetNode() *graphiteParser.Node
	//Explain(MetricBase.Backend) string
	Query(Request) ([]Response, error)
}

// SourceString represents a constant string
type SourceString struct {
	value string
	node  *graphiteParser.Node
}

func (s SourceString) String() string                { return s.value }
func (s SourceString) GetNode() *graphiteParser.Node { return s.node }
func (s SourceString) Query(q Request) ([]Response, error) {
	return nil, fmt.Errorf("Query: Internal error turning string into metrics.")
}

// SourceNumber represents a constant floating point number
type SourceNumber struct {
	value float64
	node  *graphiteParser.Node
}

func NewSourceNumber(n *graphiteParser.Node) (*SourceNumber, error) {
	// Parse it as a number
	i, err := strconv.ParseFloat(n.Name, 64)

	// TODO: Pretty-print something about non-number @ such and such
	if err != nil {
		return nil, err
	}

	return &SourceNumber{value: i, node: n}, nil
}

func (s SourceNumber) String() string                { return fmt.Sprintf("%v", s.value) }
func (s SourceNumber) GetNode() *graphiteParser.Node { return s.node }
func (m SourceNumber) Query(q Request) ([]Response, error) {
	res := NewResponse()
	res.Meta["name"] = fmt.Sprintf("%v", m.value)
	go func() {
		res.Data <- metrics.MetricValue{
			Time:  q.From,
			Value: m.value,
		}
		res.Data <- metrics.MetricValue{
			Time:  q.To,
			Value: m.value,
		}
		close(res.Data)
	}()

	return []Response{res}, nil
}

// SourceMetric is a regex-like expanded lists of metrics
type SourceMetric struct {
	metrics string
	node    *graphiteParser.Node
}

func (m SourceMetric) String() string                { return m.metrics }
func (m SourceMetric) GetNode() *graphiteParser.Node { return m.node }
func (m SourceMetric) Query(q Request) ([]Response, error) {
	responses := make([]Response, 0)
	// Loop over the known list of metrics
	metricsChan := make(chan string, 10)
	q.Backend.GetMetricsList(metricsChan)

	for metricName := range metricsChan {
		// TODO: Do some regex-like things.
		// Shell expansion to regex (or go package)
		if metricName != m.metrics {
			continue
		}

		// Build new Response
		res := NewResponse()
		res.Meta["name"] = metricName
		q.Backend.GetRawData(metricName, q.From, q.To, res.Data)
		responses = append(responses, res)
	}

	return responses, nil
}

// FunctionScale is a function that scales a given set of metrics with a given set of numbers
type FunctionScale struct {
	scale float64
	args  []astNode
	node  *graphiteParser.Node
}

func (f FunctionScale) String() string {
	args := make([]string, len(f.args))
	for i := range f.args {
		args[i] = fmt.Sprintf("%v", f.args[i])
	}
	return fmt.Sprintf("scale(%v, %v)", strings.Join(args, ", "), f.scale)
}
func (f FunctionScale) GetNode() *graphiteParser.Node { return f.node }
func (f FunctionScale) Query(q Request) ([]Response, error) {
	// Execute all arguments and execute those
	queries := make([]Response, 0)
	for i := range f.args {
		// Execute Query on each of them and start a go-routine to change the data
		results, err := f.args[i].Query(q)

		// Clean up on errors and report upstream
		if err != nil {
			// Close all result channels - or sink them?

			// Return the error upstream?
			return nil, err
		}

		// Pass on results
		for _, result := range results {
			queries = append(queries, result)
		}
	}

	// Start a go-routine that processes the data coming through each chan
	for i := range queries {
		// Switch the channels around
		srcChan := queries[i].Data
		queries[i].Data = make(chan metrics.MetricValue, cap(srcChan))

		go func(i int) {
			// Read all data from old channel, scale it and write to the new one
			for data := range srcChan {
				data.Value = data.Value * f.scale
				queries[i].Data <- data
			}
			close(queries[i].Data)
		}(i)
	}

	return queries, nil
}

func NewFunctionScale(args []astNode, n *graphiteParser.Node) (astNode, error) {
	f := FunctionScale{
		args:  make([]astNode, 0),
		scale: 1,
	}

	scaleSet := false

	for i := range args {
		switch arg := args[i].(type) {
		case *SourceString:
			return nil, errors.New("Parse error: scale() does not accept strings.")
		case *SourceNumber:
			if scaleSet {
				return nil, errors.New("Parse error: scale() only accepts one number")
			}
			scaleSet = true
			f.scale = arg.value
		default:
			f.args = append(f.args, args[i])
		}
	}

	return f, nil
}

func LookupFunction(name string, args []astNode, n *graphiteParser.Node) (astNode, error) {
	switch name {
	case "scale":
		return NewFunctionScale(args, n)
	default:
		return nil, errors.New("Unknown function '" + name + "'.")
	}
}

// Recursively convert a parsed query to a tree of something executable
func convertNodeToAST(n *graphiteParser.Node) (astNode, error) {
	switch n.Type {
	case graphiteParser.NODE_STRING:
		return &SourceString{value: n.Name, node: n}, nil
	case graphiteParser.NODE_NUMBER:
		return NewSourceNumber(n)
	case graphiteParser.NODE_METRIC:
		return &SourceMetric{metrics: n.Name, node: n}, nil
	case graphiteParser.NODE_FUNCTION:
		// Convert all arguments
		args := make([]astNode, len(n.Args))

		for i := range n.Args {
			arg, err := convertNodeToAST(n.Args[i])
			if err != nil {
				return nil, err
			}
			args[i] = arg
		}

		// Build function
		return LookupFunction(n.Name, args, n)
	default:
		return nil, errors.New(fmt.Sprintf("Did not understand input at char %d (%v)", n.Char, n))
	}
}

func ParseGraphiteQuery(q string) (astNode, error) {
	node, err := graphiteParser.Parse(q)
	if err != nil {
		return nil, err
	}

	// Walk the returned tree to get something useful out
	return convertNodeToAST(node)
}
