package main

import (
    "crypto/rand"
    "crypto/tls"
    "log"
    "net"
    "crypto/x509"
    "fmt"
    "bufio"
    "os"
)

func main() {
    cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
    if err != nil {
        log.Fatalf("server: loadkeys: %s", err)
    }
    config := tls.Config{Certificates: []tls.Certificate{cert}}
    config.Rand = rand.Reader
    service := "0.0.0.0:8000"
    listener, err := tls.Listen("tcp", service, &config)
    if err != nil {
        log.Fatalf("server: listen: %s", err)
    }
    log.Print("server: listening")
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("server: accept: %s", err)
            break
        }
        defer conn.Close()
        log.Printf("server: accepted from %s", conn.RemoteAddr())
        tlscon, ok := conn.(*tls.Conn)
        if ok {
            log.Print("ok")
            state := tlscon.ConnectionState()
            for _, v := range state.PeerCertificates {
                log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
            }
        }
        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()
    log.Println("server: waiting for messages ... ")
        scanner := bufio.NewScanner(conn)
        onChar := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
                for i := 0; i < len(data); i++ {
                    if data[i] == '>' {
			log.Println(i+1, string(data[:i]))
                        return i + 1, data[:i], nil
                    }
                }
                if !atEOF {
                        return 0, nil, nil
                }
                return 0, data, bufio.ErrFinalToken
        }
        scanner.Split(onChar)
        for scanner.Scan() {
                //fmt.Printf("%s \n", string(scanner.Text()))
        }
        if err := scanner.Err(); err != nil {
                fmt.Fprintln(os.Stderr, "reading input:", err)
        }
        log.Println("server: conn: closed")
}
