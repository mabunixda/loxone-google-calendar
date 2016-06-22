# README #

This is a microservice written in go providing weather data from forecast.io
To use this service you must create a developer account on forecast.io to get an API Key.

### How do I get set up? ###

* Clone the repo
* Set the GOPATH to the cloned repo directory
* From bitbucket-pipelines run the commans under the "script" node
* Register a development account on [google|https://console.developers.google.com/start/api?id=calendar] 
You might follow the [Go Quickstart Guide|https://developers.google.com/google-apps/calendar/quickstart/go] from google because you must also configure access for this api secrets to the dedicated calendars
* Run the calendar service using the json secret file 

```
#!bash
LoxoneGoGoogleCalendar -clientSecret $YOUR_SECRET_FILE

```
* First requests are served by http://localhost:8080/calendar
* On the first request google api will provide a link where you approve access to the actual calendar.

Within Loxone you can use 1 Virtual HTTP Input to query all the data with a single http query and parse it afterwards
into seperate variables

The default behavior is that your personal calendar is queried. You can also define a special calendar e.g. a seperate trash calendar:
http://localhost:8080/calendar?**calendarId=$GOOGLE_CALENDAR_ID**

Additional Parameters:
countDown: If this parameter is set, then a duration till the event startdate is shown. If you specify countDown=days only the amount of days is shown.
show: If this parameter is set to all, then recurring events are not filtered by their first occurance

### Example Output ###

```
#!csv
Upcoming events:
Paper:2016-06-29
Trash:2016-08-02
```

### Stuff todo ###
* Caching data if the internet connection goes down

### Who do I talk to? ###

* Repo owner or admin
* Other community or team contact