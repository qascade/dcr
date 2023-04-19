// This package will contain transformation types
package transformation

import (
	"github.com/flosch/pongo2/v6"
)

type Transformation interface{}

type MLModel interface {
	Transformation
}

// All Supported Differentially Private Queries must comply to this interface.
// We may also define some DPQuery templates.
// But we also need to give support for general DP Query.
type SqlQuery interface {
	Transformation
	Compile()
}

type NonPrivateQuery interface {
	SqlQuery
}
type PrivateQuery interface {
	SqlQuery
}

type GenericPrivateQuery struct {
	query    string
	template pongo2.Template
	inputs   map[string]interface{}
}
