# BookClub API

REST API para gerenciamento de clubes de leitura, desenvolvida em Go com Clean Architecture.  
Atividade 03 — Desenvolvimento Web II.

---

## Domínio

O sistema permite que usuários criem clubes de leitura, participem como membros, sugiram livros, sorteiem temas mensais e agendem reuniões.

---

## Tecnologias

- **Go 1.22**
- **Gin** — HTTP router
- **GORM** — ORM com PostgreSQL
- **JWT** — autenticação via `golang-jwt/jwt/v5`
- **Docker Compose** — PostgreSQL 16

---

## Como rodar

```bash
cp .env.example .env
docker compose up -d
go run ./cmd/api
```

O servidor sobe em `http://localhost:8080`.  
O banco de dados é migrado automaticamente na primeira execução.  
O PostgreSQL expõe a porta **15432** no host para evitar conflito com instâncias locais.

Variáveis de ambiente relevantes:

| Variável           | Descrição                                |
|--------------------|------------------------------------------|
| `DB_HOST`          | Host do Postgres                         |
| `DB_PORT`          | Porta do Postgres (padrão: 15432)        |
| `DB_USER`          | Usuário do Postgres                      |
| `DB_PASSWORD`      | Senha do Postgres                        |
| `DB_NAME`          | Nome do banco                            |
| `JWT_SECRET`       | Chave HMAC-SHA256 para assinar tokens    |
| `ADMIN_EMAIL`      | E-mail do admin semeado na inicialização |
| `ADMIN_PASSWORD`   | Senha do admin semeado                   |
| `VISITOR_EMAIL`    | E-mail do visitor semeado                |
| `VISITOR_PASSWORD` | Senha do visitor semeado                 |

---

## Postman

A coleção está em `docs/postman/reader-club.postman_collection.json`.  
**Não é necessário importar arquivo de environment** — todas as variáveis estão embutidas na coleção.

**Ordem de execução:**

1. **Setup / Initialize Test Data** — gera `run_id` único e define e-mail, nome do clube etc.
2. Auth / Login Admin
3. Auth / Register Member
4. Auth / Login Member
5. Auth / Login Visitor
6. Book Clubs / Create Club as Admin
7. Memberships / Join Club as Member
8. Book Suggestions / Suggest Book as Member
9. Monthly Themes / Draw Monthly Theme as Admin
10. Meetings / Create Meeting as Admin
11. Authorization Checks *(valida 403, 401 e 404)*
12. **Book Clubs / Delete Club as Admin (Run Last)**

Tokens e IDs são salvos automaticamente pelos scripts de teste. A coleção é repetível: o e-mail do membro e o nome do clube incluem um sufixo de timestamp gerado pelo Setup.

---

## Estrutura do projeto

```
cmd/api/              — entrada da aplicação e injeção de dependências
internal/
  domain/
    entity/           — entidades do domínio
    repository/       — interfaces de repositório (sem dependência de framework)
  application/
    apperr/           — erros de domínio
    dto/              — contratos de request/response
    mapper/           — conversão entidade <-> DTO
    logger/           — interface de log
    usecase/          — casos de uso (1 arquivo = 1 caso de uso)
  infrastructure/
    auth/             — serviço JWT
    persistence/      — implementações GORM dos repositórios
  interfaces/http/
    controller/       — um controller por recurso
    httperr/          — struct APIError e helpers de resposta
    middleware/       — Authenticate, RequireGlobalRoles, RequireClubRole
    router.go         — definição das rotas
```

---

## Entidades

- **User** — usuário do sistema; possui `global_role`
- **BookClub** — clube de leitura; possui nome, descrição e dono
- **Membership** — associação entre usuário e clube, com papel no clube
- **BookSuggestion** — sugestão de livro feita por um membro
- **MonthlyTheme** — livro sorteado como tema do mês para um clube
- **Meeting** — reunião agendada vinculada a um tema mensal
- **AuditLog** — registro de operações de escrita

---

## Roles

### Global (campo `global_role` no usuário)

| Role           | Atribuição                       |
|----------------|----------------------------------|
| `ROLE_ADMIN`   | Semeado na inicialização         |
| `ROLE_MEMBER`  | Padrão para auto-registro        |
| `ROLE_VISITOR` | Semeado na inicialização         |

O papel global é embutido no JWT no campo `"roles"`.

### Escopo de clube (tabela `memberships`)

Armazenado por par `(user_id, club_id)`. Verificado pelo middleware `RequireClubRole`.

---

## Matriz de permissões

| Endpoint                              | Autenticado | Papel global  | Papel no clube                |
|---------------------------------------|:-----------:|:-------------:|:-----------------------------:|
| `GET /api/v1/info`                    | Não         | —             | —                             |
| `POST /api/v1/auth/register`          | Não         | —             | —                             |
| `POST /api/v1/auth/login`             | Não         | —             | —                             |
| `GET /api/v1/clubs`                   | Sim         | —             | —                             |
| `POST /api/v1/clubs`                  | Sim         | `ROLE_ADMIN`  | —                             |
| `GET /api/v1/clubs/:id`               | Sim         | —             | —                             |
| `DELETE /api/v1/clubs/:id`            | Sim         | —             | `ROLE_ADMIN`                  |
| `POST /api/v1/clubs/:id/join`         | Sim         | —             | —                             |
| `POST /api/v1/clubs/:id/suggestions`  | Sim         | —             | `ROLE_ADMIN` ou `ROLE_MEMBER` |
| `GET /api/v1/clubs/:id/suggestions`   | Sim         | —             | Qualquer membro               |
| `POST /api/v1/clubs/:id/themes/draw`  | Sim         | —             | `ROLE_ADMIN`                  |
| `POST /api/v1/clubs/:id/meetings`     | Sim         | —             | `ROLE_ADMIN`                  |

---

## Exemplos de requisição e resposta

### GET /api/v1/info
```json
{ "name": "BookClub API", "version": "1.0.0", "status": "ok" }
```

### POST /api/v1/auth/register
```json

{ "name": "Alice", "email": "alice@example.com", "password": "securepassword" }


{
  "id": "uuid",
  "name": "Alice",
  "email": "alice@example.com",
  "global_role": "ROLE_MEMBER",
  "created_at": "2026-05-04T10:00:00Z"
}
```

### POST /api/v1/auth/login
```json

{ "email": "alice@example.com", "password": "securepassword" }


{
  "token": "eyJhbGci...",
  "user": { "id": "uuid", "name": "Alice", "email": "alice@example.com", "global_role": "ROLE_MEMBER", "created_at": "..." }
}
```

### POST /api/v1/clubs
```json

{ "name": "Sci-Fi Readers", "description": "We love sci-fi" }


{
  "id": "uuid",
  "name": "Sci-Fi Readers",
  "description": "We love sci-fi",
  "owner": { "id": "uuid", "name": "Admin", "email": "admin@example.com", "global_role": "ROLE_ADMIN", "created_at": "..." },
  "created_at": "2026-05-04T10:00:00Z"
}
```

### POST /api/v1/clubs/:club_id/join
```json

{ "role": "ROLE_MEMBER" }


{ "id": "uuid", "user": { "..." }, "role": "ROLE_MEMBER", "joined_at": "2026-05-04T10:00:00Z" }
```

### POST /api/v1/clubs/:club_id/suggestions
```json

{ "title": "Dune", "author": "Frank Herbert", "description": "Classic sci-fi epic" }

```

### POST /api/v1/clubs/:club_id/themes/draw
```json

{ "year": 2026, "month": 5 }

```

### POST /api/v1/clubs/:club_id/meetings
```json

{
  "theme_id": "uuid",
  "scheduled_at": "2026-05-20T19:00:00Z",
  "location": "Online — Google Meet",
  "notes": "Bring your favourite quotes"
}

```

---

## Códigos de status HTTP

| Situação             | Status                       |
|----------------------|------------------------------|
| Recurso criado       | `201 Created`                |
| Deletado             | `204 No Content`             |
| Não autenticado      | `401 Unauthorized`           |
| Sem permissão        | `403 Forbidden`              |
| Não encontrado       | `404 Not Found`              |
| Duplicado            | `409 Conflict`               |
| Regra de negócio     | `422 Unprocessable Entity`   |
| Erro inesperado      | `500 Internal Server Error`  |

Formato de erro padrão:
```json
{ "code": "not_found", "message": "club not found" }
```

Códigos comuns: `validation_error` · `unauthorized` · `forbidden` · `not_found` · `conflict` · `unprocessable_entity` · `internal_error`
