package main

import (
	"./lib/gorest"
	"./lib/redis"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ApiService struct {
	gorest.RestService `root:"/api/" consumes:"application/json" produces:"application/json"`

	// Hosts
	apiAddHost  gorest.EndPoint `method:"POST" path:"/hosts/" postdata:"HostDefinition"`
	apiGetHosts gorest.EndPoint `method:"GET" path:"/hosts/" postdata:"[]HostDefinition"`
}

func (serv ApiService) ApiAddHost(h HostDefinition) {
	c, err := redis.NewSynchClientWithSpec(getConnection(true).connspec)
	if err != nil {
		log.Err(err.Error())
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Error connecting to backend"))
		return
	}

	// Check for existing host definition
	exists, err := c.Sismember(HOSTS_LIST, []byte(HOST_PREFIX+":"+h.Name))
	if err != nil {
		log.Err(err.Error())
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Error connecting to backend"))
		return
	}
	if exists {
		serv.ResponseBuilder().SetResponseCode(409).WriteAndOveride([]byte("Host already exists"))
		return
	}

	// All clear, add the host
	k := HOST_PREFIX + ":" + h.Name
	c.Sadd(HOSTS_LIST, []byte(k))
	c.Hset(k, "name", []byte(h.Name))
	c.Hset(k, "address", []byte(h.Address))

	serv.ResponseBuilder().Created("/api/hosts/" + k)
}

func (serv ApiService) ApiGetHosts() (r []HostDefinition) {
	c, err := redis.NewSynchClientWithSpec(getConnection(true).connspec)
	if err != nil {
		log.Err(err.Error())
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Error connecting to backend"))
		return
	}

	hmembers, e := c.Smembers(HOSTS_LIST)
	if e != nil {
		log.Err(e.Error())
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Error connecting to backend"))
		return
	}

	ret := make([]HostDefinition, len(hmembers))

	for i := 0; i < len(hmembers); i++ {
		hmember := string(hmembers[i])

		hdef := HostDefinition{}

		// Grab full info from member
		h, e := c.Hgetall(hmember)
		if e != nil {
			log.Err(e.Error())
		} else {
			for j := 0; j < len(h); j += 2 {
				k := string(h[j])
				v := string(h[j+1])
				switch k {
				case "name":
					{
						hdef.Name = v
					}
				case "address":
					{
						hdef.Address = v
					}
				default:
					{
						log.Debug("Unknown key " + k + " sighted in host " + hmember)
					}
				}
			}

			// Get list of checks
			h, e := c.Hgetall(hmember + ":checks")
			if e != nil {
				log.Err(e.Error())
			} else {
				hdef.Checks = make([]string, len(h)/2)
				for j := 0; j < len(h); j += 2 {
					k := string(h[j])
					hdef.Checks[j/2] = strings.Replace(k, CHECK_PREFIX+":", "", -1)
				}
			}

			ret[i] = hdef
		}
	}

	return ret
}

func threadWeb() {
	log.Info("Starting web thread")
	gorest.RegisterService(new(ApiService))
	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", *webPort),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	http.Handle("/api/", gorest.Handle())
	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("web"))))
	httpServer.ListenAndServe()
}
