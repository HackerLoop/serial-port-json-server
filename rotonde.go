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

/**
 * This translates the SPJS API to the rotonde API
 * SPJS uses a mix of text commands as actions, and JSON objects as events.
 *
 * The idea is to first create definitions for all available actions and events,
 * and then match them with what is going through.
 *
 * For text command actions, we just match a rotonde action to a SPJS command,
 * each fields of the action will be listed as command parameters.
 *
 * For a json event, it is a little more complicated as SPJS objects doesn't provide a type field.
 * so the idea is to match each json event to definitions based on the fields.
 * This is the reason some definitions will not have all the fields, I just added those that were
 * always present. Or just enough to match.
 */

/**
 * This struct is made to match SPJS command names to a rotonde definition
 */
type cmdActionDefinition struct {
	name       string
	definition rotonde.Definition
}

// Generates a cmd from an object received from rotonde
func (c *cmdActionDefinition) toCmd(params map[string]interface{}) []byte {
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

// oops, this one is useless, keeping it just in case..
func (c *cmdActionDefinition) toAction(cmd []byte) *rotonde.Action {
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

/*
 * actions definition store
 */

type cmdActionDefinitions []*cmdActionDefinition

func (definitions cmdActionDefinitions) getCmdForActionIdentifier(identifier string) (*cmdActionDefinition, error) {
	for _, definition := range definitions {
		if definition.definition.Identifier == identifier {
			return definition, nil
		}
	}
	return nil, errors.New(fmt.Sprint(identifier, " Not found"))
}

func (definitions cmdActionDefinitions) getCmdForName(name string) (*cmdActionDefinition, error) {
	for _, definition := range definitions {
		if definition.name == name {
			return definition, nil
		}
	}
	return nil, errors.New(fmt.Sprint(name, " Not found"))
}

func addCmdAction(name string, definition rotonde.Definition) {
	cmdAction := cmdActionDefinition{name, definition}
	cmdActions = append(cmdActions, &cmdAction)
}

/**
 * wraps a definition, adds a isDefinitionFor method which matches based on fields.
 */

type serialEvent struct {
	rotonde.Definition
}

func (s *serialEvent) isDefinitionFor(obj map[string]interface{}) bool {
	for _, field := range s.Fields {
		if _, ok := obj[field.Name]; ok == false {
			return false
		}
	}
	return true
}

func addEvent(definition *rotonde.Definition) {
	s := serialEvent{*definition}
	events = append(events, &s)
}

var cmdActions cmdActionDefinitions
var events []*serialEvent

/**
 * Initializes all events and actions that will go through
 */
func init() {

	// actions

	cmdActions = make([]*cmdActionDefinition, 0, 10)

	list := &rotonde.Definition{Identifier: "SERIAL_LIST", Type: "action"}
	addCmdAction("list", *list)

	open := &rotonde.Definition{Identifier: "SERIAL_OPEN", Type: "action"}
	open.PushField("port", "string", "")
	open.PushField("baudrate", "number", "")
	open.PushField("buffer", "string", "")
	addCmdAction("open", *open)

	sendJson := &rotonde.Definition{Identifier: "SERIAL_SENDJSON", Type: "action"}
	sendJson.PushField("json", "string", "")
	addCmdAction("sendjson", *sendJson)

	send := &rotonde.Definition{Identifier: "SERIAL_SEND", Type: "action"}
	send.PushField("port", "string", "")
	send.PushField("data", "string", "")
	addCmdAction("send", *send)

	sendNoBuf := &rotonde.Definition{Identifier: "SERIAL_SENDNOBUF", Type: "action"}
	sendNoBuf.PushField("port", "string", "")
	sendNoBuf.PushField("data", "string", "")
	addCmdAction("sendnobuf", *sendNoBuf)

	clos := &rotonde.Definition{Identifier: "SERIAL_CLOSE", Type: "action"}
	clos.PushField("port", "string", "")
	addCmdAction("close", *clos)

	bufAlgo := &rotonde.Definition{Identifier: "SERIAL_BUFFERALGORITHMS", Type: "action"}
	bufAlgo.PushField("port", "string", "")
	addCmdAction("bufferalgorithms", *bufAlgo)

	baudRates := &rotonde.Definition{Identifier: "SERIAL_BAUDRATES", Type: "action"}
	addCmdAction("baudrates", *baudRates)

	restart := &rotonde.Definition{Identifier: "SERIAL_RESTART", Type: "action"}
	addCmdAction("restart", *restart)

	exit := &rotonde.Definition{Identifier: "SERIAL_EXIT", Type: "action"}
	addCmdAction("exit", *exit)

	fro := &rotonde.Definition{Identifier: "SERIAL_FRO", Type: "action"}
	fro.PushField("port", "string", "")
	fro.PushField("mult", "number", "")
	addCmdAction("fro", *fro)

	memstats := &rotonde.Definition{Identifier: "SERIAL_MEMSTATS", Type: "action"}
	addCmdAction("memstats", *memstats)

	version := &rotonde.Definition{Identifier: "SERIAL_VERSION", Type: "action"}
	addCmdAction("version", *version)

	hostname := &rotonde.Definition{Identifier: "SERIAL_HOSTNAME", Type: "action"}
	addCmdAction("hostname", *hostname)

	program := &rotonde.Definition{Identifier: "SERIAL_PROGRAM", Type: "action"}
	program.PushField("port", "string", "")
	program.PushField("arch", "string", "")
	program.PushField("path", "string", "")
	addCmdAction("program", *program)

	programFromUrl := &rotonde.Definition{Identifier: "SERIAL_PROGRAMFROMURL", Type: "action"}
	programFromUrl.PushField("port", "string", "")
	programFromUrl.PushField("arch", "string", "")
	programFromUrl.PushField("url", "string", "")
	addCmdAction("programfromurl", *programFromUrl)

	// events

	events = make([]*serialEvent, 0, 10)

	commands := &rotonde.Definition{Identifier: "SERIAL_COMMANDS", Type: "event"}
	commands.PushField("Commands", "", "")
	addEvent(commands)

	err := &rotonde.Definition{Identifier: "SERIAL_ERROR", Type: "event"}
	err.PushField("Error", "", "")
	addEvent(err)

	bufFlowDebug := &rotonde.Definition{Identifier: "SERIAL_BUFFLOWDEBUG", Type: "event"}
	bufFlowDebug.PushField("BufFlowDebug", "", "")
	addEvent(bufFlowDebug)

	memoryStats := &rotonde.Definition{Identifier: "SERIAL_MEMSTATS", Type: "event"}
	memoryStats.PushField("Alloc", "number", "")
	memoryStats.PushField("TotalAlloc", "number", "")
	addEvent(memoryStats)

	hostnameE := &rotonde.Definition{Identifier: "SERIAL_HOSTNAME", Type: "event"}
	hostnameE.PushField("Hostname", "", "")
	addEvent(hostnameE)

	versionE := &rotonde.Definition{Identifier: "SERIAL_VERSION", Type: "event"}
	versionE.PushField("Version", "", "")
	addEvent(versionE)

	gc := &rotonde.Definition{Identifier: "SERIAL_GC", Type: "event"}
	gc.PushField("gc", "", "")
	addEvent(gc)

	exitE := &rotonde.Definition{Identifier: "SERIAL_EXIT", Type: "event"}
	exitE.PushField("Exiting", "", "")
	addEvent(exitE)

	restartE := &rotonde.Definition{Identifier: "SERIAL_RESTART", Type: "event"}
	restartE.PushField("Restarting", "", "")
	addEvent(restartE)

	restartedE := &rotonde.Definition{Identifier: "SERIAL_RESTARTED", Type: "event"}
	restartedE.PushField("Restarted", "", "")
	addEvent(restartedE)

	broadcast := &rotonde.Definition{Identifier: "SERIAL_BROADCAST", Type: "event"}
	broadcast.PushField("Cmd", "", "")
	broadcast.PushField("Msg", "", "")
	addEvent(broadcast)

	programmerStatus := &rotonde.Definition{Identifier: "SERIAL_PROGRAMMERSTATUS", Type: "event"}
	programmerStatus.PushField("ProgrammerStatus", "", "")
	programmerStatus.PushField("Url", "", "")
	addEvent(programmerStatus)

	cmdComplete := &rotonde.Definition{Identifier: "SERIAL_CMDCOMPLETE", Type: "event"}
	cmdComplete.PushField("Cmd", "", "")
	cmdComplete.PushField("Id", "", "")
	cmdComplete.PushField("P", "", "")
	cmdComplete.PushField("BufSize", "", "")
	cmdComplete.PushField("D", "", "")
	addEvent(cmdComplete)

	portMessage := &rotonde.Definition{Identifier: "SERIAL_PORTMESSAGE", Type: "event"}
	portMessage.PushField("P", "", "")
	portMessage.PushField("D", "", "")
	addEvent(portMessage)

	output := &rotonde.Definition{Identifier: "SERIAL_OUTPUT", Type: "event"}
	output.PushField("Cmd", "", "")
	output.PushField("Desc", "", "")
	output.PushField("Port", "", "")
	addEvent(output)

	portList := &rotonde.Definition{Identifier: "SERIAL_PORTLIST", Type: "event"}
	portList.PushField("SerialPorts", "", "")
	addEvent(portList)
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

			cmdAction, err := cmdActions.getCmdForActionIdentifier(action.Identifier)
			if err != nil {
				log.Println(err)
				continue
			}
			h.broadcast <- cmdAction.toCmd(action.Data)
		}

	}
	c.ws.Close()
}

func (c *rotondeConnection) writer() {
loop:
	for message := range c.send {

		e := make(map[string]interface{})
		err := json.Unmarshal(message, &e)
		if err != nil {
			log.Println("Unmarshal error", err, string(message))
		}

		for _, event := range events {
			if event.isDefinitionFor(e) {
				ee := rotonde.Event{event.Identifier, e}
				packet := rotonde.Packet{"event", ee}
				json, err := json.Marshal(packet)
				if err != nil {
					log.Println(err)
					continue
				}

				err = c.ws.WriteMessage(websocket.TextMessage, json)
				if err != nil {
					log.Println(err)
				}
				continue loop
			}
		}
		log.Println("unknown event", string(message))

	}
	c.ws.Close()
}

func sendDefinition(ws *websocket.Conn, definition *rotonde.Definition) {
	packet := rotonde.Packet{"def", definition}

	json, err := json.Marshal(packet)
	if err != nil {
		log.Println(err)
		return
	}

	err = ws.WriteMessage(websocket.TextMessage, json)
	if err != nil {
		log.Println(err)
	}
}

func sendDefinitions(ws *websocket.Conn) {
	for _, cmdAction := range cmdActions {
		sendDefinition(ws, &cmdAction.definition)
	}

	for _, event := range events {
		sendDefinition(ws, &event.Definition)
	}

}

func startRotondeClient(serverUrl string) {
	log.Println("startRotondeClient")
	u, err := url.Parse(serverUrl)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := net.Dial("tcp", u.Host)
		if err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}
		ws, response, err := websocket.NewClient(conn, u, http.Header{}, 10000, 10000)
		if err != nil {
			log.Println(err)
			log.Println(response)
			time.Sleep(2 * time.Second)
			continue
		}
		sendDefinitions(ws)
		c := &rotondeConnection{connection{send: make(chan []byte, 256*10), ws: ws}}
		h.register <- &c.connection
		defer func() { h.unregister <- &c.connection }()
		go c.writer()
		c.reader()
	}
}
