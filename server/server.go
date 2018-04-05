package server

import (
	"errors"
	"io"
	"sync"

	"github.com/nat-ventura/message-party/server/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	port          = 9090
	ErrRecvStream = errors.New("Error receiving stream")
)

type broadcaster struct {
	listenerMu sync.RWMutex
	listeners  map[string]chan<- string
}

type ChatService struct {
	b broadcaster
}

func (b *broadcaster) Add(name string, listener chan<- string) error {
	b.listenerMu.Lock()
	defer b.listenerMu.Unlock()
	if b.listeners == nil {
		b.listeners = map[string]chan<- string{}
	}

	if _, ok := b.listeners[name]; ok {
		return status.Errorf(codes.AlreadyExists, "username %q is already in use by someone", name)
	}
	b.listeners[name] = listener
	return nil
}

func (b *broadcaster) Remove(name string) {
	b.listenerMu.Lock()
	defer b.listenerMu.Unlock()
	if c, ok := b.listeners[name]; ok {
		close(c)
		delete(b.listeners, name)
	}
}

func (b *broadcaster) Broadcast(ctx context.Context, msg string) {
	b.listenerMu.RLock()
	defer b.listenerMu.RUnlock()
	for _, listener := range b.listeners {
		select {
		case listener <- msg:
		case <-ctx.Done():
			return
		}
	}
}

func (chat *ChatService) WebChat(srv proto.ChatService_WebChatServer) error {
	msg, err := srv.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	name := msg.GetName()
	if name == "" {
		return status.Error(codes.FailedPrecondition, "first message should be name of user")
	}

	chat.b.Broadcast(srv.Context(), name+" has joined the chat")

	listener := make(chan string)
	err = chat.b.Add(name, listener)
	if err != nil {
		return err
	}
	defer func() {
		chat.b.Remove(name)
		chat.b.Broadcast(context.Background(), name+" has left the chat")
	}()

	sendErrChan := make(chan error)
	go func() {
		for {
			select {
			case msg, ok := <-listener:
				if !ok {
					return
				}
				err = srv.Send(&proto.WebResponse{Message: msg})
				if err != nil {
					sendErrChan <- err
					return
				}
			case <-srv.Context().Done():
				return
			}
		}
	}()

	recvErrChan := make(chan error)
	go func() {
		for {
			msg, err := srv.Recv()
			if err == io.EOF {
				close(recvErrChan)
				return
			}
			if err != nil {
				recvErrChan <- err
				return
			}
			chat.b.Broadcast(srv.Context(), name+": "+msg.GetMessage())
		}
	}()

	select {
	case err, ok := <-recvErrChan:
		if !ok {
			return nil
		}
		return err
	case err := <-sendErrChan:
		return err
	case <-srv.Context().Done():
		return srv.Context().Err()
	}
}

// type chatServer struct {
// 	streams []pb.Chat_ConnectServer
// }

// func (s *chatServer) Connect(stream pb.Chat_ConnectServer) error {
// 	s.streams = append(s.streams, stream)
// 	for {
// 		msg, err := stream.Recv()
// 		fmt.Printf("%s\n", msg.Text)
// 		if err == io.EOF {
// 			return nil
// 		}
// 		if err != nil {
// 			return ErrRecvStream
// 		}

// 		go func() {
// 			for i, stream_client := range s.streams {
// 				if err := stream_client.Send(msg); err != nil {
// 					s.streams = append(s.streams[:i], s.streams[i+1:]...)
// 				}
// 			}
// 		}()
// 	}
// }

// func main() {
// 	s := new(chatServer)

// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
// 	if err != nil {
// 		grpclog.Fatalf("failed to listen: %v", err)
// 	}
// 	var opts []grpc.ServerOption
// 	grpcServer := grpc.NewServer(opts...)
// 	pb.RegisterChatServer(grpcServer, s)
// 	grpcServer.Serve(lis)
// }
