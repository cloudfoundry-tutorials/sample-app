#!/bin/bash

version=$(cat version/version)
pushd source 
	zip -r sample-app.zip ./* -x manifest.yml
popd
cp source/sample-app.zip rc/sample-app-${version}.zip
cp source/manifest.yml rc/manifest-${version}.yml