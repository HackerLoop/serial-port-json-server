package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/HackerLoop/rotonde/shared"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

type cmdDefinition struct {
	name       string
	definition rotonde.Definition
}

func (c *cmdDefinition) toCmd(params map[string]interface{}) []byte {
	buf := bytes.NewBufferString(c.name)
	for _, field := range c.definition.Fields {
		value, ok := params[field.Name]
		if ok == false {
			log.Println("Missing field", field.Name)
			continue
		}
		buf.WriteString(" ")
		switch v := value.(type) {
		case string:
			buf.WriteString(v)
		case float64:
			buf.WriteString(strconv.FormatInt(int64(v), 10))
		case bool:
			buf.WriteString(strconv.FormatBool(v))
		}
	}
	return buf.Bytes()
}

func (c *cmdDefinition) toAction(cmd []byte) *rotonde.Action {
	s := string(cmd[:])
	args := strings.Split(strings.TrimSpace(s), " ")

	result := make(map[string]interface{})
	for i, field := range c.definition.Fields {
		if i >= len(args)-1 {
			break
		}
		value := args[i+1]
		switch field.Type {
		case "string":
			result[field.Name] = value
		case "number":
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			result[field.Name] = v
		case "boolean":
			v, err := strconv.ParseBool(value)
			if err != nil {
				log.Println(err)
				continue
			}
			result[field.Name] = v
		}
	}
	action := rotonde.Action{c.definition.Identifier, result}
	return &action
}

type cmdDefinitions []*cmdDefinition

func (definitions cmdDefinitions) getCmdForActionIdentifier(identifier string) (*cmdDefinition, error) {
	for _, definition := range definitions {
		if definition.definition.Identifier == identifier {
			return definition, nil
		}
	}
	return nil, errors.New(fmt.Sprint(identifier, " Not found"))
}

func (definitions cmdDefinitions) getCmdForName(name string) (*cmdDefinition, error) {
	for _, definition := range definitions {
		if definition.name == name {
			return definition, nil
		}
	}
	return nil, errors.New(fmt.Sprint(name, " Not found"))
}

func addCmd(name string, definition rotonde.Definition) {
	cmd := cmdDefinition{name, definition}
	cmds = append(cmds, &cmd)
}

var cmds cmdDefinitions

func init() {
	cmds = make([]*cmdDefinition, 0, 10)

	list := &rotonde.Definition{Identifier: "SERIAL_LIST", Type: "action"}
	addCmd("list", *list)

	open := &rotonde.Definition{Identifier: "SERIAL_OPEN", Type: "action"}
	open.PushField("device", "string", "")
	open.PushField("baudrate", "number", "")
	open.PushField("buffer", "string", "")
	addCmd("open", *open)

	sendJson := &rotonde.Definition{Identifier: "SERIAL_SENDJSON", Type: "action"}
	sendJson.PushField("json", "string", "")
	addCmd("sendjson", *sendJson)

	send := &rotonde.Definition{Identifier: "SERIAL_SEND", Type: "action"}
	send.PushField("port", "string", "")
	send.PushField("data", "string", "")
	addCmd("send", *send)

	sendNoBuf := &rotonde.Definition{Identifier: "SERIAL_SENDNOBUF", Type: "action"}
	sendNoBuf.PushField("port", "string", "")
	sendNoBuf.PushField("data", "string", "")
	addCmd("sendnobuf", *sendNoBuf)

	clos := &rotonde.Definition{Identifier: "SERIAL_CLOSE", Type: "action"}
	clos.PushField("port", "string", "")
	addCmd("close", *clos)

	bufAlgo := &rotonde.Definition{Identifier: "SERIAL_BUFFERALGORITHMS", Type: "action"}
	bufAlgo.PushField("port", "string", "")
	addCmd("bufferalgorithms", *bufAlgo)

	baudRates := &rotonde.Definition{Identifier: "SERIAL_BAUDRATES", Type: "action"}
	addCmd("baudrates", *baudRates)

	restart := &rotonde.Definition{Identifier: "SERIAL_RESTART", Type: "action"}
	addCmd("restart", *restart)

	exit := &rotonde.Definition{Identifier: "SERIAL_EXIT", Type: "action"}
	addCmd("exit", *exit)

	fro := &rotonde.Definition{Identifier: "SERIAL_FRO", Type: "action"}
	fro.PushField("port", "string", "")
	fro.PushField("mult", "number", "")
	addCmd("fro", *fro)

	memstats := &rotonde.Definition{Identifier: "SERIAL_MEMSTATS", Type: "action"}
	addCmd("memstats", *memstats)

	version := &rotonde.Definition{Identifier: "SERIAL_VERSION", Type: "action"}
	addCmd("version", *version)

	hostname := &rotonde.Definition{Identifier: "SERIAL_HOSTNAME", Type: "action"}
	addCmd("hostname", *hostname)

	program := &rotonde.Definition{Identifier: "SERIAL_PROGRAM", Type: "action"}
	program.PushField("port", "string", "")
	program.PushField("arch", "string", "")
	program.PushField("path", "string", "")
	addCmd("program", *program)

	programFromUrl := &rotonde.Definition{Identifier: "SERIAL_PROGRAM", Type: "action"}
	programFromUrl.PushField("port", "string", "")
	programFromUrl.PushField("arch", "string", "")
	programFromUrl.PushField("url", "string", "")
	addCmd("programfromurl", *programFromUrl)
}

/**
 * rotonde connection sent to hub
 */

type rotondeConnection struct {
	connection
}

func (c *rotondeConnection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		p := rotonde.Packet{}
		err = json.Unmarshal(message, &p)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if p.Type == "action" {
			action := rotonde.Action{}
			mapstructure.Decode(p.Payload, &action)

			if strings.HasPrefix(action.Identifier, "SERIAL_") == false {
				log.Println("actions should start with SERIAL_")
				continue
			}

			cmd, err := cmds.getCmdForActionIdentifier(action.Identifier)
			if err != nil {
				log.Println(err)
				continue
			}
			h.broadcast <- cmd.toCmd(action.Data)
		}

	}
	c.ws.Close()
}

func (c *rotondeConnection) writer() {
	for message := range c.send {
		s := string(message[:])
		args := strings.Split(strings.TrimSpace(s), " ")

		if len(args) == 0 {
			log.Println("zero-length command")
			continue
		}

		cmdName := args[0]
		cmd, err := cmds.getCmdForName(cmdName)
		if err != nil {
			log.Println(err)
			continue
		}

		action := cmd.toAction(message)
		packet := rotonde.Packet{"action", action}
		json, err := json.Marshal(packet)
		if err != nil {
			log.Println(err)
			continue
		}

		err = c.ws.WriteMessage(websocket.TextMessage, json)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

func startRotondeClient(serverUrl string) {
	log.Println("startRotondeClient")
	u, err := url.Parse(serverUrl)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := net.Dial("tcp", serverUrl)
		if err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}
		ws, response, err := websocket.NewClient(conn, u, http.Header{}, 10000, 10000)
		if err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}
		log.Println(response)
		c := &rotondeConnection{connection{send: make(chan []byte, 256*10), ws: ws}}
		h.register <- &c.connection
		defer func() { h.unregister <- &c.connection }()
		go c.writer()
		c.reader()
	}
}
