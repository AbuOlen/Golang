import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    vus: 1000,
    duration: '30s',
};

const channels = ['chat1', 'chat2', 'chat3'];

const tokens = [
    'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIwMDAxTTNORVYwS0FKM0JDTjU2VFhRMVNKSiIsImV4cCI6MTc0OTE2MjA3NywiaWF0IjoxNzQ5MTQ0MDc3fQ.F6SW7geiZwQv3wek4dHCc0XQVnsic2R2g84fbLNluT4',
    'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIwMDAxTTNORVY3WEVUQU43WVNIM1BBMVdGMCIsImV4cCI6MTc0OTE2MzEzMCwiaWF0IjoxNzQ5MTQ1MTMwfQ.J1jKN9EOnVITxbSwwIKRLBgMMPINqMNDXbU3LjjK_uc',
    'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIwMDAxTTNORVZERTRUOVQxWVgyWFA1UjRNUiIsImV4cCI6MTc0OTE2MzE0NiwiaWF0IjoxNzQ5MTQ1MTQ2fQ.BLn_JqRSzFU3zpCpf4rgyEaqyFhzvtRVgzYcSfi1t68',
];

export default function () {
    const channel = channels[Math.floor(Math.random() * channels.length)];
    const token = tokens[Math.floor(Math.random() * tokens.length)];

    const headers = {
        'Content-Type': 'application/json',
        'X-User-Token': token,
    };

    const messagePayload = JSON.stringify({
        channel: channel,
        message: `Hello from ${__VU}`,
    });

    // Send message
    let sendRes = http.post('http://localhost:8080/api/message/send', messagePayload, { headers });
    check(sendRes, { 'send status was 200': (r) => r.status === 200 });

    let lastReceivedId = "";

    // Try to extract last_received_id from response
    if (sendRes.status === 200) {
        try {
            const resp = JSON.parse(sendRes.body);
            if (resp && resp.id) {
                //lastReceivedId = resp.id;
            }
        } catch (e) {
            console.error('Failed to parse response:', e);
        }
    }

    const listPayload = JSON.stringify({
        channel: channel,
        last_received_id: lastReceivedId,
    });

    // List messages
    let listRes = http.post('http://localhost:8080/api/message/list', listPayload, { headers });
    check(listRes, { 'list status was 200': (r) => r.status === 200 });

    sleep(0.1);
}
