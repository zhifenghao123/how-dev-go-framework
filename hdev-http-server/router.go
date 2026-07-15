package ihttp

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

const (
	DefaultMethod       = "get"
	MethodDelimiter     = "@"
	MiddlewareDelimiter = "|"
)

// struct of config for router
type routeConf struct {
	Prefix      string                 `yaml:"prefix"`
	Middleware  string                 `yaml:"middleware"`
	Groups      []*routeConf           `yaml:"groups"`
	Actions     map[string]*actionConf `yaml:"actions"`
	Import      []string               `yaml:"import"`
	parentGroup []*routeConf
}

// struct of config for action
type actionConf struct {
	Method     string `yaml:"method"`
	Uses       string `yaml:"uses"`
	Middleware string `yaml:"middleware"`
}

// Router struct
type Router struct {
	conf []*routeConf
}

// Init router by config file
func (r *Router) Init(file string) {
	r.loadConf(file)
}

// load and extract router config file
func (r *Router) loadConf(file string, gl ...*routeConf) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln("route file read error : ", err)
	}

	var routeConf *routeConf
	err = yaml.Unmarshal(dat, &routeConf)
	if err != nil {
		log.Fatalln("route file extract err : ", err)
	}
	routeConf.parentGroup = gl
	r.conf = append(r.conf, routeConf)
	for _, i := range routeConf.Import {
		r.loadConf(i, append(gl, routeConf)...)
	}
}

// Action get all actions by config
func (r *Router) Action() []*action {
	return r.groupAction(r.conf)
}

// extract actions from config
func (r *Router) groupAction(gs []*routeConf, gl ...*routeConf) (actions []*action) {
	for _, g := range gs {
		gc := g
		if len(gc.Groups) > 0 {
			func(gs []*routeConf) {
				for _, action := range r.groupAction(gs, append(gl, append(gc.parentGroup, gc)...)...) {
					actions = append(actions, action)
				}
			}(gc.Groups)
		}

		for url, act := range gc.Actions {
			func(u string, a *actionConf) {
				if action := extract(u, a, append(gl, append(gc.parentGroup, gc)...)...); action != nil {
					actions = append(actions, action)
				}
			}(url, act)
		}
	}

	return
}

// extract action info
func extract(url string, act *actionConf, groupList ...*routeConf) *action {
	if act.Uses == "" {
		return nil
	}

	useArr := strings.Split(act.Uses, MethodDelimiter)
	if len(useArr) != 2 {
		return nil
	}
	method, controller, function := DefaultMethod, useArr[0], useArr[1]

	urlArr := strings.Split(url, " ")
	switch len(urlArr) {
	case 1:
		url = urlArr[0]
		break
	case 2:
		method, url = urlArr[0], urlArr[1]
	default:
		return nil
	}

	if act.Method != "" {
		method = act.Method
	}

	// append group prefix and middleware
	var middleware []string
	if len(groupList) > 0 {
		var prefix string
		for _, g := range groupList {
			if g == nil {
				continue
			}
			prefix = prefix + g.Prefix
			if g.Middleware == "" {
				continue
			}
			middleware = append(middleware, strings.Split(g.Middleware, MiddlewareDelimiter)...)
		}
		url = prefix + url
	}

	if act.Middleware != "" {
		middleware = append(middleware, strings.Split(act.Middleware, MiddlewareDelimiter)...)
	}

	return &action{
		Uri:        url,
		Method:     strings.ToLower(method),
		Controller: controller,
		Action:     function,
		Middleware: middleware,
	}
}

// DefaultRouter get instance of router for default
func DefaultRouter() IRouter {
	return &Router{}
}
