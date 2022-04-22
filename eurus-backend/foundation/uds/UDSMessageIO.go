package uds

import (
	"encoding/binary"
	"encoding/json"
	"eurus-backend/foundation"
	"eurus-backend/foundation/ws/ws_message"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

func ResponseUDSErrorMessage(conn net.Conn, code int64, message string, nonce string, logger *logrus.Logger) {

	res := new(ws_message.ResponseMessage)
	res.ReturnCode = code
	res.Message = message
	res.Nonce = nonce
	res.Timestamp = time.Now().Unix()

	SendUDSJsonMessage(conn, res, logger)
}

func SendUDSJsonMessage(conn net.Conn, data interface{}, logger *logrus.Logger) error {
	resData, err := json.Marshal(data)
	if err != nil {
		logger.Errorln("Unable to marshal response message: ", err)
		return err
	}
	var header []byte = make([]byte, 8)
	header[0] = 0xB
	header[1] = 0xa
	header[2] = 0xb
	header[3] = 0xd

	binary.BigEndian.PutUint32(header[4:], uint32(len(resData)))
	_, err = conn.Write(header)
	if err != nil {
		return err
	}

	_, err = conn.Write(resData)
	return err
}

func readConnection(conn net.Conn, messageHandler chan (*MessageFrame), isServerConnection bool, logger *logrus.Logger) error {
	var maxLen uint32 = 1024
	var buffer []byte = make([]byte, maxLen)
	var writePos uint32 = 0

	var totalLen uint32 = 0

	for {
		count, err := conn.Read(buffer[writePos:])

		if count > 0 {
			totalLen = totalLen + uint32(count)

			for {
				frameHeaderPos, frameLen := parseFrame(conn, buffer, totalLen, messageHandler, isServerConnection, logger)

				if frameLen > 0 {
					remainLen := totalLen - (frameHeaderPos + frameLen)

					if remainLen > 0 {
						copy(buffer, buffer[frameHeaderPos+frameLen:frameHeaderPos+frameLen+remainLen])
						totalLen = remainLen
						writePos = remainLen

						continue
					} else {
						totalLen = 0
						writePos = 0
					}
					count = 0
				}
				break
			}

			writePos += uint32(count)
		}

		if err != nil {
			messageHandler <- &MessageFrame{
				Conn:                 conn,
				Frame:                nil,
				Error:                err,
				MasterRequestMessage: nil,
				ResponseMessage:      nil}
			return err
		}

		if writePos >= uint32(float64(maxLen)*0.8) {
			oldBuffer := buffer
			maxLen *= 2
			buffer = make([]byte, maxLen)
			copy(buffer, oldBuffer)
		}

	}

}

func parseFrame(conn net.Conn, buffer []byte, totalLen uint32, messageHandler chan (*MessageFrame), isServerConnection bool, logger *logrus.Logger) (uint32, uint32) {
	var i uint32
	for i = 0; i < totalLen-8; i++ {
		if buffer[i] == 0xB && buffer[i+1] == 0xA && buffer[i+2] == 0xB && buffer[i+3] == 0xD { //Parse the header

			frameLen := binary.BigEndian.Uint32(buffer[i+4:])

			if totalLen-i-8 >= frameLen {
				//Frame found
				if isServerConnection {
					req := new(ws_message.MasterRequestMessage)
					err := json.Unmarshal(buffer[i+8:uint32(i+8)+frameLen], &req)
					if err != nil {
						messageHandler <- &MessageFrame{
							Conn:                 conn,
							Frame:                buffer[i+8 : uint32(i+8)+frameLen],
							Error:                err,
							MasterRequestMessage: nil,
						}
						ResponseUDSErrorMessage(conn, int64(foundation.RequestMalformat), foundation.RequestMalformat.String(),
							"", logger)
					} else {
						messageHandler <- &MessageFrame{
							Conn:                 conn,
							Frame:                buffer[i+8 : uint32(i+8)+frameLen],
							Error:                nil,
							MasterRequestMessage: req,
						}
					}
				} else {
					res := new(ws_message.ResponseMessage)

					err := res.UnmarshalJSON(buffer[i+8 : uint32(i+8)+frameLen])
					if err != nil {
						messageHandler <- &MessageFrame{
							Conn:  conn,
							Frame: buffer[i+8 : uint32(i+8)+frameLen],
							Error: err,
						}
					} else {

						messageHandler <- &MessageFrame{
							Conn:            conn,
							Frame:           buffer[i+8 : uint32(i+8)+frameLen],
							Error:           nil,
							ResponseMessage: res,
						}
					}
				}
				return i, frameLen + 8
			}
		}
	}

	return 0, 0
}
