##ProtoGen

###### ProtoGen is a utility to generate mapping functions that easily allow encoding a ProtoBuf structure into a []byte.


The definition of a protocol file should contain an enum with definitions of packet ids. The enum value names should map to various messages defined in the protofile. Value names should match proto messages with the suffix 'Id'. Also the field name should have a comment `//@@export@@`

A proto file should look like the following   

```proto 
file: packets.proto


enum Packet {
    SrvReturnCodeId = 0; //@@export@@
}

message SrvReturnCode {
    uint32 event = 1;
    uint32 correlationId = 2;
}

```




```bash
usage: protogen protofile enumNamePacket output_go_file

example: go run protogen /home/alexj/projects/rooms/submodules/protocol-definitions/centralserver/css1erver.proto Packet /home/alexj/projects/rooms/proto/cs/mapping.go

```

The following go code will be written to the output file
```go 
var Mapper ProtobufMessageMapper = &mapper{}

type ProtobufMessageMapper interface {
	MapIDToProto(method uint32) (interface{}, error)
	MapProtoMessageToID(msg interface{}) (uint32, error)
}

type mapper struct{}

func (mapper) MapIDToProto(method uint32) (interface{}, error)     { return MapIDToProto(method) }
func (mapper) MapProtoMessageToID(msg interface{}) (uint32, error) { return MapProtoMessageToID(msg) }

func MapIDToProto(method uint32) (interface{}, error) {...}
func MapProtoMessageToID(msg interface{}) (uint32, error) {...} 
func Parse(data []byte) (interface{}, uint32, error) {...} 
func EncodeMessage(param interface{}) ([]byte, uint32, error) {...} 
func (x Packet) Value() uint32 {...} 
}
```

  

## Building
```bash 
go build -o ~/bin/protogen .

or 

make protogen
```   
