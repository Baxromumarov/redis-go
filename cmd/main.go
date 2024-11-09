package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	dtb "github.com/baxromumarov/redis-go/database"
)

const (
	QUIT = iota
)

func main() {
	fmt.Println()
	db := dtb.Init()

	listener, err := net.Listen("tcp", ":6060")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 6379...")

	ch := make(chan int)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, db, &ch)
		val := <-ch
		fmt.Println(">>>>>>> ", val)
		if val == QUIT {
			return
		}
	}
}

func handleConnection(conn net.Conn, db *dtb.Database, ch *chan int) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("here")
				*ch <- QUIT
				return
			}
			fmt.Println("Error reading message:", err)
			return
		}

		message = strings.TrimSpace(message)
		response, err := dtb.HandleCommand(db, message)
		if err != nil {
			conn.Write([]byte("ERROR" + err.Error() + "\n"))
			continue
		}

		conn.Write([]byte(response + "\n"))
	}
}
