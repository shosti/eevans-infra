.PHONY: awx bootstrap

awx:
	ansible-playbook awx.yml

deps:
	ansible-galaxy role install -r requirements.yml && \
	ansible-galaxy collection install -r requirements.yml
