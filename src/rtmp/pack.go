package rtmp

import "bytes"
import "encoding/binary"

func GetString(value string) []byte {

    bytesBuffer := bytes.NewBuffer([]byte{})

    binary.Write(bytesBuffer, binary.BigEndian, int8(AMF_STRING))
    binary.Write(bytesBuffer, binary.BigEndian, int16(len(value)))
    binary.Write(bytesBuffer, binary.BigEndian, []byte(value))

    return bytesBuffer.Bytes() 
}

func GetStringAsObjectName(value string) []byte {

    bytesBuffer := bytes.NewBuffer([]byte{})

    binary.Write(bytesBuffer, binary.BigEndian, int16(len(value)))
    binary.Write(bytesBuffer, binary.BigEndian, []byte(value))

    return bytesBuffer.Bytes() 
}

func GetBoolean(value bool) []byte {
    bytesBuffer := bytes.NewBuffer([]byte{})

    binary.Write(bytesBuffer, binary.BigEndian, int8(AMF_BOOLEAN))
    if value {
        binary.Write(bytesBuffer, binary.BigEndian, byte(1))
    } else {
        binary.Write(bytesBuffer, binary.BigEndian, byte(0))
    }

    return bytesBuffer.Bytes() 
    
}

func GetNumber(value int64) []byte {
    bytesBuffer := bytes.NewBuffer([]byte{})

    binary.Write(bytesBuffer, binary.BigEndian, int8(AMF_NUMBER))
    binary.Write(bytesBuffer, binary.LittleEndian, value)

    return bytesBuffer.Bytes()     
}

func GetNull() []byte {

    return append([]byte{}, 0X05)
}

func GetObjectBegin() []byte {

    return append([]byte{}, AMF_OBJECT)
}

func GetObjectEnd() []byte {

    return append([]byte{}, 0X00, 0X00, 0X09)
}
