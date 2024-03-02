package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func checkName(name string) bool {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://test-api.ghole.xyz/check_name?repo_name="+name, nil)
	if err != nil {
		return false
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	return resp.StatusCode == 200
}

type Repo struct {
	UserNpub string `json:"user_npub"`
	RepoName string `json:"repo_name"`
}

func createRepo(repo Repo) bool {
	client := &http.Client{}
	b, err := json.Marshal(repo)
	if err != nil {
		fmt.Println(err)
		return false
	}
	req, err := http.NewRequest("POST", "http://test-api.ghole.xyz/deploy", bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Printf("%#v", resp.Body)
	fmt.Println(resp.StatusCode)
	return resp.StatusCode == 200
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to our API!")
	})
	router.HandleFunc("/checkname/{name}", checkNameRoute).Methods("GET")
	//router.HandleFunc("/createrepo/{name}", createRepoRoute).Methods("GET")
	router.HandleFunc("/new", createRepoPostRoute).Methods("POST")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	http.ListenAndServe(":8090", handler)
}

func createRepoPostRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var repo Repo
	err := json.NewDecoder(r.Body).Decode(&repo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n\n", repo)
	json.NewEncoder(w).Encode(createRepo(repo))

}

func checkNameRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	fmt.Printf("%#v", params["name"])
	if len(params["name"]) > 0 {
		json.NewEncoder(w).Encode(checkName(params["name"]))
		return
	}
	json.NewEncoder(w).Encode(false)
}

func createRepoRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	fmt.Printf("%#v", params["name"])
	if len(params["name"]) > 0 {
		if checkName(params["name"]) {
			repo := Repo{
				UserNpub: "npub1mygerccwqpzyh9pvp6pv44rskv40zutkfs38t0hqhkvnwlhagp6s3psn5p",
				RepoName: "test_repo_01",
			}
			json.NewEncoder(w).Encode(createRepo(repo))
			return
		}
	}
	json.NewEncoder(w).Encode(false)
}