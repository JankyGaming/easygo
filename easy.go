package easygo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//BuildRequest build your http request in one easy, convenient line, just return after this method.
func BuildRequest(domain string, port string, method string, endpoint string, headers map[string]string, params map[string]string, object interface{}) ([]byte, int, string, error) {
	var url string

	if port != "" {
		if strings.HasSuffix(domain, "/") {
			url = domain[:len(domain)-len("/")]
			url = url + ":" + port + "/" + endpoint
		} else {
			url = domain + ":" + port + "/" + endpoint
		}
	} else {
		if strings.HasSuffix(domain, "/") {
			url = domain + endpoint
		} else {
			url = domain + "/" + endpoint
		}
	}

	httpClient := &http.Client{}

	var payload []byte
	var err error

	if object != nil {
		payload, err = json.Marshal(object)
		if err != nil {
			return nil, 0, "", err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, 0, "", err
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Set(key, val)
			fmt.Println(req.Header)
		}
	}

	if params != nil {
		q := req.URL.Query()
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, "", err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, "", err
	}

	return resBody, res.StatusCode, res.Status, nil
}

//Respond used in handler functions to build responses with one line by passing in statusCode, body, and headers
func Respond(w http.ResponseWriter, r *http.Request, statusCode int, body interface{}, headers map[string]string) error {
	w.WriteHeader(statusCode)

	if body != nil {
		res, err := json.Marshal(body)
		if err != nil {
			return err
		}
		w.Write(res)
	}

	if headers != nil {
		// TODO: Set Headers
	}

	return nil
}

//RespondBasic supports: 200|ok, 201|created, 400|bad request, 401|unauthorized, 403|forbidden, 404|resource missing, 405|method not allowed, 500|internal server error
func RespondBasic(w http.ResponseWriter, r *http.Request, statusCode int) {
	msg := map[string]interface{}{}
	switch statusCode {
	case 200:
		msg["message"] = "ok"
	case 201:
		msg["message"] = "created"
	case 400:
		msg["message"] = "bad request"
	case 401:
		msg["message"] = "unauthorized"
	case 403:
		msg["message"] = "forbidden"
	case 404:
		msg["message"] = "resource missing"
	case 405:
		msg["message"] = "method not allowed"
	case 500:
		msg["message"] = "internal server error"
	default:
		msg["timestamp"] = time.Now()
		res, err := json.Marshal(msg)
		if err != nil {
			fmt.Println(err)
		}
		w.WriteHeader(statusCode)
		w.Write(res)
		return
	}
	msg["timestamp"] = time.Now()

	res, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(statusCode)
	w.Write(res)

}

// AppendStringNoDuplicates When adding a string to an array, this will check that the string doesn't all ready exist in the array
func AppendStringNoDuplicates(arr []string, str string) []string {
	doIt := true

	for _, arrStr := range arr {
		if arrStr == str {
			doIt = false
		}
	}

	if doIt {
		arr = append(arr, str)
	}

	fmt.Println("")
	return arr
}

// AppendStringSliceNoDuplicates when combing string arrays, merges while preventing duplicates
func AppendStringSliceNoDuplicates(slice1 []string, slice2 []string) []string {
	for _, str1 := range slice1 {
		doIt := true

		for _, str2 := range slice2 {
			if str1 == str2 {
				doIt = false
			}
		}

		if doIt {
			slice2 = append(slice2, str1)
		}
	}
	return slice2
}

// ResponseObject A neat little response object for simple responses
type ResponseObject struct {
	Error    bool   `json:"error"`
	Message  string `json:"message"`
	Function string `json:"function"`
}

// GetDateString takes a time object and returns the date in string form. YYYY-MM-DD
func GetDateString(tme time.Time) string {

	yr, mth, day := tme.Date()

	var dayNumber int
	dayNumber = int(day)
	var dayString string
	if dayNumber < 10 {
		dayString = strconv.Itoa(dayNumber)
		dayString = "0" + dayString
	} else {
		dayString = strconv.Itoa(dayNumber)
	}
	strconv.Itoa(dayNumber)

	var mthNumber int
	mthNumber = int(mth)
	var mthString string
	if mthNumber < 10 {
		mthString = strconv.Itoa(mthNumber)
		mthString = "0" + mthString
	} else {
		mthString = strconv.Itoa(mthNumber)
	}
	strconv.Itoa(dayNumber)

	concatStr := strconv.Itoa(yr) + "-" + mthString + "-" + dayString

	return concatStr
}

//Log Prints to stdout
func Log(msg string) {
	pc, _, line, _ := runtime.Caller(1)
	print := time.Now().Format(time.RFC3339) + " | " + "Log from (" + runtime.FuncForPC(pc).Name() + " on line " + strconv.Itoa(line) + ") | " + msg
	fmt.Println(print)
}

//LogErr Prints to stdout
func LogErr(err error) {
	pc, _, line, _ := runtime.Caller(1)
	print := time.Now().Format(time.RFC3339) + " | " + "Error in (" + runtime.FuncForPC(pc).Name() + " on line " + strconv.Itoa(line) + ") | " + err.Error()
	fmt.Println(print)
}

//DecodeBody reads and unmarshalls the body of a a request b into the struct at dest and returns the RAW body and unmarshalled body
func DecodeBody(b io.Reader, dest interface{}) ([]byte, error) {
	byt, err := ioutil.ReadAll(b)
	if err != nil {
		return byt, err
	}

	err = json.Unmarshal(byt, &dest)
	if err != nil {
		return byt, err
	}

	return byt, nil
}
