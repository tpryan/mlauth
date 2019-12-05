BASEDIR = $(shell pwd)
PROJECT = $(mlauth_project)
BUCKET = $(mlauth_bucket)

env:
	gcloud config set project $(PROJECT)

init: env
	gsutil mb gs://$(BUCKET)
	gsutil -m cp $(BASEDIR)/vision/testdata/* gs://$(BUCKET)/vision 
	gsutil -m acl -r ch -u AllUsers:R gs://$(BUCKET)/