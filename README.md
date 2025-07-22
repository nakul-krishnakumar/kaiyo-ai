
### Database Migrations
- Database migrations are handled by goose
- Make sure goose is installed and available, to install goose run:
    ```bash
    go install github.com/pressly/goose/v3/cmd/goose@latest
    ```

- Run ```source ./scripts/migrate-env.sh``` for setting up goose environment variables.