//go:build tools

package main

import _ "github.com/matryer/moq"

///go:build tools는 Go 빌드 태그(build tag)로, 이 파일이 특정 빌드 조건에서만 포함되도록 설정하는 주석입니다.
///Go의 go build, go install 등의 명령어를 사용할 때 -tags=tools 옵션을 추가하면 이 파일이 포함됩니다.
