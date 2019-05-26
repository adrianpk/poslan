#!/bin/zsh

# Vars
HOST="localhost"
PORT="8080"
PATH="signin"

# Pre
# Curl installed using nix not found if path not appropriately set
# curlcmd="$(which curl)"
# alias curl=$curlcmd

post () {

  echo "POST $1"
  /usr/bin/curl -X POST $1 --header 'Content-Type: application/json' -d @resources/rest/signin.json
}

post "http://$HOST:$PORT/$PATH"