# Poslan
A redundant mail delivery system.

## Test
```bash
$ make test-mailer
```

## Staging
App is running on GKE on ip 35.246.nnn.mmm 😉

## Pending
Some certificate problem with SendGrid:

```bash
{"msg":"Post https://api.sendgrid.com/v3/mail/send: x509: certificate signed by unknown authority"}
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
