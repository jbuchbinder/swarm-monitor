// SWARM Distributed Monitoring System
// https://github.com/jbuchbinder/swarm-monitor
//
// vim: tabstop=2:softtabstop=2:shiftwidth=2:noexpandtab

package main

import (
	"fmt"
	"github.com/fromkeith/gorest"
	redis "github.com/jbuchbinder/go-redis"
	"github.com/jbuchbinder/swarm-monitor/config"
	"net/http"
	"strings"
	"time"
)

type ApiService struct {
	gorest.RestService `root:"/api/" consumes:"application/json" produces:"application/json"`

	// Hosts
	apiAddHost  gorest.EndPoint `method:"POST" path:"/hosts/" postdata:"HostDefinition"`
	apiGetHosts gorest.EndPoint `method:"GET" path:"/hosts/" postdata:"[]HostDefinition"`

	// Contacts
	apiAddContact  gorest.EndPoint `method:"POST" path:"/contacts/" postdata:"Contact"`
	apiGetContacts gorest.EndPoint `method:"GET" path:"/contacts/" postdata:"[]Contact"`

	// Checks/statuses
	apiGetStatus gorest.EndPoint `method:"GET" path:"/status/" postdata:"[]CheckStatus"`
}

func (serv ApiService) ApiAddHost(h HostDefinition) {
	c, err := redis.NewSynchClientWithSpec(getConnection(REDIS_READWRITE).connspec)
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
	c, err := redis.NewSynchClientWithSpec(getConnection(REDIS_READONLY).connspec)
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

func (serv ApiService) ApiAddContact(d Contact) {
	c, err := redis.NewSynchClientWithSpec(getConnection(REDIS_READWRITE).connspec)
	if err != nil {
		log.Err(err.Error())
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Error connecting to backend"))
		return
	}

	// Check for existing host definition
	exists, err := c.Sismember(CONTACT_LIST, []byte(CONTACT_PREFIX+":"+d.Name))
	if err != nil {
		log.Err(err.Error())
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Error connecting to backend"))
		return
	}
	if exists {
		serv.ResponseBuilder().SetResponseCode(409).WriteAndOveride([]byte("Contact already exists"))
		return
	}

	// All clear, add the contact
	k := CONTACT_PREFIX + ":" + d.Name
	c.Sadd(CONTACT_LIST, []byte(k))
	c.Hset(k, "name", []byte(d.Name))
	c.Hset(k, "display_name", []byte(d.DisplayName))
	c.Hset(k, "email", []byte(d.EmailAddress))

	serv.ResponseBuilder().Created("/api/contacts/" + k)
}

func (serv ApiService) ApiGetContacts() (r []Contact) {
	c, err := redis.NewSynchClientWithSpec(getConnection(REDIS_READONLY).connspec)
	if err != nil {
		log.Err(err.Error())
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Error connecting to backend"))
		return
	}

	cmembers, e := c.Smembers(CONTACT_LIST)
	if e != nil {
		log.Err(e.Error())
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Error connecting to backend"))
		return
	}

	ret := make([]Contact, len(cmembers))

	for i := 0; i < len(cmembers); i++ {
		cmember := string(cmembers[i])

		cdef := Contact{}

		// Grab full info from member
		h, e := c.Hgetall(cmember)
		if e != nil {
			log.Err(e.Error())
		} else {
			for j := 0; j < len(h); j += 2 {
				k := string(h[j])
				v := string(h[j+1])
				switch k {
				case "name":
					{
						cdef.Name = v
					}
				case "display_name":
					{
						cdef.DisplayName = v
					}
				case "email":
					{
						cdef.EmailAddress = v
					}
				default:
					{
						log.Debug("Unknown key " + k + " sighted in contact " + cmember)
					}
				}
			}

			ret[i] = cdef
		}
	}

	return ret
}

func (serv ApiService) ApiGetStatus() (r []CheckStatus) {
	c, err := redis.NewSynchClientWithSpec(getConnection(REDIS_READONLY).connspec)
	if err != nil {
		log.Err(err.Error())
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Error connecting to backend"))
		return
	}

	cmembers, e := c.Smembers(CHECKS_LIST)
	if e != nil {
		log.Err(e.Error())
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Error connecting to backend"))
		return
	}

	ret := make([]CheckStatus, len(cmembers))

	// TODO: FIXME: XXX: pull statuses to return!

	return ret
}

func threadWeb() {
	log.Info("Starting web thread")
	gorest.RegisterService(new(ApiService))
	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Config.Port),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	http.Handle("/api/", gorest.Handle())
	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("web"))))
	httpServer.ListenAndServe()
}
