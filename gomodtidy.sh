#!/usr/bin/env bash

go env -w GOPRIVATE=*
go mod vendor
