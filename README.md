# net-cat >^.^<

- Recreate NetCat in server mode

- group chat: TCP connection between server and multiple clients 
 - start TCP server, listen and accept connections
 - listens on a specified port/ defaults 8989
 - max 10 connections
 - broadcast msg from client (but not empty msg)
  - log timestamp and username
  - [2020-01-20 15:48:41][client.name]:[client.message]
 - sent previous msg when new client joins
 - notify all clients when other clients join or leave

- client
 - name is required for client

must have:
- go concurrency:
 - go routines
 - channels
- mutexes
- must respect good practices
- handle errors from server side and client side

recommended:
- unit testing

Bonus:
- terminal UI
- change client name
- logs client activites into a file
- more netcet flags implemented