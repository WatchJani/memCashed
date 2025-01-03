package client

import (
	"hash"
	"hash/fnv"
	"log"
	"net"

	"github.com/WatchJani/memCashed/client/internal/types"
	p "github.com/WatchJani/memCashed/client/parser"
)

// Connection struct represents a client driver responsible for managing connections.
type Connection struct {
	Addr               string // Address to connect to.
	NumberOfConnection int    // Number of concurrent connections to establish.
	// AsynchronousMode   bool              // Flag indicating whether to use asynchronous mode.
	PayloadCh chan Communicator // Channel used for sending payloads for communication.
}

type Driver struct {
	hash.Hash32
	Conn []Connection
}

// Communicator struct represents a payload and response channel for communication.
type Communicator struct {
	payload  []byte      // Payload data to be sent.
	response chan []byte // Channel to receive the response from the server.
}

// NewCommunicator creates and returns a new Communicator with the specified payload and response channel.
func NewCommunicator(payload []byte, response chan []byte) Communicator {
	return Communicator{
		payload:  payload,
		response: response,
	}
}

func New() (*Driver, error) {
	configuration := types.LoadConfiguration()
	connections := make([]Connection, len(configuration.Server))

	for index, connection := range configuration.Server {
		con, err := NewConnection(connection.IpAddr, connection.NumberOfConnection)
		if err != nil {
			return nil, err
		}

		connections[index] = con
	}

	return &Driver{fnv.New32a(), connections}, nil
}

// NewConnection creates and returns a new Driver instance with the provided address and number of connections.
func NewConnection(addr string, numberConnection int) (Connection, error) {
	d := Connection{
		Addr:               addr,                    // Set address.
		NumberOfConnection: numberConnection,        // Set the number of connections.
		PayloadCh:          make(chan Communicator), // Create a channel for sending payloads.
	}

	return d, d.Init()
}

// Init initializes the Driver by creating a specified number of SingleConnection instances
// and starting the Worker goroutines for each connection.
func (d *Connection) Init() error {
	// Create and initialize each single connection.
	for range d.NumberOfConnection {
		singleConnection, err := NewSingleConnection(d.PayloadCh, d.Addr)
		if err != nil {
			return err // Return error if connection creation fails.
		}

		// Start the Worker goroutine for each connection.
		go singleConnection.Worker()
	}

	return nil // Initialization successful.
}

// SingleConnection represents an individual network connection and its associated communication channel.
type SingleConnection struct {
	communicatorCh chan Communicator // Channel for communicating with the Driver.
	net.Conn                         // The network connection (TCP, etc.).
}

// NewSingleConnection creates and returns a new SingleConnection instance.
func NewSingleConnection(communicatorCh chan Communicator, addr string) (*SingleConnection, error) {
	conn, err := net.Dial("tcp", addr) // Establish a TCP connection to the provided address.
	if err != nil {
		return nil, err // Return error if the connection fails.
	}

	return &SingleConnection{
		communicatorCh: communicatorCh, // Assign the provided communication channel.
		Conn:           conn,           // Assign the established network connection.
	}, nil
}

// Worker listens for incoming payloads from the communicator channel and processes them asynchronously.
func (s *SingleConnection) Worker() {
	readBuffer := make([]byte, 1024*1024+10) // Buffer for receiving data from the server.
	for payload := range s.communicatorCh {  // Loop through incoming payloads.
		// Write the payload to the connection.
		_, err := s.Conn.Write(payload.payload)
		if err != nil {
			log.Println(err) // Log the error if writing fails.
			continue
		}

		// Read the response from the server.
		n, err := s.Conn.Read(readBuffer)
		if err != nil {
			log.Println(err) // Log the error if reading fails.
			continue
		}

		response := make([]byte, n)
		copy(response, readBuffer[:n])
		// Send the received data back through the response channel.
		payload.response <- response
	}
}

// SetReq sends a request to set a key-value pair with a TTL (Time-To-Live) on the server.
func (d *Driver) SetReq(key, value []byte, ttl int) (<-chan []byte, error) {
	n, err := d.Write(key)
	if err != nil {
		return nil, err
	}

	payload, err := p.Set(key, value, ttl)
	return d.OperationReq(payload, n%len(d.Conn), err)
}

// GetReq sends a request to get a value by key from the server
func (d *Driver) GetReq(key []byte) (<-chan []byte, error) {
	n, err := d.Write(key)
	if err != nil {
		return nil, err
	}

	payload, err := p.Get(key)
	return d.OperationReq(payload, n%len(d.Conn), err)
}

// DeleteReq sends a request to delete a key-value pair from the server.
func (d *Driver) DeleteReq(key []byte) (<-chan []byte, error) {
	n, err := d.Write(key)
	if err != nil {
		return nil, err
	}

	payload, err := p.Delete(key)
	return d.OperationReq(payload, n%len(d.Conn), err)
}

// OperationReq sends the payload request to the Driver's PayloadCh and returns a response channel.
func (d *Driver) OperationReq(payload []byte, route int, err error) (<-chan []byte, error) {
	if err != nil {
		return nil, err // Return error if the operation failed.
	}

	// Create a new response channel.
	newResponse := make(chan []byte)

	// Send the payload and response channel to the PayloadCh channel.
	d.Conn[route].PayloadCh <- NewCommunicator(payload, newResponse)

	return newResponse, nil // Return the response channel.
}
