# GO Expert-ratelimiter
Projeto do Desafio Técnico "Rate Limiter" do treinamento GoExpert(FullCycle).

## O desafio
Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.
- Endereço IP: O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.
- Token de Acesso: O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato:
API_KEY: <TOKEN>
- As configurações de limite do token de acesso devem se sobrepor as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.

## Características

- Limitação de taxa por IP e Token de API
- Configurável via variáveis de ambiente
- Armazenamento em Redis
- Middleware pronto para uso
- Tempo de bloqueio configurável
- Design flexível com padrão Strategy para armazenamento
- Testes automatizados

## Requisitos

- Go 1.19 ou superior
- Docker e Docker Compose
- Redis

## Instalação

1. Clone o repositório:

```bash
git clone https://github.com/diogocardoso/go-rateLimiter.git
cd goexpert-ratelimiter
```
2. Instale as dependências:

```bash
go mod tidy
```
3. Configure as variáveis de ambiente:

```bash
cp .env.example .env
```
4. Inicie a aplicação:

```bash
# Construir as imagens
docker-compose build

# Iniciar os serviços
docker-compose up -d

# Ver logs
docker-compose logs -f

# Parar os serviços 
docker-compose down
```

## Uso

### Para testar a aplicação:

```bash
# Teste simples
curl -i http://localhost:8080/api

# Teste com várias requisições
for i in {1..5}; do curl -i http://localhost:8080/api; done
```

## Requisitos: implementação
- O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web
- O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.
- O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.
- As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “.env” na pasta raiz.
- Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.
- O sistema deve responder adequadamente quando o limite é excedido:
    - Código HTTP: 429
    - Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame
- Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.
- Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.
- A lógica do limiter deve estar separada do middleware.

## Requisitos: entrega
- O código-fonte completo da implementação.
- Documentação explicando como o rate limiter funciona e como ele pode ser configurado.
- Testes automatizados demonstrando a eficácia e a robustez do rate limiter.
- Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
- O servidor web deve responder na porta 8080.
