package transformation

import (
	"fmt"
	"testing"

	"github.com/flosch/pongo2/v6"
	"github.com/stretchr/testify/require"
)

func TestCompile(t *testing.T) {
	template := `{{hello}} {{world}}`
	context := pongo2.Context{
		"hello": "Hello",
		"world": "World",
	}
	compilation, err := ExecuteSqlTemplate(template, context)
	fmt.Println(compilation)
	require.NoError(t, err)
}
