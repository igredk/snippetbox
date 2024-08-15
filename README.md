# snippetbox
To run application use follow command from root directory:
```
go run ./cmd/web
```
Run and debug application via VSCode. Put follow config in **snippetbox/.vscode/launch.json** file:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "run app",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/web",
            "cwd": "${workspaceFolder}",
        }
    ]
}
```
