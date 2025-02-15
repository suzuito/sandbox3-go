.PHONY: mac-init
mac-init:

.PHONY: google-cloud-auth-init
google-cloud-auth-init:
	gcloud auth login
	gcloud auth application-default login
