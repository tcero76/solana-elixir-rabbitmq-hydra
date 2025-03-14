package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	hydra "github.com/ory/hydra-client-go/v2"
	"golang.org/x/oauth2"
)

const (
	hydraAdminURL  = "http://hydra:4445"
	hydraPublicURL = "http://hydra:4444"
	clientID       = "657d0bb0-2314-4c8b-b649-1525af797d72" // Reemplaza con el ID del cliente
	clientSecret   = "Wfnum75JFGpSEUTLqLiMbw1gad"           // Reemplaza con el secreto
)

var hydraAdminClient *hydra.APIClient

func main() {
	// Configurar cliente de Hydra
	cfg := hydra.NewConfiguration()
	cfg.Servers = hydra.ServerConfigurations{{URL: hydraAdminURL}}
	hydraAdminClient = hydra.NewAPIClient(cfg)

	e := echo.New()

	// Ruta para iniciar el login
	e.GET("/login", handleLogin)
	e.GET("/consent", handleConsent)
	// Ruta para manejar el callback de Hydra
	e.GET("/callback", handleCallback)
	// Ruta protegida para el frontend
	e.GET("/getToken", getToken)

	e.Start(":" + os.Getenv("PORT"))
}

func handleLogin(c echo.Context) error {
	loginChallenge := c.QueryParam("login_challenge")
	if loginChallenge == "" {
		return c.String(http.StatusBadRequest, "No login challenge provided")
	}

	// Aquí deberías implementar tu lógica de autenticación (p.ej., formulario de login)
	// Por simplicidad, asumimos que el usuario ya está autenticado
	_, _, err := hydraAdminClient.OAuth2API.GetOAuth2LoginRequest(c.Request().Context()).LoginChallenge(loginChallenge).Execute()
	if err != nil {
		return err
	}

	// Aceptar la solicitud de login
	acceptBody := hydra.AcceptOAuth2LoginRequest{
		Subject: "test-user", // ID del usuario autenticado
	}
	resp, _, err := hydraAdminClient.OAuth2API.AcceptOAuth2LoginRequest(c.Request().Context()).
		LoginChallenge(loginChallenge).
		AcceptOAuth2LoginRequest(acceptBody).
		Execute()
	// fmt.Printf("LOGIN: La respuesta es: %s \n", resp)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, resp.GetRedirectTo())
}

func handleConsent(c echo.Context) error {
	consentChallenge := c.QueryParam("consent_challenge")
	if consentChallenge == "" {
		return c.String(http.StatusBadRequest, "No consent_challenge provided")
	}

	// Obtener la solicitud de consentimiento
	_, _, err := hydraAdminClient.OAuth2API.GetOAuth2ConsentRequest(c.Request().Context()).
		ConsentChallenge(consentChallenge).
		Execute()
	if err != nil {
		fmt.Printf("Error getting consent request: %v\n", err)
		return c.String(http.StatusInternalServerError, "Error processing consent request")
	}
	var boolPtr = func(b bool) *bool {
		return &b
	}
	var int64Ptr = func(i int64) *int64 {
		return &i
	}
	acceptConsentBody := hydra.AcceptOAuth2ConsentRequest{
		GrantScope:  []string{"openid", "offline"}, // Scopes otorgados
		Remember:    boolPtr(true),
		RememberFor: int64Ptr(3600),
	}
	resp, _, err := hydraAdminClient.OAuth2API.AcceptOAuth2ConsentRequest(c.Request().Context()).
		ConsentChallenge(consentChallenge).
		AcceptOAuth2ConsentRequest(acceptConsentBody).
		Execute()

	if err != nil {
		fmt.Printf("Error accepting consent request: %v\n", err)
		return c.String(http.StatusInternalServerError, "Error accepting consent")
	}

	// fmt.Printf("Consent accepted, redirecting to: %s\n", resp.GetRedirectTo())
	return c.Redirect(http.StatusFound, resp.GetRedirectTo())
}

type debugTransport struct{}

func (t *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Println("Request URL:", req.URL)
	fmt.Println("Headers:", req.Header)
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println("Error leyendo el body:", err)
		} else {
			fmt.Println("Body:", string(bodyBytes))
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}
	return http.DefaultTransport.RoundTrip(req)
}

var token *oauth2.Token

func handleCallback(c echo.Context) error {
	var redirect = "http://" + os.Getenv("HOST_EXTERNAL") + ":" + os.Getenv("PORT_EXTERNAL") + os.Getenv("RedirectURL")
	code := c.QueryParam("code")
	if code == "" {
		return c.String(http.StatusBadRequest, "No code provided")
	}
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirect,
		Endpoint: oauth2.Endpoint{
			AuthURL:  hydraPublicURL + "/oauth2/auth",
			TokenURL: hydraPublicURL + "/oauth2/token",
		},
		Scopes: []string{"openid", "offline"},
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
		Transport: &debugTransport{},
	})
	var err error
	token, err = conf.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("Error intercambiando código por token: %v", err)
	}
	log.Printf("Access Token: %s", token.AccessToken)
	return c.Redirect(http.StatusFound, "http://"+os.Getenv("HOST_EXTERNAL")+":"+os.Getenv("PORT_EXTERNAL")+"/home")
}

func getToken(c echo.Context) error {
	return c.JSON(http.StatusOK, token.AccessToken)
}
