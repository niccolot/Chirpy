#! usr/bin/bash

JWT_SECRET=$(openssl rand -base64 64)

POLKA_API_KEY="f271c81ff7084ee5b99a5091b42d486e"

cat <<EOF > .env
JWT_SECRET=$JWT_SECRET
POLKA_API_KEY=$POLKA_API_KEY
EOF

