# Sample gRPC backend for CoreDNS

#### Generate the gRPC stubs

This is only needed if [dns.proto](https://github.com/coredns/coredns/blob/master/pb/dns.proto)
is updated:

```
cd proto
protoc dns.proto --go_out=plugins=grpc:.
```

#### Run the server

This will start the backend on port 8053 (udp/tcp):

```
go build -o main
./main
```

#### Start CoreDNS

This will start coredns using the [Corefile](./Corefile) on port 1053 (udp/tcp)
and proxy requests to the backend over gRPC:

```
coredns
```

#### Try it out

```
$ dig +short @localhost -p 1053 A foo.example.com
127.0.0.1

$ dig +short @localhost -p 1053 AAAA foo.example.com
::1
```
