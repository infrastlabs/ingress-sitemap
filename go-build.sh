# go get -u github.com/go-bindata/go-bindata/... 
go-bindata -fs -pkg asset -prefix "static/" -ignore='(\.gitignore$$|\.map$$)' -o asset/swaggerui-data.go static/...


go build ./
