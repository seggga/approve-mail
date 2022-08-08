## Тестовые запросы

аутентификация по basic-auth, ответ в куки (сервис auth)
``` bash
rm -rf /tmp/cookie.txt && \
AUTH=$(echo -ne "user_1:password_1" | base64 --wrap 0) && \
curl \
  --verbose \
  --request POST \
  --header "Content-Type: application/json" \
  --header "Authorization: Basic $AUTH" \
  --cookie-jar /tmp/cookie.txt \
  http://localhost:3001/login

```
  
аутентификация по логину и паролю в JSON, ответ в JSON (сервис auth)
``` bash
curl \
    --verbose \
    --request POST \
    --header "Content-Type: application/json" \
    --data '{"login":"user_1","password":"password_1"}' \
    http://localhost:3001/postman/login

```

аутентификация по JWT из куки, ответ в JSON (сервис auth)
``` bash
curl \
    --verbose \
    --request POST \
    --header "Content-Type: application/json" \
    --cookie /tmp/cookie.txt \
    --cookie-jar /tmp/cookie.txt \
    http://localhost:3001/postman/validate

```