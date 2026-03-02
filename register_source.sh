#!/bin/bash

NAME=$1

if [ -z "$NAME" ]; then
    echo "Usage: $0 <source_name>"
    exit 1
fi

curl -X POST http://localhost:8080/sources 
     -H "Content-Type: application/json" 
     -d "{"name": "$NAME"}"
