package test

import (
	"crypto/md5"
	"fmt"
	"testing"

	c "consul/consul"
)

// go test -v test/tags_test.go -test.run TestModifyServiceTag
func TestModifyServiceTag(t *testing.T) {
	id := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("p2p%s||%s:%d", "262", "139.196.77.26", 9082))))
	fmt.Printf("id: %s\n", id)
	c.ModifyServiceTagByID(id, "7")
}

// go test -v test/tags_test.go -test.run TestFilterByTag
func TestFilterByTag(t *testing.T) {
	c.FilterTag("p2pserver", "6")
}

// go test -v test/tags_test.go -test.run TestListTags
func TestListTags(t *testing.T) {
	c.ListTagsByName("p2pserver")
}
