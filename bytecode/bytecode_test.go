package graphql

import (
	"bytes"
	"strings"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func parseQuery(query string) ([]byte, []error) {
	i := parserCtx{
		res:    []byte{},
		query:  []byte(query),
		errors: []error{},
	}
	i.parseQueryToBytecode()
	return i.res, i.errors
}

func formatHumanReadableQuery(result string) string {
	result = strings.TrimSpace(result)
	result = strings.ReplaceAll(result, " ", "")
	result = strings.ReplaceAll(result, "\t", "")
	result = strings.ReplaceAll(result, "\r", "")
	lines := strings.Split(result, "\n")
	for i, line := range lines {
		lines[i] = strings.Split(line, "//")[0]
	}
	return strings.Join(lines, "\n")
}

func formatResToHumandReadable(result []byte) string {
	result = bytes.Join(bytes.Split(result, []byte{0}), []byte{'\n'})
	return strings.TrimSpace(string(result))
}

func parseQueryAndExpectResult(t *testing.T, query, expectedResult string) {
	res, errs := parseQuery(query)
	for _, err := range errs {
		panic(err.Error())
	}
	Equal(t, formatHumanReadableQuery(expectedResult), formatResToHumandReadable(res))
}

func TestParseSimpleQuery(t *testing.T) {
	parseQueryAndExpectResult(t, `{}`, `
		oq // [operator] [query]
		e  // [end of operator]
	`)
}

func TestParseSimpleQueryWrittenOut(t *testing.T) {
	parseQueryAndExpectResult(t, `query {}`, `
		oq // operator of type query
		e  // end of operator
	`)
}

func TestParseSimpleMutation(t *testing.T) {
	parseQueryAndExpectResult(t, `mutation {}`, `
		om // operator of type mutation
		e  // end of operator
	`)
}

func TestParseSimpleSubscription(t *testing.T) {
	parseQueryAndExpectResult(t, `subscription {}`, `
		os // operator of type subscription
		e  // end of operator
	`)
}

func TestParseQueryWithName(t *testing.T) {
	parseQueryAndExpectResult(t, `query banana {}`, `
		oqbanana // operator of type query with name banana
		e        // end of operator
	`)
}

func TestParseMultipleSimpleQueries(t *testing.T) {
	parseQueryAndExpectResult(t, `{}{}`, `
		oq // operator 1
		e  // end of operator 1
		oq // operator 2
		e  // end of operator 2
	`)
}

func TestParseMultipleQueries(t *testing.T) {
	parseQueryAndExpectResult(t, `
		query a {}
		mutation b {}
	`, `
		oqa // query operator 1
		e   // end of operator 1
		omb // mutation operator 2
		e   // end of operator 2
	`)
}

func TestParseQueryWithField(t *testing.T) {
	parseQueryAndExpectResult(t, `query {
		some_field
	}`, `
		oq          // query operator
		fsome_field // field with name some_field
		e           // end field
		e           // end operator
	`)
}

func TestParseQueryWithMultipleFields(t *testing.T) {
	parseQueryAndExpectResult(t, `query {
		some_field
		other
	}`, `
		oq          // query operator
		fsome_field // field with name some_field
		e           // end field with name some_field
		fother      // field with name other
		e           // end field with name other
		e           // end operator
	`)
}

func TestParseQueryWithFieldWithSelectionSet(t *testing.T) {
	parseQueryAndExpectResult(t, `query {
		some_field {
			foo
			bar
		}
	}`, `
		oq          // query operator
		fsome_field // field with name some_field
		ffoo        // field with name foo
		e           // end of foo
		fbar        // field with name bar
		e           // end of bar
		e           // end of some_field
		e           // end operator
	`)
}

func TestParseQueryWithFieldWithFragmentSpread(t *testing.T) {
	parseQueryAndExpectResult(t, `query {
		some_field {
			foo
			... baz
			bar
		}
	}`, `
		oq
		fsome_field
		ffoo
		e
		sfbaz // fragment spread pointing to fragment with name baz
		fbar
		e
		e
		e
	`)

	// A query that starts with "on" should parse as a fragment pointer
	parseQueryAndExpectResult(t, `query {
		some_field {
			foo
			... online
			bar
		}
	}`, `
		oq
		fsome_field
		ffoo
		e
		sfonline // fragment spread pointing to fragment with name online
		fbar
		e
		e
		e
	`)
}

func TestParseQueryWithFieldWithInlineFragmentSpread(t *testing.T) {
	parseQueryAndExpectResult(t, `query {
		some_field {
			foo
			... on baz {
				bazField
			}
			bar
		}
	}`, `
		oq
		fsome_field
		ffoo
		e
		stbaz     // fragment spread with typename baz
		fbazField // fragment field
		e         // end of fragment field
		e         // end of inline fragment
		fbar
		e
		e
		e
	`)
}

func TestParseFragment(t *testing.T) {
	parseQueryAndExpectResult(t, `fragment Foo on Bar {}`, `
		FFoo // fragment with name Foo
		Bar  // fragment type name
		e    // end of fragment
	`)
}

func TestParseFragmentWithFields(t *testing.T) {
	parseQueryAndExpectResult(t, `fragment Foo on Bar {
		fieldA
		bField
	}`, `
		FFoo    // fragment with name Foo
		Bar     // fragment type name
		ffieldA
		e
		fbField
		e
		e       // end of fragment
	`)
}
