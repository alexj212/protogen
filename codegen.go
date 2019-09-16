package main

import (
	"bytes"
	"github.com/pkg/errors"
	"go/format"
	"text/template"
)

type Packet struct {
	PacketName string
	PacketId   string
}

type MessageMapper struct {
	Date          string
	ParserName    string
	ProtoFile     string
	PacketEnum    string
	PackageName   string
	GoPackageName string
	EventList     []*Packet
}

func Generate(d *MessageMapper, formatCode bool) ([]byte, error) {
	t := template.Must(template.New("mapping").Parse(messageMapperTemplate))

	var tpl bytes.Buffer
	var err error

	if err = t.Execute(&tpl, d); err != nil {
		return nil, errors.Wrap(err, "unable to execute template")
	}

	// var config = printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}

	content := tpl.Bytes()

	if formatCode {
		content, err = format.Source(tpl.Bytes())
		if err != nil {
			return nil, errors.Wrap(err, "unable to format source")
		}
	}

	return content, err
}

var messageMapperTemplate = `
// ----------------------------------------------------------------------------        
// ----------------------------------------------------------------------------
// --- THIS FILE IS PROGRAMMATICALLY GENERATED DO NO EDIT----------------------
// --- ALL EDITS WILL BE LOST - YOU HAVE BEEN WARNED---------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------        
// File Generated: {{.Date}}
// Source: {{.ProtoFile}}
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

package {{.GoPackageName}}

import (
    "encoding/binary"
    "errors"

    "github.com/golang/protobuf/proto"
    "github.com/potakhov/loge"
)


var Mapper ProtobufMessageMapper = &mapper{}

type ProtobufMessageMapper interface {
    MapIDToProto(method uint32) (interface{}, error)
    MapProtoMessageToID(msg interface{}) (uint32, error)
}

type mapper struct{}

func (mapper) MapIDToProto(method uint32) (interface{}, error)     { return MapIDToProto(method) }
func (mapper) MapProtoMessageToID(msg interface{}) (uint32, error) { return MapProtoMessageToID(msg) }


// MapIDToProto maps all possible packet IDs to their corresponding packet types
func MapIDToProto(method uint32) (interface{}, error) {
    switch Packet(method) {

{{range $val := .EventList}}
    case {{$val.PacketId}}:
       return &{{$val.PacketName}}{}, nil
{{end}}

    }
    return nil, errors.New("unknown protocol method received")
}


// MapIDToProto maps all possible packet IDs to their corresponding packet types
func MapProtoMessageToID(msg interface{}) (uint32, error) {
    switch msg.(type) {

{{range $val := .EventList}}
    case *{{$val.PacketName}}:
       return {{$val.PacketId}}.Value(), nil
{{end}}


    }
    return 0, errors.New("unknown protocol method received")
}

func Parse(data []byte) (interface{}, uint32, error) {
    if len(data) < 4 {
        loge.Error("Receiving invalid packet")
        return nil, 0, errors.New("received invalid packet len")
    }

    idSlice := data[:4]
    packetID := binary.LittleEndian.Uint32(idSlice)

    msg, err := MapIDToProto(packetID)

    if err != nil {
        return nil, 0, err
    }

    err = proto.Unmarshal(data[4:], msg.(proto.Message))

    if err != nil {
        return nil, 0, err
    }

    return msg, packetID, nil
}

func EncodeMessage(param interface{}) ([]byte, uint32, error) {
    if param == nil {
        return nil, 0, errors.New("nil packet is not allowed")
    }

    id, err := MapProtoMessageToID(param)
    if err != nil {
        return nil, 0, err
    }

    serialized, err := proto.Marshal(param.(proto.Message))
    if err != nil {
        return nil, 0, err
    }

    idSlice := make([]byte, 4)
    binary.LittleEndian.PutUint32(idSlice, id)

    payload := append(idSlice, serialized...)

    return payload, id, nil
}

func (x {{.PacketEnum}}) Value() uint32 {
    return uint32(x)
}

`
