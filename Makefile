BASEDIR = $(shell pwd)
PROJECT = $(mlauth_project)
BUCKET = $(mlauth_bucket)

env:
	gcloud config set project $(PROJECT)

init: env
	-gsutil mb gs://$(BUCKET)
	gsutil -m cp $(BASEDIR)/vision/testdata/* gs://$(BUCKET)/vision 
	gsutil -m cp $(BASEDIR)/speech/testdata/* gs://$(BUCKET)/speech 
	gsutil -m acl -r ch -u AllUsers:R gs://$(BUCKET)/


services:
	gcloud services enable vision.googleapis.com	
	gcloud services enable speech.googleapis.com
	gcloud services enable language.googleapis.com


clean: 
	-gsutil rb gs://$(BUCKET)