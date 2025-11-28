package binding

import (
	"testing"

	"github.com/rusmanplatd/goravelframework/contracts/binding"
	"github.com/stretchr/testify/assert"
)

func TestDependencies(t *testing.T) {
	dependencies := Dependencies([]string{
		binding.Orm,
		binding.DB,
		binding.Schema,
		binding.Seeder,
	}...)

	assert.ElementsMatch(t, []string{
		binding.Config,
		binding.Log,
	}, dependencies)
}
