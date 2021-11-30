# rss-to-ical-go

Do you have an RSS feed of events you want to import into a calendar manager as .ics / iCal content and for some reason don't want to use [existing services](https://zapier.com/apps/google-calendar/integrations/rss/1449/create-google-calendar-events-from-rss-feed-items) to solve this problem??

Yeah, me neither.

_BUT_ if you do, you can use this Google Cloud Function to do this content massaging for you. Just swap in the URL of your RSS feed below, optionally overriding the event duration that defaults to 60 minutes:

`https://us-central1-absolute-pulsar-301701.cloudfunctions.net/convert-rss-to-ical-go?rssUrl={RSS URL}&eventDuration={duration}`


(e.g. [`https://us-central1-absolute-pulsar-301701.cloudfunctions.net/convert-rss-to-ical-go?rssUrl=https://demo.theeventscalendar.com/events/feed/&eventDuration=50`](https://us-central1-absolute-pulsar-301701.cloudfunctions.net/convert-rss-to-ical-go?rssUrl=https://demo.theeventscalendar.com/events/feed/&eventDuration=50))


Credit to the following packages doing about 99.9% of the work:
 - https://github.com/mmcdole/gofeed
 - https://github.com/arran4/golang-ical
