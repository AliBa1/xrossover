# xrossover
3D Multiplayer Basketball Video Game (hopefully)

## Current State
As of now I am using raylib to create the game and Golang on the client and server. My goal now is to have a simple state where clients can connect to the server, move their box, and see other clients on the server. I implemented FlatBuffers to handle serialization of data being sent between the server and client. Now it has a ball that updates it's position over the server.

## Run Server
- cd /server
- make
- ./xrossover-server [username] [port]

## Run Client
- cd /client
- make
- ./xrossover-client [username] [port]