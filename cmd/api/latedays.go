package main

import (
	"net/http"

	"github.com/biboyqg/late_days/internal/store"
	"github.com/go-chi/chi/v5"
)

//	@Summary		Get late days by student ID
//	@Description	Get late days by student ID
//	@Tags			Late Days
//	@Accept			json
//	@Produce		json
//	@Param			studentID	path		string	true	"Student ID"
//	@Success		200			{object}	store.LateDay
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/late-days/{studentID} [get]
func (app *application) getLateDaysByStudentID(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "studentID")

	lateDay, err := app.store.LateDays.GetByStudentID(r.Context(), studentID)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, lateDay); err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err.Error())
	}
}

type createLateDaysPayload struct {
	StudentID    string `json:"student_id" validate:"required,max=255"`
	StudentEmail string `json:"student_email" validate:"required,email"`
	Days         *int   `json:"days" validate:"required"`
}

//	@Summary		Create late days
//	@Description	Create late days
//	@Tags			Late Days
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		createLateDaysPayload	true	"Late days payload"
//	@Success		201		{object}	store.LateDay
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/late-days [post]
func (app *application) createLateDays(w http.ResponseWriter, r *http.Request) {
	var payload createLateDaysPayload

	if err := app.readJSON(w, r, &payload); err != nil {
		app.errorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.errorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if payload.Days == nil {
		app.errorJSON(w, http.StatusBadRequest, "days is required")
		return
	}

	lateDay := store.LateDay{
		StudentID:    payload.StudentID,
		StudentEmail: payload.StudentEmail,
		Days:         *payload.Days,
	}

	if err := app.store.LateDays.Create(r.Context(), &lateDay); err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, lateDay); err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err.Error())
	}
}

type updateLateDaysPayload struct {
	StudentID    *string `json:"student_id" validate:"omitempty,max=255"`
	StudentEmail *string `json:"student_email" validate:"omitempty,email"`
	Days         *int    `json:"days" validate:"omitempty"`
}

//	@Summary		Update late days
//	@Description	Update late days
//	@Tags			Late Days
//	@Accept			json
//	@Produce		json
//	@Param			studentID	path		string					true	"Student ID"
//	@Param			payload		body		updateLateDaysPayload	true	"Late days payload"
//	@Success		200			{object}	store.LateDay
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/late-days/{studentID} [patch]
func (app *application) updateLateDays(w http.ResponseWriter, r *http.Request) {

	studentID := chi.URLParam(r, "studentID")

	lateDay, err := app.store.LateDays.GetByStudentID(r.Context(), studentID)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if lateDay == nil {
		app.errorJSON(w, http.StatusNotFound, "late days not found")
		return
	}

	var payload updateLateDaysPayload

	if err := app.readJSON(w, r, &payload); err != nil {
		app.errorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.errorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if payload.StudentID != nil {
		lateDay.StudentID = *payload.StudentID
	}

	if payload.StudentEmail != nil {
		lateDay.StudentEmail = *payload.StudentEmail
	}

	if payload.Days != nil {
		lateDay.Days = *payload.Days
	}

	if err := app.store.LateDays.Update(r.Context(), lateDay); err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, lateDay); err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err.Error())
	}
}
