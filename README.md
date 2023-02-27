Ethereum Watcher is a simple application that watch ethereum blockchain for subscribed address. Whenever there's incoming/outgoing transactions of subscribed addresses from latest block, it will store and do notifying. It also acts as http server for user to subscribes and fetch transactions.

Current limitation is the implementation of storage is using memory storage and the notifying is logging the transaction. It doesn't use any external libraries other than from golang native.

# How to run

## Natively
1. 
```
go build -o bin/ethereum-watcher
```
2. 
```
./bin/ethereum-watcher
```

## Docker
1.
```
docker build --tag ethereum-watcher:latest
```

2.
```
docker run --rm ethereum-watcher
```
