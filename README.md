Directory structure: 
https://github.com/golang-standards/project-layout

GIN Middleware:
https://gin-gonic.com/docs/examples/using-middleware/

# Snippits

Pretty-print struct:
```Go
jsonData, err := json.MarshalIndent(cfg, "", "    ")
fmt.Println(string(jsonData))
```

Add GIN to the project:
```bash
go get -u github.com/gin-gonic/gin
```

Create a new public/private key pair for JWT:
```bash
openssl ecparam -genkey -name prime256v1 -noout -out private.pem
openssl ec -in private.pem -pubout -out public.pem
```