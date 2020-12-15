package main

import (
	"fmt"
)

// type request struct {
// 	Input *[]string `json:"input"`
// }

// type response struct {
// 	TotalStars   int            `json:"totalStars,omitempty"`
// 	InvalidRepos []string       `json:"invalidRepos,omitempty"`
// 	ValidRepos   map[string]int `json:"validRepos,omitempty"`
// }

// func respondWithError(w http.ResponseWriter, code int, message string) {
// 	respondWithJSON(w, code, map[string]string{"error": message})
// }

// func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
// 	response, _ := json.Marshal(payload)

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	w.Write(response)
// }

// func getStars(w http.ResponseWriter, r *http.Request) {
// 	// b082d2cd0e5e4202f31a
// 	// 76b2b47503154546fc0393b70c3b488d3f6d66a1

// 	// should be in app.go
// 	ctx := context.Background()
// 	// ts := oauth2.StaticTokenSource(
// 	// 	&oauth2.Token{AccessToken: "08dcd33b87acba14d8630efbf2ae2736d885ad53"},
// 	// )
// 	//tc := oauth2.NewClient(ctx, ts)

// 	//client := github.NewClient(tc)
// 	client := github.NewClient(nil)

// 	opt := &github.RepositoryListByOrgOptions{
// 		ListOptions: github.ListOptions{PerPage: 30},
// 	}

// 	// parse payload
// 	var p request
// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&p); err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
// 		return
// 	}
// 	defer r.Body.Close()
// 	log.Printf("request: %v", *p.Input)

// 	totalCount := 0
// 	validRepos := make(map[string]int)
// 	invalidRepos := make([]string, 0)
// 	var resp response
// 	seenOrgs := make(map[string][]*github.Repository)
// 	for _, input := range *p.Input {
// 		fmt.Printf("Processing input %v\n", input)
// 		token := strings.Split(input, "/")
// 		allRepos, ok := seenOrgs[token[0]]
// 		if !ok {
// 			fmt.Printf("%s is not in the map\n", token[0])
// 			results, err := listAllReposForAnOrg(ctx, token[0], client, opt)
// 			if err == nil {
// 				fmt.Printf("Find all repos for %s\n", token[0])
// 				allRepos = results
// 				seenOrgs[token[0]] = results

// 			} else {
// 				// invalid org
// 				fmt.Printf("%s is not a valid org\n", token[0])
// 				invalidRepos = append(invalidRepos, input)
// 				continue
// 			}

// 		}

// 		if count := getStarsForRepo(allRepos, token[1]); count != -1 {
// 			validRepos[input] = count
// 			totalCount += count
// 		} else {
// 			// invalid repo
// 			fmt.Printf("%s is not a valid repo\n", token[1])
// 			invalidRepos = append(invalidRepos, input)
// 		}

// 	}

// 	resp.TotalStars = totalCount
// 	resp.ValidRepos = validRepos
// 	resp.InvalidRepos = invalidRepos
// 	respondWithJSON(w, http.StatusOK, resp)
// }

// func listAllReposForAnOrg(ctx context.Context, orgName string, client *github.Client, opt *github.RepositoryListByOrgOptions) ([]*github.Repository, error) {
// 	var allRepos []*github.Repository
// 	for {
// 		repos, resp, err := client.Repositories.ListByOrg(ctx, orgName, opt)
// 		if err != nil {
// 			fmt.Printf("the organization %s does not exist %v\n", orgName, err)
// 			return nil, errors.New("the organization does not exist")
// 		}
// 		allRepos = append(allRepos, repos...)
// 		if resp.NextPage == 0 {
// 			break
// 		}
// 		opt.Page = resp.NextPage
// 	}
// 	return allRepos, nil

// }

// func getStarsForRepo(allRepos []*github.Repository, repoName string) int {
// 	//fmt.Printf("repos %v", allRepos)

// 	for _, repo := range allRepos {
// 		if repo.GetName() == repoName {
// 			count := repo.GetStargazersCount()
// 			fmt.Printf("Found: %s\n", repoName)
// 			return count
// 		}
// 	}
// 	return -1
// }

func main() {
	fmt.Println("... starts server")
	// router := mux.NewRouter().StrictSlash(true)
	// router.HandleFunc("/get-stars", getStars).Methods("POST")
	// //router.HandleFunc("/getStars", getStars)
	// log.Fatal(http.ListenAndServe(":8080", router))
	server := Server{}
	server.StartServer()
}
