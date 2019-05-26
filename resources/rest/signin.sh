#!/bin/zsh

# Pre
# Curl installed using nix not found in some configurations.
curlcmd="$(which curl)"
alias curl=$curlcmd

# Vars
HOST="localhost"
PORT="8080"
PATH="signin"

post () {
  echo "POST $1"
  curl -X POST $1 --header 'Content-Type: application/json' -d @resources/rest/signin.json
}

post "http://$HOST:$PORT/$PATH"