// -*-go-*-

package graphiteParser

import (
	"errors"
	"fmt"
)

%% machine rpn;
%% write data;

func Parse(data string) (res *Node, err error) {
	//fmt.Printf("Parsing '%s'\n", data)
	data = data + " " // The generated state-machine works better with some trailing junk
	cs, p, pe := 0, 0, len(data)
	ts, te := 0, 0
	act, eof := 0, 0

	// Make an empty root function
	rootTree := NewNode("", 0, NODE_ROOT)

	%%{
		action number {
			rootTree.addArgument(NewNode(data[ts:te], ts, NODE_NUMBER))
		}

		action metric {
			rootTree.addArgument(NewNode(data[ts:te], ts, NODE_METRIC))
		}

		action function_start {
			// Create new tree, add it as an argument to the parent and declare
			// it the new root.
			newRootTree := NewNode(data[ts:te-1], ts, NODE_FUNCTION)
			newRootTree.parent = rootTree
			rootTree.addArgument(newRootTree)
			rootTree = newRootTree
		}

		action function_end {
			rootTree = rootTree.parent
		}

		action quoted_string {
			rootTree.addArgument(NewNode(data[ts:te], ts, NODE_STRING))
		}

		metricName = [a-z.*]+;
		digits = '-'? [0-9.]+;

		main := |*
			metricName [] => metric;
			digits [] => number;
			alnum+ '(' => function_start;
			'"' alnum+ '"' => quoted_string;
			')' => function_end;
			',' => { /*fmt.Println(",")*/ };
			space;
		*|;

		write init;
		write exec;
	}%%

	_ = act;
	_ = eof;

	if cs < rpn_first_final {
		if p == pe {
			return nil, errors.New("unexpected eof")
		} else {
			return nil, errors.New(fmt.Sprintf("error at position %d", p))
		}
	}

	// The root node may only have one sibling element
	if rootTree.Type != NODE_ROOT || len(rootTree.Args) != 1 {
		return nil, errors.New("Only supposed to have one root element")
	}

	return rootTree.Args[0], nil
}
