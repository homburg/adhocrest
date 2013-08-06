// Package for delivering a development rest server.
package adhocrest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"

	"github.com/nu7hatch/gouuid"
)

var address string

// A resource is just json
type resource string

// Resources, indexed by uuid
type resources map[string]resource

// The global resource cache
type resourceCache map[string]resources

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var html string

func respond(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.Method)
}

func serveError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	data, _ := json.Marshal(map[string]string{"status": "500", "message": err.Error()})
	w.Write(data)
}

// Get new handler with the given root
func NewHandler(root string) *Handler {
	return &Handler{root, resourceCache{}}
}

type Handler struct {
	root string
	db   resourceCache
}

type route struct {
	Resource, Id string
}

func getRoute(path string) route {
	path = strings.TrimLeft(path, "/")
	pathParts := strings.Split(path, "/")

	if len(pathParts) == 1 {
		return route{pathParts[0], ""}
	}

	if len(pathParts) >= 2 {
		return route{pathParts[0], pathParts[1]}
	}

	return route{}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	r.URL.Path = strings.TrimPrefix(r.URL.Path, h.root)
	route := getRoute(r.URL.Path)
	switch strings.ToUpper(r.Method) {
	case "POST":
		// Create!
		// Return the location to the created resource

		id, err := uuid.NewV4()
		if nil != err {
			serveError(w, r, err)
			return
		}
		newResourceId := id.String()

		// Create it, the resource list?
		if _, ok := h.db[route.Resource]; !ok {
			h.db[route.Resource] = make(resources)
		}

		bodyStr, err := ioutil.ReadAll(r.Body)
		if err != nil {
			serveError(w, r, err)
			return
		}

		// Save resource
		bodyData := make(map[string]string)
		err = json.Unmarshal(bodyStr, &bodyData)
		if nil != err {
			serveError(w, r, err)
			return
		}
		bodyData["id"] = newResourceId
		dataStr, err := json.Marshal(bodyData)
		if nil != err {
			serveError(w, r, err)
			return
		}
		h.db[route.Resource][newResourceId] = resource(dataStr)
		if nil != err {
			serveError(w, r, err)
			return
		}

		w.Header().Set("location", "/"+route.Resource+"/"+newResourceId)
		w.WriteHeader(http.StatusCreated)
	case "PUT":
		// Update
		log.Printf(r.Method)
		respond(w, r)
	case "DELETE":
		// Drop
		if _, ok := h.db[route.Resource]; !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if _, ok := h.db[route.Resource][route.Id]; ok {
			delete(h.db[route.Resource], route.Id)
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	case "GET":
		// Read
		if "" == route.Id {
			// List
			jsonList := "["
			if _, ok := h.db[route.Resource]; ok {
				for _, jsonItem := range h.db[route.Resource] {
					jsonList += string(jsonItem) + ","
				}
			}
			jsonList = strings.TrimRight(jsonList, ",")
			jsonList += "]"
			w.Write([]byte(jsonList))
		} else {
			// Get
			if _, ok := h.db[route.Resource]; !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			log.Println(h.db[route.Resource])

			if resource, ok := h.db[route.Resource][route.Id]; ok {
				w.Write([]byte(resource))
				return
			} else {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
	case "HEAD":
		log.Println("Unsupported method")
		w.WriteHeader(http.StatusMethodNotAllowed)
	default:
		log.Printf("Unknown request method \"%s\"\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
