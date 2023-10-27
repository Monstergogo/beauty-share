protoc:
	cd api/protobuf-spec && $ protoc --go_out=. --go_opt=paths=source_relative \
                       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
                       *.proto

wire:
	cd internal/injector && wire

ver = 'v1.0.0'
release:
	npm run release -- --release-as $(ver)

push:
	git push --follow-tags origin main