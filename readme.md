# Register  
```
curl --location --request POST 'http://localhost:10000/register' \
--header 'Content-Type: text/plain' \
--data-raw '{
    "username": "user",
    "password": "user"
}'
```

# Check User by Username

```
curl --location --request GET 'http://localhost:10000/search-user/michael' \
--data-raw ''
```