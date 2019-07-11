default: deps
	go install

deps:
	go mod tidy
	go mod vendor

run: default
	cd ${GOPATH}/bin;./go-live

docker: deps
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
	chmod +x ./main
	sudo docker build -t go-live .

clean:
	rm -rf vendor/
	rm -rf go.sum

.PHONY: go test
