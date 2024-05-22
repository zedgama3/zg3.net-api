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
```
go get -u github.com/gin-gonic/gin
```
