import http from 'k6/http';
import { check } from 'k6';
import { sleep } from 'k6';

export let options = {
    vus: 100, // Number of VUs
    duration: '10s', // Duration of the test
};

export default function () {
    const url = 'http://localhost:1234/test';
    const payload = JSON.stringify({
        "version": "1.0",
        "eventTimeUtc": "2020-02-11T08:24:06.336Z",
        "systemId": 85,
        "requestId": 555,
        "dateTimeRangeUtc":{
            "from": "2020-04-06",
            "to": "2020-04-11"
        }
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const response = http.post(url, payload, params);

    check(response, {
        'is status 200': (r) => r.status === 200,
    });
    
}