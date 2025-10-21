# Backend API
The backend is a REST API that will communicate to the client. 

For the POST endpoints, make sure to include the following header in your HTTP request:
```
Content-Type: application/json
```
*This is not really necessary, but good practice!*

## Authentication
Some API calls will require authentication to be completed. These requests will be marked `authorization required` in this documentation. To include authentication, add the following header to your HTTP request:
```
Authorization: Bearer <token>
```

### Register Account
For a user to be registered in the backend, they first need to register. The registration request should be formatted like this:

`/auth/register` | `POST` | `authorization required`

#### Request
```
{
    "userId": string, 
    "email": string,
    "username": string,
    "joinDate": string*,
    ...
}
```
\* The `joinDate` field should follow ISO 8601 (YYYY-MM-DD).

#### Response
```
{
    "status": "success" | "fail",
    "error": NULL | string
}
```


## User

### Update User Profile

### Retreive User Profile

## Community

## Post

