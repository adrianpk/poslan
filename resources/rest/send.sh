#!/bin/zsh

# Vars
HOST="localhost"
PORT="8080"
PATH="send"
TOKEN=$(/usr/bin/curl -s -X POST -H 'Accept: application/json' -H 'Content-Type: application/json' --data '{"clientID":"dd74cb9cfb5a4f1cac4d","secret":"a5ee54c8a21a4c61820f88f14c30fa5b","rememberMe":false}' http://localhost:8080/signin | /Users/adrian/.nix-profile/bin/jq -r '.token')

# Pre
# Curl and jq installed using nix not found if path not appropriately set
# curlcmd="$(which curl)"
# alias curl=$curlcmd
# jqcmd="$(which curl)"
# alias jq=$jqcmd

post () {
  echo "POST $1"
  AUTH_HEADER="Authorization: Bearer $2"
  echo $AUTH_HEADER
  /usr/bin/curl -X POST $1 --header 'Content-Type: application/json' --header $AUTH_HEADER -d @resources/rest/send.json
}

post "http://$HOST:$PORT/$PATH" "$TOKEN"