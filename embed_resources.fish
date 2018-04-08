#!/usr/bin/env fish

# Embed WebApp
echo "Embedding Webapp"
cd server
rice embed-go

# Embed Fritz lib fonts
echo "Embedding Fritz fonts"
cd ../fritz
rice embed-go
cd ..