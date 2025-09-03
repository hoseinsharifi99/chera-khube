package app

import (
	"chera_khube/handler/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func (a *application) InitRouter(ctrl *controller, logger *zap.Logger) *gin.Engine {
	r := gin.New()

	if a.config.System.DevelopMode {
		a.wrapRouter(r)
	}

	healthR := router.NewHealthRouter()
	healthR.HandleRoutes(r, a.config)

	userR := router.NewUserRouter(ctrl.userController)
	userR.HandleRoutes(r, a.config)

	paymentR := router.NewPaymentRouter(ctrl.paymentController)
	paymentR.HandleRoutes(r, a.config)

	postR := router.NewPostRouter(ctrl.postController)
	postR.HandleRoutes(r, a.config)

	pricingR := router.NewPricingRouter(ctrl.pricingController)
	pricingR.HandleRoutes(r, a.config)

	divarR := router.NewDivarRouter(ctrl.divarController)
	divarR.HandleRoutes(r, a.config)

	adonsR := router.NewAddonsRouter(ctrl.addonsController)
	adonsR.HandleRoutes(r, a.config)

	return r
}

type Skipper func(ctx *gin.Context) bool

type (
	// CORSConfig defines the config for CORS middleware.
	CORSConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// AllowOrigin defines a list of origins that may access the resource.
		// Optional. Default value []string{"*"}.
		AllowOrigins []string `yaml:"allow_origins"`

		// AllowOriginFunc is a custom function to validate the origin. It takes the
		// origin as an argument and returns true if allowed or false otherwise. If
		// an error is returned, it is returned by the handler. If this option is
		// set, AllowOrigins is ignored.
		// Optional.
		AllowOriginFunc func(origin string) (bool, error) `yaml:"allow_origin_func"`

		// AllowMethods defines a list methods allowed when accessing the resource.
		// This is used in response to a preflight request.
		// Optional. Default value DefaultCORSConfig.AllowMethods.
		AllowMethods []string `yaml:"allow_methods"`

		// AllowHeaders defines a list of request headers that can be used when
		// making the actual request. This is in response to a preflight request.
		// Optional. Default value []string{}.
		AllowHeaders []string `yaml:"allow_headers"`

		// AllowCredentials indicates whether or not the response to the request
		// can be exposed when the credentials flag is true. When used as part of
		// a response to a preflight request, this indicates whether or not the
		// actual request can be made using credentials.
		// Optional. Default value false.
		AllowCredentials bool `yaml:"allow_credentials"`

		// ExposeHeaders defines a whitelist headers that clients are allowed to
		// access.
		// Optional. Default value []string{}.
		ExposeHeaders []string `yaml:"expose_headers"`

		// MaxAge indicates how long (in seconds) the results of a preflight request
		// can be cached.
		// Optional. Default value 0.
		MaxAge int `yaml:"max_age"`
	}
)

// The `const` block is declaring two constants HeaderOrigin and HeaderVary. These constants are
// used as keys for accessing specific headers in HTTP requests and responses.
const (
	HeaderOrigin = "Origin"
	HeaderVary   = "Vary"
)

// The DefaultSkipper function always returns false.
func DefaultSkipper(ctx *gin.Context) bool {
	return false
}

// The `var` block is declaring a variable named DefaultCORSConfig of type CORSConfig. It is
// assigning a default value to this variable, which is a default configuration for the CORS
// (Cross-Origin Resource Sharing) middleware.
var (
	// DefaultCORSConfig is the default CORS middleware config.
	DefaultCORSConfig = CORSConfig{
		Skipper:      DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
		},
		AllowCredentials: true,
	}
)

// The cors function is a method of the `application` struct that returns a `gin.HandlerFunc`
// function. This function is used as a middleware for handling Cross-Origin Resource Sharing (CORS) in
// a Go application using the Gin framework.
func (a *application) cors() gin.HandlerFunc {
	return a.corsWithConfig(DefaultCORSConfig)
}

// The noContent function is a helper function that sets the HTTP response status code to the
// provided `code` parameter and returns an empty response body. It is used to send a response with no
// content to the client.
func (a *application) noContent(c *gin.Context, code int) {
	c.Writer.WriteHeader(code)
	return
}

// The matchScheme function is used to compare the scheme (protocol) of two URLs. It takes two
// parameters, `domain` and `pattern`, which are both strings representing URLs.
func (a *application) matchScheme(domain, pattern string) bool {
	didx := strings.Index(domain, ":")
	pidx := strings.Index(pattern, ":")
	return didx != -1 && pidx != -1 && domain[:didx] == pattern[:pidx]
}

// matchSubdomain compares authority with wildcard
func (a *application) matchSubdomain(domain, pattern string) bool {
	if !a.matchScheme(domain, pattern) {
		return false
	}
	didx := strings.Index(domain, "://")
	pidx := strings.Index(pattern, "://")
	if didx == -1 || pidx == -1 {
		return false
	}
	domAuth := domain[didx+3:]
	// to avoid long loop by invalid long domain
	if len(domAuth) > 253 {
		return false
	}
	patAuth := pattern[pidx+3:]

	domComp := strings.Split(domAuth, ".")
	patComp := strings.Split(patAuth, ".")
	for i := len(domComp)/2 - 1; i >= 0; i-- {
		opp := len(domComp) - 1 - i
		domComp[i], domComp[opp] = domComp[opp], domComp[i]
	}
	for i := len(patComp)/2 - 1; i >= 0; i-- {
		opp := len(patComp) - 1 - i
		patComp[i], patComp[opp] = patComp[opp], patComp[i]
	}

	for i, v := range domComp {
		if len(patComp) <= i {
			return false
		}
		p := patComp[i]
		if p == "*" {
			return true
		}
		if p != v {
			return false
		}
	}
	return false
}

// The above code is implementing a middleware function for handling Cross-Origin Resource Sharing
// (CORS) in a Go application using the Gin framework.
func (a *application) corsWithConfig(config CORSConfig) gin.HandlerFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultCORSConfig.Skipper
	}
	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = DefaultCORSConfig.AllowOrigins
	}
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = DefaultCORSConfig.AllowMethods
	}

	var allowOriginPatterns []string
	for _, origin := range config.AllowOrigins {
		pattern := regexp.QuoteMeta(origin)
		pattern = strings.Replace(pattern, "\\*", ".*", -1)
		pattern = strings.Replace(pattern, "\\?", ".", -1)
		pattern = "^" + pattern + "$"
		allowOriginPatterns = append(allowOriginPatterns, pattern)
	}

	allowMethods := strings.Join(config.AllowMethods, ",")
	allowHeaders := strings.Join(config.AllowHeaders, ",")
	exposeHeaders := strings.Join(config.ExposeHeaders, ",")
	maxAge := strconv.Itoa(config.MaxAge)

	// return func(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Skipper(c) {
			c.Next()
		}

		req := c.Request
		res := c.Writer
		origin := req.Header.Get(HeaderOrigin)
		allowOrigin := ""

		preflight := req.Method == http.MethodOptions
		res.Header().Add(HeaderVary, HeaderOrigin)

		// No Origin provided
		if origin == "" {
			if !preflight {
				c.Next()
			}

			a.noContent(c, http.StatusNoContent)
			return
		}

		if config.AllowOriginFunc != nil {
			allowed, err := config.AllowOriginFunc(origin)
			if err != nil {
				return
			}
			if allowed {
				allowOrigin = origin
			}
		} else {
			// Check allowed origins
			for _, o := range config.AllowOrigins {
				if o == "*" && config.AllowCredentials {
					allowOrigin = origin
					break
				}
				if o == "*" || o == origin {
					allowOrigin = o
					break
				}
				if a.matchSubdomain(origin, o) {
					allowOrigin = origin
					break
				}
			}

			// Check allowed origin patterns
			for _, re := range allowOriginPatterns {
				if allowOrigin == "" {
					didx := strings.Index(origin, "://")
					if didx == -1 {
						continue
					}
					domAuth := origin[didx+3:]
					// to avoid regex cost by invalid long domain
					if len(domAuth) > 253 {
						break
					}

					if match, _ := regexp.MatchString(re, origin); match {
						allowOrigin = origin
						break
					}
				}
			}
		}

		// Origin not allowed
		if allowOrigin == "" {
			if !preflight {
				c.Next()
			}
			a.noContent(c, http.StatusNoContent)
		}

		// Simple request
		if !preflight {
			res.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			if config.AllowCredentials {
				res.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			if exposeHeaders != "" {
				res.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
			}
			c.Next()
		}

		// Preflight request
		res.Header().Add(HeaderVary, "Access-Control-Request-Method")
		res.Header().Add(HeaderVary, "Access-Control-Request-Method")
		res.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		res.Header().Set("Access-Control-Allow-Methods", allowMethods)
		if config.AllowCredentials {
			res.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if allowHeaders != "" {
			res.Header().Set("Access-Control-Allow-Credentials", allowHeaders)
		} else {
			h := req.Header.Get("Access-Control-Request-Headers")
			if h != "" {
				res.Header().Set("Access-Control-Allow-Headers", h)
			}
		}
		if config.MaxAge > 0 {
			res.Header().Set("Access-Control-Max-Age", maxAge)
		}
		a.noContent(c, http.StatusNoContent)
	}
	// }
}

func (a *application) wrapRouter(router *gin.Engine) {
	router.Use(gin.Recovery()) // Recovery with output

	router.Use(gin.Logger()) // Optional logger middleware.

	router.Use(a.cors())
}
