package k0yote3web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRewriteMeta(t *testing.T) {
	var (
		ipfsImageBaseURL = ""
	)

	helper, err := newMetaRewriter(ipfsImageBaseURL, "", "")
	assert.NoError(t, err)
	assert.NoError(t, helper.rewrite())
}
