# Generate a Certificate Authority (CA)
mkdir certs
openssl genrsa -out ./certs/rootCA.key 2048
openssl req -x509 -new -nodes -key ./certs/rootCA.key -sha256 -days 365 -out ./certs/rootCA.crt -subj "/CN=MyLocalCA"

# Create Server Certificate Files
openssl genrsa -out ./certs/server.key 2048
openssl req -new -key ./certs/server.key -out ./certs/server.csr -subj "/CN=localhost"
openssl x509 -req -in ./certs/server.csr -CA ./certs/rootCA.crt -CAkey ./certs/rootCA.key -CAcreateserial -out ./certs/server.crt -days 365 -sha256

# Create Client Certificate Files
openssl genrsa -out ./certs/client.key 2048
openssl req -new -key ./certs/client.key -out ./certs/client.csr -subj "/CN=postgres"
openssl x509 -req -in ./certs/client.csr -CA ./certs/rootCA.crt -CAkey ./certs/rootCA.key -CAcreateserial -out ./certs/client.crt -days 365 -sha256

chown 999:999 ./certs/server.key
chown 999:999 ./certs/rootCA.crt
chmod 600 ./certs/rootCA.crt
chmod 600 ./certs/server.key
# chmod 400 ./certs/.*-key.pem
# chmod 444 ./certs/.*-cert.pem
