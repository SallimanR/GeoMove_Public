# sudo openssl genrsa -out our-site.org.key 2048
# openssl req -nodes -new -key our-site.org.key -out ca.csr
# openssl x509 -req -days 365 -in our-site.org.csr -signkey our-site.org.key -out our-site.org.crt
mkdir certs
# cp our-site.org.crt our-site.org.key our-site.org.csr ./certs
sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./certs/nginx-selfsigned.key -out ./certs/nginx-selfsigned.crt

chown -R 100999:100999 ./certs
