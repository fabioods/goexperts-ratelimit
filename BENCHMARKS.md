# Go Expert Challenge - Rate Limiter Benchmarks

Resultados dos testes de carga realizados com o [Grafana k6](https://k6.io/) para avaliar o comportamento da solução desenvolvida sob pressão. Em todos os casos, a aplicação estava disponível a partir de um container Docker respondendo na porta `8080` e o k6 foi executado a partir de outro container.

## Hardware

- **CPU**: Macbook pro M3 - 11 núcleos (5 desempenho e 6 eficiência)
- **RAM**: 36GB
- **OS**: Mac OS 14.7
- **Runtime**: Docker (alpine)

## Smoke test

Teste com fim de validar se o serviço está respondendo corretamente.

- **Target**: 5 usuários
- **Duração**: 1 minuto

Comando para execução:

```sh
make test_k6_smoke
```

### Resultado

```plaintext
         /\      Grafana   /‾‾/
    /\  /  \     |\  __   /  /
   /  \/    \    | |/ /  /   ‾‾\
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/

     execution: local
        script: /scripts/smoke/smoke.test.js
 web dashboard: http://127.0.0.1:5665
        output: -

     scenarios: (100.00%) 1 scenario, 5 max VUs, 1m30s max duration (incl. graceful stop):
              * default: 5 looping VUs for 1m0s (gracefulStop: 30s)

     data_received..................: 45 MB  742 kB/s
     data_sent......................: 14 MB  240 kB/s
     http_req_blocked...............: min=375ns    med=542ns    avg=785ns  p(90)=750ns   p(95)=1µs     max=1.02ms   count=169498
     http_req_connecting............: min=0s       med=0s       avg=6ns    p(90)=0s      p(95)=0s      max=534.25µs count=169498
     http_req_duration..............: min=252.7µs  med=328.54µs avg=1.75ms p(90)=6.38ms  p(95)=10.95ms max=72.08ms  count=169498
       { expected_response:true }...: min=295.45µs med=360.95µs avg=2.11ms p(90)=7.79ms  p(95)=11.98ms max=48.81ms  count=6600
     ✓ { status:200 }...............: min=295.45µs med=360.95µs avg=2.11ms p(90)=7.79ms  p(95)=11.98ms max=48.81ms  count=6600
     ✓ { status:429 }...............: min=252.7µs  med=327.41µs avg=1.74ms p(90)=6.31ms  p(95)=10.91ms max=72.08ms  count=162898
     ✓ { status:500 }...............: min=0s       med=0s       avg=0s     p(90)=0s      p(95)=0s      max=0s       count=0
     http_req_failed................: 96.10% 162898 out of 169498
     http_req_receiving.............: min=3.12µs   med=5µs      avg=7.57µs p(90)=6.16µs  p(95)=7.25µs  max=12.46ms  count=169498
     http_req_sending...............: min=916ns    med=1.29µs   avg=2.06µs p(90)=1.75µs  p(95)=2µs     max=10.98ms  count=169498
     http_req_tls_handshaking.......: min=0s       med=0s       avg=0s     p(90)=0s      p(95)=0s      max=0s       count=169498
     http_req_waiting...............: min=247.41µs med=321.5µs  avg=1.74ms p(90)=6.36ms  p(95)=10.94ms max=72.07ms  count=169498
     http_reqs......................: 169498 2822.839147/s
     iteration_duration.............: min=539.79µs med=690.5µs  avg=3.53ms p(90)=11.51ms p(95)=15.76ms max=72.78ms  count=84749
     iterations.....................: 84749  1411.419574/s
     vus............................: 5      min=5                max=5
     vus_max........................: 5      min=5                max=5


running (1m00.0s), 0/5 VUs, 84749 complete and 0 interrupted iterations
default ✓ [======================================] 5 VUs  1m0s
```
