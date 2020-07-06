param(
  [string]$target="all"
)

$REPO="github.com/bitmaelum/bitmaelum-suite"

$APPS="bm-server","bm-client","bm-config","bm-client-ui"
$TOOLS="hash-address","jwt","proof-of-work","readmail"

$BUILD_DATE=Get-Date
$COMMIT=git rev-parse HEAD
$PKG=go list ./core
$GO_BUILD_FLAGS="-X '$PKG.BuildDate=$BUILD_DATE' -X '$PKG.GitCommit=$COMMIT'"

Write-Host -NoNewLine "Compiling [" 

foreach ($app in $APPS)
{
  go build -ldflags $GO_BUILD_FLAGS -o release/windows/$app.exe $REPO/cmd/$app
  Write-Host -NoNewLine "."
}

foreach ($tool in $TOOLS)
{
  go build -ldflags $GO_BUILD_FLAGS -o release/windows/$tool.exe $REPO/tools/$tool
  Write-Host -NoNewLine "."
}


Write-Host "]"