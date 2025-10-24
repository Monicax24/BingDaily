# Backend API
The backend is a REST API that will communicate with the client. 
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
`/auth/register` | `POST` | `authorization required`

For a user to be registered in the backend, they first need to register. The registration request should be formatted like this:

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
    "message": NULL | string
}
```


## User

### Update User Profile
`/user/profile/update` | `POST` | `authorization required`

To update a user's profile, a request should be sent like this:

#### Request
```
{
    "userId": string, // this field is not modifiable
    "username": string*,
    "email": string*,
    ...
}
```
\* These fields are optional... only include the fields that need to be updated.

#### Response
```
{
    "status": "success" | "fail",
    "message": NULL | string
}
```

### Retrieve User Profile Data

`/user/profile` | `GET` | `authorization required`

#### Request
```
{
    "userId": string
}
```

#### Response
```
{
    "status": "success" | "fail",
    "message": NULL | string,

    "data": UserObject*
}
```

* The `UserObject` structure will follow this format:
```
{
    "userId": string,
    "email": string,
    "username": string,
    "joinDate": string*,
    "communities": string[]**,
}
```
\* `joinDate` will be in `YYYY-MM-DD` format.

\*\* `communities` will be an array of `communityIds`.


## Community

### Retrieve Core Community Data

`/community` | `POST` | `authorization required`

#### Request
```
{
    "communityId": string
}
```

#### Response
```
{
    "status": "success" | "fail",
    "message": NULL | string,

    "data": CommunityObject*
}
```
\* The `CommunityObject` structure will follow this format:
```
{
    "communityId": string,
    "name": string,
    "description': string,
    "prompt": string,

    "memberCnt": int,
    ...
}
```

### Get Community Posts
`/community/posts` | `POST` | `authorization required`

This endpoint is to retrieve all the posts from a community. This request will only be fulfilled if the user requesting the data has already posted within that community.

#### Request
```
{
    "communityId": string
}
```

#### Response
```
{
    "status": "success" | "fail",
    "message": NULL | string,

    "data": PostObject*
}
```

\* The `PostObject` structure will follow this format:
```
{
    "postId": string,
    "communityId": string,
    "userId": string, // in the DB this is listed as author rn
    
    "caption": string,
    "timePosted": string, // will follow YYYY-DD-MMTHH:MM:SS (ISO 8601)
    ...
}
```

### Upload Post to Community
`/community/posts/upload` | `POST` | `authorization required`

Can only upload 1 post per prompt, so the database will check to see if the posting user has not posted yet.

#### Request
```
{
    "communityId": string,
    "userId": string,

    "caption": string,
    ...
}
```

#### Response
```
{
    "status": "success" | "fail",
    "message": NULL | string
}
```

