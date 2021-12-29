// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package graphql

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/graph-gophers/graphql-go"

	"github.com/AlayaNetwork/Alaya-Go/node"

	"github.com/AlayaNetwork/Alaya-Go/internal/ethapi"
	"github.com/AlayaNetwork/Alaya-Go/log"
	"github.com/AlayaNetwork/Alaya-Go/p2p"
	"github.com/AlayaNetwork/Alaya-Go/rpc"

	graphqlEth "github.com/AlayaNetwork/graphql-go"

	json2 "github.com/AlayaNetwork/Alaya-Go/common/json"
)

type handler struct {
	Schema    *graphql.Schema
	SchemaEth *graphqlEth.Schema
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.URL.Path == "/graphql" || r.URL.Path == "/graphql/" {
		response := h.SchemaEth.Exec(r.Context(), params.Query, params.OperationName, params.Variables)
		responseJSON, err := json2.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(response.Errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	} else {
		response := h.Schema.Exec(r.Context(), params.Query, params.OperationName, params.Variables)
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(response.Errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	}
}

// Service encapsulates a GraphQL service.
type Service struct {
	endpoint string           // The host:port endpoint for this service.
	cors     []string         // Allowed CORS domains
	vhosts   []string         // Recognised vhosts
	timeouts rpc.HTTPTimeouts // Timeout settings for HTTP requests.
	backend  ethapi.Backend   // The backend that queries will operate on.
	handler  http.Handler     // The `http.Handler` used to answer queries.
	listener net.Listener     // The listening socket.
}

// New constructs a new GraphQL service instance.
func New(backend ethapi.Backend, endpoint string, cors, vhosts []string, timeouts rpc.HTTPTimeouts) (*Service, error) {
	return &Service{
		endpoint: endpoint,
		cors:     cors,
		vhosts:   vhosts,
		timeouts: timeouts,
		backend:  backend,
	}, nil
}

// Protocols returns the list of protocols exported by this service.
func (s *Service) Protocols() []p2p.Protocol { return nil }

// APIs returns the list of APIs exported by this service.
func (s *Service) APIs() []rpc.API { return nil }

// Start is called after all services have been constructed and the networking
// layer was also initialized to spawn any goroutines required by the service.
func (s *Service) Start(server *p2p.Server) error {
	var err error
	s.handler, err = newHandler(s.backend)
	if err != nil {
		return err
	}
	if s.listener, err = net.Listen("tcp", s.endpoint); err != nil {
		return err
	}
	// create handler stack and wrap the graphql handler
	handler := node.NewHTTPHandlerStack(s.handler, s.cors, s.vhosts)
	// make sure timeout values are meaningful
	node.CheckTimeouts(&s.timeouts)
	// create http server
	httpSrv := &http.Server{
		Handler:      handler,
		ReadTimeout:  s.timeouts.ReadTimeout,
		WriteTimeout: s.timeouts.WriteTimeout,
		IdleTimeout:  s.timeouts.IdleTimeout,
	}
	go httpSrv.Serve(s.listener)
	log.Info("GraphQL endpoint opened", "url", fmt.Sprintf("http://%s", s.endpoint))
	return nil
}

// newHandler returns a new `http.Handler` that will answer GraphQL queries.
// It additionally exports an interactive query browser on the / endpoint.
func newHandler(backend ethapi.Backend) (http.Handler, error) {
	q := Resolver{backend}

	s, err := graphql.ParseSchema(schema, &q)
	if err != nil {
		return nil, err
	}

	sEth, err := graphqlEth.ParseSchema(schema, &q)
	if err != nil {
		return nil, err
	}

	h := handler{Schema: s, SchemaEth: sEth}

	mux := http.NewServeMux()
	mux.Handle("/", GraphiQL{})
	mux.Handle("/alaya", GraphiQL{}) // for alaya

	mux.Handle("/graphql", h)
	mux.Handle("/graphql/", h)

	// for alaya
	mux.Handle("/alaya/graphql", h)
	mux.Handle("/alaya/graphql/", h)

	return mux, nil
}

// Stop terminates all goroutines belonging to the service, blocking until they
// are all terminated.
func (s *Service) Stop() error {
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
		log.Info("GraphQL endpoint closed", "url", fmt.Sprintf("http://%s", s.endpoint))
	}
	return nil
}
