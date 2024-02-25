// Implement kubernetes probes.
package probes

import (
	"net/http"
	"sync"

	"fortio.org/log"
)

type state struct {
	started bool
	live    bool
	ready   bool
	// mutex
	mutex sync.Mutex
}

var State = state{}

func (s *state) SetLive(live bool) {
	s.mutex.Lock()
	s.live = live
	s.mutex.Unlock()
}

func (s *state) SetReady(ready bool) {
	s.mutex.Lock()
	s.ready = ready
	s.mutex.Unlock()
}

func (s *state) SetStarted(started bool) {
	s.mutex.Lock()
	s.started = started
	s.mutex.Unlock()
}

func (s *state) IsLive() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.live
}

func (s *state) IsReady() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.ready
}

func (s *state) IsStarted() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.started
}

func StartupProbe(w http.ResponseWriter, _ *http.Request) {
	if State.IsStarted() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func LivenessProbe(w http.ResponseWriter, _ *http.Request) {
	if State.IsLive() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func ReadinessProbe(w http.ResponseWriter, _ *http.Request) {
	if State.IsReady() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func Setup(mux *http.ServeMux) {
	mux.HandleFunc("/startup", log.LogAndCall("startup", StartupProbe))
	mux.HandleFunc("/ready", log.LogAndCall("readiness", ReadinessProbe))
	mux.HandleFunc("/health", log.LogAndCall("liveness", LivenessProbe))
}
