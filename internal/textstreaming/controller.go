package textstreaming

import (
	"bufio"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type ExternalController struct {
	manager *Manager
}

func NewController(m *Manager) *ExternalController {
	return &ExternalController{
		manager: m,
	}
}

func (c *ExternalController) MountRoutes(router chi.Router) {
	router.With().Route("/text", func(r chi.Router) {
		r.MethodFunc(http.MethodPost, "/health/check", c.HealthCheck)
		r.MethodFunc(http.MethodPost, "/generate", c.TextGenerate)
	})
}

// Api for the cloud scheduler to check the status of each inference providers
func (c *ExternalController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	c.manager.HealthCheck()
	c.RespondWith(w, "Health check for all the inference providers in progress")
}

// Api for text streaming
func (c *ExternalController) TextGenerate(w http.ResponseWriter, r *http.Request) {
	var payload Request
	reader := c.manager.TextStreaming(payload)

	// Use a JSON encoder to write to the response writer
	encoder := json.NewEncoder(w)

	// Set headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")

	// Stream data from the pipe to the response writer using the encoder
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		// Read from the pipe
		line := scanner.Text()

		// Write the read data to the response using the encoder
		encoder.Encode(Response{Result: line})

		// Flush the response to ensure the client receives the data immediately
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}

func (c *ExternalController) ErrorWith(w http.ResponseWriter, statusCode int, err error) {
	errMsg := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}
	b, err := json.Marshal(errMsg)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(statusCode)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(b)
}

func (c *ExternalController) RespondWith(w http.ResponseWriter, res interface{}) error {
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		return err
	}

	return nil
}
