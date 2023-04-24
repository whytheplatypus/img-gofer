package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/oauth2"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var (
	clientID     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	conf         = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/photoslibrary.readonly"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
		RedirectURL: "http://localhost",
	}
)

// login uses the Google OAuth client credentials flow to
// retrieve a token and create an authenticated HTTP client.
// The HTTP client returned by login will automatically refresh
// the token as necessary.
func login(ctx context.Context, conf *oauth2.Config) *http.Client {

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	return conf.Client(ctx, tok)
}

type mediaItem struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	BaseUrl     string `json:"baseUrl"`
	MimeType    string `json:"mimeType"`
	Filename    string `json:"filename"`
}

type page struct {
	Items         []mediaItem `json:"mediaItems"`
	NextPageToken string      `json:"nextPageToken"`
}

type library struct {
	Items []mediaItem
}

func fetchLibrary(client *http.Client) (*library, error) {
	lib := &library{
		Items: []mediaItem{},
	}

	nextPageToken := ""

	for {
		page, err := fetchPage(client, nextPageToken)
		if err != nil {
			return lib, err
		}
		lib.Items = append(lib.Items, page.Items...)
		if page.NextPageToken == "" {
			break
		}
		nextPageToken = page.NextPageToken
	}
	return lib, nil
}

func fetchPage(client *http.Client, pageToken string) (*page, error) {
	photosListUrl, err := url.Parse("https://photoslibrary.googleapis.com/v1/mediaItems")
	if err != nil {
		return nil, err
	}
	photosListUrl.RawQuery = fmt.Sprintf("pageToken=%s", pageToken)
	resp, err := client.Get(photosListUrl.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// decode body into a library struct
	p := &page{}
	if err := json.Unmarshal(body, p); err != nil {
		return p, err
	}
	log.Println(p)
	return p, nil
}

func main() {
	ctx := context.Background()
	client := login(ctx, conf)
	lib, err := fetchLibrary(client)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Downloading %d images\n", len(lib.Items))
	for _, item := range lib.Items {
		log.Printf("Downloading %s\n", item.Filename)
		resp, err := client.Get(fmt.Sprintf("%s=d", item.BaseUrl))
		if err != nil {
			log.Fatal(err)
		}
		// save response body to file
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile(item.Filename, body, 0644); err != nil {
			log.Fatal(err)
		}
	}
}
