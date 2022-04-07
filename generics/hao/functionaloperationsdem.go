package hao

import (
	"crypto/tls"
	"time"
)

/* type Server struct {
	Addr     string
	Port     int
	Protocol string
	Timeout  time.Duration
	MaxConns int
	TLS      *tls.Config
}

func NewDefaultServer(addr string, port int) (*Server, error) {
	return &Server{addr, port, "tcp", 30 * time.Second, 100, nil}, nil
}

func NewTLSServer(addr string, port int, tls *tls.Config) (*Server, error) {
	return &Server{addr, port, "tcp", 30 * time.Second, 100, tls}, nil
}

func NewServerWithTomeout(addr string, port int, timeout time.Duration) (*Server, error) {
	return &Server{addr, port, "tcp", timeout, 100, nil}, nil
}

func NewTLSSererWithMaxConnAndTimeout(addr string, port int, maxconns int, timeout time.Duration, tls *tls.Config) (*Server, error) {
	return &Server{addr, port, "tcp", 30 * time.Second, maxconns, tls}, nil
} */

/*type Config struct {
	Protocol string
	Timeout  time.Duration
	Maxconns int
	TLS      *tls.Config
}

type Server struct {
	Addr string
	Port int
	Conf *Config
}

func NewServer(addr string, port int, conf *Config) (*Server, error) {
	// ...

	return &Server{addr, port, conf}, nil
}

func RunFunctional() {
	srv1, _ := NewServer("localhost", 9000, nil)

	fmt.Println(srv1)
	fmt.Println("functional")

	conf := Config{Protocol: "tcp", Timeout: 60 * time.Second}
	srv2, _ := NewServer("localhost", 9000, &conf)

	fmt.Println(srv2)
} */

// 使用一个builder类来做包装

type Server struct {
	Addr     string
	Port     int
	Protocol string
	Timeout  time.Duration
	MaxConns int
	TLS      *tls.Config
}

/*
//使用一个builder类来做包装
type ServerBuilder struct {
	Server
}

func (sb *ServerBuilder) Create(addr string, port int) *ServerBuilder {
	sb.Server.Addr = addr
	sb.Server.Port = port
	//其它代码设置其它成员的默认值
	return sb
}

func (sb *ServerBuilder) WithProtocol(protocol string) *ServerBuilder {
	sb.Server.Protocol = protocol
	return sb
}

func (sb *ServerBuilder) WithMaxConn(maxconn int) *ServerBuilder {
	sb.Server.MaxConns = maxconn
	return sb
}

func (sb *ServerBuilder) WithTimeOut(timeout time.Duration) *ServerBuilder {
	sb.Server.Timeout = timeout
	return sb
}

func (sb *ServerBuilder) WithTLS(tls *tls.Config) *ServerBuilder {
	sb.Server.TLS = tls
	return sb
}

func (sb *ServerBuilder) Build() Server {
	return sb.Server
}

func RunFunctionalBuild() int {

	sb := ServerBuilder{}
	server, err := sb.Create("127.0.0.1", 8080).WithProtocol("udp").WithMaxConn(1024).WithTimeOut(30 * time.Second).Build()

	// fmt.Println(server)
	return 1
}
*/

type Option func(*Server)

func Protocol(p string) Option {
	return func(s *Server) {
		s.Protocol = p
	}
}
