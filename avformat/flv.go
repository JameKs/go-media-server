package avformat

import (
	"io"
	"github.com/songshenyi/go-media-server/logger"
	"encoding"
	"bytes"
)


//func (h *FLVHeader)SetAudioFlag(flag bool){
//	f := 0x04 * flag
//	h.Data[4] = (( h.Data[4] & 0x01) |  f)
//}
//
//func (h *FLVHeader)SetVideoFlag(flag bool){
//	f := 0x01 * flag
//	h.Data[4] = (( h.Data[4] & 0x04) |  f)
//}
//
//func (h *FLVHeader)GetAudioFlag()(flag bool){
//	return h.Data[4] & 0x04
//}
//
//func (h *FLVHeader)GetVideoFlag()(flag bool){
//	return h.Data[4] & 0x01
//}

type Marshaler interface {
	encoding.BinaryUnmarshaler
	encoding.BinaryMarshaler
	Size() int
}

type GMSUint8 uint8


type GMSUint24 uint32


type GMSUint64 uint64

type RtmpMessageType GMSUint8

const (
	RtmpMsgAmf0DataMessage RtmpMessageType = 18 // 0x12
	RtmpMsgAmf3DataMessage RtmpMessageType = 15 // 0x0F

	RtmpMsgAudioMessage RtmpMessageType = 8 // 0x08
	RtmpMsgVideoMessage RtmpMessageType = 9 // 0x09
)

type FlvTagTimestamp uint64

type FlvTagUint24 uint32

type FlvMessage struct {
	Tag *FlvTag
	Header *FlvHeader

	MetaData bool
	VideoSequenceHeader bool
	AudioSequenceHeader bool
}

const (
	FlvHeaderSize int = 13
	FlvTagHeaderSize int = 11
	FlvPreTagLenSize int = 4
)



type FlvHeader struct {
	Signature [3]byte
	Version uint8
	EnableAudio bool
	EnableVideo bool
	Offset uint32
}

func ReadFlvHeader(r io.Reader)(h *FlvHeader, err error){
	h = &FlvHeader{}
	var buf bytes.Buffer
	if _, err := io.CopyN(&buf, r, int64(FlvHeaderSize)) ; err != nil {
		logger.Warn("read flv headerfailed")
		return h, err
	}

	data := buf.Bytes()
	copy(h.Signature[:], data)
	if string(h.Signature[:]) != "FLV" {
		logger.Warnf("flv header Signature is wrong, %s", string(h.Signature[:]))
		return h, err
	}

	if h.Version = data[3]; h.Version != 1 {
		logger.Warn("flv header Version invalid")
		return h, err
	}

	AVFlag := data[4]
	h.EnableAudio = (AVFlag & 0x04) != 0
	h.EnableVideo = (AVFlag & 0x01) != 0

	logger.Info("read flv header")
	return
}


type FlvTag struct{
	TagType   RtmpMessageType
	DataSize  uint32
	TimeStamp uint64
	StreamId  uint32
	Payload   []byte
}

func FlvGetDataSize(data []byte)(size uint32, err error){
	size = uint32(data[2]) | uint32(data[1])<<8 | uint32(data[0])<<16
	return
}

func FlvGetTimestamp(data []byte)(size uint32, err error){
	return
}

func FlvGetStreamId(data []byte)(size uint32, err error){
	return
}

func ReadFlvTag(r io.Reader)(tag *FlvTag, err error){
	tag = &FlvTag{}
	var buf bytes.Buffer
	if _, err := io.CopyN(&buf, r, int64(FlvTagHeaderSize)) ; err != nil {
		logger.Warn("read flv headerfailed")
		return tag, err
	}

	data := buf.Bytes()

	tag.TagType = RtmpMessageType(data[0])
	tag.DataSize, err = FlvGetDataSize(data[1:4])
	timeStamp, err := FlvGetTimestamp(data[4:8])
	tag.TimeStamp = uint64(timeStamp)
	tag.StreamId, err = FlvGetStreamId(data[8:10])

	buf.Reset()
	written, err := io.CopyN(&buf, r, int64(tag.DataSize))
	if (written != int64(tag.DataSize) || err != nil) {
		logger.Warn("read flv tag Data failed")
		return tag, err
	}

	tag.Payload = make([]byte, tag.DataSize)
	copy(tag.Payload, buf.Bytes())

	buf.Reset()
	io.CopyN(&buf, r, 4)

	logger.Info("read flv tag")
	return
}


// SoundFormat UB [4]
// Format of SoundData. The following values are defined:
//     0 = Linear PCM, platform endian
//     1 = ADPCM
//     2 = MP3
//     3 = Linear PCM, little endian
//     4 = Nellymoser 16 kHz mono
//     5 = Nellymoser 8 kHz mono
//     6 = Nellymoser
//     7 = G.711 A-law logarithmic PCM
//     8 = G.711 mu-law logarithmic PCM
//     9 = reserved
//     10 = AAC
//     11 = Speex
//     14 = MP3 8 kHz
//     15 = Device-specific sound
// Formats 7, 8, 14, and 15 are reserved.
// AAC is supported in Flash Player 9,0,115,0 and higher.
// Speex is supported in Flash Player 10 and higher.
type RtmpCodecAudio uint8

const (
	RTMPLinearPCMPlatformEndian RtmpCodecAudio = iota
	RTMPADPCM
	RTMPMP3
	RTMPLinearPCMLittleEndian
	RTMPNellymoser16kHzMono
	RTMPNellymoser8kHzMono
	RTMPNellymoser
	RTMPReservedG711AlawLogarithmicPCM
	RTMPReservedG711MuLawLogarithmicPCM
	RTMPReserved
	RTMPAAC
	RTMPSpeex
	RTMPReserved1CodecAudio
	RTMPReserved2CodecAudio
	RTMPReservedMP3_8kHz
	RTMPReservedDeviceSpecificSound
	RTMPReserved3CodecAudio
	RTMPDisabledCodecAudio
)

// AACPacketType IF SoundFormat == 10 UI8
// The following values are defined:
//     0 = AAC sequence header
//     1 = AAC raw
type RtmpAacType uint8

const (
	RTMPAacSequenceHeader RtmpAacType = iota
	RTMPAacRawData
	RTMPAacReserved
)

// E.4.3.1 VIDEODATA
// CodecID UB [4]
// Codec Identifier. The following values are defined:
//     2 = Sorenson H.263
//     3 = Screen video
//     4 = On2 VP6
//     5 = On2 VP6 with alpha channel
//     6 = Screen video version 2
//     7 = AVC
type RtmpCodecVideo uint8

const (
	RTMPReservedCodecVideo RtmpCodecVideo = iota
	RTMPReserved1CodecVideo
	RTMPSorensonH263
	RTMPScreenVideo
	RTMPOn2VP6
	RTMPOn2VP6WithAlphaChannel
	RTMPScreenVideoVersion2
	RTMPAVC
	RTMPDisabledCodecVideo
	RTMPReserved2CodecVideo
)

// E.4.3.1 VIDEODATA
// Frame Type UB [4]
// Type of video frame. The following values are defined:
//     1 = key frame (for AVC, a seekable frame)
//     2 = inter frame (for AVC, a non-seekable frame)
//     3 = disposable inter frame (H.263 only)
//     4 = generated key frame (reserved for server use only)
//     5 = video info/command frame
type RtmpAVCFrame uint8

const (
	RTMPReservedAVCFrame RtmpAVCFrame = iota
	RTMPKeyFrame
	RTMPInterFrame
	RTMPDisposableInterFrame
	RTMPGeneratedKeyFrame
	RTMPVideoInfoFrame
	RTMPReserved1AVCFrame
)

// AVCPacketType IF CodecID == 7 UI8
// The following values are defined:
//     0 = AVC sequence header
//     1 = AVC NALU
//     2 = AVC end of sequence (lower level NALU sequence ender is
//         not required or supported)
type RtmpVideoAVCType uint8

const (
	RTMPSequenceHeader RtmpVideoAVCType = iota
	RTMPNALU
	RTMPSequenceHeaderEOF
	RTMPReservedAVCType
)



func (v *FlvMessage) isVideoSequenceHeader() bool {
	// TODO: FIXME: support other codecs.
	if len(v.Tag.Payload) < 2 {
		return false
	}

	b := v.Tag.Payload

	// sequence header only for h264
	codec := RtmpCodecVideo(b[0] & 0x0f)
	if codec != RTMPAVC {
		return false
	}

	frameType := RtmpAVCFrame((b[0] >> 4) & 0x0f)
	avcPacketType := RtmpVideoAVCType(b[1])
	return frameType == RTMPKeyFrame && avcPacketType == RTMPSequenceHeader
}

func (v *FlvMessage) isAudioSequenceHeader() bool {
	// TODO: FIXME: support other codecs.
	if len(v.Tag.Payload) < 2 {
		return false
	}

	b := v.Tag.Payload

	soundFormat := RtmpCodecAudio((b[0] >> 4) & 0x0f)
	if soundFormat != RTMPAAC {
		return false
	}

	aacPacketType := RtmpAacType(b[1])
	return aacPacketType == RTMPAacSequenceHeader
}

func (v RtmpMessageType) isAudio() bool {
	return v == RtmpMsgAudioMessage
}

func (v RtmpMessageType) isVideo() bool {
	return v == RtmpMsgVideoMessage
}

func (v RtmpMessageType) isData() bool {
	return v.isAmf0Data() || v.isAmf3Data()
}

func (v RtmpMessageType) isAmf0Data() bool {
	return v == RtmpMsgAmf0DataMessage
}

func (v RtmpMessageType) isAmf3Data() bool {
	return v == RtmpMsgAmf3DataMessage
}
