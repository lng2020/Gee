package geerpc

import (
	"encoding/json"
	"errors"
	"goTinyToys/geerpc/codec"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

type Option struct {
	MagicNumber int
	CodecType   codec.Type
}

var DefaultOption = &Option{
	MagicNumber: 0x3bef5c,
	CodecType:   codec.GobType,
}

type Server struct {
	serviceMap sync.Map
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

type request struct {
	h      *codec.Header
	argv   interface{}
	replyv interface{}
}

func (s *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		go s.ServeConn(conn)
	}
}

func Accept(lis net.Listener) {
	DefaultServer.Accept(lis)
}

func (s *Server) ServeConn(conn net.Conn) {
	defer func() { _ = conn.Close() }()

	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error:", err)
		return
	}
	if opt.MagicNumber != DefaultOption.MagicNumber {
		log.Printf("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}

	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return
	}

	serverCodec := f(conn)
	s.ServeCodec(serverCodec)
}

func (s *Server) ServeCodec(cc codec.Codec) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := s.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			s.sendResponse(cc, req.h, nil, sending)
			continue
		}
		wg.Add(1)
		go s.handleRequest(cc, req, sending, wg)
	}
	wg.Wait()
	_ = cc.Close()
}

// readRequest reads a request from the connection.
func (s *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := s.readRequestHeader(cc)
	if err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}

	req := &request{h: h}
	if err := cc.ReadBody(req.argv); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return req, err
}

// readRequestHeader reads a request header from the connection.
func (s *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

// sendResponse sends a response for the request to the connection.
func (s *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()

	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

// handleRequest handles the request.
func (s *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	// find service
	serviceMethod := req.h.ServiceMethod
	svc, mtype, err := s.findService(serviceMethod)
	if err != nil {
		req.h.Error = err.Error()
		s.sendResponse(cc, req.h, nil, sending)
		return
	}

	// call service
	returns := svc.call(mtype.method.Name, reflect.ValueOf(req.argv), reflect.ValueOf(req.replyv))
	s.sendResponse(cc, req.h, returns, sending)
}

// findService finds the service registered with the given name.
func (s *Server) findService(serviceMethod string) (svc *service, mtype *methodType, err error) {
	s.serviceMap.Range(func(key, value interface{}) bool {
		svc = value.(*service)
		mtype = svc.methods[serviceMethod]
		if mtype == nil {
			err = errors.New("method not found: " + serviceMethod)
			return false
		}
		return true
	})
	return
}
