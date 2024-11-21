# Botzilla Client

`botzillaclient` is a library designed to facilitate communication with Botzilla. For more information about Botzilla itself, refer to the official documentation: Botzilla Documentation.

`Note`: botzillaclient requires Botzilla to be actively running on the network to function properly.

# Functions

```
Start(name string, port string) error
```
- `Description`: Establishes a connection to the Botzilla server.

- `Parameters`:
	- `name`: A unique identifier for each component connecting to the server.

	- `port`: The TCP port on which the server will listen.

- `Returns`: An error if the connection fails, otherwise nil.

- `Usage`: This function must be called by every component that wishes to connect to the Botzilla server. It opens a TCP server on the specified port.

<br />

```
SendCommand(follower string, body string) (string, error)
```
- `Description`: Sends a command to a specified follower and waits for a response.

- `Parameters`:
	- `follower`: The unique identifier of the follower to which the command is being sent.

	- `body`: The body of the command to send.
- `Returns`: The response from the follower and any potential error.

- `Usage`: This function is used to send commands to a specific follower and receive their responses. It is designed to handle one-way communication from the sender to the follower.


<br />

```
BroadcastMessage(followers []string, body string) error
```
- `Description`: Sends a message to multiple followers.

- `Parameters`:
	- `followers`: A list of follower identifiers to which the message will be sent.

	- `body`: The body of the message.
- `Returns`: An error if the message could not be delivered to any of the followers.
- `Usage`: This function broadcasts a message to multiple followers. It is useful for sending the same message to a group of followers.

<br />

## Receiving Commands and Messages

`CommandListener`

To handle incoming commands, you must implement the CommandListener interface:

```
type CommandListener interface {
	OnReceive(body string) (string, error)
}
```
- `Description`: This interface is used to process incoming commands. The OnReceive function is invoked each time a new command is received.

- `Returns`: When a command is received and processed via the CommandListener, the return value of the OnReceive function is automatically sent back to the sender.

- `Execution`: OnReceive is executed in a new goroutine for each incoming command, allowing for concurrent handling of multiple commands.

<br />

`MessageListener`

To handle incoming messages, you must implement the MessageListener interface:

```
type MessageListener interface {
	OnReceive(body string)
}
```
- `Description`: This interface is used to process incoming messages. The OnReceive function is invoked each time a new message is received.

- `Execution`: OnReceive is executed in a new goroutine for each incoming message.

### Setting Up Listeners
Once you've implemented the required interfaces, you can register listeners for commands and messages using the following functions:

`Start Command Listener`
```
func StartCommandListener( listener CommandListener )
```

<br />

`Start Message Listener`
```
func StartMessageListener( listener MessageListener )
```

## Example Usage
```
// Define a command listener
type MyCommandListener struct{}

func (m *MyCommandListener) OnReceive(body string) (string, error) {
    // Process the command
    return "Command processed", nil
}

// Define a message listener
type MyMessageListener struct{}

func (m *MyMessageListener) OnReceive(body string) error {
    // Process the message
    fmt.Println("Received message:", body)
    return nil
}

func main() {
    // Start the server
    err := Start("MyComponent", "8080")
    if err != nil {
        log.Fatal(err)
    }

    // Register listeners
    err = RegisterCommandListener(&MyCommandListener{})
    if err != nil {
        log.Fatal(err)
    }

    err = RegisterMessageListener(&MyMessageListener{})
    if err != nil {
        log.Fatal(err)
    }

    // Send a command to a follower
    response, err := SendCommand("Follower1", "Some command")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Received response:", response)

    // Broadcast a message to multiple followers
    err = BroadcastMessage([]string{"Follower1", "Follower2"}, "Hello followers!")
    if err != nil {
        log.Fatal(err)
    }
}
```