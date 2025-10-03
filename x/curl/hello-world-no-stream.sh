

curl -s localhost:11435/api/generate -d '{
    "model": "llama3.2:latest",
    "prompt": "hello",
    "stream": false
}'
