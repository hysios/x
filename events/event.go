package events

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/hysios/x/events/common"
)

func NewMessage(payload interface{}) *common.Message {
	b, _ := json.Marshal(payload)
	return message.NewMessage(watermill.NewUUID(), b)
}
