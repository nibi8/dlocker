# Distributed locker

Tiny distributed locker (distributed mutex).
You need to implement persistent lock storage in order to use it (or use the provided mongodb storage provider).

## Scheme

Lock "unique_lock_name" with 5 min duration:

```
Instace 1: [set lock success] [running within the duration of the lock (5 minutes)]
```

Must complete execution before the lock expires.

```
Instace 2:       [set lock fail] [sleep during lock (5 min)                        ] [try get lock]
```

Execution of other locks is not affected.

## Usage

Create storage provider. Implement `StorageProvider` interface for you persistent storage or use default implementation (mongosp package has a mongodb implementation)

```go
sp, err := mongosp.NewStorageProvider(ctx, db, "lockerTest")
```

Create locker

```go
locker := dlocker.NewLocker(sp)
```

Create lock

```go
lock, err := models.NewLock("unique_lock_name", 10, 5)
...
lockCtx, _, err := locker.LockWithWait(ctx, lock)
```

For more details see examples.
