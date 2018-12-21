package mock

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Mock struct {
	server *httptest.Server
	mux    *http.ServeMux
	lock   sync.Mutex
}

type Router struct {
	AcceptedHttpMethods []string
	Path                string
	RespHttpCode        int
	RespBody            string
	ReqJsonBodyStruct   interface{}
	FormParams          []FormParameter
}

type FormParameter struct {
	Name       string
	IsRequired bool
}

func NewMock(closeAfter time.Duration) *Mock {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		//	mock router as default
		resp.WriteHeader(http.StatusNotFound)

		writeStringToResp(resp, fmt.Sprintf("%s not found in mock routers", req.URL))
	})

	if closeAfter != 0 {
		time.AfterFunc(closeAfter, server.Close)
	}

	return &Mock{
		mux:    mux,
		server: server,
	}
}

func NewMockWithRoute(router *Router) *Mock {
	mockServer := NewMock(0)

	mockServer.AddRouter(router)

	return mockServer
}

func (m *Mock) Close() {
	m.server.Close()
}

func (m *Mock) GetServerURL() string {
	return m.server.URL
}

func (m *Mock) AddHandler(path string, handler func(resp http.ResponseWriter, req *http.Request)) {
	m.mux.HandleFunc(path, handler)
}

func (m *Mock) AddRouter(routers ...*Router) {
	if len(routers) == 0 {
		return
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	for _, router := range routers {
		if router == nil {
			continue
		}

		m.AddHandler(
			router.Path,
			func(resp http.ResponseWriter, req *http.Request) {
				if len(router.AcceptedHttpMethods) != 0 {
					currentMethod := req.Method

					isAccepted := false

					for _, acceptedMethod := range router.AcceptedHttpMethods {
						if acceptedMethod == currentMethod {
							isAccepted = true
							break
						}
					}

					if !isAccepted {
						resp.WriteHeader(http.StatusMethodNotAllowed)
						return
					}
				}

				// Checking json body
				if router.ReqJsonBodyStruct != nil {
					bodyBytes, err := ioutil.ReadAll(req.Body)

					if len(bodyBytes) == 0 || err != nil {
						resp.WriteHeader(http.StatusBadRequest)

						return
					}

					jsonStructType := reflect.TypeOf(router.ReqJsonBodyStruct)

					newObj := reflect.New(jsonStructType).Elem().Addr()

					err = json.Unmarshal(bodyBytes, newObj.Interface())

					if err != nil {
						resp.WriteHeader(http.StatusBadRequest)

						writeStringToResp(resp, err.Error())

						return
					}
				} else if len(router.FormParams) != 0 { // If not json maybe are there form parameters ?
					errors := make([]string, 0, 1)

					for _, param := range router.FormParams {
						val := req.Form[param.Name]

						if param.IsRequired && len(val) == 0 {
							errors = append(errors, fmt.Sprintf("value of %s is blank", param.Name))
						}
					}

					if len(errors) != 0 {
						resp.WriteHeader(http.StatusBadRequest)

						writeStringToResp(resp, strings.Join(errors, ", "))

						return
					}
				}

				if router.RespHttpCode != 0 {
					resp.WriteHeader(router.RespHttpCode)
				}

				if len(router.RespBody) != 0 {
					writeStringToResp(resp, router.RespBody)
				}
			},
		)
	}
}

func writeStringToResp(resp http.ResponseWriter, str string) {
	n, err := io.WriteString(resp, str)

	if n != len(str) || err != nil {
		fmt.Printf("failed within writing response body [str=%s, wroteCount=%d, err=%v]\n", str, n, err)
	}
}
