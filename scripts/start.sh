#! usr/bin/bash

if [ "$PG_STATUS" != "active" ]; then
    echo "PostgreSQL is not running. Starting PostgreSQL..."
    sudo systemctl start postgresql
    
    # Check again if it started successfully
    if [ "$(systemctl is-active postgresql)" == "active" ]; then
        echo "PostgreSQL started successfully."
    else
        echo "Failed to start PostgreSQL."
    fi
else
    echo "PostgreSQL is already running."
fi

go build -o out && ./out  
