package client

import (
	"crypto/rand"
	"fmt"
	"log"
	"net"
)

type Driver struct {
	Addr               string
	NumberOfConnection int
	AsynchronousMode   bool
	PayloadCh          chan Communicator
}

type Communicator struct {
	payload  []byte
	response chan []byte
	reqID    string
}

func NewCommunicator(payload []byte, response chan []byte, reqID string) Communicator {
	return Communicator{
		payload:  payload,
		response: response,
		reqID:    reqID,
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
	reqStore       map[string]chan []byte
	net.Conn
}

func NewSingleConnection(communicatorCh chan Communicator, addr string) (*SingleConnection, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &SingleConnection{
		reqStore:       make(map[string]chan []byte),
		communicatorCh: communicatorCh,
		Conn:           conn,
	}, nil
}

func (s *SingleConnection) Worker() {
	for payload := range s.communicatorCh {
		s.reqStore[payload.reqID] = payload.response

		_, err := s.Conn.Write(payload.payload)
		if err != nil {
			log.Println(err)
			continue
		}

		//
	}
}

func (d *Driver) SetReq(key, value []byte, ttl int) (<-chan []byte, error) {
	newResponse, reqId := make(chan []byte), GenerateReqID()

	payload, err := SetDriver(key, value, reqId, ttl)
	if err != nil {
		return nil, err
	}

	d.PayloadCh <- NewCommunicator(payload, newResponse, string(reqId))

	return newResponse, nil
}

func GenerateReqID() []byte {
	id := make([]byte, 12)
	_, err := rand.Read(id)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate random ID: %v", err))
	}

	return id
}
