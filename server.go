package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ahmetb/coredns-grpc-backend-sample/pb"
	"github.com/miekg/dns"
	"google.golang.org/grpc"
)

type dnsServer struct{}

func (d *dnsServer) Query(ctx context.Context, in *pb.DnsPacket) (*pb.DnsPacket, error) {
	m := new(dns.Msg)
	if err := m.Unpack(in.Msg); err != nil {
		return nil, fmt.Errorf("failed to unpack msg: %v", err)
	}
	for _, q := range m.Question {
		// TODO: query database and add answers here
		hdr := dns.RR_Header{Name: q.Name,
			Rrtype: q.Qtype,
			Class:  q.Qclass}

		switch q.Qtype {
		case dns.TypeA:
			m.Answer = append(m.Answer, &dns.A{
				Hdr: hdr,
				A:   net.IPv4(127, 0, 0, 1)}) // TODO use a real IP
		case dns.TypeAAAA:
			m.Answer = append(m.Answer, &dns.AAAA{
				Hdr:  hdr,
				AAAA: net.IPv6loopback})
		default:
			return nil, fmt.Errorf("only A/AAAA supported, got qtype=%d", q.Qtype)
		}
	}

	if len(m.Answer) == 0 {
		m.Rcode = dns.RcodeNameError
	} else {
		m.Response = true
	}

	out, err := m.Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack msg: %v", err)
	}
	return &pb.DnsPacket{Msg: out}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterDnsServiceServer(grpcServer, &dnsServer{})
	panic(grpcServer.Serve(lis))
}
