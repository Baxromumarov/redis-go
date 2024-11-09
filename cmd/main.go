package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	dtb "github.com/baxromumarov/redis-go/database"
	er "github.com/baxromumarov/redis-go/errors"
)

type Connection struct {
	conn              net.Conn
	connType          string
	activeConnections int
	maxConnections    int
	totalConnections  int
}

func main() {
	fmt.Println()
	db := dtb.Init()

	// In every two seconds, the cleaner will check for expired keys
	// it will work in background
	db.StartExpirationCleaner(2 * time.Second)

	listener, err := net.Listen("tcp", ":6060")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Redis-go server listening on port 6060...")

	ch := make(chan error)
	defer close(ch)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, db, &ch)
		val := <-ch
		switch val {

		case er.ErrConnClosed:
			fmt.Println("Connection closed")
		default:
			fmt.Print("Error reading message:", val)
		}

	}
}

func handleConnection(conn net.Conn, db *dtb.Database, ch *chan error) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	fmt.Println(conn.RemoteAddr().Network())
	fmt.Println(conn.RemoteAddr().String())
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {

				*ch <- er.ErrConnClosed
				return
			}
			*ch <- err

		}

		message = strings.TrimSpace(message)
		response, err := dtb.HandleCommand(db, message)
		if err != nil {
			conn.Write([]byte("ERROR: " + err.Error() + "\n"))
			continue
		}

		conn.Write([]byte(response + "\n"))
	}
}
