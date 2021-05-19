// Package iotserver provides server interface for ESP8266 devices
package iotserver

import (
	"bufio"
	"io"
	"log"
	"net"
	"strconv"
)

const bufferSize int = 100

type iotDevice struct {
	id         string
	lightValue string
	serverCh   chan string // Server to IOT device
	clientCh   chan string // IOT device to server
	conn       net.Conn
	r          io.Reader
	w          io.Writer
	buff       []byte
}

func NewIotDevice(conn net.Conn, upStrem chan string) *iotDevice {
	return &iotDevice{
		conn:     conn,
		serverCh: upStrem,
		clientCh: make(chan string, 5),
		r:        bufio.NewReader(conn),
		w:        bufio.NewWriter(conn),
		buff:     make([]byte, bufferSize),
	}
}

func (dev *iotDevice) reader() {
ILOOP:
	for {
		n, err := dev.r.Read(dev.buff)
		dev.buff[n] = 0
		data := dev.buff[:n+1]

		switch err {

		case io.EOF:
			break ILOOP

		case nil:
			dev.serverCh <- string(data)
			//log.Printf("Receive: %s", data)

		default:
			log.Fatalf("Receive data failed:%s", err)
			return
		}

	}

	log.Println("Break from reader loop, dev id:", dev.id)
}

type IotServer struct {
	port        int
	cons        []*iotDevice
	devUpstream chan string
}

func NewIotServer(port int) *IotServer {
	return &IotServer{
		port:        port,
		devUpstream: make(chan string, 5),
	}
}

func (srv *IotServer) Start() error {

	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(srv.port))

	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", srv.port, err)
	}

	defer listen.Close()

	log.Printf("Begin listen port: %d", srv.port)

	go func() {
		for {
			select {
			case msg := <-srv.devUpstream:
				log.Println(msg)
			}
		}
	}()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go srv.handler(conn)
	}

}

func (srv *IotServer) handler(conn net.Conn) {
	device := NewIotDevice(conn, srv.devUpstream)
	//go device.writer()
	go device.reader()
	srv.cons = append(srv.cons, device)
}
