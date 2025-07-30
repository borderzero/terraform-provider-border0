.PHONY: init apply destroy

init:
	(cd ../../; make install)
	rm -rf .terraform*
	terraform init

plan:
	terraform plan -parallelism=100

apply:
	terraform apply -parallelism=20

destroy:
	terraform destroy
