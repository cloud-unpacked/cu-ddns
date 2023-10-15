#!/usr/bin/env bash

set -e

rm -rf ./completions
mkdir ./completions

for sh in bash fish zsh; do
	go run ./ddns/main.go completion "$sh" >"completions/cu-ddns.$sh"
done
