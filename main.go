package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

// IMPORTANT: Use the domain where your account/key belongs
const providerURL = "https://ninbvnportal.com.ng/api/nin-verification"

var (
	apiKey   string
	ninRegex = regexp.MustCompile(`^\d{11}$`)
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	apiKey = os.Getenv("CHECKMYNIN_API_KEY")
	if apiKey == "" {
		log.Fatal("Missing CHECKMYNIN_API_KEY (set in .env or environment)")
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/lookup", lookupHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func jsonErr(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func lookupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		NIN     string `json:"nin"`
		Consent bool   `json:"consent"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonErr(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if !ninRegex.MatchString(req.NIN) {
		jsonErr(w, "Invalid NIN — must be exactly 11 digits", http.StatusBadRequest)
		return
	}
	if !req.Consent {
		jsonErr(w, "User consent is required", http.StatusBadRequest)
		return
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"nin":     req.NIN,
		"consent": true,
	})

	providerReq, err := http.NewRequest(http.MethodPost, providerURL, bytes.NewBuffer(payload))
	if err != nil {
		jsonErr(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	providerReq.Header.Set("Content-Type", "application/json")
	providerReq.Header.Set("x-api-key", apiKey)

	resp, err := http.DefaultClient.Do(providerReq)
	if err != nil {
		jsonErr(w, "Provider unreachable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		jsonErr(w, "Failed to read provider response", http.StatusInternalServerError)
		return
	}

	// Debug logging (only shown in terminal)
	log.Printf("Provider response [%d]: %s", resp.StatusCode, string(body))

	// Structure exactly matching provider's successful response
	var providerResp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Firstname        string `json:"firstname"`
			Middlename       string `json:"middlename"`
			Surname          string `json:"surname"`
			Telephoneno      string `json:"telephoneno"`
			ResidenceState   string `json:"residence_state"`
			ResidenceTown    string `json:"residence_town"`
			ResidenceAddress string `json:"residence_address"`
			ResidenceLGA     string `json:"residence_lga"`
			Birthcountry     string `json:"birthcountry"`
			Birthstate       string `json:"birthstate"`
			Birthlga         string `json:"birthlga"`
			Gender           string `json:"gender"`
			Nin              string `json:"nin"`
			Birthdate        string `json:"birthdate"`
			Photo            string `json:"photo"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &providerResp); err != nil {
		jsonErr(w, "Bad response from provider", http.StatusInternalServerError)
		return
	}

	if providerResp.Status != "success" {
		msg := providerResp.Message
		if msg == "" {
			msg = "Lookup failed"
		}
		jsonErr(w, msg, http.StatusBadRequest)
		return
	}

	d := providerResp.Data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"firstname":  d.Firstname,
		"surname":    d.Surname,
		"middlename": d.Middlename,
		"fullname":   d.Firstname + " " + d.Middlename + " " + d.Surname,
		"phone":      d.Telephoneno,
		"gender":     d.Gender,
		"birthdate":  d.Birthdate,
		"birthstate": d.Birthstate,
		"birthlga":   d.Birthlga,
		"address":    d.ResidenceAddress + ", " + d.ResidenceTown + ", " + d.ResidenceState,
		"photo":      d.Photo,
	})
}
