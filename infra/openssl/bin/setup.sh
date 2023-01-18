#!/bin/bash

# SITE="kyle.speer.infra.challenge.s3-website-us-east-1.amazonaws.com"
SITE=$1

echo "Creating Certificate Authority (CA)..."
openssl req -x509 \
            -sha256 -days 356 \
            -nodes \
            -newkey rsa:2048 \
            -subj "/CN=${SITE}/C=US/L=Philadelphia" \
            -keyout rootKey.pem -out rootCA.pem

echo "Creating private key..."
openssl genrsa -out serverKey.pem 2048

echo "Configuring CSR (Certificate Signing Request)..."
cat > csr.conf <<EOF
[ req ]
default_bits = 2048
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[ dn ]
C = US
ST = Pennsylvania
L = Philadelphia
O = speer
OU = kylespeer
CN = $SITE

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = $SITE

EOF

echo "Generating CSR..."
openssl req -new -key serverKey.pem -out server.csr -config csr.conf


echo "Configuring SSL (Secure Socket Layer)..."
cat > cert.conf <<EOF

authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = $SITE

EOF

openssl x509 -req \
    -in server.csr \
    -CA rootCA.pem -CAkey rootKey.pem \
    -CAcreateserial -out server.pem \
    -days 365 \
    -sha256 -extfile cert.conf
