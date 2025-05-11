package orchestrator

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ZolotarevAlexandr/yl_final/db"
)

// handlePing handles GET /api/v1/ping healthcheck endpoint.
// It should always return 200 with message "orchestrator is up and running"
func handlePing(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "orchestrator is up and running"})
}

// handleCalculate processes POST /api/v1/calculate to add a new expression.
func handleCalculate(w http.ResponseWriter, r *http.Request) {
	var req struct{ Expression string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
		http.Error(w, "invalid data", http.StatusUnprocessableEntity)
		return
	}
	userID := UserIDFromContext(r.Context()) // from JWT middleware

	exprID, err := CreateExpression(req.Expression, userID)
	if err != nil {
		http.Error(w, "error processing expression", http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": exprID})
}

// handleListExpressions returns a list of all expressions.
func handleListExpressions(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	var list []db.Expression
	if err := db.DB.Where("user_id = ?", userID).Find(&list).Error; err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	// map to your JSON shape…
	json.NewEncoder(w).Encode(map[string]any{"expressions": list})
}

// handleGetExpression returns a specific expression by its id.
func handleGetExpression(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	userID := UserIDFromContext(r.Context())
	var expr db.Expression
	if err := db.DB.Where("id = ? AND user_id = ?", id, userID).
		First(&expr).Error; err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"expression": expr})
}
