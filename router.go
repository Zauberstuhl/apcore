// apcore is a server framework for implementing an ActivityPub application.
// Copyright (C) 2019 Cory Slep
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package apcore

import (
	"net/http"

	"github.com/go-fed/activity/pub"
	"github.com/gorilla/mux"
)

type Router struct {
	router            *mux.Router
	actor             pub.Actor
	errorHandler      http.Handler
	badRequestHandler http.Handler
}

func (r *Router) wrap(route *mux.Route) *Route {
	return &Route{
		route:             route,
		actor:             r.actor,
		errorHandler:      r.errorHandler,
		badRequestHandler: r.badRequestHandler,
	}
}

func (r *Router) ActorPostInbox(path string) *Route {
	return r.wrap(r.router.HandleFunc(path,
		func(w http.ResponseWriter, req *http.Request) {
			isApRequest, err := r.actor.PostInbox(req.Context(), w, req)
			if err != nil {
				ErrorLogger.Errorf("Error in ActorPostInbox: %s", err)
				r.errorHandler.ServeHTTP(w, req)
			} else if !isApRequest {
				r.badRequestHandler.ServeHTTP(w, req)
			}
		}))
}

func (r *Router) ActorPostOutbox(path string) *Route {
	return r.wrap(r.router.HandleFunc(path,
		func(w http.ResponseWriter, req *http.Request) {
			isApRequest, err := r.actor.PostOutbox(req.Context(), w, req)
			if err != nil {
				ErrorLogger.Errorf("Error in ActorPostOutbox: %s", err)
				r.errorHandler.ServeHTTP(w, req)
			} else if !isApRequest {
				r.badRequestHandler.ServeHTTP(w, req)
			}
		}))
}

func (r *Router) ActorGetInbox(path string, web func(http.ResponseWriter, *http.Request)) *Route {
	return r.wrap(r.router.HandleFunc(path,
		func(w http.ResponseWriter, req *http.Request) {
			isApRequest, err := r.actor.GetInbox(req.Context(), w, req)
			if err != nil {
				ErrorLogger.Errorf("Error in ActorGetInbox: %s", err)
				r.errorHandler.ServeHTTP(w, req)
			} else if !isApRequest {
				web(w, req)
			}
		}))
}

func (r *Router) ActorGetOutbox(path string, web func(http.ResponseWriter, *http.Request)) *Route {
	return r.wrap(r.router.HandleFunc(path,
		func(w http.ResponseWriter, req *http.Request) {
			isApRequest, err := r.actor.GetOutbox(req.Context(), w, req)
			if err != nil {
				ErrorLogger.Errorf("Error in ActorGetOutbox: %s", err)
				r.errorHandler.ServeHTTP(w, req)
			} else if !isApRequest {
				web(w, req)
			}
		}))
}

func (r *Router) ActivityPubOnlyHandleFunc(path string, apHandler pub.HandlerFunc) *Route {
	// TODO: construct pub.HandlerFunc in here instead
	return r.wrap(r.router.HandleFunc(path,
		func(w http.ResponseWriter, req *http.Request) {
			isASRequest, err := apHandler(req.Context(), w, req)
			if err != nil {
				ErrorLogger.Error(err)
			}
			if !isASRequest && r.router.NotFoundHandler != nil {
				r.router.NotFoundHandler.ServeHTTP(w, req)
			}
		}))
}

func (r *Router) ActivityPubAndWebHandleFunc(path string, apHandler pub.HandlerFunc, f func(http.ResponseWriter, *http.Request)) *Route {
	// TODO: construct pub.HandlerFunc in here instead
	return r.wrap(r.router.HandleFunc(path,
		func(w http.ResponseWriter, req *http.Request) {
			isASRequest, err := apHandler(req.Context(), w, req)
			if err != nil {
				ErrorLogger.Error(err)
			}
			if !isASRequest {
				f(w, req)
			}
		}))
}

func (r *Router) Get(name string) *Route {
	return r.wrap(r.router.Get(name))
}

func (r *Router) WebOnlyHandle(path string, handler http.Handler) *Route {
	return r.wrap(r.router.Handle(path, handler))
}

func (r *Router) WebOnlyHandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.wrap(r.router.HandleFunc(path, f))
}

func (r *Router) Headers(pairs ...string) *Route {
	return r.wrap(r.router.Headers(pairs...))
}

func (r *Router) Host(tpl string) *Route {
	return r.wrap(r.router.Host(tpl))
}

func (r *Router) Methods(methods ...string) *Route {
	return r.wrap(r.router.Methods(methods...))
}

func (r *Router) Name(name string) *Route {
	return r.wrap(r.router.Name(name))
}

func (r *Router) NewRoute() *Route {
	return r.wrap(r.router.NewRoute())
}

func (r *Router) Path(tpl string) *Route {
	return r.wrap(r.router.Path(tpl))
}

func (r *Router) PathPrefix(tpl string) *Route {
	return r.wrap(r.router.PathPrefix(tpl))
}

func (r *Router) Queries(pairs ...string) *Route {
	return r.wrap(r.router.Queries(pairs...))
}

func (r *Router) Schemes(schemes ...string) *Route {
	return r.wrap(r.router.Schemes(schemes...))
}

func (r *Router) Use(mwf ...mux.MiddlewareFunc) {
	r.router.Use(mwf...)
}

func (r *Router) Walk(walkFn mux.WalkFunc) error {
	return r.router.Walk(walkFn)
}

type Route struct {
	route             *mux.Route
	actor             pub.Actor
	errorHandler      http.Handler
	badRequestHandler http.Handler
}

// TODO: move Router methods to Route, have Router delegate to Route. No code dupe
