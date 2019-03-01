# Deceptions - Project Proposal

## Project Description
1. Who is your target audience?  Who do you envision using your application? Depending on the domain of your application, there may be a variety of audiences interested in using your application.  You should hone in on one of these audiences.
⋅⋅⋅ We have chosen to create the recommended UW Hangouts Live Video Chat App for our final project in order to target anyone that is a student or faculty member at the University of Washington. We were even thinking of possibly adding an additional link to canvas so students could automatically see/share their course content instead of having to pull up another window on their browser if feasible. If we went with this scenario, then our audience base would be specified to the students/teachers/TAs within a certain class. But overall we want to provide this video conferencing product within the UW academic community. We understand that our audience base is very broad, but as we pursue this idea in the future we can narrow our stakeholders/user base.
	
2. Why does your audience want to use your application? Please provide some sort of reasoning. 
⋅⋅⋅ Our audience would like to use this application because our project creates a tight knit community within our chat session where only UW students have access to certain resources. Not to mention it would help with working remote If integrated with canvas, we could certainly provide a unique experience where users don’t need to access multiple sites to share one file.

3. Why do you as developers want to build this application?
⋅⋅⋅ As developers we have never tackled creating a video conference so we believe this project would give us the opportunity to learn something new. Also with the inclusion of a messaging service within the Hangouts session, we can utilize most of the content taught in this course in a different environment. Moreover, we wanted to choose an application that could be completed in the time frame allotted to us.

## Technical Description
1. Include an architectural diagram mapping 1) all server and database components, 2) flow of data, and its communication type (message? REST API?).

2. A summary table of user stories in the following format (P0 P1 means priority. These classifications do not factor into grading. They are more for your own benefit to think about what would be top priorities if you happened to run out of time)

| Use Case | Priority     | User          | Description   | Technology |
| :------  | :----------  | :-----------  | :-----------  | :-------- |
|   1      | P0           | As a User     | I want to create a Video Conference between 4 people where I can listen and speak to other users | RedisStore, MySQL |
|   2      | P1           | As a User     | I want to send messages to other people within the conference call   | RabbitMQ, RedisStore, HTML/CSS/JS
|   3      | P0           | As a User     | I want to be authenticated to use this product as a UW member | RedisStore, MySQL |
|   4      | P1           | As a User     | I want to add more people to a UW Hangout | MySQL |
|   5      | P2           | As a User     | I want to know who joined the conference call | RabbitMQ |
|   6      | P3           | As a Developer | I want to retain HD video during the hangout | WebRTC |
|   7      | P3           | As a Developer | I want to maintain fast messaging updates across different users based on message timestamps | RabbitMQ, MySQL |

2. For each of your user story, describe in 2-3 sentences what your technical implementation strategy is. Explicitly note in **bold** which technology you are using (if applicable):

| Number | Strategy |
| :----- | :------- |
| 1      | We will be using **Redis** to store current and previous sessions of conferences. **MySQL** will store our user information to call. |
| 2      | **MySQL** will store the user messages for each session to provide context of the previous chat. We will use **HTML/CSS/JS** to provide client feedback and to parse the message data pulled from **MySQL**. |
| 3      | **MySQL** will be used to support multiple users in each UW Hangout. This is through updating a **MySQL** column that contains a list of users in the call. |
| 4      | **MySQL** will authenticate a valid UW Husky ID and check **Redis** for all  the  multiple user sessions and to populate the user interface with all these sessions. |
| 5      | **RabbitMQ** is used to notify all the users that are in the same conference call. The information will be displayed as a popup saying “_user_ has joined this call”. |
| 6      | WebRTC offers HD video configurations by consuming their API and prompting users to give access to their webcam for the conference. Video display size can be customized. |
| 7      | **MySQL** will be used to store message timestamps along with the message through a “createdAt” value. **RabbitMQ** will be used to notify users of a new message created. |