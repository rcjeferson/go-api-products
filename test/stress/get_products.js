import { check } from 'k6';
import http from 'k6/http';

export const options = {
  stages: [
    { duration: '1m', target: 20 },
    { duration: '1m', target: 100 },
    { duration: '1m', target: 200 },
  ]
}

export default function() {
  const res = http.get('http://localhost:8000/products');
  check(res, { 'status was 200': (r) => r.status == 200 });
}