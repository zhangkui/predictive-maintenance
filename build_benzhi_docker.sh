#!/usr/bin/env sh

set -eu

image_name="${IMAGE_NAME:-predictive-maintenance}"
image_tag="${IMAGE_TAG:-benzhi}"

docker build \
  --file benzhi.Dockerfile \
  --tag "${image_name}:${image_tag}" \
  .
