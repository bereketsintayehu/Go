# Task Manager API

This is a task management system with user authentication, role-based access control, and task CRUD operations.

## Getting Started

### Prerequisites

- Go 1.16 or later
- MongoDB Atlas or local MongoDB instance
- Git

## [Postaman API Doc](https://documenter.getpostman.com/view/20213080/2sA3rzLD74)

### Installation

1. **Clone the repository**

```bash
git clone https://github.com/yourusername/task-manager.git
cd task-manager
```

2. **Set up environment variables**

Create a `.env` file in the root directory and add the following variables:

```env
MONGO_URI="url"
JWT_SECRET="ursecret"
SUPER_ADMIN_EMAIL="admin@admin.com"
SUPER_ADMIN_PASSWORD="admin"
```

3. **Install dependencies**

Ensure you have Go modules enabled, then run:

```bash
go mod tidy
```

4. **Run the server**

```bash
go run main.go
```

The server will start on `http://localhost:8080`.

### Creating a Super Admin

Upon running the server for the first time, a super admin user will be created with the credentials specified in the `.env` file.

## API Endpoints

### Authentication

- **Register a User**
  - **POST** `/register`
  - Protected by: No protection, public access
  - Body: `{"name": "John Doe", "email": "john.doe@example.com", "password": "password123", "role": 0}`

- **Login**
  - **POST** `/login`
  - Protected by: No protection, public access
  - Body: `{"email": "john.doe@example.com", "password": "password123"}`

### User Management

- **Get All Users**
  - **GET** `/users`
  - Protected by: Admin, Super Admin
  - Headers: `Authorization: Bearer <JWT_TOKEN>`

- **Get User by ID**
  - **GET** `/users/:id`
  - Protected by: User, Admin, Super Admin
  - Headers: `Authorization: Bearer <JWT_TOKEN>`

- **Get Current User**
  - **GET** `/me`
  - Protected by: User, Admin, Super Admin
  - Headers: `Authorization: Bearer <JWT_TOKEN>`

- **Update User**
  - **PUT** `/users/:id`
  - Protected by: User (self), Admin, Super Admin
  - Headers: `Authorization: Bearer <JWT_TOKEN>`
  - Body: `{"name": "John Doe Updated", "email": "john.doe@example.com", "role": 2}`

- **Delete User**
  - **DELETE** `/users/:id`
  - Protected by: User (self), Admin, Super Admin
  - Headers: `Authorization: Bearer <JWT_TOKEN>`

### Task Management

- **Get All Tasks**
  - **GET** `/tasks`
  - Protected by: User, Admin, Super Admin
  - Headers: `Authorization: Bearer <JWT_TOKEN>`

- **Get Task by ID**
  - **GET** `/tasks/:id`
  - Protected by: User, Admin, Super Admin
  - Headers: `Authorization: Bearer <JWT_TOKEN>`

- **Create Task**
  - **POST** `/tasks`
  - Protected by: User, Admin, Super Admin
  - Headers: `Authorization: Bearer <JWT_TOKEN>`
  - Body: `{"title": "New Task", "description": "Description of the new task", "status": "Pending"}`

- **Update Task**
  - **PUT** `/tasks/:id`
  - Protected by: User (creator), Admin, Super Admin
  - Headers: `Authorization: Bearer <JWT_TOKEN>`
  - Body: `{"title": "Updated Task", "description": "Updated description for the task", "status": "Completed"}`

- **Delete Task**
  - **DELETE** `/tasks/:id`
  - Protected by: User (creator), Admin, Super Admin
  - Headers: `Authorization: Bearer <JWT_TOKEN>`


## Built With

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [MongoDB](https://www.mongodb.com/) - NoSQL database
- [Go](https://golang.org/) - Programming language
