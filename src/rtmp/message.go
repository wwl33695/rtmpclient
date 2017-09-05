package rtmp

import "encoding/binary"

func (self *RTMPClient) SendRTMPMessage(rtmpheader *RTMPHeader, bodybuf []byte) error {

    header := RTMPChunkHeader {
        flags : (rtmpheader.chunkfmt<<6)|rtmpheader.chunkstreamid,
        msg_type : rtmpheader.msg_type,
    }
    msglen := uint32(len(bodybuf))

    var intbuf [4]byte
    binary.BigEndian.PutUint32(intbuf[:], uint32(msglen))
    copy(header.msg_len[:], intbuf[1:])
    binary.LittleEndian.PutUint32(intbuf[:], uint32(rtmpheader.msg_streamid))
    copy(header.msg_streamid[:], intbuf[:])
    binary.BigEndian.PutUint32(intbuf[:], uint32(rtmpheader.timestamp))
    copy(header.timestamp[:], intbuf[1:])
//    println("msglen=", msglen, intbuf[0], intbuf[1], intbuf[2], intbuf[3])

    if msglen <= self.chunklen {

        self.buffer.Reset()
        binary.Write(&self.buffer, binary.BigEndian, &header)
        binary.Write(&self.buffer, binary.BigEndian, bodybuf)

//        buf := self.buffer.Bytes()
//        println("connectlen", self.buffer.Len(), buf[0],buf[1],buf[2],buf[3],buf[4],buf[5],buf[6],buf[7],buf[8],buf[9],buf[10],buf[11])
        _, err := self.conn.Write(self.buffer.Bytes())
        if err != nil {
            return err
        }        
    } else {

        self.buffer.Reset()
        binary.Write(&self.buffer, binary.BigEndian, &header)
        binary.Write(&self.buffer, binary.BigEndian, bodybuf[:self.chunklen])

        _, err := self.conn.Write(self.buffer.Bytes())
        if err != nil {
            return err
        }        

        pos := self.chunklen
        for msglen > pos {

            self.buffer.Reset()
            binary.Write(&self.buffer, binary.BigEndian, uint8(0xC3))
            if pos + self.chunklen > msglen {
                binary.Write(&self.buffer, binary.BigEndian, bodybuf[pos:msglen])
            } else {
                binary.Write(&self.buffer, binary.BigEndian, bodybuf[pos:pos+self.chunklen])
            }

            _, err := self.conn.Write(self.buffer.Bytes())
            if err != nil {
                return err
            }        

            pos += self.chunklen
        }
    }

    return nil
}

func (self *RTMPClient) RecvMessage() error{

    buf := make([]byte, self.chunklen * 2)
    err := recvbuffer(self.conn, buf[:1])
    if err != nil {
        return err
    }
    flag := buf[0]
//    println("flag = ", flag, uint8(flag)>>6, uint8(flag)&0x3f)

    self.messages[flag&0x3f].header.chunkfmt = flag>>6
    self.messages[flag&0x3f].header.chunkstreamid = flag&0x3f

    headerlen := getHeaderLength(flag)
    err = recvbuffer(self.conn, buf[:headerlen-1])
    if err != nil {
        return err
    }
//    println("1", uint8(flag)>>6, flag&0x3f, headerlen)
    if headerlen >= 4 {
        var intbuf [4]byte
        copy(intbuf[1:], buf[0:3])
        timestamp := binary.BigEndian.Uint32(intbuf[:])

        self.messages[flag&0x3f].header.timestamp = timestamp
    }
    if headerlen >= 8 {
        var intbuf [4]byte
        copy(intbuf[1:], buf[3:6])
        msglen := binary.BigEndian.Uint32(intbuf[:])

        self.messages[flag&0x3f].header.msg_len = msglen
        self.messages[flag&0x3f].header.msg_type = buf[6]
//        println("8", uint8(flag)>>6, headerlen, msglen)
    }
    if headerlen >= 12 {

        streamid := binary.LittleEndian.Uint32(buf[7:11])
        self.messages[flag&0x3f].header.msg_streamid = streamid
//        println("12", uint8(flag)>>6, headerlen, msglen)
//        println("12", uint8(flag)>>6, headerlen, buf[0], buf[1], buf[2], buf[3], buf[4], buf[5], buf[6] )
    }

    readlen := self.messages[flag&0x3f].header.msg_len - uint32(self.messages[flag&0x3f].buf.Len())
    if readlen > self.chunklen {
        readlen = self.chunklen
    }

//    println("readlen", readlen)
    err = recvbuffer(self.conn, buf[:readlen])
    if err != nil {
        return err
    }

    self.messages[flag&0x3f].buf.Write(buf[:readlen])
    if uint32(self.messages[flag&0x3f].buf.Len()) == self.messages[flag&0x3f].header.msg_len {
        self.HandleMessage(&self.messages[flag&0x3f])
        self.messages[flag&0x3f].buf.Reset()
    }

    return nil
}

func (self *RTMPClient) HandleMessage(msg *RTMPMessage) error{

//    println("HandleMessage", msg.header.msg_type, msg.header.chunkstreamid, msg.header.msg_len)
        //,msg.header.timestamp, msg.header.msg_streamid)

    if msg.header.msg_type == MSG_INVOKE {
        self.HandleInvoke(msg)
    } else if msg.header.msg_type == MSG_NOTIFY {
        self.HandleNotify(msg)
    } else if msg.header.msg_type == MSG_SET_CHUNK {
        self.HandleSetChunkSize(msg)
    } else if msg.header.msg_type == MSG_VIDEO {
        self.HandleVideoData(msg)
    } else if msg.header.msg_type == MSG_AUDIO {
        self.HandleAudioData(msg)
    }

    return nil
}

func (self *RTMPClient) HandleInvoke(msg *RTMPMessage) error{


//    println("HandleInvoke", msg.header.msg_streamid, msg.header.chunkstreamid)

    var buffer [2048]byte
    binary.Read(&msg.buf, binary.BigEndian, &buffer[0])

    var length uint16
    binary.Read(&msg.buf, binary.BigEndian, &length)
    binary.Read(&msg.buf, binary.BigEndian, buffer[:length])
    method := string(buffer[:length])

    binary.Read(&msg.buf, binary.BigEndian, &buffer[0])
    var transactionid uint64 = 0
    binary.Read(&msg.buf, binary.LittleEndian, &transactionid)
    println("HandleInvoke method", method, "transactionid", uint64(transactionid), 
                msg.header.chunkstreamid, msg.header.msg_streamid)

    if method == "_result" {

        if self.times == 0 {
            println("connect result")
            self.SetWindowAcknowledgementSize()

            self.CreateStream()            
        } else if self.times == 1 {

            binary.Read(&msg.buf, binary.BigEndian, &buffer[0])
//            println("CreateStream result", buffer[0], msg.buf.Len())
            binary.Read(&msg.buf, binary.BigEndian, &buffer[0])
            var streamid float64 = 0
            binary.Read(&msg.buf, binary.BigEndian, &streamid)

            println("CreateStream result", buffer[0], uint64(streamid), msg.buf.Len())

//            self.Play("test", uint32(streamid))    
            self.Play("hks", uint32(streamid))    
        }
        self.times++
    }

    return nil
}

func (self *RTMPClient) HandleNotify(msg *RTMPMessage) error{

//    println("HandleInvoke", msg.header.msg_streamid, msg.header.chunkstreamid)

    var buffer [2048]byte
    binary.Read(&msg.buf, binary.BigEndian, &buffer[0])

    var length uint16
    binary.Read(&msg.buf, binary.BigEndian, &length)
    binary.Read(&msg.buf, binary.BigEndian, buffer[:length])
    method := string(buffer[:length])
    println("HandleNotify method", method, msg.buf.Len())

    return nil
}

func (self *RTMPClient) HandleSetChunkSize(msg *RTMPMessage) error{

    var chunklen uint32 = 0
    binary.Read(&msg.buf, binary.BigEndian, &chunklen)
//    println("chunklen", chunklen)
    self.chunklen = chunklen

    return nil
}

func (self *RTMPClient) HandleVideoData(msg *RTMPMessage) error{

    buf := msg.buf.Bytes()

    if buf[1] == 0x01 {

        var i uint32 = 5
        for i < uint32(len(buf)) {
            framelength := binary.BigEndian.Uint32(buf[i:i+uint32(self.avccmediaheaderlen)])
            i+=uint32(self.avccmediaheaderlen)

            self.file264.Write(h264startcode)
            self.file264.Write(buf[i:i+framelength])                

            i += framelength
        }

 //       println("nalu ", buf[9]&0x1f)

    } else if buf[1] == 0x00 {        
        println("key frame", len(buf))

        data := buf[5:]
        profile := data[1]
        levelid := data[3]
        self.avccmediaheaderlen = (data[4]&0x03)+1
        println("self.avccmediaheaderlen", self.avccmediaheaderlen)
        numSps := data[5] & 0x1f
        spsLen := (uint(data[6]))<<8 + uint(data[7])
        sps := data[8 : 8+spsLen]
        idx := 8 + spsLen
        numPps := data[idx]
        println("profile levelid numSps numPps", profile, levelid, numSps, numPps)
        idx++
        ppsLen := (uint(data[idx]))<<8 + uint(data[idx+1])
        idx += 2
        pps := data[idx : idx+ppsLen]
        //------------------
        self.file264.Write(h264startcode)
        self.file264.Write(sps)
        self.file264.Write(h264startcode)
        self.file264.Write(pps)
    }

//    println("nalusize", frametype, nalu, nalusize, msg.buf.Len())

//    println("nalusize", msg.buf.Bytes()[0],msg.buf.Bytes()[1],
//        msg.buf.Bytes()[2],msg.buf.Bytes()[3])


    return nil
}

func (self *RTMPClient) HandleAudioData(msg *RTMPMessage) error{

    buf := msg.buf.Bytes()
    if buf[0] == 0xaf {
        if buf[1] == 0x00 {
            self.adts = NewAdts(buf[2:4])
            println("aac config buf", len(buf))
        } else {
            data := buf[2:]
            blen := uint32(len(data))
            hdr := self.adts.ToAdts(int(blen))
            self.fileaac.Write(hdr)
            self.fileaac.Write(data)
        }
    }

    return nil
}