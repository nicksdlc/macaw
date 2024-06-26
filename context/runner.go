package context

import (
	"log"

	"github.com/nicksdlc/macaw/communicators"
	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/responder"
	"github.com/nicksdlc/macaw/sender"
)

var inconsistentRunnerMessage = "Inconsistent runner type. Want %s, got %s"
var receiverName = "receiver"
var senderName = "sender"

var runners = make(map[string]runner)

type runner func(communicator communicators.Communicator, cfg config.Configuration)

func init() {
	runners[receiverName] = runReceiver
	runners[senderName] = runSender
}

func get(communicator communicators.Communicator, cfg config.Configuration) runner {
	r, ok := runners[cfg.Mode]
	if !ok {
		log.Panicf("Not recognized mode")
	}

	return r
}

func runSender(communicator communicators.Communicator, cfg config.Configuration) {
	if cfg.Mode == senderName {
		sender := sender.NewMessageSender(communicator, cfg.Requests)
		resp, err := sender.SendWithResponse()
		if err != nil {
			log.Fatalf("Error sending message: %s", err)
		}
		for r := range resp {
			log.Printf(" [runner] Received response: %s", r.Body)
		}
	} else {
		log.Fatalf(inconsistentRunnerMessage, senderName, cfg.Mode)
	}
}

func runReceiver(communicator communicators.Communicator, cfg config.Configuration) {
	if cfg.Mode == receiverName {
		listener := responder.NewMessageResponder(communicator, cfg.Responses)
		listener.Listen()

		var forever chan struct{}

		log.Printf(" [runner] Waiting for messages. To exit press CTRL+C")
		<-forever
	} else {
		log.Fatalf(inconsistentRunnerMessage, receiverName, cfg.Mode)
	}
}
