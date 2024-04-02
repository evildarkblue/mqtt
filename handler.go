package mqtt

import (
	"log"
	"strings"

	"github.com/94peter/mqtt/config"
	"github.com/94peter/mqtt/trans"

	"github.com/eclipse/paho.golang/paho"
)

// handler is a simple struct that provides a function to be called when a message is received. The message is parsed
// and the count followed by the raw message is written to the file (this makes it easier to sort the file)

type Handler interface {
	Close()
	handle(msg *paho.PublishReceived)
}

type handler struct {
	trans  map[string]trans.Trans
	logger *log.Logger
}

// NewHandler creates a new output handler and opens the output file (if applicable)
func NewHandler(conf *config.Config, tsmap map[string]trans.Trans) Handler {
	return &handler{
		trans:  tsmap,
		logger: conf.Logger,
	}
}

// Close closes the file
func (o *handler) Close() {
	for _, t := range o.trans {
		t.Close()
	}
}

func (o *handler) sendMsg(t trans.Trans, topic string, data []byte) {
	err := t.Send(topic, data)
	if err != nil {
		o.println("send fail: " + err.Error())
		o.println(string(data))
	}
}

// handle is called when a message is received
func (o *handler) handle(msg *paho.PublishReceived) {
	for k, t := range o.trans {
		if k == msg.Packet.Topic {
			o.sendMsg(t, msg.Packet.Topic, msg.Packet.Payload)
		} else if strings.Contains(k, "#") {
			prefix := strings.Split(k, "#")[0]
			if strings.HasPrefix(msg.Packet.Topic, prefix) {
				o.sendMsg(t, msg.Packet.Topic, msg.Packet.Payload)
			}
		}
	}
}

func (o *handler) println(a ...any) {
	if o.logger == nil {
		return
	}
	o.logger.Println(a...)
}
