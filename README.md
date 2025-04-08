# URL Shortener

A robust URL shortening service built with Go. It offers secure user authentication, detailed URL management, and analytics capabilities.

## Features

- **URL Shortening**: Generate shortened URLs that redirect to original destinations.
- **User Authentication**: Secure Firebase-based authentication system.
- **User Management**: Create, manage, and verify user accounts.
- **URL Management**: Create, list, and track shortened URLs per user with pagination.
- **Analytics**: Gather insights on URL usage and performance.
- **PostgreSQL Database**: Reliable data persistence powered by GORM ORM.
- **RESTful API**: Well-structured API endpoints for all operations.

## Technologies

- **Backend**: Go (Golang) with the Gin web framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: Firebase Authentication
- **Environment Management**: godotenv for configuration

## Prerequisites

- Go 1.23 or higher
- PostgreSQL database
- Firebase project with authentication enabled

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/drako02/url-shortener.git
   cd url-shortener
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Create a `.env` file in the project root with the following variables:

   ```
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=your_db_name
   ```

4. Place your Firebase admin SDK credentials in a secure location and update the path in `config/firebase-admin.go`.

5. Run the application:

   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8080` by default.

## API Documentation

### URL Endpoints

#### Create a shortened URL
- **POST** `/create`
- **Body**: 
    ```json
    {
        "url": "https://example.com/long-url",
        "uid": "user_id"
    }
    ```
- **Response**: Returns the created shortened URL with a corresponding short code.

#### Access a shortened URL
- **GET** `/:shortCode`
- **Response**: Redirects to the original URL.

#### Get user's URLs
- **POST** `/user-urls`
- **Body**:
    ```json
    {
        "uid": "user_id",
        "limit": 10,
        "offset": 0
    }
    ```
- **Response**: Provides a paginated list of the user's shortened URLs.

### User Endpoints

#### Create a new user
- **POST** `/users`
- **Body**: User details including UID from Firebase.
- **Response**: Returns the created user details.

#### Check if user exists
- **POST** `/users/exists`
- **Body**: Email to check.
- **Response**: Boolean indicating if the user exists.

#### Get user details
- **GET** `/users/:uid`
- **Response**: Provides detailed user information.
- **Note**: This endpoint requires authentication.

## Project Structure

```
url_shortener/
├── config/           # Database and Firebase configuration files
├── handlers/         # Request handlers for API endpoints
├── middlewares/      # Authentication and other middleware functions
├── models/           # Database models
├── repositories/     # Database operations and data access layers
├── routes/           # API route definitions
├── .env              # Environment variables (not tracked in git)
├── go.mod            # Go module dependencies
├── go.sum            # Go module checksums
└── main.go           # Application entry point
```

## Configuration

### Database Configuration

Database connection is configured in `config/database.go`. The application uses PostgreSQL with the following environment variables:

- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name

### Firebase Configuration

Firebase authentication is set up in `config/firebase-admin.go`. You need to provide the path to your Firebase Admin SDK credentials file.

## Development

To run the application in development mode:

```bash
go run main.go
```

The server will automatically launch on `http://localhost:8080`.

## Contributing

1. Fork the repository.
2. Create a feature branch: `git checkout -b feature-name`
3. Commit your changes: `git commit -am 'Add feature'`
4. Push the branch: `git push origin feature-name`
5. Submit a pull request.

## License

[MIT License](LICENSE)

## Contact

Project maintained by [drako02](https://github.com/drako02)