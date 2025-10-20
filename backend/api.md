# Backend API
The backend is a REST API that will communicate to the client via JSON post requests. Due to the data transmitted being in JSON format, make sure to add the following header to your HTTP request:
```json
Content-Type: application/json
```

## Authentication
Some API calls will require authentication to be completed. These requests will be marked `authorization required` in this documenation. To include authentication in a request add the following header to your HTTP request:
```json
Authorization: Bearer <token>

```

### Registering an Account
In order for a user to be registerd in the backend, they will need to send a registration request. The registration request should be formatted as such:

#### Request
```json
{
    "user-id": string, 
    "email": string,
    "username": string,
    "join-date": string*,
    ...
}
```
\* For the `join-date` field, follow ISO 8601 (YYYY-MM-DD).

#### Response
```json
{
    "status": "success" or "fail",
    "error": NULL or string
}
```


## Fetch Data



## Upload Data

