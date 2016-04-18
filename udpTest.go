package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"
)

type udpServer struct {
	conn *net.UDPConn
}

type netEvent struct {
	x, y int8
}

func main() {
	udpConnection := initUDP()

	netEventChan := make(chan netEvent)
	go udpConnection.udpReadLoop(netEventChan)

	ticker := time.Tick(1 * time.Second)

	for {
		select {
		case nEvent := <-netEventChan:
			handleNetEvent(&nEvent)
		case <-ticker:
			udpConnection.sendData()
		}
	}
}

func (server *udpServer) sendData() {

	b := []byte("Hello, UDP world!")
	server.conn.Write(b)

}

func handleNetEvent(e *netEvent) {
	log.Print("Recieved X = ", e.x, "   Y = ", e.y)
}

func initUDP() udpServer {
	log.Printf("Opening UDP server...")

	ListenAddr, err := net.ResolveUDPAddr("udp", "192.168.0.199:8000")
	if err != nil {
		log.Printf("S Error: ", err)
	}
	log.Print("Listening on    ", ListenAddr)

	RemoteAddr, err := net.ResolveUDPAddr("udp", "192.168.0.199:6000")
	if err != nil {
		log.Printf("S Error: ", err)
	}
	log.Print("Sending to      ", RemoteAddr)

	conn, err := net.DialUDP("udp", ListenAddr, RemoteAddr)
	if err != nil {
		log.Printf("D Error: ", err)
	}
	log.Print("Connection OK!")

	server := udpServer{conn: conn}
	return server
}

func (server *udpServer) udpReadLoop(c chan netEvent) {

	input := make([]byte, 2, 2)

	for {
		for i, _ := range input {
			input[i] = 0
		}

		nBytes, err := server.conn.Read(input)

		if nBytes != 2 {
			log.Printf("Recieved unexpected ammount of values")
			continue
		}

		byteReader := bytes.NewReader(input)

		event := new(netEvent)

		binary.Read(byteReader, binary.LittleEndian, &event.x)
		err = binary.Read(byteReader, binary.LittleEndian, &event.y)

		if err != nil {
			log.Fatal("binary.Read failed:", err)
		}
		/*
			        for i, _ := range input {
					    err = binary.Read(byteReader, binary.LittleEndian, &event.val[i])

			            if err != nil {
						    log.Fatal("binary.Read failed:", err)
					    }
			        }
		*/

		//log.Print("I have a net event: %+v", event)

		c <- *event
	}
}
