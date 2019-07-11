module go-live

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/go-sql-driver/mysql v1.4.1
	github.com/google/go-cmp v0.3.0 // indirect
	github.com/jinzhu/gorm v1.9.8
	github.com/joho/godotenv v1.3.0
	github.com/julienschmidt/httprouter v1.2.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/orcaman/concurrent-map v0.0.0-20190314100340-2693aad1ed75
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.3.0
	google.golang.org/appengine v0.0.0-00010101000000-000000000000 // indirect
)

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.26.0
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190123085648-057139ce5d2b
	golang.org/x/lint => github.com/golang/lint v0.0.0-20181026193005-c67002cb31c3
	golang.org/x/net => github.com/golang/net v0.0.0-20190311183353-d8887717615a
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20180821212333-d2e6202438be
	golang.org/x/sync => github.com/golang/sync v0.0.0-20181108010431-42b317875d0f
	golang.org/x/sys => github.com/golang/sys v0.0.0-20180905080454-ebe1bf3edb33
	golang.org/x/text => github.com/golang/text v0.3.0
	golang.org/x/time => github.com/golang/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190328211700-ab21143f2384
	google.golang.org/appengine => github.com/golang/appengine v1.1.0
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20180817151627-c66870c02cf8
	google.golang.org/grpc => github.com/grpc/grpc-go v1.20.1
)
