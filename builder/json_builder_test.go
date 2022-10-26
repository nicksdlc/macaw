package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleSerialize(t *testing.T) {
	input :=
		`{
			"version": "1.0",
			"eventTimeUtc": "2020-02-11T08:24:06.336Z",
			"systemId": 85,
			"policyId": 555,
			"dateTimeRangeUtc":{
				"from": "2020-04-06",
				"to": "2020-04-11"
		   },
		   "requiredFields" : [ "siteId", "segmentId", "contactId", "mediaType", "segmentStartTime", "segmentEndTime", "duration", "callDirection", "recorderId", "inLitigationHold", "customerId", "caseId", "contactStartTime", "archiveId", "setNumber", "fsArchiveClass", "esmArchiveClass", "fsClusterId", "esmClusterId", "isDeleted", "isUnavailable", "isFsArchived", "isEsmArchived" ]	   
		}`

	intermidiate := serialize(input)

	result := deserialize(intermidiate)

	assert.Contains(t, result, "\"policyId\":555")
}

func TestShouldReplaceWithRandomNumber(t *testing.T) {
	input := `{number: {{(sum 2 1)}}}`

	result := build(input)

	assert.Equal(t, result, "{number: 3}")
}
