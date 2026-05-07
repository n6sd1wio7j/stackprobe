package export

import (
	"net/http"
)

// NewHandler returns an http.Handler that serves health-check records
// as JSON (default) or CSV when ?format=csv is present.
func NewHandler(exp *Exporter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		format := r.URL.Query().Get("format")

		switch format {
		case "csv":
			w.Header().Set("Content-Type", "text/csv")
			w.Header().Set("Content-Disposition", `attachment; filename="stackprobe-export.csv"`)
			if err := exp.WriteCSV(w); err != nil {
				http.Error(w, "failed to generate CSV", http.StatusInternalServerError)
			}
		default:
			w.Header().Set("Content-Type", "application/json")
			if err := exp.WriteJSON(w); err != nil {
				http.Error(w, "failed to generate JSON", http.StatusInternalServerError)
			}
		}
	})
}
