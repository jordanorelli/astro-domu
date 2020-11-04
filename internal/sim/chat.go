package sim

import (
	"github.com/jordanorelli/astro-domu/internal/wire"
)

type SendChatMessage struct {
	Text string `json:"text"`
}

func (SendChatMessage) NetTag() string { return "chat/send-msg" }

func (m *SendChatMessage) exec(w *world, r *room, p *player, seq int) result {
	for _, p2 := range r.players {
		p2.outbox <- wire.Response{Body: ChatMessage{From: p.name, Text: m.Text}}
	}
	return result{reply: wire.OK{}}
}

type ChatMessage struct {
	From string `json:"from"`
	Text string `json:"text"`
}

func (ChatMessage) NetTag() string { return "chat/msg" }

func init() {
	wire.Register(func() wire.Value { return new(SendChatMessage) })
	wire.Register(func() wire.Value { return new(ChatMessage) })
}
