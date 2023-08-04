package main

import (
	"chess-profile/apis"
	"chess-profile/calc"
	"chess-profile/format"
	"chess-profile/types"
	"net/http"
	"time"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

var apiMap = map[string]types.ChessApi{
	"lichess": apis.NewLichess(httpClient),
}

type HTTPServer struct{}

func (hh HTTPServer) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	if len(params["site"]) == 0 || len(params["user"]) == 0 {
		writer.WriteHeader(400)
		return
	}

	site := params["site"][0]
	user := params["user"][0]
	api, ok := apiMap[site]
	if !ok {
		writer.Write([]byte("Bad site!"))
		return
	}
	games, err := api.GetGames(user)
	if err != nil {
		writer.Write([]byte("Failed to read user games!\n" + string(err.Error())))
		return
	}
	summary := calc.NewCalculator(2, 3).Calc(user, games)
	html := format.Format(summary)
	writer.Header()["Content-Type"] = []string{"text/html"}
	writer.Write([]byte(html))
}

func main() {
	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      HTTPServer{},
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
