package main
import "net/http"

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request){
	if cfg.platform != "dev"{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	cfg.db.DeleteAllUsers(r.Context())
    w.Write([]byte("Hits reset to 0"))
}