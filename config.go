package config

import (
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// MARK: Types

type Config struct {
	all   []Var
	flags map[string]Var
	env   map[string]Var
}

// MARK: Public Functions

func (c *Config) ReadAll() (err error) {
	err = c.readEnvVars()
	flag.Parse()
	return err
}

func (c *Config) readEnvVars() error {
	i, errs := 0, make([]error, len(c.env))
	for k, v := range c.env {
		if s := os.Getenv(k); s != "" {
			errs[i] = v.set(s)
		}
		i++
	}
	return errors.Join(errs...)
}

func (c *Config) Register(v ...Var) {
	if c.env == nil {
		c.env = make(map[string]Var)
	}
	if c.flags == nil {
		c.flags = make(map[string]Var)
	}
	for _, cVar := range v {
		if cVar == nil {
			return
		}
		c.all = append(c.all, cVar)
		if n := cVar.EnvName(); n != "" {
			c.env[n] = cVar
		}
		if n := cVar.FlagName(); n != "" {
			c.flags[n] = cVar
		}
	}
}

func (c *Config) NewAny(cfg VarCfg[any]) (*any, Var) {
	return newConfigVar(c, cfg)
}

func (c *Config) NewBool(cfg VarCfg[bool]) (*bool, Var) {
	return newConfigVar(c, cfg)
}

func (c *Config) NewString(cfg VarCfg[string]) (*string, Var) {
	return newConfigVar(c, cfg)
}

func (c *Config) NewInt(cfg VarCfg[int]) (*int, Var) {
	return newConfigVar(c, cfg)
}

func (c *Config) NewUint(cfg VarCfg[uint]) (*uint, Var) {
	return newConfigVar(c, cfg)
}

func (c *Config) NewFloat64(cfg VarCfg[float64]) (*float64, Var) {
	return newConfigVar(c, cfg)
}

func (c *Config) NewDuration(cfg VarCfg[time.Duration]) (*time.Duration, Var) {
	return newConfigVar(c, cfg)
}

func (c *Config) Flags() map[string]any {
	m := make(map[string]any)
	for k, v := range c.flags {
		m[k] = v.Value()
	}
	return m
}
func (c *Config) Env() map[string]any {
	m := make(map[string]any)
	for k, v := range c.env {
		m[k] = v.Value()
	}
	return m
}
func (c *Config) HandleHTTP() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", c.getConfigVars).Methods(http.MethodGet)
	r.HandleFunc("/", c.postConfigVars).Methods(http.MethodPost)
	r.HandleFunc("/", c.putConfigVars).Methods(http.MethodPut)
	return r
}

// MARK: Private Functions

func newConfigVar[T any](c *Config, cfg VarCfg[T]) (*T, Var) {
	v, i := NewVar(cfg)
	c.Register(i)
	return v, i
}

func (c *Config) getConfigVars(w http.ResponseWriter, _ *http.Request) {
	vars := c.Env()
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(vars)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("failed to marshal global to JSON"))
		return
	}
	_, _ = w.Write(data)
}

func (c *Config) postConfigVars(w http.ResponseWriter, r *http.Request) {
	var buf []byte
	n, err := r.Body.Read(buf)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if n == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	params := make(map[string]string)
	err = json.Unmarshal(buf, &params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	i, errs := 0, make([]error, len(params))
	for k, v := range params {
		if cVar, found := c.env[k]; found {
			errs[i] = wrap(cVar.set(v), k)
			_ = os.Setenv(k, v)
		} else if cVar, found = c.flags[k]; found {
			errs[i] = wrap(cVar.set(v), k)
		}
		i++
	}

	err = errors.Join(errs...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("error setting config variables"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *Config) putConfigVars(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	i, errs := 0, make([]error, len(params))
	for k, v := range params {
		v := strings.Join(v, ",")
		if cVar, found := c.env[k]; found {
			errs[i] = wrap(cVar.set(v), k)
			_ = os.Setenv(k, v)
		} else if cVar, found = c.flags[k]; found {
			errs[i] = wrap(cVar.set(v), k)
		}
		i++
	}
	err := errors.Join(errs...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("error setting config variables"))
		return
	}
	w.WriteHeader(http.StatusCreated)
}
