# ✔ JUST DO IT – *A Full-Stack Task Manager*  

I created this project to learn backend technologies by extending a React TODO app I was building to explore **React Context, Reducers, and Memoization**. I then wanted to gain hands-on experience with backend development, so I decided to implement a **Golang backend** for authentication and data management.  

This project primarily serves as a learning experience for backend development, with a focus on **Golang** and **MongoDB**.

🚧 **Note:** *The frontend has not been integrated into this project yet.* 🚧

##  Features  
- User authentication with **JWT**  
- Secure password storage using **bcrypt**  
- Rate-limiting with **ratelimit**  
- CRUD operations with **MongoDB**  
- Backend uses **Air** for live reloading during development  

## 🛠️ Technical Overview  

### ◽ Backend  
- [Golang](https://go.dev/) – Backend language  
- [Gin](https://gin-gonic.com/) – Web framework
- [Air](https://github.com/cosmtrek/air) – Live reloading for development  
- [jwt-go](https://github.com/golang-jwt/jwt) – JWT-based authentication  
- [bcrypt](https://golang.org/x/crypto/bcrypt) – Password hashing  
- [ratelimit](https://github.com/juju/ratelimit) – Rate limiting  
- [MongoDB](https://www.mongodb.com/) – Database  

### ◽ Frontend (Planned)  
- [React](https://reactjs.org/) – UI Library  
- [TypeScript](https://www.typescriptlang.org/) – Static typing  

## ⚡ Getting Started  

### Prerequisites  
- Node.js & npm installed (for frontend, once integrated)  
- Golang installed  
- MongoDB running  

### Installation  

```sh
# Clone the repository
git clone git@github.com:suitableDev/myTaskManager.git
cd myTaskManager

# Install backend dependencies
cd server
go mod tidy

# Run the backend with live reloading
air
```

### Environment Variables  
Before running the backend, rename `example.env` to `.env` and update the following values:  

```sh
MONGO_URI="your_mongodb_connection_string"
SECRET_KEY="your_secret_key_for_jwt_hasing"
```

- **MONGO_URI** – Set this to your MongoDB connection string.  
- **SECRET_KEY** – Set this to any secure string for signing JWT tokens.  

This will ensure the backend can connect to the database and properly handle authentication. 
