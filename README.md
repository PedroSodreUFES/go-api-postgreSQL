![Go](https://img.shields.io/badge/Language-Go-00ADD8?logo=go)

## API com Go Language e PostgreSQL

### Contexto
CRUD de usuários contendo id, nome, sobrenome e biografia.

### Rotas:
+ DELETE /user/:id - deleta um usuário pelo id.
+ GET /user/:id - lê um usuário pelo id.
+ GET /users - lista todos os usuários.
+ POST /user - cria um usuário.
+ PUT /user/:id - edita um usuário pelo id.

### Tecnologias usadas
+ Go Language
+ Chi framework
+ Pgx
+ Google UUID
+ Pacote json
+ Pacote errors
+ Pacote log
+ Pacote http

### Como rodar o programa
```bash
docker-compose up -d
go run .
```
Ou
```bash
docker-compose up -d
go build 
./main
```

### Como parar o container Docker e excluir os dados
```bash
docker-compose down -v
```