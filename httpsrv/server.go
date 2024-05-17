package httpsrv

import (
	"context"
	"encoding/json"
	"net/http"

	"rift/authz/client"
	"rift/memdb"
)

func New(
	db *memdb.DB,
	authzC *client.Client,
) *http.ServeMux {
	mux := http.NewServeMux()

	// authn middleware
	memberId := "member"
	orgId := "org"

	// this is set during registration/invite/accept/change role/...
	if err := authzC.WriteOrganizationAdmin(
		context.Background(),
		orgId,
		memberId,
	); err != nil {
		panic(err)
	}

	mux.HandleFunc("POST /offdays", func(w http.ResponseWriter, r *http.Request) {
		if err := authzC.CanCreateOrganizationOffDay(r.Context(), orgId, memberId); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		var offDay memdb.OffDay
		if err := json.NewDecoder(r.Body).Decode(&offDay); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		db.AddOffDay(&offDay)
		authzC.WriteOffDayOrganization(r.Context(), offDay.ID, orgId)
	})

	mux.HandleFunc("GET /offdays", func(w http.ResponseWriter, r *http.Request) {
		offDayIds, err := authzC.ListViewOffDays(r.Context(), memberId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		offDays := db.GetOffDays(offDayIds...)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(offDays); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("GET /offdays/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := authzC.CanViewOffDay(r.Context(), id, memberId); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		offDay := db.GetOffDay(id)
		if offDay == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(offDay); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("PATCH /offdays/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := authzC.CanEditOffDay(r.Context(), id, memberId); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		// parse request body
		var offDay memdb.OffDay
		if err := json.NewDecoder(r.Body).Decode(&offDay); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		db.AddOffDay(&offDay)
	})

	mux.HandleFunc("DELETE /offdays/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := authzC.CanDeleteOffDay(r.Context(), id, memberId); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		db.DeleteOffDay(id)

		if err := authzC.DeleteOffDayOrganization(r.Context(), id, orgId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	return mux
}
