#!/bin/sh


if [ !  -f ~/build/src/eurus-backend/doc/EurusAPI.yaml ]; then
    echo "File not found"
    exit 1
fi

rm ~/eurus/api/EurusAPI.yaml
cp -f ~/build/src/eurus-backend/doc/$1 ~/eurus/api/

echo "Deployed"
