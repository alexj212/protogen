package main

import "C"
import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/emicklei/proto"
)

func Parse(file string, enumName, fieldPrefix string) (*MessageMapper, error) {
	reader, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()

	if err != nil {
		return nil, err
	}

	mapper := &MessageMapper{}
	mapper.PacketEnum = enumName
	mapper.ProtoFile = file
	mapper.ProtoGenVer = DATE
	mapper.CommandLine = strings.Join(os.Args, " ")
	mapper.ParserName = os.Args[0]
	currentTime := time.Now()
	mapper.Date = currentTime.Format("2006.01.02 15:04:05")

	if strings.Contains(mapper.ParserName, "/") {
		mapper.ParserName = mapper.ParserName[strings.LastIndex(mapper.ParserName, "/")+1:]
	}

	for _, each := range definition.Elements {
		switch v := each.(type) {
		case *proto.Package:
			mapper.PackageName = v.Name
		}
	}

	proto.Walk(definition,
		proto.WithEnum(func(enum *proto.Enum) {
			handleEnum(mapper, enum, fieldPrefix)
		}),
		proto.WithOption(func(enum *proto.Option) {
			handleOption(mapper, enum)
		}),
		proto.WithService(func(service *proto.Service) {
			handleService(mapper, service)
		}),
		proto.WithMessage(func(message *proto.Message) {
			handleMessage(mapper, message)
		}))

	if mapper.GoPackageName == "" {
		mapper.GoPackageName = mapper.PackageName
	}

	if strings.Contains(mapper.GoPackageName, "/") {
		mapper.GoPackageName = mapper.GoPackageName[strings.LastIndex(mapper.GoPackageName, "/")+1:]
	}

	return mapper, nil
}

func handleOption(mapper *MessageMapper, s *proto.Option) {
	// fmt.Printf("handleOption: %v   / %v\n", s.Name, s.Constant.Source)
	if s.Name == "go_package" {
		mapper.GoPackageName = s.Constant.Source
	}
}

func handleEnum(mapper *MessageMapper, s *proto.Enum, fieldPrefix string) {
	// fmt.Printf("handleEnum: %v\n", s.Name)

	for _, each := range s.Elements {

		enumField, ok := each.(*proto.EnumField)
		if ok && enumField.InlineComment != nil{
			// fmt.Printf("handleEnum[%v]: %v  %v\n", i, enumField.Name, enumField.Integer, )
			messageName := messageNameExtractor(enumField.InlineComment.Message())

			if messageName != "" {
				// fmt.Printf("    %v\n", enumField.InlineComment.Message() )

				packetField := &Packet{}
				packetField.PacketId = fmt.Sprintf("%v_%v", fieldPrefix, enumField.Name)
				packetField.PacketId = strings.ToUpper(string(packetField.PacketId[0])) + packetField.PacketId[1:]

				packetField.PacketName = fixupPacketName(messageName)
				mapper.EventList = append(mapper.EventList, packetField)
			}
		}
	}
}

func handleService(mapper *MessageMapper, s *proto.Service) {
	// fmt.Printf("handleService: %v\n", s.Name)
}

func handleMessage(mapper *MessageMapper, m *proto.Message) {
	// fmt.Printf("handleMessage: %v\n",m.Name)
}

func fixupPacketName(s string) string {

	pieces := strings.Split(s, "_")
	var name string
	for _, piece := range pieces {
		name = name + strings.Title(piece)
	}
	return name
}

func messageNameExtractor(comment string) string {
	// comment := "  sddafdasfa   @@protogen:pkt_server_registration@@  // asfasdfasd"
	re := regexp.MustCompile("@@protogen:(.*?)@@")
	matches := re.FindStringSubmatch(comment)

	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
