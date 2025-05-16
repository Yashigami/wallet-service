import http from 'k6/http';
import { check } from 'k6';

export let options = {
    stages: [
        { duration: '10s', target: 1000 }, // выход на 1000 RPS
        { duration: '30s', target: 1000 }, // поддержка 1000 RPS
        { duration: '10s', target: 0 },    // спад нагрузки
    ],
    thresholds: {
        http_req_failed: ['rate<0.01'],    // не более 1% ошибок
        http_req_duration: ['p(95)<500'],  // 95% запросов быстрее 500мс
    },
};

const walletId = '9e1c9d04-5c62-4f8e-9d7f-ec6d2f2c1d2a';
const baseUrl = 'http://localhost:8080/api/v1';

export default function () {
    // Случайная операция: DEPOSIT, WITHDRAW, GET
    const choice = Math.floor(Math.random() * 3);

    let res;

    if (choice === 0) {
        res = http.post(`${baseUrl}/wallet`, JSON.stringify({
            walletId: walletId,
            operationType: 'DEPOSIT',
            amount: 1000
        }), {
            headers: { 'Content-Type': 'application/json' },
        });
    } else if (choice === 1) {
        res = http.post(`${baseUrl}/wallet`, JSON.stringify({
            walletId: walletId,
            operationType: 'WITHDRAW',
            amount: 500
        }), {
            headers: { 'Content-Type': 'application/json' },
        });
    } else {
        res = http.get(`${baseUrl}/wallets/${walletId}`);
    }

    check(res, {
        'status is 200': (r) => r.status === 200,
    });
}
