$esc=$([char]27)

Write-Output "$esc[36;1m*** Format check$esc[0m"
gofmt -l .

Write-Output "$esc[36;1m*** Vet check$esc[0m"
go vet ./...

Write-Output"$esc[36;1m*** Lint check$esc[0m"
go/bin/golint ./...

Write-Output "$esc[36;1m*** Tests$esc[0m"
go test ./...
