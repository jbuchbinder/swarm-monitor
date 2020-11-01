// SWARM Distributed Monitoring System
// https://github.com/jbuchbinder/swarm-monitor
//
// vim: tabstop=2:softtabstop=2:shiftwidth=2:noexpandtab

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fromkeith/gorest"
	"github.com/go-redis/redis/v8"
	"github.com/jbuchbinder/swarm-monitor/config"
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
	c := redis.NewClient(getConnection(REDIS_READWRITE).connspec)

	// Check for existing host definition
	exists, err := c.SIsMember(ctx, HOSTS_LIST, []byte(HOST_PREFIX+":"+h.Name)).Result()
	if err != nil {
		log.Printf("ERR: " + err.Error())
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Error connecting to backend"))
		return
	}
	if exists {
		serv.ResponseBuilder().SetResponseCode(409).WriteAndOveride([]byte("Host already exists"))
		return
	}

	// All clear, add the host
	k := HOST_PREFIX + ":" + h.Name
	c.SAdd(ctx, HOSTS_LIST, []byte(k))
	c.HSet(ctx, k, "name", []byte(h.Name))
	c.HSet(ctx, k, "address", []byte(h.Address))

	serv.ResponseBuilder().Created("/api/hosts/" + k)
}

func (serv ApiService) ApiGetHosts() (r []HostDefinition) {
	c := redis.NewClient(getConnection(REDIS_READWRITE).connspec)

	hmembers, e := c.SMembers(ctx, HOSTS_LIST).Result()
	if e != nil {
		log.Printf("ERR: " + e.Error())
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Error connecting to backend"))
		return
	}

	ret := make([]HostDefinition, len(hmembers))

	for i := 0; i < len(hmembers); i++ {
		hmember := string(hmembers[i])

		hdef := HostDefinition{}

		// Grab full info from member
		h, e := c.HGetAll(ctx, hmember).Result()
		if e != nil {
			log.Printf("ERR: " + e.Error())
		} else {
			for k, v := range h {
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
						log.Printf("DEBUG: Unknown key " + k + " sighted in host " + hmember)
					}
				}
			}

			// Get list of checks
			h, e := c.HGetAll(ctx, hmember+":checks").Result()
			if e != nil {
				log.Printf("ERR: " + e.Error())
			} else {
				hdef.Checks = make([]string, len(h)/2)
				for _, k := range h {
					hdef.Checks = append(hdef.Checks, strings.Replace(k, CHECK_PREFIX+":", "", -1))
				}
			}

			ret[i] = hdef
		}
	}

	return ret
}

func (serv ApiService) ApiAddContact(d Contact) {
	c := redis.NewClient(getConnection(REDIS_READWRITE).connspec)

	// Check for existing host definition
	exists, err := c.SIsMember(ctx, CONTACT_LIST, []byte(CONTACT_PREFIX+":"+d.Name)).Result()
	if err != nil {
		log.Printf("ERR: " + err.Error())
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Error connecting to backend"))
		return
	}
	if exists {
		serv.ResponseBuilder().SetResponseCode(409).WriteAndOveride([]byte("Contact already exists"))
		return
	}

	// All clear, add the contact
	k := CONTACT_PREFIX + ":" + d.Name
	c.SAdd(ctx, CONTACT_LIST, []byte(k))
	c.HSet(ctx, k, "name", []byte(d.Name))
	c.HSet(ctx, k, "display_name", []byte(d.DisplayName))
	c.HSet(ctx, k, "email", []byte(d.EmailAddress))

	serv.ResponseBuilder().Created("/api/contacts/" + k)
}

func (serv ApiService) ApiGetContacts() (r []Contact) {
	c := redis.NewClient(getConnection(REDIS_READWRITE).connspec)

	cmembers, e := c.SMembers(ctx, CONTACT_LIST).Result()
	if e != nil {
		log.Printf("ERR: " + e.Error())
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Error connecting to backend"))
		return
	}

	ret := make([]Contact, len(cmembers))

	for i := 0; i < len(cmembers); i++ {
		cmember := string(cmembers[i])

		cdef := Contact{}

		// Grab full info from member
		h, e := c.HGetAll(ctx, cmember).Result()
		if e != nil {
			log.Printf("ERR: " + e.Error())
		} else {
			for k, v := range h {
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
						log.Printf("DEBUG: Unknown key " + k + " sighted in contact " + cmember)
					}
				}
			}

			ret[i] = cdef
		}
	}

	return ret
}

func (serv ApiService) ApiGetStatus() (r []CheckStatus) {
	c := redis.NewClient(getConnection(REDIS_READWRITE).connspec)

	cmembers, e := c.SMembers(ctx, CHECKS_LIST).Result()
	if e != nil {
		log.Printf("ERR: " + e.Error())
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Error connecting to backend"))
		return
	}

	ret := make([]CheckStatus, len(cmembers))

	// TODO: FIXME: XXX: pull statuses to return!

	return ret
}

func threadWeb() {
	log.Printf("INFO: Starting web thread")
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
