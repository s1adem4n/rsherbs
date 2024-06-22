cd frontend || exit
bun install
bun run build
cd ..
go build -o build/main main.go
