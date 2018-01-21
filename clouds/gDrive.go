package clouds

import (
	"log"
	"fmt"
	"os/user"
	"path/filepath"
	"os"
	"net/url"
	"encoding/json"
	"io/ioutil"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
	"golang.org/x/oauth2"
	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
	"strconv"
)
type GDriveHandler struct{
	Name string
	FileList map[int] string
}

var srv *drive.Service
func (g *GDriveHandler) Init() string{
	ctx := context.Background()

	b, err := ioutil.ReadFile("resources/client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/fullDrive.json
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client, authUrl := getClient(ctx, config)

	srv, err = drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)
	}

	return authUrl
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) (*http.Client, string) {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	var authUrl string
	if err != nil {
		tok, authUrl = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok), authUrl
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, string) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)


	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok, authURL
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("fullDrive.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func (g *GDriveHandler) GetFileList() map[int]string {
	fileList := make(map[int] string)
	r, err := srv.Files.List().Do()
	if err != nil{
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	for _, i := range r.Files {
		d, err := strconv.Atoi(i.Id)
		if err != nil{
			log.Fatalf("Could not convert file id")
		}
		fileList[d] = i.Name + i.FileExtension
	}
	g.FileList = fileList
	return fileList
}

func (g *GDriveHandler) DownloadById(id int){
	idx := strconv.Itoa(id)
	out, err := srv.Files.Get(idx).Download()
	if err != nil{
		panic(err.Error())
	}
	final, err := os.Create(g.FileList[id])
	if err != nil{
		panic(err.Error())
	}
	defer final.Close()
	io.Copy(final, out.Body)
}

func (g *GDriveHandler) DownloadByName(name string) {
	var id int
	for i,j := range g.FileList{
		if j == name{
			id = i
			break
		}
	}
	idx := strconv.Itoa(id)
	out, err := srv.Files.Get(idx).Download()
	if err != nil{
		panic(err.Error())
	}
	final, err := os.Create(g.FileList[id])
	if err != nil{
		panic(err.Error())
	}
	defer final.Close()
	io.Copy(final, out.Body)
}

func (g * GDriveHandler) Upload(path, name string){
	in, err := os.Open(path)
	if err != nil{
		panic(err.Error())
	}
	_, err = srv.Files.Create(&drive.File{Name: name}).Media(in).Do()
	if err != nil{
		panic(err.Error())
	}
}