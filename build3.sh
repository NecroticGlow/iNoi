#!/usr/bin/env bash
set -e

appName="iNoi"
outDir="build"
zipName="${appName}-windows-386.zip"

builtAt="$(date +'%F %T %z')"
gitCommit=$(git rev-parse --short HEAD || echo unknown)
gitAuthor="The iNoi Projects Contributors <inoi@peifeng.li>"
version="windows-386"

echo "== Build Windows 386 =="
echo "commit: $gitCommit"

# ldflags（不依赖 tag）
ldflags="\
-w -s \
-X 'github.com/OpenListTeam/OpenList/v4/internal/conf.BuiltAt=$builtAt' \
-X 'github.com/OpenListTeam/OpenList/v4/internal/conf.GitAuthor=$gitAuthor' \
-X 'github.com/OpenListTeam/OpenList/v4/internal/conf.GitCommit=$gitCommit' \
-X 'github.com/OpenListTeam/OpenList/v4/internal/conf.Version=$version' \
"

# 准备目录
rm -rf "$outDir"
mkdir -p "$outDir/tmp"

# 构建 Windows 386
export GOOS=windows
export GOARCH=386
export CGO_ENABLED=0

go build -o "$outDir/tmp/$appName.exe" -ldflags="$ldflags" -tags=jsoniter .

# 打包 zip
cd "$outDir/tmp"
zip "../$zipName" "$appName.exe"
cd ../..

rm -rf "$outDir/tmp"

echo "== Done =="
echo "Output: $outDir/$zipName"
