package impl

import "net/http"

func (a *api) PingDB(w http.ResponseWriter, r *http.Request) error {
	status := http.StatusOK
	ctx := r.Context()

	err := a.mc.PingDB(ctx)
	if err != nil {
		return err
	}

	w.WriteHeader(status)
	return nil
}
