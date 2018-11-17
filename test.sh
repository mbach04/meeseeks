#!/bin/bash
#Use ansible to test meeseeks application responses

#generate certs
mkdir ./certs
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout certs/key.pem -out certs/cert.pem

#run tests
ansible-playbook -i ansible/inventory ansible/test_login.yml
ansible-playbook -i ansible/inventory ansible/test_ls.yml

#cleanup
rm -rf ./certs

