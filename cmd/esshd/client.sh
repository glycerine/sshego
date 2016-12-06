#!/bin/bash

ssh-add -D
ssh-add ./id_rsa_alice

echo "login as alice. The first password is 'answer1'."

echo

echo "The 2FA/goole-authenticator code is given by the QR-code in example.otp-qrcode.png / example.otp"

echo

ssh alice@localhost -p 2022
