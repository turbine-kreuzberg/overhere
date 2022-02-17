package overhere_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/lixiangzhong/dnsutil"
	overhere "github.com/turbine-kreuzberg/overhere/pkg"
)

func TestNewServer(t *testing.T) {
	port := 12353
	verbose := false

	// startup
	srv, err := overhere.NewServer("1.2.3.4", port, verbose)
	if err != nil {
		t.Errorf("NewServer() error = %v", err)
		return
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			t.Errorf("srv.ListenAndServe() error = %v", err)
		}
	}()

	// dns client
	dig := dnsutil.Dig{}
	dig.Retry = 3

	err = dig.SetDNS(fmt.Sprintf("localhost:%d", port))
	if err != nil {
		t.Errorf("dig.SetDNS() error = %v", err)
		return
	}

	// verify google.com
	got, err := dig.A("google.com")
	if err != nil {
		t.Errorf("dig.A(\"google.com\") error = %v", err)
		return
	}

	if len(got) == 0 {
		t.Errorf("missing result for google.com resolution")
		return
	}

	if got[0].A.Equal(net.ParseIP("1.2.3.4")) {
		t.Errorf("google.com resolved to fallback ip")
		return
	}

	// verify development-domain
	got, err = dig.A("development-domain")
	if err != nil {
		t.Errorf("dig.A(\"development-domain\") error = %v", err)
		return
	}

	if len(got) == 0 {
		t.Errorf("missing result for development-domain resolution")
		return
	}

	if !got[0].A.Equal(net.ParseIP("1.2.3.4")) {
		t.Errorf("development-domain does not resolved to fallback ip")
		return
	}

	// shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = srv.ShutdownContext(ctx)
	if err != nil {
		t.Errorf("dig.A(\"google.com\") error = %v", err)
		return
	}
}
