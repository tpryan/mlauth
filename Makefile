BASEDIR = $(shell pwd)
PROJECT = $(mlauth_project)
BUCKET = $(mlauth_bucket)

env:
	gcloud config set project $(PROJECT)

init: env bucket serviceaccount services key

serviceaccount:
	gcloud iam service-accounts create mlauth

key:
	gcloud iam service-accounts keys create creds/creds.json \
  --iam-account mlauth@$(PROJECT).iam.gserviceaccount.com

bucket:
	-gsutil mb gs://$(BUCKET)
	gsutil -m cp $(BASEDIR)/vision/testdata/* gs://$(BUCKET)/vision 
	gsutil -m cp $(BASEDIR)/speech/testdata/* gs://$(BUCKET)/speech 
	gsutil -m acl -r ch -u AllUsers:R gs://$(BUCKET)/

services:
	gcloud services enable vision.googleapis.com	
	gcloud services enable speech.googleapis.com
	gcloud services enable language.googleapis.com

clean: 
	-gsutil rm -rf gs://$(BUCKET)/*
	-gsutil rb gs://$(BUCKET)
	-gcloud iam service-accounts delete mlauth@$(PROJECT).iam.gserviceaccount.com -q

test:
	cd language && go test
	cd speech && go test
	cd vision && go test