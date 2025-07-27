# idm

### Adminka
http://localhost:9990/admin/

### сайт Keycloak
https://www.keycloak.org/app/


## получение токена
curl --location 'http://localhost:9990/realms/idm/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
-d 'client_id=idmapp' \
-d 'username=idm19090' \
-d 'password=12345' \
-d 'grant_type=password' \
-d 'client_secret=bZdU9YfqMkMTPX8zXNUoIs0uT88MRbVw' \
-d 'scope=openid'

## создание сотрудника
curl --location 'https://localhost:8080/api/v1/employees' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer xxxxxx' \
--data '{
"name": "Vava Viva"
}'