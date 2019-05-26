#!/bin/sh
# Build
make build

# Free ports
killall -9 main

# Set environment variables
# App
export LIVENESS_PING_PORT=18080
# Probe
export PROBE_PING_PORT=18081
export PROBE_MAX_INACTIVE_SECS=900
export PROBE_MAX_INACTIVE_ATTEMPTS=3
# Mailers
# Amazon SES
export PROVIDER_NAME_1=amazon
export PROVIDER_TYPE_1=amazon-ses
export PROVIDER_ENABLED_1=true
export PROVIDER_TESTONLY_1=false
export PROVIDER_SENDER_NAME_1=SendMailTest
export PROVIDER_SENDER_EMAIL_1=sendmailtest@sharkslasers.com
# Sendgrid
export PROVIDER_NAME_2=sendgrid
export PROVIDER_TYPE_2=sendgrid
export PROVIDER_ENABLED_2=true
export PROVIDER_TESTONLY_2=false
export PROVIDER_SENDER_NAME_2=SendMailTest
export PROVIDER_SENDER_EMAIL_2=sendmailtest@sharkslasers.com

# Start
# Ref.: Fresh - https://github.com/gravityblast/fresh
# go get github.com/pilu/fresh
# fresh
go run --race main.go