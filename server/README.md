
### Database Migrations
- Database migrations are handled by goose
- Make sure goose is installed and available, to install goose run:
    ```bash
    go install github.com/pressly/goose/v3/cmd/goose@latest
    ```

- Run ```source ./scripts/migrate-env.sh``` for setting up goose environment variables.

### How to use run backend

* **For development:**

  ```bash
  docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
  ```
* **For production:**

  ```bash
  docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
  ```

---

https://platform.openai.com/docs/guides/function-calling