package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"cats-go-api/internal/model"
	"cats-go-api/internal/service"
)

type CatHandler struct {
	svc *service.CatService
}

func NewCatHandler(svc *service.CatService) *CatHandler {
	return &CatHandler{svc: svc}
}

func writeJSONError(w http.ResponseWriter, status int, detail string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"detail": detail})
}

func (h *CatHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	cats, err := h.svc.List(r.Context())
	if err != nil {
		log.Printf("error listing cats: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error al obtener la lista de gatos")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cats)
}

func (h *CatHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	idStr := r.PathValue("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "id required")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	cat, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("error getting cat: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error al obtener el gato")
		return
	}
	if cat == nil {
		writeJSONError(w, http.StatusNotFound, "Gato no encontrado")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cat)
}

func (h *CatHandler) GetByCatID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	catIDStr := r.PathValue("cat_id")
	if catIDStr == "" {
		writeJSONError(w, http.StatusBadRequest, "cat_id required")
		return
	}

	cat, err := h.svc.GetByCatID(r.Context(), catIDStr)
	if err != nil {
		log.Printf("error getting cat: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error al obtener el gato")
		return
	}
	if cat == nil {
		writeJSONError(w, http.StatusNotFound, "Gato no encontrado")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cat)
}

func (h *CatHandler) Save(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var input model.SaveCatInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json")
		return
	}

	cat, err := h.svc.Save(r.Context(), input)
	if errors.Is(err, service.ErrDuplicateCat) {
		writeJSONError(w, http.StatusConflict, "Ya existe un gato con este ID en la base de datos")
		return
	}
	if err != nil {
		log.Printf("error saving cat: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error al guardar el gato")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cat)
}

func (h *CatHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	idStr := r.PathValue("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "id required")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input model.SaveCatInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json")
		return
	}

	cat, err := h.svc.Update(r.Context(), input, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, http.StatusNotFound, "Gato no encontrado")
			return
		}
		log.Printf("error updating cat: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error al actualizar el gato")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cat)
}

func (h *CatHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	idStr := r.PathValue("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "id required")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.svc.Delete(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		writeJSONError(w, http.StatusNotFound, "Gato no encontrado")
		return
	}
	if err != nil {
		log.Printf("error deleting cat: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error al eliminar el gato")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
