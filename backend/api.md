# Backend API
The backend is a REST API that will communicate to the client via JSON post requests. Due to the data transmitted being in JSON format, make sure to add the following header to your HTTP request:
```
Content-Type: application/json
```

## Authentication
Some API calls will require authentication to be completed. These requests will be marked `authorization required` in this documenation. To include authentication in a request add the following header to your HTTP request:
```
Authorization: Bearer <token>

```

### Registering an Account
`/auth/register`
`authorization required`
For a user to be registered in the backend, they will need to send a registration request. The registration request should be formatted like this:

#### Request
```
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
```
{
    "status": "success" or "fail",
    "error": NULL or string
}
```


## Fetch / Update Data
For all of these requests

### User



## Upload Data

