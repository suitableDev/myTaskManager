# ✔ JUST DO IT – *A Full-Stack Task Manager*  

I created this project to learn backend technologies by extending a React TODO app I was building to explore **React Context, Reducers, and Memoization**. I then wanted to gain hands-on experience with backend development, so I decided to implement a **Golang backend** for authentication and data management.  

This project primarily serves as a learning experience for backend development, with a focus on **Golang** and **MongoDB**.

🚧 **Note:** *The frontend has not been integrated into this project yet.* 🚧

## Features  
- User authentication with **JWT**  
- Secure password storage using **bcrypt**  
- Rate-limiting with **ratelimit**  
- CRUD operations with **MongoDB**  
- Email functionality with **Postmark** (for account verification, password reset, etc.)  
- Backend uses **Air** for live reloading during development  

# 🛠️ Technical Overview  

### ◽ Backend  
- [Golang](https://go.dev/) – Backend language  
- [Gin](https://gin-gonic.com/) – Web framework  
- [Air](https://github.com/cosmtrek/air) – Live reloading for development  
- [jwt-go](https://github.com/golang-jwt/jwt) – JWT-based authentication  
- [bcrypt](https://golang.org/x/crypto/bcrypt) – Password hashing  
- [ratelimit](https://github.com/juju/ratelimit) – Rate limiting  
- [MongoDB](https://www.mongodb.com/) – Database  
- [Postmark](https://postmarkapp.com/) – Email service for sending verification and notification emails

### ◽ Frontend (Planned)  
- [React](https://reactjs.org/) – UI Library  
- [TypeScript](https://www.typescriptlang.org/) – Static typing  

# 🏁 Getting Started  

### Prerequisites  
- Node.js & npm installed (for frontend, once integrated)  
- Golang installed  
- MongoDB running – Either:  
  - **MongoDB Atlas** account (for a cloud-hosted database)  
  - **Local MongoDB instance** installed manually  
  - **Docker-based MongoDB instance** (see Docker setup instructions)  
- Postmark account (for email functionality)  


# Installation  
### Setup: 
```sh
# Clone the repository
git clone git@github.com:suitableDev/myTaskManager.git
cd myTaskManager

# Install backend dependencies (the tools our code needs)
cd server
go mod tidy
```

### Environment Variables
To ensure the backend can connect to the database, properly handle authentication, and integrate email functionalities, rename either:
-  `example.env` to `.env` 
 
- `example.docker-compose.env` to `docker-compose.env` for Dockerised version

then update the following values:
```sh
# Leave this unless you're running mongodb locally WITHOUT docker
MONGO_URI=mongodb://mongo:27017 
# Change these
SECRET_KEY="your_secret_key_for_jwt_hashing"
POSTMARK_API_TOKEN="your-postmark-api-token"
POSTMARK_SENDER_EMAIL="your-verified-email@example.com"
POSTMARK_EMAIL_LINK_ADDRESS="your_site_address"
```
- MONGO_URI – *Set this to your MongoDB connection string.*
- SECRET_KEY – *Set this to any secure string for signing JWT tokens.*
- POSTMARK_API_TOKEN – *Set this to your Postmark API token for email sending.*
- POSTMARK_SENDER_EMAIL – *Set this to the email address you have verified with Postmark.*
- POSTMARK_EMAIL_LINK_ADDRESS – *Set this to the base URL for your site (used for email link generation).*



# Option 1: Run Locally (Without Docker) 🏠

```sh
# Install 'air' for live reloading (makes development easier!)
go install github.com/cosmtrek/air@latest

# Run the backend with live reloading
air
```
# Option 2: Run Locally With Docker 🐳
*From main app folder. Ensure you have updated the example docker-compose.env*

Build and start the backend and database using Docker Compose

```sh
docker-compose up --build
```

 Access the application at: http://localhost:8080


```sh
# Stops the containers
docker-compose down
# Stops the containers and removes persitant data
docker-compose down -v
```

## Health check
```sh
### Check the server is up and running
http://localhost:8080/health
```
