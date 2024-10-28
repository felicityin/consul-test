package consul

import (
	"context"
	"fmt"
	"time"

	consul "github.com/hashicorp/consul/api"
)

// DefaultWaitTime is the maximum wait time allowed by Consul
const DefaultWaitTime = 10 * time.Minute

const SleepTime = 5 * time.Second

// Watcher is a wrapper around the Consul client that watches for changes to a keys and directories
type Watcher struct {
	consul *consul.Client
}

// New returns a new Watcher
func New(consulClient *consul.Client) *Watcher {
	return &Watcher{
		consul: consulClient,
	}
}

// WatchKey watches for changes to a key and emits a key value pair
func (w *Watcher) WatchKey(ctx context.Context, key string) (<-chan *consul.KVPair, error) {
	out := make(chan *consul.KVPair)
	kv := w.consul.KV()
	var waitIndex uint64

	opts := &consul.QueryOptions{
		AllowStale:        true,
		RequireConsistent: false,
		UseCache:          true,
		WaitIndex:         waitIndex,
		WaitTime:          DefaultWaitTime,
	}

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			kvPair, meta, err := kv.Get(key, opts.WithContext(ctx))
			if kvPair != nil {
				fmt.Printf("k: %s, v: %v\n", kvPair.Key, kvPair.Value)
			}
			if err != nil {
				if consul.IsRetryableError(err) {
					waitIndex = 0
				}
				continue
			}

			// if we have the same index, then we didn't find any new values
			if waitIndex == meta.LastIndex {
				time.Sleep(SleepTime)
				continue
			}

			waitIndex = meta.LastIndex

			if kvPair != nil {
				out <- kvPair
			}
		}
	}()

	return out, nil
}
