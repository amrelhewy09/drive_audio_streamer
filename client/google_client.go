package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

func NewGoogleDriveClient() RemoteFileClient {

	return &DriveClient{
		client: GoogleConnect(createConfig()),
	}
}

func GoogleConnect(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func (c *DriveClient) ConfigureRemoteFileDestination() error {
	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithHTTPClient(c.client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	c.driveService = srv
	return nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func createConfig() *oauth2.Config {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, drive.DriveReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	config.RedirectURL = "http://localhost:8080/callback"
	return config
}

func handleCallback(w http.ResponseWriter, r *http.Request, cancel context.CancelFunc) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Authorization code: %v\n", code)

	cancel()
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	go bootServerForAuthCode()

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return token
}

func (d *DriveClient) List() ([]*File, error) {
	query := "mimeType='audio/mpeg' or mimeType='audio/wav' or mimeType='audio/x-wav' or mimeType='audio/ogg'"
	r, err := d.driveService.Files.List().Q(query).Fields("items(id,title,fileSize)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	files := make([]*File, len(r.Items))
	for i, item := range r.Items {
		files[i] = &File{Name: item.Title, Id: item.Id, Size: item.FileSize}
	}

	return files, nil
}

func (d *DriveClient) Stream(id string, start int64, end int64) (io.ReadCloser, error) {
	req := d.driveService.Files.Get(id)
	if start >= 0 && end > start {
		req.Header().Set("Range", fmt.Sprintf("bytes=%d-%d", start, end-1))
	} else {
		return nil, fmt.Errorf("invalid range")
	}
	resp, err := req.Download()
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func bootServerForAuthCode() {
	ctx, cancel := context.WithCancel(context.Background())
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		handleCallback(w, r, cancel)
	})
	server := &http.Server{Addr: ":8080"}
	go cancelContext(ctx, server)
	server.ListenAndServe()

}

func cancelContext(ctx context.Context, server *http.Server) {
	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
