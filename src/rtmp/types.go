package rtmp

import "net"
import "bytes"
import "os"

type RTMPChunkHeader struct {
    flags uint8
    timestamp [3]byte
    msg_len [3]byte
    msg_type uint8
    msg_streamid [4]byte
};

type RTMPHeader struct {
    chunkfmt uint8
    chunkstreamid uint8
    timestamp uint32
    msg_len uint32
    msg_type uint8
    msg_streamid uint32
};

type RTMPMessage struct {
    header RTMPHeader
        
    buf bytes.Buffer
};

const RANDOM_LEN = 1536 -8
const DEFAULT_CHUNK_LEN = 128
var h264startcode []byte = []byte{0x0,0x0,0x0,0x1}

type RTMPHandshake struct {
    time1 uint32
    time2 uint32
    random [RANDOM_LEN]byte
};

const (
    AMF_NUMBER = 0x00
    AMF_BOOLEAN = 0x01
    AMF_STRING = 0x02
    AMF_OBJECT = 0x03
)

// Chunk stream ID
const (
    CS_ID_PROTOCOL_CONTROL = uint32(2)
    CS_ID_COMMAND          = uint32(3)
    CS_ID_USER_CONTROL     = uint32(4)
)

const (
    EVENT_STREAM_BEGIN       = uint16(0)
    EVENT_STREAM_EOF         = uint16(1)
    EVENT_STREAM_DRY         = uint16(2)
    EVENT_SET_BUFFER_LENGTH  = uint16(3)
    EVENT_STREAM_IS_RECORDED = uint16(4)
    EVENT_PING_REQUEST       = uint16(6)
    EVENT_PING_RESPONSE      = uint16(7)
    EVENT_REQUEST_VERIFY     = uint16(0x1a)
    EVENT_RESPOND_VERIFY     = uint16(0x1b)
    EVENT_BUFFER_EMPTY       = uint16(0x1f)
    EVENT_BUFFER_READY       = uint16(0x20)
)

const (
    BINDWIDTH_LIMIT_HARD    = uint8(0)
    BINDWIDTH_LIMIT_SOFT    = uint8(1)
    BINDWIDTH_LIMIT_DYNAMIC = uint8(2)
)

const (

    MSG_SET_CHUNK = uint8(0x1)
    MSG_BYTES_READ = uint8(0x3)
    MSG_USER_CONTROL = uint8(0x4)
    MSG_RESPONSE = uint8(0x5)
    MSG_REQUEST = uint8(0x6)
    MSG_AUDIO = uint8(0x8)
    MSG_VIDEO = uint8(0x9)
    MSG_INVOKE3 = uint8(0x11)    /* AMF3 */
    MSG_NOTIFY = uint8(0x12)
    MSG_OBJECT = uint8(0x13)
    MSG_INVOKE = uint8(0x14)    /* AMF0 */
    MSG_FLASH_VIDEO = uint8(0x16)
)

const (
    MAX_TIMESTAMP                       = uint32(2000000000)
    AUTO_TIMESTAMP                      = uint32(0XFFFFFFFF)
    DEFAULT_HIGH_PRIORITY_BUFFER_SIZE   = 2048
    DEFAULT_MIDDLE_PRIORITY_BUFFER_SIZE = 128
    DEFAULT_LOW_PRIORITY_BUFFER_SIZE    = 64
    DEFAULT_CHUNK_SIZE                  = uint32(128)
    DEFAULT_WINDOW_SIZE                 = 2500000
    DEFAULT_CAPABILITIES                = int64(0x2e40)
    DEFAULT_AUDIO_CODECS                = int64(0xeea840)
    DEFAULT_VIDEO_CODECS                = int64(0x806f4000)
    DEFAULT_VIDEO_FUNCTION              = int64(0xf03f)
     FMS_CAPBILITIES                     = uint32(255)
    FMS_MODE                            = uint32(2)
    SET_PEER_BANDWIDTH_HARD             = byte(0)
    SET_PEER_BANDWIDTH_SOFT             = byte(1)
    SET_PEER_BANDWIDTH_DYNAMIC          = byte(2)
)

type RTMPClient struct {
    conn net.Conn
    buffer bytes.Buffer

    messages [64]RTMPMessage

    times int
    chunklen uint32
    adts *AdtsInfo

    //debug
    file264 *os.File
    fileaac *os.File
}