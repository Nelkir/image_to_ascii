#!/bin/bash
go build -o bin/image_to_ascii . &&
sudo cp bin/image_to_ascii /usr/local/bin/
