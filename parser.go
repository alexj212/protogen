package main

import "C"
import (
    "fmt"
    "os"
    "strings"
    "time"

    "github.com/emicklei/proto"
)

func Parse(file string, enumName string) (*MessageMapper, error) {
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
    mapper.ParserName = os.Args[0]
    currentTime := time.Now()
    mapper.Date = currentTime.Format("2006.01.02 15:04:05")

    if strings.Contains(mapper.ParserName, "/") {
        mapper.ParserName = mapper.ParserName[strings.LastIndex(mapper.ParserName, "/")+1:]
    }

    for _, each := range definition.Elements {
        // fmt.Printf("each: %T  %v\n", each, each)
        switch v := each.(type) {
        case *proto.Package:
            mapper.PackageName = v.Name
        }
    }

    proto.Walk(definition,
        proto.WithEnum(func(enum *proto.Enum) {

            if enum.Name == enumName {
                handleEnum(mapper, enum)
            }
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

func handleEnum(mapper *MessageMapper, s *proto.Enum) {
    // fmt.Printf("handleEnum: %v\n", s.Name)

    for _, each := range s.Elements {

        enumField, ok := each.(*proto.EnumField)
        if ok {
            // fmt.Printf("handleEnum[%v]: %v  %v\n", i, enumField.Name, enumField.Integer, )
            if enumField.InlineComment != nil {
                // fmt.Printf("    %v\n", enumField.InlineComment.Message() )

                if strings.Contains(enumField.InlineComment.Message(), "@@export@@") {

                    packetField := &Packet{}
                    packetField.PacketId = fmt.Sprintf("%v_%v", mapper.PacketEnum, enumField.Name)

                    if strings.HasSuffix(enumField.Name, "Id") {
                        packetField.PacketName = enumField.Name[:len(enumField.Name)-2]
                    } else {
                        packetField.PacketName = enumField.Name
                    }

                    mapper.EventList = append(mapper.EventList, packetField)
                }
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
