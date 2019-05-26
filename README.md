# Poslan
A redundant mail delivery system.

## Status
  * [Status](docs/status/index.md)

## Pending
  * Integration tests.
  * Simple SPA client.
  * Installation, deployment & use docs.
  * Deploy to GKE.
  * ...

## Make
**Main commands**

```bash
all
build
test
clean
run
deps
build-stage
```

See `makefile` for new additions.

## Test

```bash
$ ./resources/rest/signin.sh
$ ./resources/rest/send.sh
```

## Notes
**This is a kind of PoC, in a real world app**
* Container environment variables should be available in control version repository.
