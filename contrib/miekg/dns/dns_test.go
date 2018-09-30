package dns

import (
	"net"
	"testing"
	"time"

	"github.com/lonnblad/dd-trace-go/ddtrace/ext"
	"github.com/lonnblad/dd-trace-go/ddtrace/mocktracer"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestDNS(t *testing.T) {
	mux := dns.NewServeMux()
	mux.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		w.WriteMsg(m)
	})
	addr := getFreeAddr(t).String()

	// start the server
	go func() {
		err := ListenAndServe(addr, "udp", mux)
		if err != nil {
			t.Fatal(err)
		}
	}()
	waitTillUDPReady(t, addr)

	mt := mocktracer.Start()
	defer mt.Stop()

	m := new(dns.Msg)
	m.SetQuestion("miek.nl.", dns.TypeMX)

	_, err := Exchange(m, addr)
	assert.NoError(t, err)

	spans := mt.FinishedSpans()
	assert.Len(t, spans, 2)
	for _, s := range spans {
		assert.Equal(t, "dns.request", s.OperationName())
		assert.Equal(t, "dns", s.Tag(ext.SpanType))
		assert.Equal(t, "dns", s.Tag(ext.ServiceName))
		assert.Equal(t, "QUERY", s.Tag(ext.ResourceName))
	}
}

func getFreeAddr(t *testing.T) net.Addr {
	li, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := li.Addr()
	li.Close()
	return addr
}

func waitTillUDPReady(t *testing.T, addr string) {
	deadline := time.Now().Add(time.Second * 10)
	for time.Now().Before(deadline) {
		m := new(dns.Msg)
		m.SetQuestion("miek.nl.", dns.TypeMX)
		_, err := dns.Exchange(m, addr)
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}
