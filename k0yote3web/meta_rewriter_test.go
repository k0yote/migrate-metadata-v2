package k0yote3web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRewriteMeta(t *testing.T) {
	var (
		ipfsImageBaseURL = "https://cloudflare-ipfs.com/ipfs/QmNTE1Uvz8JqkL81sJ6UMxw2fwfhtfU8LBtA9SY8t1yD7E/"
	)

	helper := newRewriteHelper(ipfsImageBaseURL)
	err := helper.rewrite()
	assert.NoError(t, err)
}
