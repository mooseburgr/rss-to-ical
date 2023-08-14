package p

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const validRssUrl = "https://demo.theeventscalendar.com/events/feed/"

// const validRssUrl = "https://feeds.twit.tv/twit.xml"

func TestHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rss-to-ical?rssUrl="+validRssUrl, nil)
	w := httptest.NewRecorder()

	HandleRequest(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	assert.NotEmpty(t, string(data))
}

func TestSuccessfulConvert(t *testing.T) {

	cal, err := doConvert(validRssUrl, 50)

	assert.Nil(t, err)
	calString := cal.Serialize()

	assert.True(t, strings.HasPrefix(calString, "BEGIN:VCALENDAR"))
	assert.Contains(t, calString, "BEGIN:VEVENT")
	assert.Contains(t, calString, "END:VEVENT")
	assert.Contains(t, calString, "END:VCALENDAR")
}

func TestFailedConvert(t *testing.T) {

	cal, err := doConvert("https://dne.com/rss", 45)

	assert.Nil(t, cal)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "404 Not Found")
}
