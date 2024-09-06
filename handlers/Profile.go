package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/EupravaProjekat/court/Models"
	"github.com/EupravaProjekat/court/Repo"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"log"
	"mime"
	"net/http"
	"strconv"
	"time"
)

type Courthandler struct {
	l    *log.Logger
	repo *Repo.Repo
}

func NewCourthandler(l *log.Logger, r *Repo.Repo) *Courthandler {
	return &Courthandler{l, r}

}

func (h *Courthandler) CheckIfUserExists(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	res := ValidateJwt(r, h.repo)
	if res == nil {
		err := errors.New("user doesnt exist")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)

}
func (h *Courthandler) NewUser(w http.ResponseWriter, r *http.Request) {

	res := ValidateJwt2(r, h.repo)

	rt, err := DecodeBodyUser(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}
	if res != rt.Email {
		err := errors.New("user doesnt exist")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	newUUID := uuid.New().String()
	rt.Uuid = newUUID
	rt.Role = "Guest"
	err = h.repo.NewUser(rt)
	if err != nil {
		log.Printf("Operation Failed: %v\n", err)
		w.WriteHeader(http.StatusNotAcceptable)
		_, err := w.Write([]byte("Profile not found"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (h *Courthandler) NewCase(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON body of the request
	var payload Models.Case
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate JWT and retrieve the user
	user := ValidateJwt(r, h.repo)
	if user == nil {
		http.Error(w, "User doesn't exist", http.StatusForbidden)
		return
	}

	// Generate a new unique ID for the case
	caseID := uuid.New().String()

	// Create a new Case instance
	newCase := Models.Case{
		ID:           caseID,
		Type:         payload.Type,
		Status:       "Open",                          // Default status is "Open"
		FilingDate:   time.Now().Format(time.RFC3339), // Set the filing date to the current time
		HearingDates: payload.HearingDates,            // Use the provided hearing dates
		Judge:        payload.Judge,
		Plaintiff:    payload.Plaintiff,
		Defendant:    payload.Defendant,
		Lawyers:      payload.Lawyers,
	}

	// Create a new Request instance to associate the case with the user
	newRequest := Models.Request{
		ID:          uuid.New().String(), // Generate a unique ID for the request
		Type:        "CaseCreation",      // Define the type of request
		Status:      "Pending",           // Default status
		Case:        newCase.ID,          // Store the case ID as a string reference
		Description: "New case created",
		CreatedAt:   time.Now().Format(time.RFC3339), // Set the creation time
	}

	// Append the new request to the user's Requests slice
	user.Requests = append(user.Requests, newRequest)

	// Persist the new case in the database (optional, depending on your repo implementation)
	err = h.repo.Update(user)
	if err != nil {
		http.Error(w, "Failed to save the case", http.StatusInternalServerError)
		return
	}

	// Return the created case as a JSON response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCase)
}
func (h *Courthandler) GetAllCases(w http.ResponseWriter, r *http.Request) {

	res := ValidateJwt(r, h.repo)
	if res == nil {
		err := errors.New("user doesnt exist")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if res.Role != "Operator" {
		err := errors.New("role error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	response, err := h.repo.GetAllCases()
	if err != nil {
		log.Printf("Operation Failed: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Cases not found"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, response)
}
func (h *Courthandler) GetallRequests(w http.ResponseWriter, r *http.Request) {

	res := ValidateJwt(r, h.repo)
	if res == nil {
		err := errors.New("user doesnt exist")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if res.Role != "Operator" {
		err := errors.New("role error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	response, err := h.repo.GetAllRequest()
	if err != nil {
		log.Printf("Operation Failed: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Requests not found"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, response)
}
func (h *Courthandler) GetProfile(w http.ResponseWriter, r *http.Request) {

	emaila := mux.Vars(r)["email"]
	ee := new(protos.ProfileRequest)
	ee.Email = emaila
	res := ValidateJwt(r, h.repo)
	if res == nil {
		err := errors.New("user doesnt exist")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	response, err := h.repo.GetByEmail(ee.Email)
	if err != nil || response == nil {
		log.Printf("Operation Failed: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Profile not found"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, response)
}

func (h *Courthandler) NewRequest(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	// Decode request body into GetRequest struct
	var req Models.GetRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Assign a new UUID to the request and set status
	newUUID := uuid.New().String()
	rt := Models.Request{
		ID:     newUUID,
		Status: "received",
		Case:   req.Uuid, // Assuming the case information comes from the UUID
	}

	// Prepare data to send to the external service (e.g., case-checking service)
	userData := []byte(`{"case_id":"` + req.Uuid + `"}`)

	// Placeholders for external service URL and secret
	apiUrl := "YOUR_EXTERNAL_SERVICE_URL" ///////////////////////////
	secretCode := "prosecution-service-secret-code"

	// Create a new HTTP POST request
	externalReq, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(userData))
	if err != nil {
		http.Error(w, "Failed to create external request", http.StatusInternalServerError)
		return
	}
	externalReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	externalReq.Header.Set("jwt", r.Header.Get("jwt")) // Forward JWT from the original request
	externalReq.Header.Set("secret", secretCode)

	// Send the request to the external service
	client := &http.Client{}
	externalResp, err := client.Do(externalReq)
	if err != nil {
		http.Error(w, "Error in external service communication", http.StatusInternalServerError)
		return
	}
	defer externalResp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(externalResp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from external service", http.StatusInternalServerError)
		return
	}

	// Unmarshal response into a Response struct
	var resp Models.Response
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		http.Error(w, "Error parsing external service response", http.StatusInternalServerError)
		return
	}

	// Update the request's `CaseStatus` field based on external response
	rt.Description = "Case status: " + strconv.FormatBool(resp.CaseStatus)

	// Validate user (you may have a function like `ValidateJwt`)
	user := ValidateJwt(r, h.repo)
	if user == nil {
		http.Error(w, "User doesn't exist", http.StatusForbidden)
		return
	}

	// Update user with the new request
	user.Requests = append(user.Requests, rt)
	err = h.repo.Update(user)
	if err != nil {
		log.Printf("Failed to update user data: %v\n", err)
		http.Error(w, "Couldn't add request", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Request successfully processed"))
}
func (h *Courthandler) GetRequest(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	rt, err := DecodeBody2(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := ValidateJwt(r, h.repo)
	if res == nil {
		err := errors.New("user doesnt exist")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	respon, err := h.repo.GetRequest(rt.Uuid)
	if err != nil {
		log.Printf("Operation failed: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("couldn't find request"))
		if err != nil {
			return
		}
		return
	}
	RenderJSON(w, respon)
	w.WriteHeader(http.StatusOK)
}
func (h *Courthandler) CheckIfPersonIsProsecuted(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request body to get the email or UUID
	var req struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare the request payload to send to the external /prosecution service
	payload, err := json.Marshal(map[string]string{"email": req.Email})
	if err != nil {
		http.Error(w, "Failed to prepare request", http.StatusInternalServerError)
		return
	}

	// Define the external service URL (this is where you call the external API)
	externalServiceUrl := "http://localhost:9199/prosecute"

	// Create the POST request to the external service
	externalReq, err := http.NewRequest("POST", externalServiceUrl, bytes.NewBuffer(payload))
	if err != nil {
		http.Error(w, "Failed to create request to external service", http.StatusInternalServerError)
		return
	}
	externalReq.Header.Set("Content-Type", "application/json")

	// Send the request to the external service
	client := &http.Client{}
	resp, err := client.Do(externalReq)
	if err != nil {
		http.Error(w, "Failed to communicate with external service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response from the external service
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from external service", http.StatusInternalServerError)
		return
	}

	// Forward the response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode) // Use the same status code from the external service
	w.Write(body)                  // Write the external response back to the client
}
