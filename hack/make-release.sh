#!/usr/bin/env bash

set -e

goBack() {
  num=$1
  echo -ne "\033[${num}D"
}

# Setup for release
rm -rf "Release"

echo
echo " ===> Creating binaries for all major platforms <=== "
echo
gox -parallel=2 -ldflags="$LDFLAGS" -tags "$TAGS" -output="Release/staging/containerenv-{{.OS}}-{{.Arch}}" -os "linux windows darwin" -arch "amd64 386" github.com/jaredallard/containerenv/cmd/...

echo 
echo " ===> Setting up staging for release generation <=== "
echo 
mkdir -p "Release/staging" "Release/binaries" || true

echo " ===> Creating releases <=== "
echo 
pushd "Release/staging" >/dev/null || exit 1
for bin in *; do
  name=$(sed 's/\.exe$//' <<< "$bin")

  echo -ne "  -> $name 0%"
  cp ../../README.md ../../LICENSE ./

  mv "$bin" "containerenv"
  chmod +x "containerenv"

  goBack "2"
  echo -n "20%"
  tar --transform 's,^,,' -cvf "$name.tar" "containerenv" "README.md" "LICENSE" >/dev/null

  goBack "3"
  echo -n "50%"
  xz -e -T 0 "$name.tar"

  goBack "3"
  echo -n "99%"
  sha256sum "$name.tar.xz" | awk '{ print $1 }' > "$name.tar.xz.sha256"
  mv "$name.tar.xz"* ..
  mv "containerenv" "../binaries/$name"

  goBack "3"
  echo -e "\033[32mdone\033[0m"
done
popd >/dev/null || exit 1

echo 
echo " ===> Release created in './Release' <==="
echo
