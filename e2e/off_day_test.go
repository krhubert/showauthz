package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"rift/assert"
	"rift/memdb"
)

func TestOffDay(t *testing.T) {
	resp, err := http.Post(tsURL+"/offdays", "application/json", bytes.NewReader([]byte(`{"ID":"of1","OrganizationID":"org"}`)))
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	resp, err = http.Get(tsURL + "/offdays")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	var offDays []*memdb.OffDay
	json.NewDecoder(resp.Body).Decode(&offDays)
	assert.Equal(t, offDays, []*memdb.OffDay{{ID: "of1", OrganizationID: "org"}})
}
