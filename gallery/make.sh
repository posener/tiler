#! /usr/bin/env bash

go run ./cmd/tiler/ -img in/cake.png -tiles tiles/circle.png \
  -out gallery/cake.png -shift 1,1 -colors 8 -scale 1,0.8,0.6,0.4,0.2

go run ./cmd/tiler/ -img in/starry-night.png -tiles tiles/circle.png \
  -out gallery/starry-night.png -colors 16 -scale 0.1

go run ./cmd/tiler/ -img in/starry-night.png -tiles tiles/circle.png \
  -out gallery/starry-night-shift-1.png -shift 1,1 -colors 4 -scale 0.1
