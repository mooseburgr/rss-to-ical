// Package p contains an HTTP Cloud Function.
package p

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/mmcdole/gofeed"
)

// Fetches RSS feed of events located at request param "rssUrl" and converts
// to iCal format with the specified "eventDuration," defaulting to 60 minutes
func HandleRequest(w http.ResponseWriter, r *http.Request) {

	// validate URL param
	rssUrl := r.URL.Query().Get("rssUrl")
	if _, err := url.ParseRequestURI(rssUrl); err != nil {
		http.Error(w, "Invalid URL: "+err.Error(), http.StatusBadRequest)
		return
	}
	// default to event duration of 60min if invalid or unspecified
	eventDuration, err := strconv.Atoi(r.URL.Query().Get("eventDuration"))
	if err != nil {
		eventDuration = 60
	}

	cal, err := doConvert(rssUrl, eventDuration)

	if err != nil {
		http.Error(w, "Error parsing feed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// write to response
	w.Header().Add("Content-Type", "text/calendar")
	w.Write([]byte(cal.Serialize()))
}

func doConvert(rssUrl string, eventDuration int) (*ical.Calendar, error) {

	// fetch RSS feed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rssUrl)
	if err != nil {
		log.Printf("Error parsing feed: %v", err)
		return nil, err
	}
	log.Print("Fetched feed: ", feed.Title)

	// convert to ical
	cal := ical.NewCalendar()

	// top level properties
	productId := "-//" + feed.Title + "//mooseburgr/rss-to-ical"
	cal.SetProductId(productId)
	cal.SetVersion("2.0")
	cal.SetCalscale("GREGORIAN")
	cal.SetMethod(ical.MethodPublish)
	cal.SetName(feed.Title)
	cal.SetXWRCalName(feed.Title)
	cal.SetXWRCalDesc(feed.Description)
	cal.SetXWRCalID(productId)
	cal.SetLastModified(*feed.UpdatedParsed)

	// copy events
	for _, item := range feed.Items {
		event := cal.AddEvent(item.GUID)
		event.SetStartAt(*item.PublishedParsed)
		event.SetEndAt(item.PublishedParsed.Add(time.Minute * time.Duration(eventDuration)))
		event.SetSummary(item.Title)
		event.SetDescription(item.Description)
		event.SetLocation(item.Link)
		event.SetURL(item.Link)
		event.SetOrganizer(authorsToOrganizer(item.Authors))
	}

	return cal, nil
}

func authorsToOrganizer(authors []*gofeed.Person) string {
	var result []string
	for _, author := range authors {
		if strings.TrimSpace(author.Email) == "" {
			result = append(result, strings.TrimSpace(author.Name))
		} else {
			result = append(result, fmt.Sprintf("%s (%s)",
				strings.TrimSpace(author.Name), strings.TrimSpace(author.Email)))
		}
	}
	return strings.Join(result, ", ")
}
