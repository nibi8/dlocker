# Distributed locker

Embed distributed locker (distributed mutex).
You need to implement persistent lock storage in order to use it (or use the provided mongodb storage provider).
With Extender you can extend lock several or infinity times.

For critical data updates (during lock), it is recommended to additionally use data versioning.

For more details see examples.
