package config

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_NewConfigTest(t *testing.T) {
	c := NewConfigTest()

	tokens := strings.Split(c.ProjectRoot, "/")
	projectRoot := tokens[len(tokens)-1]

	assert.Equal(t, "scrapygo", projectRoot)
}
