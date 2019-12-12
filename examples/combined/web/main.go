package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/tpryan/mlauth/speech"

	"github.com/pkg/errors"
	"github.com/tpryan/mlauth/language"
	"github.com/tpryan/mlauth/vision"
)

var errorAuthFalse = errors.New("the content did not pass authentication")

const keyVision = "Golden Retriever"
const keyLanguage = "location"
const keySpeech = "Brooklyn"
const positiveLanguage = true
const tokenVision = "4A56D6"
const tokenSpeech = "7C45E0"
const tokenLanguage = "6D64A5"
const secret = "Your Google Cloud credit redemption code is 4A56D6-6D64A5-7C45E0"

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/healthz", handleHealth)
	http.HandleFunc("/keys", handleKeys)
	http.HandleFunc("/auth/language", handleLanguage)
	http.HandleFunc("/auth/vision", handleVision)
	http.HandleFunc("/auth/speech", handleSpeech)
	http.HandleFunc("/auth/secret", handleSecret)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// AuthResponse is the response from a successful authentication attempt.
type AuthResponse struct {
	Auth  bool   `json:"auth"`
	Token string `json:"token"`
}

// TokenResponse is the response from a successful request for the secret.
type TokenResponse struct {
	Result bool   `json:"result"`
	Secret string `json:"secret"`
}

// String outputs the response as a json string.
func (t *TokenResponse) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	resp := string(b)
	return resp
}

//Keys is the structure that represents the list of keys for the front end app.
type Keys struct {
	Vision   string `json:"vision"`
	Speech   string `json:"speech"`
	Language string `json:"language"`
}

// String outputs the response as a json string.
func (k *Keys) String() string {
	b, err := json.Marshal(k)
	if err != nil {
		return ""
	}
	resp := string(b)
	return resp
}

func handleKeys(w http.ResponseWriter, r *http.Request) {
	resp := Keys{Vision: keyVision, Speech: keySpeech, Language: keyLanguage}

	if positiveLanguage {
		resp.Language = resp.Language + ",positive"
	} else {
		resp.Language = resp.Language + ",negative"
	}

	writeResponse(w, http.StatusOK, resp.String())
	return
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, http.StatusOK, "{\"status\":\"ok\"}")
	return
}

func handleSecret(w http.ResponseWriter, r *http.Request) {
	visionSub := r.FormValue("token_vision")
	languageSub := r.FormValue("token_language")
	speechSub := r.FormValue("token_speech")

	rep := TokenResponse{Result: false, Secret: ""}

	if visionSub == tokenVision &&
		languageSub == tokenLanguage &&
		speechSub == tokenSpeech {
		rep.Result = true
		rep.Secret = secret
		writeResponse(w, http.StatusOK, rep.String())
		return

	}
	writeResponse(w, http.StatusUnauthorized, rep.String())
	return
}

func handleLanguage(w http.ResponseWriter, r *http.Request) {
	sentence := r.FormValue("sentence")
	auth, err := authLanguage(sentence)
	if err != nil {
		writeResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	rep := AuthResponse{Auth: auth, Token: tokenLanguage}

	b, err := json.Marshal(rep)
	if err != nil {
		writeResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	resp := string(b)

	writeResponse(w, http.StatusOK, resp)
	return

}

func handleVision(w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("picture")
	if err != nil {
		tmp := errors.Wrap(err, "could not get file from form:")
		writeResponse(w, http.StatusInternalServerError, tmp.Error())
		return
	}
	defer file.Close()

	auth, err := authVision(file)
	if err != nil {
		tmp := errors.Wrap(err, "could not auth vision:")
		writeResponse(w, http.StatusUnauthorized, tmp.Error())
		return
	}
	rep := AuthResponse{Auth: auth, Token: tokenVision}

	b, err := json.Marshal(rep)
	if err != nil {
		writeResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	resp := string(b)

	writeResponse(w, http.StatusOK, resp)
	return

}

func handleSpeech(w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("audio")
	if err != nil {
		tmp := errors.Wrap(err, "could not get file from form:")
		writeResponse(w, http.StatusInternalServerError, tmp.Error())
		return
	}
	defer file.Close()

	auth, err := authSpeech(file)
	if err != nil {
		tmp := errors.Wrap(err, "could not auth vision:")
		writeResponse(w, http.StatusUnauthorized, tmp.Error())
		return
	}
	rep := AuthResponse{Auth: auth, Token: tokenSpeech}

	b, err := json.Marshal(rep)
	if err != nil {
		writeResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	resp := string(b)

	writeResponse(w, http.StatusOK, resp)
	return

}

func authLanguage(content string) (bool, error) {

	result, err := language.Auth(keyLanguage, content, true)

	if err != nil {
		return false, errors.Wrap(err, "could not authenticate content: ")
	}

	if result {
		return true, nil
	}
	return false, errorAuthFalse

}

func authVision(file io.Reader) (bool, error) {

	result, err := vision.AuthFromReader(keyVision, file)

	if err != nil {
		return false, errors.Wrap(err, "could not authenticate content: ")
	}

	if result {
		return true, nil
	}
	return false, errorAuthFalse
}

func authSpeech(file io.Reader) (bool, error) {

	result, err := speech.AuthFromReader(keySpeech, file)

	if err != nil {
		return false, errors.Wrap(err, "could not authenticate content: ")
	}

	if result {
		return true, nil
	}
	return false, errorAuthFalse
}

func writeResponse(w http.ResponseWriter, code int, msg string) {

	if code != http.StatusOK {
		log.Printf(msg)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write([]byte(msg))

	return
}
