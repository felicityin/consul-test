package test

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"

	c "consul/consul"
)

// go test -v test/kv_test.go -test.run TestPutKV
func TestPutKV(t *testing.T) {
	kv := c.GetClient(c.ConsulAddr, c.Scheme, c.Token).KV()

	key := "test-key"
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
	key := "test-key"
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

// go test -v test/kv_test.go -test.run TestWatchKV
func TestWatchKV(t *testing.T) {
	key := "test-key"
	client := c.GetClient(c.ConsulAddr, c.Scheme, c.KvToken)
	kvWatcher := c.New(client)

	ch, _ := kvWatcher.WatchKey(context.Background(), key)

	for kv := range ch {
		if kv == nil {
			t.Errorf("err: kv is nil")
		} else {
			i, err := strconv.Atoi(string(kv.Value))
			if err != nil {
				fmt.Printf("Error %v\n", err)
				return
			}
			t.Logf("new value! k: %s, v: %d\n", kv.Key, i)
		}
	}
}
