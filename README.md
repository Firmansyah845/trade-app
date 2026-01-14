### installation

how to run this project:

1. clone the repo
2. enter the directory
3. run go mod tidy
4. copy the file `configs/application.sample.yaml` to `configs/application.yaml`
5. update the `application.yaml` based on your configuration
6. run the project
   ```zsh
   go run main.go server
   ```
7. visit the url and port based on configuration on `configs/application.yaml`
    ```zsh
   http://localhost:3000/ping
   ```