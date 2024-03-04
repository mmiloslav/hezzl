package main

import (
	"encoding/json"
	"io"
	"net/http"
	"postgres"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
)

var validate *validator.Validate

const (
	errInternal     = "errors.internal"
	errGoodNotFound = "errors.good.notFound"
	errWrongParams  = "wrong.params"
)

func initValidator() {
	validate = validator.New()
}

func validateRequest(req interface{}) error {
	err := validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			logrus.Errorf("Field '%s' is %s", err.Field(), err.Tag())
		}
		return err
	}

	return nil
}

// type IRequest interface{}

type badResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func writeResponse(w http.ResponseWriter, response interface{}, statusCode int) {
	byteBody, err := json.Marshal(response)
	if err != nil {
		logrus.Errorf("couldnt marshal during writeAnswer of object [%+v] with error [%s]", response, err.Error())
		return
	}
	w.Header().Add("Content-Type", "text/json")
	w.WriteHeader(statusCode)
	io.WriteString(w, string(byteBody))
}

// PING
type pingResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func pingHandler(writer http.ResponseWriter, request *http.Request) {
	logrus.Info("health check")

	pingResponse := pingResponse{
		Success: true,
	}

	writeResponse(writer, pingResponse, 200)
}

// GOOD CREATE
type goodCreateRequest struct {
	Name string
}

type goodCreateUpdateResponse struct {
	Success     bool      `json:"success,omitempty"`
	ID          int       `json:"id"`
	ProjectID   int       `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"createdAt"`
}

func goodCreate(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("handling good create request...")
	badResponse := badResponse{
		Success: false,
		Error:   errInternal,
	}
	req := goodCreateRequest{}
	resp := goodCreateUpdateResponse{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logrus.Errorf("error decode request body [%s]", err.Error())
		writeResponse(w, badResponse, 500)
		return
	}

	err = validateRequest(req)
	if err != nil {
		logrus.Errorf("error validate request [%s]", err.Error())
		badResponse.Error = errWrongParams
		writeResponse(w, badResponse, 500)
		return
	}

	queryValues := r.URL.Query()
	projectId, err := strconv.Atoi(queryValues.Get("projectId"))
	if err != nil {
		logrus.Errorf("error convert projecIdParam [%s] to int [%s]", queryValues.Get("projectId"), err.Error())
		badResponse.Error = errWrongParams
		writeResponse(w, badResponse, 500)
		return
	}

	good := postgres.Good{
		ProjectID: projectId,
		Name:      req.Name,
	}
	err = good.Create()
	if err != nil {
		logrus.Errorf("error creating good [%s]", err.Error())
		writeResponse(w, badResponse, 500)
		return
	}

	resp.New(good)
	logrus.Infof("successfully created new good with id [%d]", good.ID)
	writeResponse(w, resp, 200)
}

func (resp *goodCreateUpdateResponse) New(good postgres.Good) {
	resp.Success = true
	resp.ID = good.ID
	resp.ProjectID = good.ProjectID
	resp.Name = good.Name
	resp.Description = good.Description
	resp.Priority = good.Priority
	resp.Removed = good.Removed
	resp.CreatedAt = good.CreatedAt
}

// GOOD UPDATE
type goodUpdateRequest struct {
	Name        string
	Description string
}

func goodUpdate(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("handling good update request...")
	badResponse := badResponse{
		Success: false,
		Error:   errInternal,
	}
	req := goodUpdateRequest{}
	resp := goodCreateUpdateResponse{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logrus.Errorf("error decode request body [%s]", err.Error())
		writeResponse(w, badResponse, 500)
		return
	}

	err = validateRequest(req)
	if err != nil {
		logrus.Errorf("error validate request [%s]", err.Error())
		badResponse.Error = errWrongParams
		writeResponse(w, badResponse, 500)
		return
	}

	queryValues := r.URL.Query()
	id, err := strconv.Atoi(queryValues.Get("id"))
	if err != nil {
		logrus.Errorf("error convert id [%s] to int [%s]", queryValues.Get("id"), err.Error())
		badResponse.Error = errWrongParams
		writeResponse(w, badResponse, 500)
		return
	}

	projectId, err := strconv.Atoi(queryValues.Get("projectId"))
	if err != nil {
		logrus.Errorf("error convert projecIdParam [%s] to int [%s]", queryValues.Get("projectId"), err.Error())
		badResponse.Error = errWrongParams
		writeResponse(w, badResponse, 500)
		return
	}

	good := postgres.Good{
		ID:        id,
		ProjectID: projectId,
	}
	err = good.Update(req.Name, req.Description)
	if err != nil {
		logrus.Errorf("error updating good with id=[%d], projectID=[%d] [%s]", id, projectId, err.Error())
		writeResponse(w, badResponse, 500)
		return
	}

	resp.New(good)
	logrus.Infof("successfully updated good with id [%d], projectID=[%d]", id, projectId)
	writeResponse(w, resp, 200)
}

// GOOD DELETE
type goodDeleteResponse struct {
	Success   bool `json:"success"`
	ID        int  `json:"id"`
	ProjectID int  `json:"projectId"` // не до конца понял что такое campaignId, но подумал что имелось ввиду projectId
	Removed   bool `json:"removed"`
}

func goodDelete(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("handling good delete request...")
	badResponse := badResponse{
		Success: false,
		Error:   errInternal,
	}
	resp := goodDeleteResponse{}

	queryValues := r.URL.Query()
	id, err := strconv.Atoi(queryValues.Get("id"))
	if err != nil {
		logrus.Errorf("error convert id [%s] to int [%s]", queryValues.Get("id"), err.Error())
		badResponse.Error = errWrongParams
		writeResponse(w, badResponse, 500)
		return
	}

	projectId, err := strconv.Atoi(queryValues.Get("projectId"))
	if err != nil {
		logrus.Errorf("error convert projecIdParam [%s] to int [%s]", queryValues.Get("projectId"), err.Error())
		badResponse.Error = errWrongParams
		writeResponse(w, badResponse, 500)
		return
	}

	good := postgres.Good{
		ID:        id,
		ProjectID: projectId,
	}
	err = good.Delete()
	if err != nil {
		logrus.Errorf("error deleting good with id=[%d], projectID=[%d] [%s]", id, projectId, err.Error())
		writeResponse(w, badResponse, 500)
		return
	}

	resp.New(good)
	logrus.Infof("successfully deleted good with id [%d], projectID=[%d]", id, projectId)
	writeResponse(w, resp, 200)
}

func (resp *goodDeleteResponse) New(good postgres.Good) {
	resp.Success = true
	resp.ID = good.ID
	resp.ProjectID = good.ProjectID
	resp.Removed = good.Removed
}

// GOODS LIST
type goodsListResponse struct {
	Success bool                       `json:"success"`
	Meta    meta                       `json:"meta"`
	Goods   []goodCreateUpdateResponse `json:"goods"`
}

type meta struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

func goodsList(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("handling goods list request...")
	badResponse := badResponse{
		Success: false,
		Error:   errInternal,
	}
	resp := goodsListResponse{}
	var err error

	queryValues := r.URL.Query()

	limit := 10
	limitParam := queryValues.Get("limit")
	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			logrus.Errorf("error convert limit [%s] to int [%s]", limitParam, err.Error())
			badResponse.Error = errWrongParams
			writeResponse(w, badResponse, 500)
			return
		}
	}

	offset := 0 // я думаю что если не указан, то имелось ввиду 0, а не 1. тк если 1, то пропускается 1ая запись
	offsetParam := queryValues.Get("offset")
	if offsetParam != "" {
		offset, err = strconv.Atoi(offsetParam)
		if err != nil {
			logrus.Errorf("error convert offset [%s] to int [%s]", offsetParam, err.Error())
			badResponse.Error = errWrongParams
			writeResponse(w, badResponse, 500)
			return
		}
	}

	goods := postgres.GoodSlice{}
	total, err := goods.Many(limit, offset)
	if err != nil {
		logrus.Errorf("error getting all goods [%s]", err.Error())
		writeResponse(w, badResponse, 500)
		return
	}

	resp.New(goods, total, limit, offset)
	logrus.Info("successfully got all goods")
	writeResponse(w, resp, 200)
}

func (resp *goodsListResponse) New(goodSlice postgres.GoodSlice, total, limit, offset int) {
	resp.Success = true

	removed := 0
	goods := []goodCreateUpdateResponse{}
	for _, good := range goodSlice {
		if good.Removed {
			removed++
		}

		goods = append(goods, goodCreateUpdateResponse{
			ID:          good.ID,
			ProjectID:   good.ProjectID,
			Name:        good.Name,
			Description: good.Description,
			Priority:    good.Priority,
			Removed:     good.Removed,
			CreatedAt:   good.CreatedAt,
		})
	}

	resp.Meta = meta{
		Total:   total,
		Removed: removed,
		Limit:   limit,
		Offset:  offset,
	}
	resp.Goods = goods
}

// GOOD REPRIORITIZE
type goodReprioritizeRequest struct {
	NewPriority int
}

type goodReprioritizeResponse struct {
	Success    bool       `json:"success"`
	Priorities []priority `json:"priorities"`
}

type priority struct {
	ID       int `json:"id"`
	Priority int `json:"priority"`
}

func goodReprioritize(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("handling good reprioritize request...")
	badResponse := badResponse{
		Success: false,
		Error:   errInternal,
	}
	req := goodReprioritizeRequest{}
	resp := goodReprioritizeResponse{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logrus.Errorf("error decode request body [%s]", err.Error())
		writeResponse(w, badResponse, 500)
		return
	}

	err = validateRequest(req)
	if err != nil {
		logrus.Errorf("error validate request [%s]", err.Error())
		badResponse.Error = errWrongParams
		writeResponse(w, badResponse, 500)
		return
	}

	queryValues := r.URL.Query()
	id, err := strconv.Atoi(queryValues.Get("id"))
	if err != nil {
		logrus.Errorf("error convert id [%s] to int [%s]", queryValues.Get("id"), err.Error())
		badResponse.Error = errWrongParams
		writeResponse(w, badResponse, 500)
		return
	}

	projectId, err := strconv.Atoi(queryValues.Get("projectId"))
	if err != nil {
		logrus.Errorf("error convert projecIdParam [%s] to int [%s]", queryValues.Get("projectId"), err.Error())
		badResponse.Error = errWrongParams
		writeResponse(w, badResponse, 500)
		return
	}

	good := postgres.Good{
		ID:        id,
		ProjectID: projectId,
	}
	goods := postgres.GoodSlice{}
	err = goods.Reprioritize(req.NewPriority, &good)
	if err != nil {
		logrus.Errorf("error reprioritizing goods from good with id=[%d], projectID=[%d] [%s]", id, projectId, err.Error())
		writeResponse(w, badResponse, 500)
		return
	}

	resp.New(goods)
	logrus.Infof("successfully reprioritized good with id [%d], projectID=[%d]. New priority=[%d]", id, projectId, req.NewPriority)
	writeResponse(w, resp, 200)
}

func (resp *goodReprioritizeResponse) New(goods postgres.GoodSlice) {
	resp.Success = true
	for _, good := range goods {
		resp.Priorities = append(resp.Priorities, priority{
			ID:       good.ID,
			Priority: good.Priority,
		})
	}
}
