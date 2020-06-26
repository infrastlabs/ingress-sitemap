# go get -u github.com/go-bindata/go-bindata/... 
echo "go-bindata..."
go-bindata -fs -pkg asset -prefix "static/" -ignore='(\.gitignore$$|\.map$$)' -o asset/swaggerui-data.go static/...

echo "go build..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64
go build -ldflags "-s -w"
# go build ./
