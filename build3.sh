#!/usr/bin/env bash
set -e

appName="iNoi"
outDir="build"
zipName="${appName}-windows-386.zip"

builtAt="$(date +'%F %T %z')"
gitCommit=$(git rev-parse --short HEAD || echo unknown)
gitAuthor="The iNoi Projects Contributors <inoi@peifeng.li>"
version="windows-386"
webVersion=$(wget -qO- -t1 -T2 "https://api.github.com/repos/li-peifeng/iNoi-Web/releases/latest" | grep "tag_name" | head -n 1 | awk -F ":" '{print $2}' | sed 's/\"//g;s/,//g;s/ //g' || echo unknown)

echo "== Build Windows 386 =="
echo "commit: $gitCommit"

ldflags="\
-w -s \
-X 'github.com/OpenListTeam/OpenList/v4/internal/conf.BuiltAt=$builtAt' \
-X 'github.com/OpenListTeam/OpenList/v4/internal/conf.GitAuthor=$gitAuthor' \
-X 'github.com/OpenListTeam/OpenList/v4/internal/conf.GitCommit=$gitCommit' \
-X 'github.com/OpenListTeam/OpenList/v4/internal/conf.Version=$version' \
-X 'github.com/OpenListTeam/OpenList/v4/internal/conf.WebVersion=$webVersion' \
"

FetchWebRelease() {
  curl -L https://github.com/li-peifeng/iNoi-Web/releases/latest/download/dist.tar.gz -o dist.tar.gz
  tar -zxvf dist.tar.gz
  rm -rf public/dist
  mv -f dist public
  rm -rf dist.tar.gz
}

rm -rf "$outDir"
mkdir -p "$outDir/tmp"

FetchWebRelease

export GOOS=windows
export GOARCH=386
export CGO_ENABLED=1
export CC=i686-w64-mingw32-gcc
export CXX=i686-w64-mingw32-g++

go build -o "$outDir/tmp/$appName.exe" -ldflags="$ldflags" -tags=jsoniter .

mkdir -p "$outDir/tmp/public"
cp -R public/dist "$outDir/tmp/public/"

cd "$outDir/tmp"
zip -r "../$zipName" "$appName.exe" "public/dist"
cd ../..

rm -rf "$outDir/tmp"

echo "== Done =="
echo "Output: $outDir/$zipName"
