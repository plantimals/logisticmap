#!/bin/bash

go clean

go build animateLogMap.go

time ./animateLogMap
