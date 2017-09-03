package rtmp

import "net"
import "encoding/binary"
import "os"

func GetRTMPClient(conn net.Conn) *RTMPClient {
    client := &RTMPClient{
        conn : conn,
        chunklen : DEFAULT_CHUNK_LEN,
    };

    file, err := os.OpenFile("111.264", os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        println(err.Error())
        return nil
    }
    client.file264 = file

    file, err = os.OpenFile("111.aac", os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        println(err.Error())
        return nil
    }
    client.fileaac = file

    return client
}

func (self *RTMPClient)Handshake() error {

    buf := make([]byte, 1536)
    buf[0] = 3
    println("c0", buf[0])
    self.conn.Write(buf[:1])
    println("c1", buf[8])
    self.conn.Write(buf)

    err := recvbuffer(self.conn, buf[:1])
    if err != nil {
        return err
    }
    println("s0 version", buf[0])

    var s1buf [1536]byte
    err = recvbuffer(self.conn, s1buf[:])
    if err != nil {
        return err
    }
    println("s1", s1buf[8])
    
    err = recvbuffer(self.conn, buf)
    if err != nil {
        return err
    }
    println("s2",buf[8])

    println("c2",s1buf[8])
    self.conn.Write(s1buf[:])

    return nil
}

func (self *RTMPClient)Connect(url string) error {

    bodybuf := append([]byte{},GetString("connect")...)
    bodybuf = append(bodybuf, GetNumber(1)...)

    bodybuf = append(bodybuf, GetObjectBegin()...)
    bodybuf = append(bodybuf, GetStringAsObjectName("app")...)
    bodybuf = append(bodybuf, GetString("live")...)

    bodybuf = append(bodybuf, GetStringAsObjectName("tcUrl")...)
    bodybuf = append(bodybuf, GetString(url)...)

    bodybuf = append(bodybuf, GetStringAsObjectName("fpad")...)
    bodybuf = append(bodybuf, GetBoolean(false)...)

    bodybuf = append(bodybuf, GetStringAsObjectName("capabilities")...)
    bodybuf = append(bodybuf, GetNumber(DEFAULT_CAPABILITIES)...)

    bodybuf = append(bodybuf, GetStringAsObjectName("audioCodecs")...)
    bodybuf = append(bodybuf, GetNumber(DEFAULT_AUDIO_CODECS)...)

    bodybuf = append(bodybuf, GetStringAsObjectName("videoCodecs")...)
    bodybuf = append(bodybuf, GetNumber(DEFAULT_VIDEO_CODECS)...)

    bodybuf = append(bodybuf, GetStringAsObjectName("videoFunction")...)
    bodybuf = append(bodybuf, GetNumber(DEFAULT_VIDEO_FUNCTION)...)

    bodybuf = append(bodybuf, GetObjectEnd()...)

    rtmpheader := RTMPHeader{
        chunkfmt : 0x0,
        chunkstreamid : 0x3,
        timestamp : 0x0,
        msg_type : 0x14,
        msg_streamid : 0x0,
    }

    return self.SendRTMPMessage(&rtmpheader, bodybuf)
}

func (self *RTMPClient) CreateStream() error {

    bodybuf := append([]byte{},GetString("createStream")...)
    bodybuf = append(bodybuf, GetNumber(2)...)

    bodybuf = append(bodybuf, GetNull()...)

    rtmpheader := RTMPHeader{
        chunkfmt : 0x0,
        chunkstreamid : 0x7,
        timestamp : 0x0,
        msg_type : 0x14,
        msg_streamid : 0,
    }

    return self.SendRTMPMessage(&rtmpheader, bodybuf)
}

func (self *RTMPClient) Play(streamid string, msgstreamid uint32) error {

    bodybuf := append([]byte{},GetString("play")...)
    bodybuf = append(bodybuf, GetNumber(3)...)

    bodybuf = append(bodybuf, GetNull()...)
    bodybuf = append(bodybuf,GetString(streamid)...)
//    bodybuf = append(bodybuf, GetNumber(0)...)

    rtmpheader := RTMPHeader{
        chunkfmt : 0x0,
        chunkstreamid : 0x8,
        timestamp : 0x0,
        msg_type : 0x14,
        msg_streamid : msgstreamid,
    }

    return self.SendRTMPMessage(&rtmpheader, bodybuf)
}

func (self *RTMPClient) SetWindowAcknowledgementSize() error {

    var intbuf [4]byte
    binary.BigEndian.PutUint32(intbuf[:], uint32(2500000))
    bodybuf := append([]byte{}, intbuf[:]...)

    rtmpheader := RTMPHeader{
        chunkfmt : 0x0,
        chunkstreamid : 0x2,
        timestamp : 0x0,
        msg_type : 0x5,
        msg_streamid : 1,//1337,
    }

    return self.SendRTMPMessage(&rtmpheader, bodybuf)
}
