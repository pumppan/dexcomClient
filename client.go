package dexcomClient

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

type DexcomClient struct {
	AuthCode    string
	DexcomToken string
	config      *Config
	oAuthToken  *Token
	logger
}

func NewClient(config *Config) *DexcomClient {
	dc := &DexcomClient{
		config: config,
		logger: &defaultLogger{config: config},
	}

	if config.IsDev {
		fmt.Println("Dev server starting on :8000")
		url := config.getBaseUrl() + "/v1/oauth2/login?client_id=" + config.ClientId + "&redirect_uri=" + config.RedirectURI + "&response_type=code&scope=offline_access"
		fmt.Println(url)
		defer dc.startDevServer(url)
	}
	return dc
}

func NewClientWithToken(config *Config, token *Token) *DexcomClient {
	return &DexcomClient{
		config:     config,
		oAuthToken: token,
		logger:     &defaultLogger{config: config},
	}
}

func (client *DexcomClient) startDevServer(url string) {
	server := &http.Server{Addr: ":8000"}

	router := mux.NewRouter()
	router.Path("/oauth").Queries("code", "{code}").HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		client.AuthCode = req.FormValue("code")
		_, err := client.GetOauthToken()
		if err != nil {
			panic(err)
		}
		server.Shutdown(context.Background())
	})

	server.Handler = router
	exec.Command("open", url)
	server.ListenAndServe()
}
