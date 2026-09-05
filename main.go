package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-http-server/internal/users"
	"go-http-server/internal/util"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const PORT = 3000

type UserData struct {
	FirstName string
	LastName  string
	Email     string
}

type Server struct {
	userManager *users.Manager
}

func main() {
	manager := users.NewManager()
	defer manager.Shutdown()

	s := Server{manager}
	mux := http.NewServeMux()
	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", PORT), Handler: mux}

	mux.HandleFunc("/{$}", HandleRootEndpoint)
	mux.HandleFunc("/goodbye", HandleGoodbyeEndpoint)
	mux.HandleFunc("/hello", HandleHelloEndpoint)
	mux.HandleFunc("/param/{name}", HandleParamEndpoint)
	mux.HandleFunc("/header", HandleHeaderEndpoint)
	mux.HandleFunc("POST /json", HandleJSONEndpoint)
	mux.HandleFunc("POST /user", s.HandleUserEndpointPOST)
	mux.HandleFunc("GET /user", s.HandleUserEndpointGET)

	go func() {
		slog.Info("Listening to", "port", PORT)
		err := httpServer.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server error", "err", err)
			os.Exit(1)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	wg.Go(func() {
		defer wg.Done()

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT)

		<-sc
		slog.Info("server shutting down")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		err := httpServer.Shutdown(shutdownCtx)

		if err != nil {
			slog.Error("error shutting down server", "err", err)
		}
	})

	wg.Wait()
	slog.Info("server shutdown comple")
}

func HandleRootEndpoint(res http.ResponseWriter, req *http.Request) {
	_, err := res.Write([]byte("Welcome!"))

	if err != nil {
		slog.Error("error writing response body:", "err", err)
		return
	}
}

func HandleGoodbyeEndpoint(res http.ResponseWriter, _ *http.Request) {
	_, err := res.Write([]byte("Goodbye!"))

	if err != nil {
		slog.Error("error writing response body:", "err", err)
		return
	}
}

func HandleHelloEndpoint(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	nameList, ok := query["name"]

	if !ok {
		http.Error(res, "You need to provide a search query named \"name\"", http.StatusBadRequest)
		return
	}

	name := nameList[0]

	util.WriteResponseBody(res, "Hello, ", name, "!")
}

func HandleParamEndpoint(res http.ResponseWriter, req *http.Request) {
	name := req.PathValue("name")
	util.WriteResponseBody(res, "Hello, ", name, "!")
}

func HandleHeaderEndpoint(res http.ResponseWriter, req *http.Request) {
	name := req.Header.Get("name")

	if name == "" {
		http.Error(res, "You must set a value for \"name\" at the request header", http.StatusBadRequest)
		return
	}

	util.WriteResponseBody(res, "Hello, ", name, "!")
}

func HandleJSONEndpoint(res http.ResponseWriter, req *http.Request) {
	bodyInBytes, err := io.ReadAll(req.Body)

	if err != nil {
		slog.Error("error reading request body:", "err", err)
		http.Error(res, "error reading request body", http.StatusInternalServerError)
		return
	}

	var unmarshalledReqBody UserData
	err = json.Unmarshal(bodyInBytes, &unmarshalledReqBody)

	if err != nil {
		slog.Error("error deserialising request body:", "err", err)
		http.Error(res, "error deserialising request body", http.StatusBadRequest)
		return
	}

	if unmarshalledReqBody.FirstName == "" {
		http.Error(res, `You must not leave the "name" field empty`, http.StatusBadRequest)
		return
	}

	util.WriteResponseBody(res, "Hello, ", unmarshalledReqBody.FirstName, "!")
}

func (s *Server) HandleUserEndpointPOST(res http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(res, fmt.Sprintf(`unsupported Content-Type header: %q`, contentType), http.StatusUnsupportedMediaType)
		return
	}

	reqBody := http.MaxBytesReader(res, req.Body, 100000) // 100 KB

	decoder := json.NewDecoder(reqBody)
	decoder.DisallowUnknownFields()

	var u UserData

	err := decoder.Decode(&u)
	if err != nil {
		slog.Error("error decoding request body to /user", "err", err)
		http.Error(res, "bad request body", http.StatusBadRequest)
		return
	}

	err = s.userManager.AddUser(u.FirstName, u.LastName, u.Email)
	if err != nil {
		http.Error(res, fmt.Sprintf("error adding user: %v", err), http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

func (s *Server) HandleUserEndpointGET(res http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	firstName, lastName := q.Get("firstName"), q.Get("lastName")

	if firstName == "" || lastName == "" {
		http.Error(res, `Netheir "firstName" or "lastName" search query can be empty`, http.StatusBadRequest)
		return
	}

	result, err := s.userManager.GetUserByName(firstName, lastName)

	if err != nil {
		if errors.Is(err, users.ErrNoUserFound) {
			http.Error(res, "No user found", http.StatusNotFound)
		} else {
			slog.Error("error retrieving user", "err", err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		return
	}

	converted := convertUserToUserData(result)
	marshalled, err := json.Marshal(converted)

	if err != nil {
		slog.Error("error marshalling user data", "err", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(marshalled)

	if err != nil {
		slog.Error("error writing response body", "err", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func convertUserToUserData(u *users.User) *UserData {
	return &UserData{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email.Address,
	}
}
