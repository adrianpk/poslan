# Poslan
A redundant mail delivery system.

## Test
```bash
$ make test-mailer
```

## SPA Client

**[Poslan-cli](https://github.com/adrianpk/poslan-cli)**

## Make
**Main commands**

```bash
all
build
tests
clean
run
deps
build-stage
```

See [`makefile`](makefile) for new additions.

## Curl Test

```bash
$ ./resources/rest/signin.sh
$ ./resources/rest/send.sh
```

## Notes
**This is a kind of PoC, in a real world app**
* Container environment variables should not be under version control.
