package mgmt

import (
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/handler"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/processor"
	"github.com/bitmaelum/bitmaelum-suite/internal/apikey"
	"net/http"
)

// FlushQueues handler will flush all the queues normally on tickers
func FlushQueues(w http.ResponseWriter, req *http.Request) {
	key := handler.GetAPIKey(req)
	if !key.HasPermission(apikey.PermFlush) {
		handler.ErrorOut(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Flush queues. Note that this means that multiple queue processing can run multiple times
	go processor.ProcessRetryQueue()
	go processor.ProcessStuckIncomingMessages()
	go processor.ProcessStuckProcessingMessages()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(handler.StatusOk("Flushing queues"))
}
