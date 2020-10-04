#!/usr/bin/env bash

clear
go build -ldflags "\
-X main.taigaUsername=your-username \
-X main.taigaPassword=your-password \
-X main.sandboxProjectSlug=my-sandbox-slug \
-X main.sandboxEpicID=1234567890 \
-X main.sandboxFileUploadPath=/tmp/bad-puns-make-me-sic.png\
" \
-o taigo-dev.exe || ECHO Failed to build the binary && exit 1

./taigo-dev
