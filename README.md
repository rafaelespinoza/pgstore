# pgstore

A session store backend for [github.com/gorilla/sessions](https://github.com/gorilla/sessions).

This is a fork of [github.com/antonlindstrom/pgstore](https://github.com/antonlindstrom/pgstore).

## Development

```
make deps
```

Integration tests against a live DB can be run on a machine with `docker` and
`docker-compose` installed. Those files live in the directory `.docker/`.

```
# Build test environment and run tests.
make -f .docker/Makefile test-up

# Teardown test environment.
make -f .docker/Makefile test-down
```

## Documentation

See [full documentation](https://pkg.go.dev/github.com/gorilla/sessions) of the
underlying interface.

See `examples/` for example usage.

## Differences from original project

These are some notable differences compared to the [original
project](https://github.com/antonlindstrom/pgstore). The schema of the
`http_sessions` DB table differs in order to:

- more efficiently lookup a session by key
- ease lookups of expired sessions in the automated Cleanup job

The `id` field was removed from the `http_sessions` table because the library
does not use it. Likewise, the `PGSession.ID` field was removed. The primary
key is now an opaque identifier, `key`, which is already generated in the
original implementation.

## Thanks

I've stolen, borrowed and gotten inspiration from the other backends available:

* the original [pgstore](https://github.com/antonlindstrom/pgstore)
* [redistore](https://github.com/boj/redistore)
* [mysqlstore](https://github.com/srinathgs/mysqlstore)
* [babou dbstore](https://github.com/drbawb/babou/blob/master/lib/session/dbstore.go)

Thank you all for sharing your code!

What makes this backend different is that it's for PostgreSQL.
