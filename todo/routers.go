package todo

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

func Router() *chi.Mux {
	r := chi.NewMux()

	r.Get("/", listItemsHandler)
	r.Get("/{itemId}", getItemHandler)
	r.Post("/", createItemHandler)
	r.Post("/done", makeItemDoneHandler)

	return r
}

func writeMessage(w http.ResponseWriter, status int, msg string) {
	var j struct {
		Msg string `json:"message"`
	}

	j.Msg = msg

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(j)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeMessage(w, status, err.Error())
}

func listItemsHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var resp struct {
		Items []TodoItem `json:"items,omitempty"`
	}

	items, err := listItems(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	resp.Items = items
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func getItemHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	itemId := chi.URLParam(req, "itemId")

	id, err := ulid.Parse(itemId)

	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var resp TodoItem

	item, err := findItem(ctx, id)

	if err != nil {
		if err == ErrTodoNotFound {
			writeMessage(w, http.StatusNotFound, "item not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	resp = item

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func createItemHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ctx := req.Context()

	title := req.FormValue("title")

	id, err := createItem(ctx, title)

	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var resp struct {
		Id string `json:"id"`
	}

	resp.Id = id.String()

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(resp)
}

func makeItemDoneHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ctx := req.Context()
	idStr := req.FormValue("id")

	id, err := ulid.Parse(idStr)

	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	err = makeItemDone(ctx, id)

	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
