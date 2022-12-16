package data

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestJsonParsingOfRequest(t *testing.T) {
	input :=
		`{
			"version": "1.0",
			"eventTimeUtc": "2020-02-11T08:24:06.336Z",
			"systemId": 85,
			"requestId": 555,
			"dateTimeRangeUtc":{
				"from": "2020-04-06",
				"to": "2020-04-11"
		   },
		   "requiredFields" : [ "siteId", "segmentId", "contactId", "mediaType", "segmentStartTime", "segmentEndTime", "duration", "callDirection", "recorderId", "inLitigationHold", "customerId", "caseId", "contactStartTime", "archiveId", "setNumber", "fsArchiveClass", "esmArchiveClass", "fsClusterId", "esmClusterId", "isDeleted", "isUnavailable", "isFsArchived", "isEsmArchived" ]	   
		}`

	var req Request
	err := json.Unmarshal([]byte(input), &req)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Equal(t, req.SystemID, 85)
}
