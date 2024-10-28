# Test tags

```
go test -v test/tags_test.go -test.run TestModifyServiceTag

go test -v test/tags_test.go -test.run TestFilterByTag
```

# Test KV

```
go test -v test/kv_test.go -test.run TestPutKV

go test -v test/kv_test.go -test.run TestDeleteKV

go test -v test/kv_test.go -test.run TestWatchKV
```