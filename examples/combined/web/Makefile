BASEDIR = $(shell pwd)
PROJECT = $(mlauth_project)

env:
	gcloud config set project $(PROJECT)

init: env
	gcloud app create --project=$(PROJECT) --region=us-central -q
	gcloud services enable cloudbuild.googleapis.com

deploy: env
	gcloud app deploy -q

dev:
	go run main.go	