#!/bin/bash

echo 'eval "$(mise activate bash)"' >> ~/.bashrc
mise trust

go install github.com/air-verse/air@latest
