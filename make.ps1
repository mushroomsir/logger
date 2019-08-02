

function test {
    go test -coverprofile="logger.coverprofile" ./...
    go tool cover -html="logger.coverprofile"
    Remove-Item *.coverprofile
}

switch ($args[0]) {
    "test" { test }
    Default { "Not support " + $args[0] }
}
