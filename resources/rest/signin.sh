#!/bin/zsh

# Vars
HOST="localhost"
PORT="8080"
PATH="signin"

# Pre
# Curl and jq installed using nix not found if path not appropriately set
# Uncomment these helper lines or replace '/usr/bin/curl' by your system values
# if curl is not included in you PATH.
# curlcmd="$(which curl)"
# alias curl=$curlcmd

post () {

  echo "POST $1"
  /usr/bin/curl -X POST $1 --header 'Content-Type: application/json' -d @resources/rest/signin.json
}

post "http://$HOST:$PORT/$PATH"