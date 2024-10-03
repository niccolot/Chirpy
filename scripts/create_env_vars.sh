#! usr/bin/bash

JWT_SECRET=$(openssl rand -base64 64)

POLKA_API_KEY="f271c81ff7084ee5b99a5091b42d486e"
DB_URL="postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable" 
PLATFORM="dev"

cat <<EOF > .env
JWT_SECRET=$JWT_SECRET
POLKA_API_KEY=$POLKA_API_KEY
DB_URL=$DB_URL
PLATFORM=$PLATFORM
EOF

