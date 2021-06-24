package dns

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

func InitDnsService() error {
	if !conf.Config.Dns.EnableService {
		log.Warnf("dns service not start")
		return nil
	}

	server := &dns.Server{Addr: ":" + strconv.Itoa(int(conf.Config.Dns.ServicePort)), Net: "udp"}
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Compress = false

		switch r.Opcode {
		case dns.OpcodeQuery:
			for _, q := range m.Question {
				switch q.Qtype {
				case dns.TypeA:
					ip := LookupHostsFastestBack(q.Name)
					if ip != "" {
						rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
						if err == nil {
							head := rr.Header()
							head.Ttl = 60
							m.Answer = append(m.Answer, rr)
						}
						log.Infof("[dns query]%v %v", q.Name, ip)
					}
				}
			}
		default:
			log.Warnf("unsupported opcode")
		}

		err := w.WriteMsg(m)
		if err != nil {
			log.Errorf("err:%v", err)
		}
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	}()
	go func() {
		for {
			select {
			case <-time.After(time.Minute * 10):
				keys := dnsCache.Keys(true)
				for _, key := range keys {
					domain := key.(string)
					LookupHostsFastestIp(domain)
				}
			case <-sigCh:
				log.Info("dns service stop")
				err := server.Shutdown()
				if err != nil {
					log.Errorf("err:%v", err)
				}
				return
			}
		}
	}()

	return nil
}
