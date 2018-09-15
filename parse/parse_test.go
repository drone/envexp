package parse

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var tests = []struct {
	Text string
	Node Node
}{

	//
	// text only
	//
	{
		Text: "text",
		Node: &TextNode{Value: "text"},
	},
	{
		Text: "}text",
		Node: &TextNode{Value: "}text"},
	},
	{
		Text: "http://github.com",
		Node: &TextNode{Value: "http://github.com"}, // should not escape double slash
	},
	{
		Text: "$${string}",
		Node: &TextNode{Value: "${string}"}, // should not escape double dollar
	},
	{
		Text: "$$string",
		Node: &TextNode{Value: "$string"}, // should not escape double dollar
	},

	//
	// variable only
	//
	{
		Text: "${string}",
		Node: &FuncNode{Param: "string"},
	},
	{
		Text: "${string_A}",
		Node: &FuncNode{Param: "string_A"},
	},
	{
		Text: "${string_A-b}",
		Node: &FuncNode{Param: "string_A-b"},
	},

	//
	// text transform functions
	//
	{
		Text: "${string,}",
		Node: &FuncNode{
			Param: "string",
			Name:  ",",
			Args:  nil,
		},
	},
	{
		Text: "${string,,}",
		Node: &FuncNode{
			Param: "string",
			Name:  ",,",
			Args:  nil,
		},
	},
	{
		Text: "${string^}",
		Node: &FuncNode{
			Param: "string",
			Name:  "^",
			Args:  nil,
		},
	},
	{
		Text: "${string^^}",
		Node: &FuncNode{
			Param: "string",
			Name:  "^^",
			Args:  nil,
		},
	},

	//
	// substring functions
	//
	{
		Text: "${string:position}",
		Node: &FuncNode{
			Param: "string",
			Name:  ":",
			Args: []Node{
				&TextNode{Value: "position"},
			},
		},
	},
	{
		Text: "${string:position:length}",
		Node: &FuncNode{
			Param: "string",
			Name:  ":",
			Args: []Node{
				&TextNode{Value: "position"},
				&TextNode{Value: "length"},
			},
		},
	},

	//
	// string removal functions
	//
	{
		Text: "${string#substring}",
		Node: &FuncNode{
			Param: "string",
			Name:  "#",
			Args: []Node{
				&TextNode{Value: "substring"},
			},
		},
	},
	{
		Text: "${string##substring}",
		Node: &FuncNode{
			Param: "string",
			Name:  "##",
			Args: []Node{
				&TextNode{Value: "substring"},
			},
		},
	},
	{
		Text: "${string%substring}",
		Node: &FuncNode{
			Param: "string",
			Name:  "%",
			Args: []Node{
				&TextNode{Value: "substring"},
			},
		},
	},
	{
		Text: "${string%%substring}",
		Node: &FuncNode{
			Param: "string",
			Name:  "%%",
			Args: []Node{
				&TextNode{Value: "substring"},
			},
		},
	},

	//
	// string replace functions
	//
	{
		Text: "${string/substring/replacement}",
		Node: &FuncNode{
			Param: "string",
			Name:  "/",
			Args: []Node{
				&TextNode{Value: "substring"},
				&TextNode{Value: "replacement"},
			},
		},
	},
	{
		Text: "${string//substring/replacement}",
		Node: &FuncNode{
			Param: "string",
			Name:  "//",
			Args: []Node{
				&TextNode{Value: "substring"},
				&TextNode{Value: "replacement"},
			},
		},
	},
	{
		Text: "${string/#substring/replacement}",
		Node: &FuncNode{
			Param: "string",
			Name:  "/#",
			Args: []Node{
				&TextNode{Value: "substring"},
				&TextNode{Value: "replacement"},
			},
		},
	},
	{
		Text: "${string/%substring/replacement}",
		Node: &FuncNode{
			Param: "string",
			Name:  "/%",
			Args: []Node{
				&TextNode{Value: "substring"},
				&TextNode{Value: "replacement"},
			},
		},
	},

	//
	// default value functions
	//
	{
		Text: "${string=default}",
		Node: &FuncNode{
			Param: "string",
			Name:  "=",
			Args: []Node{
				&TextNode{Value: "default"},
			},
		},
	},
	{
		Text: "${string:=default}",
		Node: &FuncNode{
			Param: "string",
			Name:  ":=",
			Args: []Node{
				&TextNode{Value: "default"},
			},
		},
	},
	{
		Text: "${string:-default}",
		Node: &FuncNode{
			Param: "string",
			Name:  ":-",
			Args: []Node{
				&TextNode{Value: "default"},
			},
		},
	},
	{
		Text: "${string:?default}",
		Node: &FuncNode{
			Param: "string",
			Name:  ":?",
			Args: []Node{
				&TextNode{Value: "default"},
			},
		},
	},
	{
		Text: "${string:+default}",
		Node: &FuncNode{
			Param: "string",
			Name:  ":+",
			Args: []Node{
				&TextNode{Value: "default"},
			},
		},
	},

	//
	// length function
	//
	{
		Text: "${#string}",
		Node: &FuncNode{
			Param: "string",
			Name:  "#",
		},
	},

	//
	// special characters in argument
	//
	{
		Text: "${string#$%:*{}",
		Node: &FuncNode{
			Param: "string",
			Name:  "#",
			Args: []Node{
				&TextNode{Value: "$%:*{"},
			},
		},
	},

	// text before and after function
	{
		Text: "hello ${#string} world",
		Node: &ListNode{
			Nodes: []Node{
				&TextNode{
					Value: "hello ",
				},
				&ListNode{
					Nodes: []Node{
						&FuncNode{
							Param: "string",
							Name:  "#",
						},
						&TextNode{
							Value: " world",
						},
					},
				},
			},
		},
	},

	// escaped function arguments
	{
		Text: `${string/\/position/length}`,
		Node: &FuncNode{
			Param: "string",
			Name:  "/",
			Args: []Node{
				&TextNode{
					Value: "/position",
				},
				&TextNode{
					Value: "length",
				},
			},
		},
	},
	{
		Text: `${string/\/position\\/length}`,
		Node: &FuncNode{
			Param: "string",
			Name:  "/",
			Args: []Node{
				&TextNode{
					Value: "/position\\",
				},
				&TextNode{
					Value: "length",
				},
			},
		},
	},
	{
		Text: `${string/position/\/length}`,
		Node: &FuncNode{
			Param: "string",
			Name:  "/",
			Args: []Node{
				&TextNode{
					Value: "position",
				},
				&TextNode{
					Value: "/length",
				},
			},
		},
	},
	{
		Text: `${string/position/\/length\\}`,
		Node: &FuncNode{
			Param: "string",
			Name:  "/",
			Args: []Node{
				&TextNode{
					Value: "position",
				},
				&TextNode{
					Value: "/length\\",
				},
			},
		},
	},

	// functions in functions
	{
		Text: "${string:${position}}",
		Node: &FuncNode{
			Param: "string",
			Name:  ":",
			Args: []Node{
				&FuncNode{
					Param: "position",
				},
			},
		},
	},
	{
		Text: "${string:${stringy:position:length}:${stringz,,}}",
		Node: &FuncNode{
			Param: "string",
			Name:  ":",
			Args: []Node{
				&FuncNode{
					Param: "stringy",
					Name:  ":",
					Args: []Node{
						&TextNode{Value: "position"},
						&TextNode{Value: "length"},
					},
				},
				&FuncNode{
					Param: "stringz",
					Name:  ",,",
				},
			},
		},
	},
	{
		Text: "${string#${stringz}}",
		Node: &FuncNode{
			Param: "string",
			Name:  "#",
			Args: []Node{
				&FuncNode{Param: "stringz"},
			},
		},
	},
	{
		Text: "${string=${stringz}}",
		Node: &FuncNode{
			Param: "string",
			Name:  "=",
			Args: []Node{
				&FuncNode{Param: "stringz"},
			},
		},
	},
	{
		Text: "${string//${stringy}/${stringz}}",
		Node: &FuncNode{
			Param: "string",
			Name:  "//",
			Args: []Node{
				&FuncNode{Param: "stringy"},
				&FuncNode{Param: "stringz"},
			},
		},
	},
}

func TestParse(t *testing.T) {
	for _, test := range tests {
		got, err := Parse(test.Text)
		if err != nil {
			t.Error(err)
		}

		if diff := cmp.Diff(test.Node, got.Root); diff != "" {
			t.Errorf(diff)
		}
	}
}
