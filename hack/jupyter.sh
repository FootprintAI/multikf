#!/usr/bin/env bash

docker run -it --rm \
	--gpus all \
	-p 8888:8888 \
	-v /data/jupyter-volume:/home/jovyan/data \
	quay.io/jupyter/datascience-notebook:python-3.12.10
