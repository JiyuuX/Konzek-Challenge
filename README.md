# Konzek-Challenge

P.S. This project made for a company's job application challenge project! Not personal project! <br /> <br /> <br /> 

Assignment: Building a Concurrent Web Service in Go <br /> 
Scenario: <br /> 
You are working on a project that requires building a high-performance web service in Go. The service will accept incoming HTTP requests, process data concurrently, and store it in a database. The application should be scalable, maintainable, and well-tested.


NOTE : <br />
The basic HTTP server is installed using the net/http package. Endpoints have been created that process POST, GET, DELETE, UPDATE requests. <br />
The data model is defined (Task structure). PostgreSQL database connection has been made and CRUD operations (Create, Read, Update, Delete) have been applied. <br />
HTTP handlers are processed simultaneously using Gorutins. Synchronization is achieved using sync.Mutex. <br />
Incoming data is verified, Errors are handled appropriately and meaningful error messages are sent to clients <br />
A record is kept for errors and important events. However, monitoring tools such as Prometheus and Grafana have not yet been integrated. <br />
Basic security measures have been implemented. Encryption (bcrypt), JWT authentication and basic authorization have been implemented. (at konzek-challenge-jwt) <br /> <br />

NOTE : <br />
The Challenge project consists of two independent directories: the directory with authorization (JWT) (konzek-challenge-jwt) and the directory without JWT (konzek-challenge). <br /> <br />

RUN <br />
"go mod" automatically detects its dependencies in the project and creates a file called go.mod and stores it there. This file contains the versions and dependency information of all packages used. <br />
You can start the project as a mod using the "go mod init" command. <br />
"go mod tidy" is a command used to clean up project dependencies and update the go.mod file. <br />  <br />


go mod init example.com/YOUR_PROJECT_NAME  <br />
go mod tidy <br /> <br />

go run .\main.go <br /> <br /> <br /> 

FOR KONZEK-CHALLENGE DIRECTORY <br /> 
API ENDPOINT: <br /> 
http://localhost:8080/  (ROOT)
<br /> <br /> 

REQUEST EXAMPLES : <br /> <br /> 

CREATE TASK : <br /> <br /> 
http://localhost:8080/ (POST): <br /> 
![image](https://github.com/JiyuuX/Konzek-Challenge/assets/139239394/de731f97-86d1-456c-ba27-d57c2cd78960)  <br />  <br /> 

POSTGRES (DATABASE): <br /> 
![image](https://github.com/JiyuuX/Konzek-Challenge/assets/139239394/eb0d32df-ef32-46ae-9e53-6e4d7b2c0b9b) <br /> <br />  <br /> 


UPDATE : <br /> 
http://localhost:8080/ (PUT): <br /> 
![image](https://github.com/JiyuuX/Konzek-Challenge/assets/139239394/55b1734b-5c52-468f-a0e3-bcc8e64fc430) <br /> 

POSTGRES (DATABASE):  <br /> 
![image](https://github.com/JiyuuX/Konzek-Challenge/assets/139239394/83b57e54-10d7-43a2-8368-99fe5df7545e)  <br /> <br /> <br />

GET ALL TASKS: <br /> 
http://localhost:8080/ (GET) : <br /> 
![image](https://github.com/JiyuuX/Konzek-Challenge/assets/139239394/47177fba-eede-4711-b783-0edd58e27315) <br />  <br /> 

DELETE : <br /> 
http://localhost:8080/ (DELETE) : <br /> 
![image](https://github.com/JiyuuX/Konzek-Challenge/assets/139239394/7a59e46c-ea03-4351-88b5-080ff00193cb) <br />

After Delete: <br />
![image](https://github.com/JiyuuX/Konzek-Challenge/assets/139239394/9cca5d4e-109c-4199-9fe7-c4c8b8324eb5) <br /> <br /> <br /> <br /> <br /> <br />



NOTE : <br />
Authentication (OAuth2 or JWT) is additionally requested within the scope of the project. The JWT implementation is available in the "konzek-challenge-jwt" directory. <br /> <br />
FOR KONZEK-CHALLENGE-JWT: <br /> <br /> 
A simple authorization process was implemented by generating JWT tokens, but session management or JWT time limitation (expiration) was not implemented, for example, to protect against attacks during login verification <br />  <br /> 

WITH JWT Authentication : <br />
API ENDPOINTS : <br />
/ (ROOT) <br />
/register  <br />
/login  <br /> <br /> 

User registration system was created. To make an API request, the user must first register (username, password, email). The user can register if the username is unique. The registered user is saved in the database. (table name in the database is users) (user password is hashed and saved in the database) and JTW token is automatically generated.
The user cannot make GET, POST, UPDATE, DELETE requests without registering.  <br />

When the user wants to send a request, the user must send the request by adding the JWT token information generated for him to the Authentication Key in the request header.

Before making a request, it must be logged into the system at the /login endpoint.

After logging in, GET, POST, DELETE, UPDATE operations can be performed - with JWT token. <br /> <br /> <br /> 

REQUEST EXAMPLES :  <br /> 
IF user tries to make GET request to root endpoint without using token, it's not authorized :   <br /> 
![image](https://github.com/JiyuuX/Konzek-Challenge/assets/139239394/866454ce-6da7-46dd-a59b-8ccad43c764a)  <br />  <br /> 

IF user try to send request even with the JWT token it is also not authorized. Have to send login first :  <br /> 
![image](https://github.com/JiyuuX/Konzek-Challenge/assets/139239394/9341e7e0-086f-404b-92b3-2399992210d6) <br /> <br /> 

After login user can send POST,GET, DELETE, UPDATE requests etc. <br /> 
![image](https://github.com/JiyuuX/Konzek-Challenge/assets/139239394/45fcc4ab-b84b-40d2-bb7c-aa0e5b5db268) <br /> <br /> 

http://localhost:8080/?id=<task_id> (PUT) : <br /> 
![image](https://github.com/JiyuuX/Konzek-Challenge/assets/139239394/7efcd5fc-6b42-4bdf-8e86-cd49df0eefa1) <br />  <br /> 

REGISTER : <br /> 
http://localhost:8080/register (POST) :  <br /> 
![image](https://github.com/JiyuuX/Konzek-Challenge/assets/139239394/27a6d9c2-048b-4bdc-89a2-38b4215d66bd)  <br /> <br /> 

IF user already exist by username: <br /> 
![image](https://github.com/JiyuuX/Konzek-Challenge/assets/139239394/4465066a-6ac3-47e2-ba55-31bf7e2efe0b) <br /> 

After registration, Postgres : <br /> 
![image](https://github.com/JiyuuX/Konzek-Challenge/assets/139239394/feaee57f-bb0e-43eb-8036-a5a4f03cd25a) <br />   <br />   <br />  

LAST THOUGHTS : <br /> 
In this project, builded a high-performance web service in Go. The service will accept incoming HTTP requests, process data concurrently, and store it in a database. For API documentation, Swagger can be used for auto-documentation, however, I did it via Postman. <br /> 

In the codes, I added the information required for database connection directly to the code section - since it is a project within the scope of the challenge project. This is actually not a safe way. It is safer to use environment variables rather than writing sensitive information (such as database credentials) directly into the code. This can help you build a better application from a security perspective and prevents sensitive information from leaking in the codebase. Environment variables provide configuration information to applications running on a system through variables in the environment. This information can be hidden from other applications running on the system or from system administrators. Therefore, storing sensitive data such as database credentials, APIs keys through environment variables increases the security of the code. <br /> 

Since the scope of the Challenge project is very, very broad, the project has been terminated by the author in here.



















