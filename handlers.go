package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"encoding/json"

	"github.com/rs/xid"
	"github.com/zenazn/goji/web"
)

func ensureGuidExists(guid string, w http.ResponseWriter) bool {
	exists, err := Redis.SIsMember(sessionsKey, guid).Result()
	if err != nil {
		panic(err)
	}
	if !exists {
		http.Error(w, "Invalid ID", http.StatusNotFound)
	}
	return exists
}

func Root(c web.C, w http.ResponseWriter, r *http.Request) {
	guid := xid.New().String()
	err := Redis.SAdd(sessionsKey, guid).Err()
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, hookUrl(guid), http.StatusSeeOther)
}

func GetData(c web.C, w http.ResponseWriter, r *http.Request) {
	// Obtain the GUID and check that it was created.
	guid := c.URLParams["guid"]
	if !ensureGuidExists(guid, w) {
		return
	}
	values, err := Redis.LRange(sessionKey(guid), 0, -1).Result()
	if err != nil {
		panic(err)
	}
	for i := range values {
		var h HookDatum
		err := json.Unmarshal([]byte(values[i]), &h)
		if err != nil {
			panic(err)
		}
		io.WriteString(w, fmt.Sprintf("%s\n%s\n%s\n%s\n\n",
			h.Method, h.Headers, h.Time, h.Body))
	}
}

func PostData(c web.C, w http.ResponseWriter, r *http.Request) {
	guid := c.URLParams["guid"]
	if !ensureGuidExists(guid, w) {
		return
	}
	// Read request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Create new entry.
	h := HookDatum{
		Method:  "POST",
		Headers: r.Header,
		Body:    string(body),
		Time:    time.Now().UTC(),
	}
	// Store it.
	val, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	err = Redis.LPush(sessionKey(guid), val).Err()
	if err != nil {
		panic(err)
	}
	// Redirect to GetData.
	http.Redirect(w, r, hookUrl(guid), http.StatusCreated)
}
