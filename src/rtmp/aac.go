package rtmp

//AdtsInfo an adts header info,to construct adts header
//profileID profile identifier
//sampleRateIdx to index the sample rate
//chanNum channel number
type AdtsInfo struct {
	profileID     uint8
	sampleRateIdx uint8
	chanNum       uint8
}

/*
var sampleRates []byte = {96000,
    88200, 64000, 48000, 44100, 32000,24000,
    22050, 16000, 12000, 11025, 8000, 7350}
*/

//NewAdts from config to adts info
//config  the config bytes array contains the adts info
func NewAdts(config []byte) *AdtsInfo {
	c := config
	//objType := (c[0]>>3)&0xff  //5 Bit
	sampleIdx := (c[0]&0x7)<<1 | c[1]>>7 //4 Bit
	chanNum := (c[1] >> 3) & 0xf         //4 Bit
	return &AdtsInfo{profileID: 1, sampleRateIdx: sampleIdx,
		chanNum: chanNum}
}

//ToAdts fromo AdtsInfo to bytes
//size the aac frame length
func (info *AdtsInfo) ToAdts(size int) []byte {
	b := NewBStreamWriter(7)

	b.WriteBits(0xfff, 12)
	b.WriteBits(0, 1)
	b.WriteBits(0, 2)
	b.WriteBits(1, 1)
	b.WriteBits((uint64)(info.profileID), 2)
	b.WriteBits((uint64)(info.sampleRateIdx), 4)
	b.WriteBits(0, 1)
	b.WriteBits((uint64)(info.chanNum), 3)
	b.WriteBits(0, 1)
	b.WriteBits(0, 1)
	b.WriteBits(0, 1)
	b.WriteBits(0, 1)
	b.WriteBits((uint64)(size+7), 13)
	b.WriteBits(0x7ff, 11)
	b.WriteBits(0, 2)
	return b.Bytes()
}
