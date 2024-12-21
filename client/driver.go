package client

import (
	"log"
	"net"
)

// Driver struct represents a client driver responsible for managing connections.
type Driver struct {
	Addr               string
	NumberOfConnection int
	AsynchronousMode   bool
	PayloadCh          chan Communicator
}

type Communicator struct {
	payload  []byte
	response chan []byte
}

func NewCommunicator(payload []byte, response chan []byte) Communicator {
	return Communicator{
		payload:  payload,
		response: response,
	}
}

func New(addr string, numberConnection int) *Driver {
	return &Driver{
		Addr:               addr,
		NumberOfConnection: numberConnection,
		AsynchronousMode:   true,
		PayloadCh:          make(chan Communicator),
	}
}

func (d *Driver) Init() error {
	communicator := d.PayloadCh
	for range d.NumberOfConnection {
		singleConnection, err := NewSingleConnection(communicator, d.Addr)
		if err != nil {
			return err
		}

		go singleConnection.Worker()
	}

	return nil
}

type SingleConnection struct {
	communicatorCh chan Communicator
	net.Conn
}

func NewSingleConnection(communicatorCh chan Communicator, addr string) (*SingleConnection, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &SingleConnection{
		communicatorCh: communicatorCh,
		Conn:           conn,
	}, nil
}

// AsynchronousMode
func (s *SingleConnection) Worker() {
	readBuffer := make([]byte, 1024*1024+10)
	for payload := range s.communicatorCh {
		_, err := s.Conn.Write(payload.payload)
		if err != nil {
			log.Println(err)
			continue
		}

		n, err := s.Conn.Read(readBuffer)
		if err != nil {
			log.Println(err)
			continue
		}

		payload.response <- readBuffer[:n]
	}
}

type Operation func(byte, []byte, []byte, int) ([]byte, error)

func (d *Driver) SetReq(key, value []byte, ttl int) (<-chan []byte, error) {
	return d.OperationReq(Set(key, value, ttl))
}

func (d *Driver) GetReq(key []byte) (<-chan []byte, error) {
	return d.OperationReq(Get(key))
}

func (d *Driver) DeleteReq(key []byte) (<-chan []byte, error) {
	return d.OperationReq(Delete(key))
}

func (d *Driver) OperationReq(payload []byte, err error) (<-chan []byte, error) {
	if err != nil {
		return nil, err
	}

	newResponse := make(chan []byte)

	d.PayloadCh <- NewCommunicator(payload, newResponse)

	return newResponse, nil

}
