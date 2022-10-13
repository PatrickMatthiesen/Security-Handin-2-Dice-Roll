# # 1. Generate CA's private key and self-signed certificate
# openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=DK/CN=My CA"

# echo "CA's self-signed certificate"
# openssl x509 -in ca-cert.pem -noout -text

# # 2. Generate web server's private key and certificate signing request (CSR)
# openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=DK/CN=My Server"

# # 3. Use CA's private key to sign web server's CSR and get back the signed certificate
# openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf

# echo "Server's signed certificate"
# openssl x509 -in server-cert.pem -noout -text

# openssl req -new -x509 -days 365 -nodes -out server-cert.pem -keyout server-key.pem -subj "/C=DK/ST=Denmark/O=ITU/CN=example.org/emailAddress=server@itu.dk"
# openssl req -new -x509 -days 365 -nodes -out client-cert.pem -keyout client-key.pem -subj "/C=DK/ST=Denmark/O=ITU/CN=example.org/emailAddress=server@itu.dk"

$name = "alice"
openssl req -x509 -newkey rsa:4096 -keyout "$name.key.pem" -out "$name.cert.pem" -sha256 -days 365 -nodes -addext "subjectAltName = DNS:localhost,DNS:$name"
$name = "bob"
openssl req -x509 -newkey rsa:4096 -keyout "$name.key.pem" -out "$name.cert.pem" -sha256 -days 365 -nodes -addext "subjectAltName = DNS:localhost,DNS:$name"

