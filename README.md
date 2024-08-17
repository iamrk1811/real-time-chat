### Real-Time Chat Application
This project is a real-time chat application built using Go (Golang) and `https://github.com/gorilla/websocket`. It supports direct messaging between users and group chat functionality. The application includes session-based authentication, message storage, and retrieval functionalities.

## Features
* **Real-Time Communication**: Utilizes WebSocket for instant messaging between users.
* **Direct Messaging**: Supports private messaging between two users.
* **Group Chats**: Allows users to communicate in groups with real-time message broadcasting.
* **Session Management**: Session-based authentication ensures secure communication.
* **Persistent Storage**: Messages are stored in a database for retrieval and history.

## Project setup
#### Prerequisite
1. Go installed
2. postgreSQL

### Config
1. Create a `config.yaml` file in the root of the project DIR, you can find an example in `config.yaml.example`
2. Setup the config accordingly

### Run
```
git clone https://github.com/iamrk1811/real-time-chat.git
cd real-time-chat
go mod download
make all
```

## To be added
0. Swagger
1. API to create user and groups
2. API to add users to group
3. Efficient logging
4. Caching can be implemented for active groups member