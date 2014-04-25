// line 1 "query/graphiteParser/parser.rl"
// -*-go-*-
package graphiteParser

import (
	"errors"
	"fmt"
)

// line 10 "query/graphiteParser/parser.rl"

// line 15 "query/graphiteParser/parser.go"
var _rpn_actions []byte = []byte{
	0, 1, 2, 1, 6, 1, 7, 1, 8,
	1, 9, 1, 10, 1, 11, 1, 12,
	1, 13, 2, 0, 1, 2, 3, 4,
	2, 3, 5,
}

var _rpn_key_offsets []byte = []byte{
	0, 0, 6, 13, 16, 23, 38, 42,
	45, 51, 59,
}

var _rpn_trans_keys []byte = []byte{
	48, 57, 65, 90, 97, 122, 34, 48,
	57, 65, 90, 97, 122, 46, 48, 57,
	40, 48, 57, 65, 90, 97, 122, 32,
	34, 41, 42, 44, 45, 46, 9, 13,
	48, 57, 65, 90, 97, 122, 42, 46,
	97, 122, 46, 48, 57, 42, 46, 48,
	57, 97, 122, 40, 46, 48, 57, 65,
	90, 97, 122, 40, 42, 46, 48, 57,
	65, 90, 97, 122,
}

var _rpn_single_lengths []byte = []byte{
	0, 0, 1, 1, 1, 7, 2, 1,
	2, 2, 3,
}

var _rpn_range_lengths []byte = []byte{
	0, 3, 3, 1, 3, 4, 1, 1,
	2, 3, 3,
}

var _rpn_index_offsets []byte = []byte{
	0, 0, 4, 9, 12, 17, 29, 33,
	36, 41, 47,
}

var _rpn_indicies []byte = []byte{
	0, 0, 0, 1, 2, 0, 0, 0,
	1, 3, 3, 1, 5, 6, 6, 6,
	4, 7, 8, 9, 10, 11, 12, 13,
	7, 14, 6, 15, 1, 10, 10, 10,
	16, 3, 3, 17, 10, 13, 3, 10,
	16, 5, 3, 14, 6, 6, 17, 5,
	10, 10, 6, 6, 15, 16,
}

var _rpn_trans_targs []byte = []byte{
	2, 0, 5, 7, 5, 5, 4, 5,
	1, 5, 6, 5, 3, 8, 9, 10,
	5, 5,
}

var _rpn_trans_actions []byte = []byte{
	0, 0, 5, 0, 17, 3, 0, 11,
	0, 7, 0, 9, 0, 0, 25, 22,
	13, 15,
}

var _rpn_to_state_actions []byte = []byte{
	0, 0, 0, 0, 0, 19, 0, 0,
	0, 0, 0,
}

var _rpn_from_state_actions []byte = []byte{
	0, 0, 0, 0, 0, 1, 0, 0,
	0, 0, 0,
}

var _rpn_eof_trans []byte = []byte{
	0, 0, 0, 0, 5, 0, 17, 18,
	17, 18, 17,
}

const rpn_start int = 5
const rpn_first_final int = 5
const rpn_error int = 0

const rpn_en_main int = 5

// line 11 "query/graphiteParser/parser.rl"

func Parse(data string) (res *Node, err error) {
	//fmt.Printf("Parsing '%s'\n", data)
	data = data + " " // The generated state-machine works better with some trailing junk
	cs, p, pe := 0, 0, len(data)
	ts, te := 0, 0
	act, eof := 0, 0

	// Make an empty root function
	rootTree := NewNode("", 0, NODE_ROOT)

	// line 112 "query/graphiteParser/parser.go"
	{
		cs = rpn_start
		ts = 0
		te = 0
		act = 0
	}

	// line 120 "query/graphiteParser/parser.go"
	{
		var _klen int
		var _trans int
		var _acts int
		var _nacts uint
		var _keys int
		if p == pe {
			goto _test_eof
		}
		if cs == 0 {
			goto _out
		}
	_resume:
		_acts = int(_rpn_from_state_actions[cs])
		_nacts = uint(_rpn_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _rpn_actions[_acts-1] {
			case 2:
				// line 1 "NONE"

				ts = p

				// line 144 "query/graphiteParser/parser.go"
			}
		}

		_keys = int(_rpn_key_offsets[cs])
		_trans = int(_rpn_index_offsets[cs])

		_klen = int(_rpn_single_lengths[cs])
		if _klen > 0 {
			_lower := int(_keys)
			var _mid int
			_upper := int(_keys + _klen - 1)
			for {
				if _upper < _lower {
					break
				}

				_mid = _lower + ((_upper - _lower) >> 1)
				switch {
				case data[p] < _rpn_trans_keys[_mid]:
					_upper = _mid - 1
				case data[p] > _rpn_trans_keys[_mid]:
					_lower = _mid + 1
				default:
					_trans += int(_mid - int(_keys))
					goto _match
				}
			}
			_keys += _klen
			_trans += _klen
		}

		_klen = int(_rpn_range_lengths[cs])
		if _klen > 0 {
			_lower := int(_keys)
			var _mid int
			_upper := int(_keys + (_klen << 1) - 2)
			for {
				if _upper < _lower {
					break
				}

				_mid = _lower + (((_upper - _lower) >> 1) & ^1)
				switch {
				case data[p] < _rpn_trans_keys[_mid]:
					_upper = _mid - 2
				case data[p] > _rpn_trans_keys[_mid+1]:
					_lower = _mid + 2
				default:
					_trans += int((_mid - int(_keys)) >> 1)
					goto _match
				}
			}
			_trans += _klen
		}

	_match:
		_trans = int(_rpn_indicies[_trans])
	_eof_trans:
		cs = int(_rpn_trans_targs[_trans])

		if _rpn_trans_actions[_trans] == 0 {
			goto _again
		}

		_acts = int(_rpn_trans_actions[_trans])
		_nacts = uint(_rpn_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _rpn_actions[_acts-1] {
			case 3:
				// line 1 "NONE"

				te = p + 1

			case 4:
				// line 27 "query/graphiteParser/parser.rl"

				act = 1
			case 5:
				// line 23 "query/graphiteParser/parser.rl"

				act = 2
			case 6:
				// line 31 "query/graphiteParser/parser.rl"

				te = p + 1
				{
					// Create new tree, add it as an argument to the parent and declare
					// it the new root.
					newRootTree := NewNode(data[ts:te-1], ts, NODE_FUNCTION)
					newRootTree.parent = rootTree
					rootTree.addArgument(newRootTree)
					rootTree = newRootTree
				}
			case 7:
				// line 44 "query/graphiteParser/parser.rl"

				te = p + 1
				{
					rootTree.addArgument(NewNode(data[ts:te], ts, NODE_STRING))
				}
			case 8:
				// line 40 "query/graphiteParser/parser.rl"

				te = p + 1
				{
					rootTree = rootTree.parent
				}
			case 9:
				// line 57 "query/graphiteParser/parser.rl"

				te = p + 1
				{ /*fmt.Println(",")*/
				}
			case 10:
				// line 58 "query/graphiteParser/parser.rl"

				te = p + 1

			case 11:
				// line 27 "query/graphiteParser/parser.rl"

				te = p
				p--
				{
					rootTree.addArgument(NewNode(data[ts:te], ts, NODE_METRIC))
				}
			case 12:
				// line 23 "query/graphiteParser/parser.rl"

				te = p
				p--
				{
					rootTree.addArgument(NewNode(data[ts:te], ts, NODE_NUMBER))
				}
			case 13:
				// line 1 "NONE"

				switch act {
				case 0:
					{
						cs = 0
						goto _again
					}
				case 1:
					{
						p = (te) - 1

						rootTree.addArgument(NewNode(data[ts:te], ts, NODE_METRIC))
					}
				case 2:
					{
						p = (te) - 1

						rootTree.addArgument(NewNode(data[ts:te], ts, NODE_NUMBER))
					}
				}

				// line 299 "query/graphiteParser/parser.go"
			}
		}

	_again:
		_acts = int(_rpn_to_state_actions[cs])
		_nacts = uint(_rpn_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _rpn_actions[_acts-1] {
			case 0:
				// line 1 "NONE"

				ts = 0

			case 1:
				// line 1 "NONE"

				act = 0

				// line 319 "query/graphiteParser/parser.go"
			}
		}

		if cs == 0 {
			goto _out
		}
		p++
		if p != pe {
			goto _resume
		}
	_test_eof:
		{
		}
		if p == eof {
			if _rpn_eof_trans[cs] > 0 {
				_trans = int(_rpn_eof_trans[cs] - 1)
				goto _eof_trans
			}
		}

	_out:
		{
		}
	}

	// line 63 "query/graphiteParser/parser.rl"

	_ = act
	_ = eof

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
