#!/bin/bash
hey -n 100 -c 10 -m POST \
"http://localhost:8080/api/create/random/human/with/transaction"
