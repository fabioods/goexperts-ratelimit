# Go Expert Challenge - Rate Limiter

Esse desafio tem como objetivo implementar um rate limit, que consiste em limitar o número de requests recebidas dentro um tempo pré-estabelecido.
O rate limit vai ter duas formas de funcionamento, por IP ou por access token, sendo o número para cada um deles configurável via variáveis de ambiente.

## Arquitetura

Essa aplicação é feita em go e consiste em um servidor web que recebe requisições via HTTP e tem implementado um middleware para validar as requests recebidas, validando o rate limit em caso.

Para validar o rate limit existe uma validação por IP ou por token, onde recebe-se uma `API_KEY` no header da request. Para saber se o rate limit foi atingido os dados das requests são armazenados em um banco de dados Redis, já as configurações para o númnero máximo de requests por ip e token, bem como o tempo de vida que isso será analisado é definido em um arquivo `.env`.

- `RATE_LIMITER_IP_MAX_REQUESTS`: Máximo de requets por IP
- `RATE_LIMITER_TOKEN_MAX_REQUESTS`: Máximo de requests por token
- `RATE_LIMITER_TIME_WINDOW_MILISECONDS`: Tempo de vida em ms

Exemplo de configuração com um tempo de vida de 5 minutos, 10 requests por IP e 100 requets por token:

```sh
RATE_LIMITER_IP_MAX_REQUESTS=10
RATE_LIMITER_TOKEN_MAX_REQUESTS=100
RATE_LIMITER_TIME_WINDOW_MILISECONDS=300000
```

No caso de do rate limit ser atingido será retornado o status code `429 Too Many Requests` com a mensagem `{"message":"rate limit exceeded"}`.
Outras informações são fornecidas como os cabeçalhos `X-Ratelimit-Limit`, `X-Ratelimit-Remaining` e `X-Ratelimit-Reset` com informações sobre o limite, quantidade restante e tempo de reset, respectivamente.

## Benchmarks

Foi utilizado o [Grafana k6](https://k6.io/) para realizar testes de carga do tipo [_smoke_](https://grafana.com/docs/k6/latest/testing-guides/test-types/smoke-testing/) no serviço para avaliar o comportamento da solução desenvolvida. Os resultados se encontram [aqui](./BENCHMARKS.md).

## Executando o projeto

**Obs:** é necessário ter o [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/) instalados.

1. Crie um arquivo `.env` na raiz do projeto copiando o conteúdo de `.env.example` e ajuste-o conforme necessário. Por padrão, os seguintes valores são utilizados:

```sh
LOG_LEVEL="debug" # Nível de log da aplicação
WEB_SERVER_PORT=8080 # Porta do servidor Web

# Configurações do Redis
REDIS_HOST="localhost"
REDIS_PORT=6379
REDIS_PASSWORD=""
REDIS_DB=0

RATE_LIMITER_IP_MAX_REQUESTS=10 # Número máximo de requisições por IP
RATE_LIMITER_TOKEN_MAX_REQUESTS=100 # Número máximo de requisições por token
RATE_LIMITER_TIME_WINDOW_MILISECONDS=1000 # Janela de tempo em milissegundos
```

2. Execute o comando `docker compose up -d` para iniciar a aplicação e o Redis.

### Exemplos de requisições

- **Requisição com checagem via IP com sucesso:**

```sh
$ curl -vvv http://localhost:8080
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* connect to ::1 port 8080 from ::1 port 49992 failed: Connection refused
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.7.1
> Accept: */*
>
* Request completely sent off
< HTTP/1.1 200 OK
< Accept: application/json
< Content-Type: application/json
< X-Ratelimit-Limit: 10
< X-Ratelimit-Remaining: 9
< X-Ratelimit-Reset: 1732710055
< Date: Wed, 27 Nov 2024 12:20:54 GMT
< Content-Length: 27
<
{"message":"Hello World!"}
* Connection #0 to host localhost left intact
```

- **Requisição com checagem via token com sucesso:**

```sh
$ curl -H 'API_KEY: some-api-key-123' -vvv http://localhost:8080
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* connect to ::1 port 8080 from ::1 port 50000 failed: Connection refused
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.7.1
> Accept: */*
> API_KEY: some-api-key-123
>
* Request completely sent off
< HTTP/1.1 200 OK
< Accept: application/json
< Content-Type: application/json
< X-Ratelimit-Limit: 100
< X-Ratelimit-Remaining: 99
< X-Ratelimit-Reset: 1732710093
< Date: Wed, 27 Nov 2024 12:21:32 GMT
< Content-Length: 27
<
{"message":"Hello World!"}
* Connection #0 to host localhost left intact

```

- **Requisição com checagem via IP bloqueada:**

```sh
$ curl -vvv http://localhost:8080

* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* connect to ::1 port 8080 from ::1 port 50135 failed: Connection refused
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.7.1
> Accept: */*
>
* Request completely sent off
< HTTP/1.1 200 OK
< Accept: application/json
< Content-Type: application/json
< X-Ratelimit-Limit: 10
< X-Ratelimit-Remaining: 0
< X-Ratelimit-Reset: 1732710180
< Date: Wed, 27 Nov 2024 12:23:00 GMT
< Content-Length: 27
<
{"message":"rate limit exceeded"}
```

- **Requisição com checagem via token bloqueada:**

```sh
$ curl -H 'API_KEY: some-api-key-123' -vvv http://localhost:8080

* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* connect to ::1 port 8080 from ::1 port 50142 failed: Connection refused
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.7.1
> Accept: */*
> API_KEY: some-api-key-123
>
* Request completely sent off
< HTTP/1.1 200 OK
< Accept: application/json
< Content-Type: application/json
< X-Ratelimit-Limit: 100
< X-Ratelimit-Remaining: 98
< X-Ratelimit-Reset: 1732710222
< Date: Wed, 27 Nov 2024 12:23:41 GMT
< Content-Length: 27
<
{"message":"rate limit exceeded"}
```

## Testes

### Testes de unidade

Para executar os testes de unidade e validar a cobertura, execute o comando `make test`.

### Testes de estresse

Para executar os testes de estresse com k6, siga os passos:

1. Inicie a aplicação e o Redis com o comando `docker-compose up -d`;
2. Execute o comando `make test_k6_smoke` para iniciar o teste de estresse do tipo _smoke_ (duração de 1 minuto);

É possível visualizar os resultados dos testes obtidos por mim na pasta `./scripts/k6/smoke`, tanto em formato de texto quanto em HTML.
