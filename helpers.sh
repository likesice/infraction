#!/usr/bin/sh

load_env() {
while IFS== read -r key value; do
  printf -v "$key" %s "$value" && export "$key"
done <.env
}

load_env