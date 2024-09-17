SSO project to study gRPC framework in golang.

How to init your workspace:
1. Install golang 1.22.5 or higher (learn more: https://go.dev/)
2. Install sqlite (learn more: https://www.sqlite.org/)
3. Install task util (learn more: https://taskfile.dev/)
4. Create file ./storage/sso.db or any other, but don't forget to change STORE_PATH in `.env` file
5. Install fieldalignment linter:
```bash
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
```
6. Run the application
```bash
task run
```

Additionally you are able to run tests:
1. Run the application
2. Run tests
```bash
task test
```