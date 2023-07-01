GOOS=darwin GOARCH=amd64 go build
file go-karton
mv go-karton build/go-karton
cp config.toml build/config.toml
zip -r go-karton.zip build
