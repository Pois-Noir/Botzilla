package core

import "encoding/json"

type packetType int

const (
	MessagePacket packetType = iota
	BroadcastPacket
	BadPacket
)

type Packet struct {
	Header map[string]string
	Body   string
}

func (p Packet) Encode() ([]byte,error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return data, nil 
}

func Decode(input []byte) (*Packet, error) {
	var p *Packet
	err := json.Unmarshal(input, p)
	if err != nil {
		return nil, err
	}
	return p, nil

}
