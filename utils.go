package ics

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/goark/fetch"
)

var mutex *sync.Mutex
var idCounter int

// if DeleteTempFiles is true , after we download ics and parse it , the local temp file  will be deleted
var DeleteTempFiles bool

// Describes the file path to the folder with the temp ics files
var FilePath string

// if RepeatRuleApply is true , the rrule will create new objects for the repeated events
var RepeatRuleApply bool

// max of the rrule repeat for single event
var MaxRepeats int

// unixtimestamp
const uts = "1136239445"

// ics date time format
const IcsFormat = "20060102T150405Z"

// Y-m-d H:i:S time format
const YmdHis = "2006-01-02 15:04:05"

// ics date format ( describes a whole day)
const IcsFormatWholeDay = "20060102"

// downloads the calendar before parsing it
func downloadFromUrl(url string) (string, error) {
	// split the url to get the name of the file (like basic.ics)
	tokens := strings.Split(url, "/")

	// create the name of the file
	fileName := fmt.Sprintf("%s%s_%s", FilePath, time.Now().Format(uts), tokens[len(tokens)-1])

	// creates the path
	if err := os.MkdirAll(FilePath, 0777); err != nil {
		return "", err
	}

	// creates the file in the path folder
	output, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	// close the file
	defer output.Close()

	// get the URL
	u, err := fetch.URL(url)
	if err != nil {
		return "", err
	}
	response, err := fetch.New().Get(u)
	if err != nil {
		return "", err
	}
	// close the response body
	defer response.Close()

	// copy the response from the url to the temp local file
	_, err = io.Copy(output, response.Body())
	if err != nil {
		return "", err
	}

	//return the file that contains the info
	return fileName, nil
}

func stringToByte(str string) []byte {
	return []byte(str)
}

// removes newlines and cutset from given string
func trimField(field, cutset string) string {
	re, _ := regexp.Compile(cutset)
	cutsetRem := re.ReplaceAllString(field, "")
	return strings.TrimRight(cutsetRem, "\r\n")
}

// checks if file exists
func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil
}

func parseDayNameToIcsName(day string) string {
	var dow string
	switch day {
	case "Mon":
		dow = "MO"
	case "Tue":
		dow = "TU"
	case "Wed":
		dow = "WE"
	case "Thu":
		dow = "TH"
	case "Fri":
		dow = "FR"
	case "Sat":
		dow = "ST"
	case "Sun":
		dow = "SU"
	default:
		// fmt.Println("DEFAULT :", start.Format("Mon"))
		dow = ""
	}
	return dow
}
