# Generate a Certificate Authority (CA)
# mkdir certs
# openssl genrsa -out ./certs/rootCA.key 2048
# openssl req -x509 -new -nodes -key ./certs/rootCA.key -sha256 -days 365 -out ./certs/rootCA.crt -subj "/CN=MyLocalCA"
#
# # Create Server Certificate Files
# openssl genrsa -out ./certs/server.key 2048
# openssl req -new -key ./certs/server.key -out ./certs/server.csr -subj "/CN=localhost"
# openssl x509 -req -in ./certs/server.csr -CA ./certs/rootCA.crt -CAkey ./certs/rootCA.key -CAcreateserial -out ./certs/server.crt -days 365 -sha256
#
# # Create Client Certificate Files
# openssl genrsa -out ./certs/client.key 2048
# openssl req -new -key ./certs/client.key -out ./certs/client.csr -subj "/CN=postgres"
# openssl x509 -req -in ./certs/client.csr -CA ./certs/rootCA.crt -CAkey ./certs/rootCA.key -CAcreateserial -out ./certs/client.crt -days 365 -sha256

# Set podman permissions
podman unshare mkdir -p ./certs
#
# Generate certificates inside the namespace
podman unshare sh -c '
  cd ./certs
  openssl genrsa -out rootCA.key 2048
  openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 365 -out rootCA.crt -subj "/CN=MyLocalCA"
  openssl genrsa -out server.key 2048
  openssl req -new -key server.key -out server.csr -subj "/CN=localhost"
  openssl x509 -req -in server.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out server.crt -days 365 -sha256
  openssl genrsa -out client.key 2048
  openssl req -new -key client.key -out client.csr -subj "/CN=postgres"
  openssl x509 -req -in client.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out client.crt -days 365 -sha256
  chmod 600 server.key rootCA.key client.key
  chmod 644 server.crt rootCA.crt client.crt
'

# Verify ownership (from namespace perspective)
podman unshare chown -R 999:999 ./certs
podman unshare ls -la ./certs/
