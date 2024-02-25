// Implement kubernetes probes.
package probes

import (
	"net/http"
	"sync"

	"fortio.org/dflag"
	"fortio.org/log"
)

type state struct {
	started bool
	live    bool
	// mutex
	mutex sync.Mutex
}

var (
	State     = state{}
	ReadyFlag = dflag.NewBool(false, "Initial readiness state")
)

func (s *state) SetLive(live bool) {
	s.mutex.Lock()
	s.live = live
	s.mutex.Unlock()
	log.Infof("Setting live to %v", live)
}

func (s *state) SetReady(ready bool) {
	_ = ReadyFlag.SetV(ready)
	log.Infof("Setting ready to %v", ready)
}

func (s *state) SetStarted(started bool) {
	s.mutex.Lock()
	s.started = started
	s.mutex.Unlock()
	log.Infof("Setting started to %v", started)
}

func (s *state) IsLive() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.live
}

func (s *state) IsReady() bool {
	return ReadyFlag.Get()
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
