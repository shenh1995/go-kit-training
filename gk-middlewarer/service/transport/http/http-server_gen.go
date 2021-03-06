// Autogenerated code, do not change directly.
// To make changes to this file, please modify the templates at
// go-kit-middlewarer/tmpl/*.tmpl

package http

import (
	"net/http"

	ep "github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	"gk-middlewarer/service"
	"gk-middlewarer/service/endpoint"
)

type toEndpoint func(service.StringService) ep.Endpoint

// ServerLayer is a wrapper for github.com/OahcUil94/go-kit-training/gk-middlewarer/service.StringService which returns a
// github.com/go-kit/kit/endpoint.Middleware.  This allows you to specify
// Middleware while creating HTTP Servers.
type ServerLayer func(base service.StringService, path string) ep.Middleware

func epID(ep ep.Endpoint) ep.Endpoint {
	return ep
}

func serverFactory(stringService service.StringService, config ServerConfig, path string, endp toEndpoint, dec httptransport.DecodeRequestFunc, enc httptransport.EncodeResponseFunc) *httptransport.Server {
	var middlewares []ep.Middleware
	for _, w := range config.ServerLayers {
		middlewares = append(middlewares, w(stringService, path))
	}

	middlewares = append(middlewares, config.Middlewares...)

	var options []httptransport.ServerOption
	if config.ErrorEncoder != nil {
		options = append(options, httptransport.ServerErrorEncoder(config.ErrorEncoder))
	}
	options = append(options, httptransport.ServerBefore(config.RequestFuncs...))
	options = append(options, httptransport.ServerAfter(config.ServerReponseFuncs...))
	options = append(options, config.Options...)

	server := httptransport.NewServer(
		ep.Chain(epID, middlewares...)(endp(stringService)),
		dec,
		enc,
		options...,
	)
	config.Mux.Handle(path, server)
	return server
}

// ServersForEndpoints will take the given arguments, associate all of
// the proper endpoints together, and register itself as an HTTP handler for
// github.com/OahcUil94/go-kit-training/gk-middlewarer/service.StringService.
func ServersForEndpoints(stringService service.StringService, wrappers ...ServerLayer) (servers map[string]*httptransport.Server) {
	return ServersForEndpointsWithConfig(stringService, ServerConfig{ServerLayers: wrappers})
}

// ServersForEndpointsWithOptions will take the given arguments, associate all of
// the proper endpoints together, and register itself as an HTTP handler for
// github.com/OahcUil94/go-kit-training/gk-middlewarer/service.StringService.
func ServersForEndpointsWithOptions(stringService service.StringService, wrappers []ServerLayer, options []httptransport.ServerOption) (servers map[string]*httptransport.Server) {
	return ServersForEndpointsWithConfig(stringService, ServerConfig{ServerLayers: wrappers, Options: options})
}

// ServersForEndpointsWithConfig will take the given arguments, associate
// all of the endpoints togher, and register itself as an HTTP handler for
// github.com/OahcUil94/go-kit-training/gk-middlewarer/service.StringService.
//
// The function uses the ServerConfig specification to be setup. Any properties
// can be specified within the ServerConfig structure.
func ServersForEndpointsWithConfig(stringService service.StringService, config ServerConfig) (servers map[string]*httptransport.Server) {
	if config.Mux == nil {
		config.Mux = http.DefaultServeMux
	}

	return map[string]*httptransport.Server{

		endpoint.PathUppercase: serverFactory(stringService, config, endpoint.PathUppercase, makeUppercaseEndpoint, decodeUppercaseRequest, encodeUppercaseResponse),
		endpoint.PathCount:     serverFactory(stringService, config, endpoint.PathCount, makeCountEndpoint, decodeCountRequest, encodeCountResponse),
	}
}

// Mux represents an interface abstration for a Mux. This is useful when
// wanting to use something other than the ServeMux in the default http package.
// However, due to the signature restriction on the functions, adapaters will
// likely be required for other implementations.
type Mux interface {
	// Handle registers the handler for the given pattern.
	// According to net/http.ServeMux If a handler already exists for pattern,
	// the Handle invocation panics.
	Handle(pattern string, handler http.Handler)

	// HandleFunc registers the handler function for the given pattern.
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))

	// Any Mux must also implement net/http.Handler.
	http.Handler
}

// ServerConfig represents a set of configuation options that can be passed
// and overwritten when instanciating the Handlers.  This allows for maximum
// configuration and cutomization.  If nothing is provided, then defaults will
// be used.
type ServerConfig struct {
	// Mux represents the default Mux to use.  Defaults to
	// net/http.DefaultServeMux
	Mux Mux

	// Options represents a list of potential
	// github.com/go-kit/kit/transport/http.ServerOption(s).  These options
	// allow for direct manipulation of the
	// github.com/go-kit/kit/transport/http.Server, if desired.
	// These Options will be applied after the supplied ErrorEncoder, if it is
	// provided.
	Options []httptransport.ServerOption

	// ServerLayers represents a list of potential ServerLayers. Since a
	// ServerLayer generates an Endpoint, the provided ServerLayers will be
	// invoked as a chain of middlewares, in the order provided, to the
	// generated Endpoint.
	ServerLayers []ServerLayer

	// Middlewares represetns a list of potential
	// github.com/go-kit/kit/endpoint.Middleware(s). These Middlewares will be
	// applied after any supplied ServerLayers.
	Middlewares []ep.Middleware

	// RequestFuncs represents a list of potential
	// github.com/go-kit/kit/transport/http.RequestFunc(s) that will be invoked
	// before the processing of the Endpoint.
	RequestFuncs []httptransport.RequestFunc

	// ServerResponseFuncts represetns a list of potential
	// github.com/go-kit/kit/transport/http.ServerResponseFunc(s) that will be
	// invoked before the flush of the response generated by the Endpoint.
	ServerReponseFuncs []httptransport.ServerResponseFunc

	// ErrorEncoder allows for you to overwrite the ErrorEncoder.  If nothing
	// is specified, the Default from go-kit will be used.
	//
	// If a different ErrorEncoder is needed for different endpoints, then it
	// is recommended that the returned Servers be modified with go-kit's
	// ServerOptions, externally.
	ErrorEncoder httptransport.ErrorEncoder
}
