#!/bin/bash
set -ex
jwt_token='eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzgzNjgyMTIsInVzZXJfaWQiOjF9.7FQnzUejbr5yRYuD8X9JZrCB9rQqhwut6oxJSjWnmvU'
curl_params=()

curl -X POST \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer ${jwt_token}" \
	-d '{"title":"post it","content":"Lorem ipsum dolor sit amet, consetetur\nsadipscing elitr, sed diam nonumy eirmod\n\ntempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet."}' \
	"http://localhost:18000/api/posts"

curl -X POST \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer ${jwt_token}" \
  -d '{"content":"This is my comment."}' \
	"http://localhost:18000/api/posts/1/comments"


curl -q "http://localhost:18000/api/posts" | jq '.'
curl -q "http://localhost:18000/api/posts/1/comments"  | jq '.'
