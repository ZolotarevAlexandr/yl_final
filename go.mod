module github.com/ZolotarevAlexandr/yl_final

go 1.24

toolchain go1.24.0

replace (
	github.com/ZolotarevAlexandr/yl_final/agent => ./agent
	github.com/ZolotarevAlexandr/yl_final/calculator => ./calculator
	github.com/ZolotarevAlexandr/yl_final/db => ./db
	github.com/ZolotarevAlexandr/yl_final/orchestrator => ./orchestrator
	github.com/ZolotarevAlexandr/yl_final/grpc => ./grpc
)

require (
	github.com/ZolotarevAlexandr/yl_final/agent v0.0.0-00010101000000-000000000000
	github.com/ZolotarevAlexandr/yl_final/orchestrator v0.0.0-00010101000000-000000000000
)

require (
	github.com/ZolotarevAlexandr/yl_final/calculator v0.0.0-00010101000000-000000000000 // indirect
	github.com/ZolotarevAlexandr/yl_final/db v0.0.0-00010101000000-000000000000 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/grpc v1.72.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gorm.io/driver/sqlite v1.5.7 // indirect
	gorm.io/gorm v1.26.1 // indirect
)
