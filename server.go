package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	// "os/user"
	"flag"
	"path/filepath"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"

	"github.com/sirupsen/logrus"
)

const (
	longDateForm  = "Jan 2, 2006 at 3:04pm (MST)"
	shortDateForm = "2006-01-02"
)

var (
	clientSecret string
	port         string
	currentTime  time.Time
)

type eventlist struct {
	IsDue bool
	Items []event
}
type event struct {
	Name string
	Date string
}

// forecastHandler takes a forecast.Request object
// and passes it to the forecast.io API
func calendarHandler(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	b, err := ioutil.ReadFile(clientSecret)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read client secret file: %v", err), http.StatusInternalServerError)
		return
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/calendar-go-quickstart.json
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {

		http.Error(w, fmt.Sprintf("Unable to parse client secret file to config: %v", err), http.StatusInternalServerError)
		return
	}
	client := getClient(ctx, config)

	srv, err := calendar.New(client)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve calendar Client %v", err), http.StatusInternalServerError)
		return
	}

	var calendarID = r.URL.Query().Get("calendarId")
	if calendarID == "" {
		calendarID = "primary"
	}
	currentTime = time.Now()
	t := currentTime.Format(time.RFC3339)
	events, err := srv.Events.List(calendarID).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()

	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve next ten of the user's events: %v", err), http.StatusInternalServerError)
		return
	}

	var recurringEvents = map[string]string{}

	var showEvents = r.URL.Query().Get("show")
	var countDown = r.URL.Query().Get("countDown")

	res := eventlist{}
	res.IsDue = (countDown != "")
	for _, i := range events.Items {
		if _, ok := recurringEvents[i.RecurringEventId]; ok && showEvents != "all" {
			continue
		}
		recurringEvents[i.RecurringEventId] = i.Id
		newItem := event{
			Name: i.Summary,
			Date: getCountDown(i.Start, countDown),
		}
		res.Items = append(res.Items, newItem)
	}
	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(js))
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getCountDown(start *calendar.EventDateTime, countDown string) string {
	var when string
	if start.DateTime != "" {
		when = start.DateTime
	} else {
		when = start.Date
	}

	due, err := time.Parse(longDateForm, when)
	if err != nil {
		due, err = time.Parse(shortDateForm, when)
	}
	if err != nil {
		return when
	}
	diff := due.Sub(currentTime)
	if countDown != "days" {
		when = diff.String()
	} else {
		var diffDays = diff.Hours() / 24
		if start.DateTime != "" {
			when = fmt.Sprintf("%.2f", diffDays)
		} else {
			when = fmt.Sprintf("%.0f", diffDays)
		}
	}
	return when
}

// JSONResponse is a map[string]string
// response from the web server
type JSONResponse map[string]string

// String returns the string representation of the
// JSONResponse object
func (j JSONResponse) String() string {
	str, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{
  "error": "%v"
}`, err)
	}

	return string(str)
}

func failHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, JSONResponse{
		"error": fmt.Sprintf("Not a valid endpoint: %s", r.URL.Path),
	})
	return
}

// writeError sends an error back to the requester
// and also logrus. the error
func writeError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, JSONResponse{
		"error": msg,
	})
	logrus.Printf("writing error: %s", msg)
	return
}

func init() {
	flag.StringVar(&clientSecret, "clientSecret", "", "client secret from google calendar")
	flag.StringVar(&port, "p", "8080", "port for server to run on")
	flag.Parse()

	if clientSecret == "" {
		logrus.Fatalf("You need to pass a google calendar secret file!")
	}
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	tokenCacheDir := filepath.Join(".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	cache := filepath.Join(tokenCacheDir,
		url.QueryEscape("loxonegogooglecalendar.json"))
	return cache, nil
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
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	// create mux server
	mux := http.NewServeMux()

	mux.HandleFunc("/calendar", calendarHandler) // forecast handler
	mux.HandleFunc("/", failHandler)             // everything else fail handler

	// set up the server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	logrus.Infof("Starting server on port %q", port)
	logrus.Fatal(server.ListenAndServe())
}
