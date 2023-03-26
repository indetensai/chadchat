# CHADCHAT

**CHADCHAT** is a Golang backend project that showcases my knowledge and skills in web development. It uses Fiber, a fast and lightweight web framework, to handle HTTP requests and responses. It uses pgx, a PostgreSQL driver and toolkit, to connect and interact with the database. It uses websockets, a protocol for bidirectional communication, to enable real-time chat functionality. It uses JWT, a standard for secure authentication and authorization, to protect the endpoints and verify the users. **CHADCHAT** is designed with clean architecture principles. It can be easily deployed with Docker, a platform for containerization and orchestration. **CHADCHAT** is a project that I'm proud of and I hope it will impress potential employers and clients.

## Prerequisites

Before you continue, ensure you have met the following requirements:

- You have installed the latest version of Go.
- You have installed PostgreSQL and created a database for the project.
- You have created .env file, with `POSTGRES_URL` defined.

## Installation

To install **CHADCHAT**, follow these steps:

1. Clone this repository: `git clone https://github.com/indetensai/chadchat.git`
2. Change into the project directory: `cd chadchat`
3. Install the dependencies: `go mod download`
4. Fill `.env` file with the required environment variables.
5. Build the executable: `go build -o chadchat cmd/chadchat/main.go`

## Usage

To run CHADCHAT, follow these steps:

1. Start the executable: `./chadchat`
2. The server will listen on port 8080.
3. To interact with the chat API, you can use any HTTP client or websocket client of your choice.

The chat API has the following endpoints:

- `POST /user/register`: Create a new user account.
- `POST /user/login`: Login with an existing user account and get a JWT token.
- `GET /refresh`: Refresh user's access and refresh token.
- `POST /chatroom`: Create a new chatroom. Requires authentication.
- `GET /ws/:room_id<guid>`: Connect to a chatroom by ID using websockets. Requires authentication.
- `GET /chat/:room_id<guid>/history`: Show chatroom history by ID. Requires authentication.