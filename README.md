Directory structure: 
https://github.com/golang-standards/project-layout

GIN Middleware:
https://gin-gonic.com/docs/examples/using-middleware/

# TODO
- [ ] Tidy file structure
  - [ ] Create folders and files for managing the database and move the code there.  If needed, move structs as well.

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