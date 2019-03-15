# Decepticons - Project Proposal <br>

<img src="b1b0d9e5a5af325e52b31a0a561fac94.jpg" alt="Decepticons Wallpaper"
	title="Decepticons Wallpaper" width="1200" height="550" /> <br>

## Project Description
1. Who is your target audience?  Who do you envision using your application? Depending on the domain of your application, there may be a variety of audiences interested in using your application.  You should hone in on one of these audiences.<br>
We have chosen to create the recommended UW Hangouts Live Video Chat App for our final project in order to target anyone that is a student or faculty member at the University of Washington. We were even thinking of possibly adding an additional link to canvas so students could automatically see/share their course content instead of having to pull up another window on their browser if feasible. If we went with this scenario, then our audience base would be specified to the students/teachers/TAs within a certain class. But overall we want to provide this video conferencing product within the UW academic community. We understand that our audience base is very broad, but as we pursue this idea in the future we can narrow our stakeholders/user base.
	
2. Why does your audience want to use your application? Please provide some sort of reasoning. <br>
Our audience would like to use this application because our project creates a tight knit community within our chat session where only UW students have access to certain resources. Not to mention it would help with working remote If integrated with canvas, we could certainly provide a unique experience where users don’t need to access multiple sites to share one file.

3. Why do you as developers want to build this application? <br>
As developers we have never tackled creating a video conference so we believe this project would give us the opportunity to learn something new. Also with the inclusion of a messaging service within the Hangouts session, we can utilize most of the content taught in this course in a different environment. Moreover, we wanted to choose an application that could be completed in the time frame allotted to us.

## Technical Description
1. Include an architectural diagram mapping 1) all server and database components, 2) flow of data, and its communication type (message? REST API?).<br>
Overall Architecture <br>
<img src="Decepticons - Updated Flow Diagram.jpeg" alt="Flow Diagram"
	title="Architectural Diagram"  />
Detailed Flow Diagram <br>
<img src="Flow Diagram 2nd Version.jpeg" alt="Flow Diagram 2.0"
	title="Architectural Diagram"  />

2. A summary table of user stories in the following format (P0 P1 means priority. These classifications do not factor into grading. They are more for your own benefit to think about what would be top priorities if you happened to run out of time)

| Use Case | Priority     | User          | Description   | Technology |
| :------  | :----------  | :-----------  | :-----------  | :-------- |
|   1      | P0           | As a User     | I want to create a Video Conference between 4 people where I can listen and speak to other users | RedisStore, MySQL, WebRTC |
|   2      | P0           | As a User     | I want to send messages to other people within the conference call   | RabbitMQ, RedisStore, HTML/CSS/JS, Messaging Microservice, Websocket
|   3      | P0           | As a User     | I want to be authenticated to use this product as a UW member | Gateway Microservice, RedisStore, MySQL |
|   4      | P1           | As a User     | I want to create a chatroom | MySQL, Chatroom Microservice |
|   5      | P1           | As a User     | I want to update my existing chatroom | MySQL, Chatroom Microservice|
|   6      | P1           | As a User     | I want to delete a chatroom | MySQL, Chatroom Microservice |
|   7      | P1           | As a User     | I want to view all chatrooms | MySQL, Chatroom Microservice |
|   8      | P1           | As a User     | I want to add members to chatroom | MySQL, Chatroom Microservice |
|   9      | P1           | As a User     | I want to delete members from chatroom | MySQL, Chatroom Microservice |
|   10      | P1           | As a User     | I want to add a message | MySQL, Message Microservice, RabbitMQ, Websocket |
|   11      | P1           | As a User     | I want to delete a message | MySQL, Message Microservice, RabbitMQ, Websocket |
|   12     | P1           | As a User     | I want to update a message | MySQL, Message Microservice, RabbitMQ, Websocket |

3. For each of your user story, describe in 2-3 sentences what your technical implementation strategy is. Explicitly note in **bold** which technology you are using (if applicable):

| Number | Strategy |
| :----- | :------- |
| 1      | **MySQL** will store our user information to call as well as any chatroom/message information not stored originally. **RedisStore** will be used to create a session when the user logs in. **WebRTC** will maintain the peer to peer connection so people can connect to chatroom|
| 2      | **MySQL** will be used to store message timestamps along with the message through a “createdAt” value. **RabbitMQ** will be used to notify users of a new message created. The RabbitMQ would talk directly to the client and not the gateway. The **Microservice** will handle all the API requests. **WebSockets** will be used to send information back to the client and update the UI. **HTML/CSS/JS** will be used to visualize all the conversations|
| 3      | **Gateway Microservice** will authenticate the user with the help of **RedisStore** to create a session. We will utilize **MySQL** to authorized the user information provided.|
| 4      | **MySQL** will store new chatroom information. We will use the **Chatroom Microservice** to handle the data request.|
| 5      | **MySQL** will store the updated chatroom information. We will use the **Chatroom Microservice** to handle the data request.|
| 6      | **MySQL** will delete a chatroom's information. We will use the **Chatroom Microservice** to handle the data request.|
| 7      | **MySQL** will get all the chatrooms. We will use the **Chatroom Microservice** to handle the data request.|
| 8      | **MySQL** will store new members for a specified chatroom. We will use the **Chatroom Microservice** to handle the data request.|
| 9      | **MySQL** will delete members for a specified chatroom. We will use the **Chatroom Microservice** to handle the data request.|
| 10      | **MySQL** will store a new message for a specified chatroom and only for the writer of the message. We will use the **Message Microservice** to handle the data request. **RabbitMQ and Websockets** would relay back the information from the server when the add button is clicked|
| 11      | **MySQL** will delete a message for a specified chatroom and only forthe writer of the message. We will use the **Message Microservice** to handle the data request. **RabbitMQ and Websockets** would relay back the information from the server when the delete button is clicked|
| 12     | **MySQL** will update a message for a specified chatroom and only for the writer of the message. We will use the **Message Microservice** to handle the data request. **RabbitMQ and Websockets** would relay back the information from the server when the add button is clicked|

4. Include a list of available endpoints your application will provide and what is the purpose it serves. Ex: GET /driver/{id}

**ALL OF THESE CALLS REQUIRE A USER IN THE X-USER HEADER**

### GET /v1/users/
* Gets all users
* Content-Type header should all be set to application/json
	* 200: Successfully retrieved all users
	* 401: User is not logged in
	* 500: Internal Server Error
	
### POST /v1/users
* Creates a new user account
* Content-Type header should all be set to application/json
	* 201: Successfully created user
	* 400: Request body is not a valid user
	* 415: Content-Type not application/json
	* 500: Internal Server Error

### Get /v1/chatroom/:chatid/members
* Get members for specific chatroom
* Content-Type header should all be set to application/json
	* 200: Successfully retrieved memebers from channel
	* 401: No valid user in the X-User header
	* 500: Internal Server Error
	
### POST /v1/sessions
* Creates a new session for the user
* Content-Type header should all be set to application/json
	* 201: Successfully created a session
	* 400: Request body is not a valid
	* 401: Email/Password incorrect
	* 415: Content-Type is not application/json
	* 500: Internal Server Error
	
### GET /v1/chatroom
* Will respond with all of the chatrooms that the user has stored
* Content-Type header should all be set to application/json
	* 200: Successfully retrieved the chatrooms
	* 401: No valid user in the X-User header
	* 500: Internal Server Error

### POST /v1/chatroom
* Creates a new chatroom
* Content-Type header should all be set to application/json
	* 201: Successfully created the chatroom
	* 403: Unauthorized to Make Channel
	* 500: Internal Server Error

### GET /v1/chatroom/:id
* Will respond by grabbing a specific chatroom 
* Content-Type header should all be set to application/json
	* 200: Successfully retrieved the chatrooms
	* 404: Chatroom with specific ID does not exist
	* 500: Internal Server Error

### PATCH /v1/chatroom/:id
* Updates the specific chatroom’s name and description
* Content-Type header should all be set to application/json
	* 200: Successfully changed
	* 403: User forbidden
	* 404: Chatroom with specific ID does not exist
	* 500: Internal Server Error
	
### POST /v1/chatroom/:id
* Response includes added message body
* Content-Type header should all be set to application/json
	* 201: Successfully retrieved the chatrooms
	* 403: User forbidden
	* 404: Chatroom with specific ID does not exist
	* 500: Internal Server Error

### DELETE /v1/chatroom/:id
* Deletes a specific chatroom
	* 200: Successfully deleted
	* 403: User forbidden
	* 404: Chatroom with specific ID does not exist
	* 500: Internal Server Error

### POST /v1/chatroom/:chatid/:id
* Response includes updated chatroom's name/description/private/public
* Content-Type header should all be set to application/json
	* 201: Successfully retrieved the chatrooms
	* 403: User forbidden
	* 404: Chatroom with specific ID does not exist
	* 500: Internal Server Error

### PATCH /v1/chatroom/:chatid/messages/:id
* Updates the body of a specific message
* Content-Type header should all be set to application/json
	* 200: Successfully changed
	* 403: User forbidden
	* 404: Message with specific ID does not exist
	* 500: Internal Server Error

### DELETE /v1/chatroom/:chatid/messages/:id
* Will delete a specific message
	* 200: Successfully deleted a message
	* 403: User forbidden
	* 404: Message with specific ID does not exist
	* 500: Internal Server Error

### POST /v1/chatroom/:id/members
* Updates the specific chat room’s member list
* Content-Type header should all be set to application/json
	* 200: Successfully changed
	* 401: No valid user in the X-User header
	* 403: User forbidden/Not an Actual User
	* 404: Chatroom with specific ID does not exist
	* 500: Internal Server Error

### DELETE /v1/chatroom/:id/members
* Deletes a user from the list of members in the chatroom
	* 200: Successfully changed
	* 401: No valid user in the X-User header
	* 403: User forbidden/Not an Actual User
	* 404: Chatroom with specific ID does not exist
	* 500: Internal Server Error

### GET /v1/users/:id
* Gets specific user
* Content-Type header should all be set to application/json
	* 201: Successfully created user
	* 400: User id is not valid
	* 401: User is not logged in
	* 500: Internal Server Error
	
5. Include any database schemas as appendix
### Sessions
This will contain a redis key-value store that contains sessionIDs, session startTime, and the users information

### Users
```
create table if not exists users (
   id int not null auto_increment primary key,
   email varchar(128) not null UNIQUE,
   passHash binary(60) not null,
   userName varchar(255) not null UNIQUE,
   firstName varchar(64) not null,
   lastName varchar(128) not null,
   photoURL varchar(2083) not null
);
```

### Messages
```
create table if not exists messages (
   id int not null auto_increment UNIQUE primary key,
   channelID int not null,
   body varchar(128) not null,
   createdAt datetime not null,
   creatorID int not null,
   editedAt datetime
);
```

### Users Sign In
```
create table if not exists signin (
   pKey int not null auto_increment primary key,
   id int not null,
   signingTimeDate datetime not null,
   ipAddress varchar(128) not null UNIQUE
);
```

### Chatroom
```
create table if not exists chatroom (
   id int not null auto_increment primary key,
   chatroomID int not null,
   createdAt datetime not null,
   creatorID int not null
);
```

### Chatroom Members
```
create table if not exists chatroom_members (
   id int not null auto_increment primary key,
   chatroomID int not null,
   userID int not null
);
```


