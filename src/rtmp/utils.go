package rtmp

import "net"
import "errors"

func recvbuffer(conn net.Conn, buffer []byte) error {

    pos := 0
    for pos < len(buffer) {
        c, err := conn.Read(buffer[pos:])
        if err != nil {
            return err
        }
        if c <= 0 {
            return errors.New("Read data count <= 0")
        }
        pos += c
    }
    return nil
}

func getHeaderLength(flag byte) int{
    var HEADER_LENGTH []int = []int{12, 8, 4, 1};
    return HEADER_LENGTH[flag>>6]
}