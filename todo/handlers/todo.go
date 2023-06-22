package handlers

import (
	"encoding/json"
	"mda/todo"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

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

func ListItems(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var resp struct {
		Items []todo.TodoItem `json:"items,omitempty"`
	}

	items, err := todo.ListItems(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	resp.Items = items
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func GetItem(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	itemId := chi.URLParam(req, "itemId")

	id, err := ulid.Parse(itemId)

	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var resp todo.TodoItem

	item, err := todo.FindItem(ctx, id)

	if err != nil {
		if err == todo.ErrTodoNotFound {
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

func CreateItem(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ctx := req.Context()

	title := req.FormValue("title")

	id, err := todo.CreateItem(ctx, title)

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

func MakeItemDone(w http.ResponseWriter, req *http.Request) {
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

	err = todo.MakeItemDone(ctx, id)

	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
