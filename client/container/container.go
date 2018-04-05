package container

import (
	"strings"

	"honnef.co/go/js/dom"
	"honnef.co/go/js/xhr"
	"myitcv.io/highlightjs"
	r "myitcv.io/react"

	"github.com/nat-ventura/message-party/client"
	"github.com/nat-ventura/message-party/proto"
)

//go:generate reactGen

type ContainerDef struct {
	r.ComponentDef
}

type ContainerState struct {
	client   proto.ChatServiceClient
	examples *exampleSource
}

func Container() *ContainerElem {
	return buildContainerElem()
}

func (c ContainerDef) GetInitialState() ContainerState {
	return ContainerState{
		client:   nil,
		examples: newExampleSource(),
	}
}

func (c ContainerDef) ComponentWillMount() {
	newSt := c.State()

	newSt.client = NewChatServiceClient(
		strings.TrimSuffix(dom.GetWindow().Document().BaseURI(), "/"),
	)

	c.SetState(newSt)
}

func (c ContainerDef) Render() r.Element {
	content := []r.Element{
		p.renderExample(
			exampleBookChat,
			r.Span(nil,
				r.S("bi-directional streaming chat"),
			),
			WebChat(webchat.WebChatProps{Client: c.State().client}),
		),
	}

	return r.Div(nil,
		r.Div(&r.DivProps{ClassName: "container"},
			content...,
		),
	)
}

func plainPanel(children ...r.Element) r.Element {
	return r.Div(&r.DivProps{ClassName: "panel panel-default panel-body"},
		children...,
	)
}
