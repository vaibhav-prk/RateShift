package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "ok",
			"path":    r.URL.Path,
			"method":  r.Method,
			"time":    time.Now().UTC().Format(time.RFC3339),
			"headers": r.Header,
		})
	})
	fmt.Println("Mock backend listening on :8081")
	http.ListenAndServe(":8081", nil)
}
