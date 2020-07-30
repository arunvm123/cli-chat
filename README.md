#cli-chat

    This is a project I build to learn a bit about grpc. I was planning to build a chat room where the client communicated with the server via grpc and was looking at some cli based UI in Golang to implement the client with. Thats when I came across this [repo](https://github.com/Luqqk/go-cli-chat). The UI was exactly what I had in mind, so I decided to use it but with a grpc implementation.

    When first running the client, the user has to enter their name after which the client recieves a stream through which the server communicates. It is a one way stream and only the server can stream data to the client. To send messages from the client,there is a separate broadcast function. The server recieves the message and then broadcasts it to all open streams. When the client leaves the stream,it is closed and no more data can be send on it. The UI is pretty simple and it was built with [gocui](https://github.com/jroimartin/gocui).

## How to run?

    Run the server first

    ```
    go run cmd/server/main.go
    ```

    and then the client

    ```
    go run cmd/client/main.go
    ```

    and it should be working
