#!/bin/bash

# Set the Common Name for the Certificate Authority
CA_NAME="MyCA"

# Set the Common Names for the SSL certificates
MANAGER_CERT_NAME="Manager"
WORKER_CERT_NAME="Worker"

# Set folder names for each server
MANAGER_FOLDER="manager"
WORKER_FOLDER="worker"

# Set IP and hostname information
MANAGER_IP="192.168.1.10"
MANAGER_HOSTNAME="manager.local"

WORKER_IP="192.168.1.20"
WORKER_HOSTNAME="worker.local"

# Create directories to store the CA and certificate files
mkdir -p certs/${MANAGER_FOLDER}
mkdir -p certs/${WORKER_FOLDER}

# Step 1: Generate a private key for the Certificate Authority (CA)
openssl genpkey -algorithm RSA -out certs/ca-key.pem

# Step 2: Generate a self-signed certificate for the CA
openssl req -x509 -new -key certs/ca-key.pem -out certs/ca-cert.pem -subj "/CN=${CA_NAME}"

# Copy the CA certificate to each server folder
cp certs/ca-cert.pem certs/${MANAGER_FOLDER}/
cp certs/ca-cert.pem certs/${WORKER_FOLDER}/

# Step 3: Generate a private key for the Manager SSL certificate
openssl genpkey -algorithm RSA -out certs/${MANAGER_FOLDER}/key.pem

# Step 4: Generate a Certificate Signing Request (CSR) for the Manager SSL certificate
openssl req -new -key certs/${MANAGER_FOLDER}/key.pem -out certs/${MANAGER_FOLDER}/csr.pem -subj "/CN=${MANAGER_CERT_NAME}" -addext "subjectAltName = IP:${MANAGER_IP},DNS:${MANAGER_HOSTNAME}"

# Step 5: Sign the Manager SSL certificate with the CA
openssl x509 -req -in certs/${MANAGER_FOLDER}/csr.pem -CA certs/ca-cert.pem -CAkey certs/ca-key.pem -out certs/${MANAGER_FOLDER}/cert.pem -CAcreateserial -extfile <(printf "subjectAltName = IP:${MANAGER_IP},DNS:${MANAGER_HOSTNAME}")

# Step 6: Generate a private key for the Worker SSL certificate
openssl genpkey -algorithm RSA -out certs/${WORKER_FOLDER}/key.pem

# Step 7: Generate a Certificate Signing Request (CSR) for the Worker SSL certificate
openssl req -new -key certs/${WORKER_FOLDER}/key.pem -out certs/${WORKER_FOLDER}/csr.pem -subj "/CN=${WORKER_CERT_NAME}" -addext "subjectAltName = IP:${WORKER_IP},DNS:${WORKER_HOSTNAME}"

# Step 8: Sign the Worker SSL certificate with the CA
openssl x509 -req -in certs/${WORKER_FOLDER}/csr.pem -CA certs/ca-cert.pem -CAkey certs/ca-key.pem -out certs/${WORKER_FOLDER}/cert.pem -CAcreateserial -extfile <(printf "subjectAltName = IP:${WORKER_IP},DNS:${WORKER_HOSTNAME}")

# Optional: Display information about the generated certificates
echo "Manager Certificate:"
openssl x509 -in certs/${MANAGER_FOLDER}/cert.pem -noout -text

echo "Worker Certificate:"
openssl x509 -in certs/${WORKER_FOLDER}/cert.pem -noout -text

echo "Certificates and CA generated successfully. Files are located in the 'certs' directory."
