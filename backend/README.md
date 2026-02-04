# Backend API

## Notes
### Post Request

The backend follows a REST API format. 
For POST endpoints, make sure to include the following header in your HTTP request:
```
Content-Type: application/json
```
*This is not really necessary, but good practice!*

### Responses
All endpoints will send back a response in this format:
```json
{
    "status": "success" | "fail",
    "message": string,
    "data": {},
}
```

`status` denotes if the request was successful. 

`message` includes any debug information regarding the prior request. This is useful to display/check to troubleshoot any failed requests.

`data` contains any relevant data returned by the server for the request. The specific structure of this field will be detailed in the docs (if applicable).

### Media Upload/Download

The backend uses S3 to store media (images, videos, etc.). To protect the privacy of users, some media operations require authentication. A presigned upload/download url will be sent by the backend for these operations.

For upload operations, a `PUT` request must be sent to the url.

Right now the urls are set to expire in **10 minutes**, but that is subject to change. There also isn't any functionality to extend the duration of a url. 


## Authentication
Some API calls will require authentication. These requests will be marked `authorization required`. To include authentication, add the following header to your HTTP request:
```
Authorization: Bearer <token>
```

During authentication, the server will automatically retrieve the current user based on the contents of the token. All operations performed will be executed on behalf of this user.

## User

### Register User
`/users/register` | `POST` | `authorization required`

Before making any requests, the user first needs to register. A registration request should look like this:

#### Request
```json
{
    "email": string,
    "username": string,
    "updatePicture": bool,
}
```

*Note, the email must be the same as the one used by firebase!*

#### Response
```json
"data": {
    "uploadUrl" : string*,
}
```

\* `uploadUrl` is only sent when the `updatePicture` field is set to `true`. It contains a presigned `PUT` url that can be used to upload media.


### Update User Profile
`/users/update` | `POST` | `authorization required`

To update a user's profile, send a request detailing which properties need to be modified:

#### Request
```json
{
    "username": string*,
    "updatePicture": bool*,
}
```
\* All fields are optional, only include what needs to be changed.

#### Response
```json
"data": {
    "uploadUrl": string*,
}
```

\* Only sent back if `updatePicture` is set to `true`.

*Note: there is a chance that a request may only be partially fulfilled. `success` will be set to `false` upon any operation failing. Check `message` for more details.*

### Retrieve Current User Profile Data

`/users/profile` | `GET` | `authorization required`

#### Response
```json
"data": {
    "user": UserObject*,
}
```

\* `UserObject` will follow this format:
```json
{
    "userId": string,
    "email": string,
    "username": string,
    "joinDate": string*,
    "profilePicture": string,
    "communities": string[]**,
}
```
\* The `joinDate` field follows ISO 8601 (YYYY-MM-DD).

\*\* `communities` will be an array of `communityId` strings.


## Community

### List Communites

`/communities/list` | `GET` | `authorization required`

#### Response
```json
"data": {
    "communities": CommunityObject[]*
}
```

\* `CommunityObject` will follow this format:
```json
{
    "communityId": string,
    "name": string,
    "description": string,
    "prompt": string,
    "memberCnt": int
}
```


### Retrieve Core Community Data

`/communities/<communityId>` | `GET` | `authorization required`

#### Response
```json
"data": {
    "community": CommunityObject,
}
```

### Join Community
`/communities/join/<communityId>` | `GET` | `authorization required`

Backend will attempt to add the current user to specified community.

### Leave Community
`/communities/leave/<communityId>` | `GET` | `authorization required`

Backend will attempt to remove the current user to specified community.

### Get Community Posts
`/communities/posts/<communityId>` | `GET` | `authorization required`

This endpoint is to retrieve all the posts from a community. This request will only be fulfilled if the user requesting the data has already posted within that community.

#### Response
```json
"data": {
    "posts": PostObject[]*,
}
```

\* `PostObject` will follow this format:
```json
{
    "postId": string,
    "communityId": string,
    "userId": string,
    "caption": string,
    "timePosted": string, // follows YYYY-DD-MMTHH:MM:SS (ISO 8601)
    "imageUrl": string*,
}
```
\* `imageUrl` is a signed URL that can be used to download/access media

### Upload Post to Community
`/communities/posts/upload` | `POST` | `authorization required`

A user can only have 1 post per community. The database will check to see if the user has posted yet.

#### Request
```json
{
    "communityId": string,
    "caption": string,
}
```

#### Response
```json
"data": {
    "postId": string, // newly created postId
    "uploadUrl": string, // send PUT with media here
}
```

### Delete Post from Community

`/communities/posts/delete/<communityId>` | `POST` | `authorization required`

Will delete the user's post in the specified community if it exists.
