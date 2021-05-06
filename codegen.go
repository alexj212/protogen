package main

import (
	"bytes"
	"go/format"
	"text/template"

	"github.com/pkg/errors"
)

type Packet struct {
	PacketName string
	PacketId   string
}

type MessageMapper struct {
	Date          string
	ParserName    string
	ProtoFile     string
	CommandLine   string
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
// Source        : {{.ProtoFile}}
// Command Line  : {{.CommandLine}}
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

package {{.GoPackageName}}

import (
    "encoding/binary"
    
    "google.golang.org/protobuf/proto"
	"github.com/pkg/errors"
    "github.com/potakhov/loge"
)

// Codec is the exported codec used to marshall a []byte to and from a proto.Message
var Codec ProtobufCodec = &codec{}

// ProtobufCodec is an interface of functions used to marshall a []byte to and from a proto.Message
type ProtobufCodec interface {

	// MapIDToProto maps all possible packet IDs to their corresponding packet types struct
    MapIDToProto(method uint32) (proto.Message, error)

	// MapProtoMessageToID maps all possible packet IDs to their corresponding packet types
    MapProtoMessageToID(msg proto.Message) (uint32, error)

	// EncodeMessage takes a proto.Message and will return a []byte, packetID and error. error is nil if no error encountered in conversion.
    EncodeMessage(param proto.Message) ([]byte, uint32, error)

	// Parse takes a []byte and return a mapped proto.Message, packetID and error. error is nil if no error encountered in conversion
    Parse(data []byte) (proto.Message, uint32, error)
}

type codec struct{}

// MapIDToProto maps all possible packet IDs to their corresponding packet types struct
func (codec) MapIDToProto(method uint32) (proto.Message, error)     { return MapIDToProto(method) }

// MapProtoMessageToID maps all possible packet IDs to their corresponding packet types
func (codec) MapProtoMessageToID(msg proto.Message) (uint32, error) { return MapProtoMessageToID(msg) }

// EncodeMessage takes a proto.Message and will return a []byte, packetID and error. error is nil if no error encountered in conversion.
func (codec) EncodeMessage(param proto.Message) ([]byte, uint32, error) { return EncodeMessage(param) }

// Parse takes a []byte and return a mapped proto.Message, packetID and error. error is nil if no error encountered in conversion
func (codec) Parse(data []byte) (proto.Message, uint32, error) { return Parse(data) }


// MapIDToProto maps all possible packet IDs to their corresponding packet types
func MapIDToProto(method uint32) (proto.Message, error) {
    switch {{.PacketEnum}}(method) {

{{range $val := .EventList}}
    case {{$val.PacketId}}:
       return &{{$val.PacketName}}{}, nil
{{end}}

    }
    return nil, errors.Errorf("unknown protocol method received: %v [%b]", method, method)
}


// MapProtoMessageToID maps all possible packet IDs to their corresponding packet types
func MapProtoMessageToID(msg proto.Message) (uint32, error) {
    switch msg.(type) {

{{range $val := .EventList}}
    case *{{$val.PacketName}}:
       return {{$val.PacketId}}.Value(), nil
{{end}}


    }
	return 0, errors.Errorf("unknown protocol method received msg: %T", msg) 
}

// Parse takes a []byte and return a mapped proto.Message, packetID and error. error is nil if no error encountered in conversion
func Parse(data []byte) (proto.Message, uint32, error) {
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

// EncodeMessage take a proto.Message and will return a []byte, packetID and error. error is nil if no error encountered in conversion.
func EncodeMessage(param proto.Message) ([]byte, uint32, error) {
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

// Value return the packet id as a uint32 value
func (x {{.PacketEnum}}) Value() uint32 {
    return uint32(x)
}

`
