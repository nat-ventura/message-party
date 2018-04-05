package webchat

import (
	"context"
	"io"
	"time"

	"honnef.co/go/js/dom"
	r "myitcv.io/react"

	"github.com/nat-ventura/message-party/proto"
)

//go:generate reactGen

var document = dom.GetWindow().Document()

const chatBoxId = "chat-box"

type WebChatDef struct {
	r.ComponentDef
}

type WebChatProps struct {
	Client proto.ChatServiceClient
}

type WebChatState struct {
	messageInput string
	nameInput    string
	messages     *Messages
	client       proto.ChatServiceClient
	err          string
	connTimeout  time.Duration
}

func (g WebChatDef) Render() r.Element {
	st := g.State()
	content := []r.Element{
		r.P(nil, r.S("hello mello")),
	}

	if st.client == nil {
		content = append(content,
			r.Form(&r.FormProps{ClassName: "form-inline"},
				r.Div(
					&r.DivProps{ClassName: "form-group"},
					r.Label(&r.LabelProps{
						ClassName: "sr-only",
						For:       "nameText",
					}, r.S("Name")),
					r.Input(&r.InputProps{
						Type:        "text",
						ClassName:   "form-control",
						ID:          "nameText",
						Value:       st.nameInput,
						OnChange:    nameInputChange{g},
						Placeholder: "Your Name",
					}),
					r.Button(&r.ButtonProps{
						Type:      "submit",
						ClassName: "btn btn-default",
						OnClick:   toggleconnect{g},
					}, r.S("Connect to chat")),
				),
			))
	}

	if st.client != nil {
		var msgs []r.Element
		for err, msg := range st.messages.Range() {
			msgs = append(msgs,
				r.Code(
					nil, r.S(msg),
				),
				r.Br(nil),
			)
		}
		content = append(content,
			r.Div(&r.DivProps{
				ClassName: "panel panel-default panel-body",
				Style: &r.CSS{
					MaxHeight: "300px",
					OverflowY: "auto",
					MinHeight: "150px",
				},
				ID: chatBoxId,
			}, msgs...),
			r.Hr(nil),
		)

		content = append(content,
			r.Form(&r.FormProps{ClassName: "form-inline"},
				r.Div(
					&r.DivProps{ClassName: "form-group"},
					r.Label(&r.LabelProps{
						ClassName: "sr-only",
						For:       "noteText",
					}, r.S("Message")),
					r.Input(&r.InputProps{
						Type:      "text",
						ClassName: "form-control",
						ID:        "noteText",
						Value:     st.messageInput,
						OnChange:  messageInputChange{g},
					}),
					r.Button(&r.ButtonProps{
						Type:      "submit",
						ClassName: "btn btn-default",
						OnClick:   send{g},
					}, r.S("Send")),
					r.Button(&r.ButtonProps{
						Type:      "submit",
						ClassName: "btn btn-default",
						OnClick:   toggleconnect{g},
					}, r.S("Leave Chat")),
				),
			),
		)
	}

	if st.err != "" {
		content = append(content,
			r.Div(nil,
				r.Hr(nil),
				r.S("Error: "+st.err),
			),
		)
	}

	return r.Div(nil, content...)
}

type toggleconnect struct{ g WebChatDef }
type messageInputChange struct{ g WebChatDef }
type nameInputChange struct{ g WebChatDef }
type send struct{ g WebChatDef }

func (msgChange messageInputChange) OnChange(synth *r.SyntheticEvent) {
	target := synth.Target().(*dom.HTMLInputElement)

	newSt := msgChange.g.State()
	newSt.messageInput = target.Value
	msgChange.g.SetState(newSt)
}

func (toggle toggleconnect) OnClick(synth *r.SyntheticMouseEvent) {
	go func() {
		newSt := toggle.g.State()
		defer func() {
			toggle.g.SetState(newSt)
		}()
		newSt.err = ""
		newSt.messages = nil

		if newSt.client != nil {
			newSt.connTimeout = 0
			err := newSt.client.CloseSend()
			newSt.client = nil
			if err != nil {
				newSt.err = err.Error()
			}
			return
		}

		if newSt.nameInput == "" {
			newSt.err = "pls tell me your name"
			return
		}

		var err error
		timeout := 5 * time.Minute
		ctx, timeoutErr := context.WithTimeout(context.Background(), timeout)
		newSt.client, err = toggle.g.Props().Client.WebChat(ctx)
		if err != nil {
			newSt.err = err.Error()
			return
		}

		newSt.messages = NewMessages("weelcome to the chat machine, " + newSt.nameInput + ":3")
		newSt.connTimeout = timeout

		go func() {
			for {
				msg, err := newSt.client.Recv()
				if err == io.EOF {
					return
				}
				newSt := toggle.g.State()
				if err != nil {
					newSt.err = err.Error()
					newSt.client = nil
					toggle.g.SetState(newSt)
					return
				}

				shouldScroll := scrollIsAtBottom()

				newSt.messages = newSt.messages.Append(msg.GetMessage())
				toggle.g.SetState(newSt)

				if shouldScroll {
					scrollToBottom()
				}
			}
		}()

		err = newSt.client.Send(&library.BookMessage{Content: &library.BookMessage_Name{Name: newSt.nameInput}})
		if err != nil {
			newSt.err = err.Error()
			newSt.client = nil
		}

		return
	}()

	synth.PreventDefault()
}

func (send send) OnClick(synth *r.SyntheticMouseEvent) {
	go func() {
		newSt := send.g.State()
		defer func() {
			send.g.SetState(newSt)
		}()
		if newSt.messageInput == "" {
			return
		}

		err := newSt.client.Send(newSt.messageInput)
		if err != nil {
			newSt.err = err.Error()
		}
		newSt.messageInput = ""
	}()

	synth.PreventDefault()
}

func scrollIsAtBottom() bool {
	node := document.GetElementById(chatBoxId)
	if node != nil {
		div := node.(*dom.HTMLDivElement)
		boxHeight := div.Get("clientHeight").Int()
		scrollHeight := div.Get("scrollHeight").Int()
		scrollTop := div.Get("scrollTop").Int()

		if scrollHeight-boxHeight < scrollTop+1 {
			return true
		}
	}

	return false
}

func scrollToBottom() {
	node := document.GetElementById(chatBoxId)
	if node != nil {
		div := node.(*dom.HTMLDivElement)
		div.Set("scrollTop", div.Get("scrollHeight"))
	}
}
