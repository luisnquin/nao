#!/bin/bash

last_tag=$(git tag | tail -n 1)
content=""

while read -r line;
do
    if [[ "$line" == *"Version"* ]]; then
        content+="Version string = \"$last_tag\"" && content+="\n"
    else
        content+="$line" && content+="\n"
    fi
done < ./src/constants/app.go

echo "$content" > ./src/constants/app.go

gofmt -w ./src/constants/app.go
