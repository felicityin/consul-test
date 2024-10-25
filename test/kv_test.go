package test

import (
	"bytes"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"

	c "consul/consul"
)

// go test -v test/kv_test.go -test.run TestPutKV
func TestPutKV(t *testing.T) {
	kv := c.GetClient(c.ConsulAddr, c.Scheme, c.Token).KV()

	key := "Qme8rV3sDXki9HfvK6zYRRvdfiC2ZP3Ybw8zGf6qTdaBL5"
	value := []byte("7")

	// Put the key
	p := &api.KVPair{Key: key, Flags: 42, Value: value}
	if _, err := kv.Put(p, &api.WriteOptions{Token: c.KvToken}); err != nil {
		t.Fatalf("put kv err: %s, key: %s", err.Error(), key)
	}

	time.Sleep(1)

	// Get should work
	pair, meta, err := kv.Get(key, &api.QueryOptions{Token: c.KvToken})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if pair == nil {
		t.Fatalf("expected value: %#v", pair)
	}
	if !bytes.Equal(pair.Value, value) {
		t.Fatalf("unexpected value: %#v", pair)
	}
	if pair.Flags != 42 {
		t.Fatalf("unexpected value: %#v", pair)
	}
	if meta.LastIndex == 0 {
		t.Fatalf("unexpected value: %#v", meta)
	}
}

// go test -v test/kv_test.go -test.run TestDeleteKV
func TestDeleteKV(t *testing.T) {
	kv := c.GetClient(c.ConsulAddr, c.Scheme, c.Token).KV()

	// Get a get without a key
	key := "Qme8rV3sDXki9HfvK6zYRRvdfiC2ZP3Ybw8zGf6qTdaBL5"
	_, pair, err := kv.Get(key, &api.QueryOptions{Token: c.KvToken})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if pair == nil {
		t.Fatalf("unexpected value: %#v", pair)
	}

	// Delete
	if _, err := kv.Delete(key, &api.WriteOptions{Token: c.KvToken}); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get should fail
	_, pair, err = kv.Get(key, &api.QueryOptions{Token: c.KvToken})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if pair != nil {
		t.Fatalf("unexpected value: %#v", pair)
	}
}
