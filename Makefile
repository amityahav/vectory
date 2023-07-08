.PHONEY: all

swagger:
	@swagger generate server --target gen/api --name Vectory --spec api/spec.yaml --principal string --exclude-main --keep-spec-order
	@swagger generate client -f api/spec.yaml -t pkg -c client

ent:
	go run -mod=mod entgo.io/ent/cmd/ent generate ./db/metadata/schema --target=./gen/ent