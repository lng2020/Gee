package geerpc

import (
	"fmt"
	"goTinyToys/geerpc/codec"
	"log"
	"net"
	"sync"
)

type Call struct {
	Seq           uint64
	ServiceMethod string
	Args          interface{}
	Reply         interface{}
	Error         error
	Done          chan *Call
}

var ErrShutdown = fmt.Errorf("connection is shut down")
var defaultOption = &Option{
	MagicNumber: 0x3bef5c,
	CodecType:   codec.GobType,
}

type Client struct {
	cc       codec.Codec
	opt      *Option
	sending  sync.Mutex
	header   codec.Header
	mu       sync.Mutex
	seq      uint64
	pending  map[uint64]*Call
	closing  bool
	shutdown bool
}

func NewClientWithCodec(cc codec.Codec, option ...*Option) *Client {
	opt := defaultOption
	if len(option) > 0 && option[0] != nil {
		opt = option[0]
	}
	client := &Client{
		cc:      cc,
		opt:     opt,
		pending: make(map[uint64]*Call),
	}
	go client.receive()
	return client
}

func NewClient(conn net.Conn, option ...*Option) *Client {
	f := codec.NewCodecFuncMap[defaultOption.CodecType]
	if f == nil {
		log.Panic("invalid codec type")
	}
	return NewClientWithCodec(f(conn), option...)
}

func (client *Client) Close() error {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing {
		return ErrShutdown
	}
	client.closing = true
	return client.cc.Close()
}

func (client *Client) IsAvailable() bool {
	client.mu.Lock()
	defer client.mu.Unlock()
	return !client.shutdown
}

func (client *Client) registerCall(call *Call) (uint64, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.shutdown {
		return 0, ErrShutdown
	}
	seq := client.seq
	client.seq++
	client.pending[seq] = call
	return seq, nil
}

func (client *Client) removeCall(seq uint64) *Call {
	client.mu.Lock()
	defer client.mu.Unlock()
	call := client.pending[seq]
	delete(client.pending, seq)
	return call
}

func (client *Client) terminateCalls(err error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	for seq, call := range client.pending {
		call.Error = err
		call.Done <- call
		delete(client.pending, seq)
	}
}

func (client *Client) send(call *Call) {
	client.sending.Lock()
	defer client.sending.Unlock()
	seq, err := client.registerCall(call)
	if err != nil {
		call.Error = err
		call.Done <- call
		return
	}
	client.header.Seq = seq
	client.header.ServiceMethod = call.ServiceMethod
	client.header.Error = ""
	if err := client.cc.Write(&client.header, call.Args); err != nil {
		call := client.removeCall(seq)
		if call != nil {
			call.Error = err
			call.Done <- call
		}
		client.terminateCalls(err)
	}
}

func (client *Client) receive() {
	var err error
	for err == nil {
		var header codec.Header
		if err = client.cc.ReadHeader(&header); err != nil {
			break
		}
		call := client.removeCall(header.Seq)
		switch {
		case call == nil:
			err = client.cc.ReadBody(nil)
		case header.Error != "":
			call.Error = fmt.Errorf(header.Error)
			client.cc.ReadBody(nil)
			call.Done <- call
		default:
			err = client.cc.ReadBody(call.Reply)
			if err != nil {
				call.Error = err
			}
			call.Done <- call
		}
	}
	client.sending.Lock()
	client.terminateCalls(err)
	client.shutdown = true
	client.sending.Unlock()
}

func (client *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	if done == nil {
		done = make(chan *Call, 10)
	} else if cap(done) == 0 {
		log.Println("done channel is unbuffered")
	}
	call := &Call{ServiceMethod: serviceMethod, Args: args, Reply: reply, Done: done}
	client.send(call)
	return call
}

func (client *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	call := <-client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}

func Dial(network, address string, option ...*Option) (*Client, error) {
	opt, err := parseOption(option...)
	if err != nil {
		return nil, err

	}

	cc, err := net.Dial(network, address)
	if err != nil {
		return nil, err

	}

	defer func() {
		if err != nil {
			cc.Close()
		}
	}()

	return NewClient(cc, opt), nil
}

func parseOption(option ...*Option) (*Option, error) {
	opt := defaultOption
	if len(option) > 0 && option[0] != nil {
		opt = option[0]
	}
	if opt.MagicNumber != defaultOption.MagicNumber {
		return nil, fmt.Errorf("magic number is not match")
	}
	return opt, nil
}
