curl -X POST http://localhost:8080/api/articles \
  -H 'Content-Type: application/json' \
  -d '{"title":"Hello Reflection","body":"Deep dive","tags":["go","gin"]}'
