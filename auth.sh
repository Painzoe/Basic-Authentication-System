#!/bin/sh

# Simple CLI client for the Go Auth Service
# Requires: curl

API_URL="http://localhost:8080/api"

show_usage() {
    echo "Usage: $0 [register|login] <username> <password>"
    exit 1
}

# Check for required arguments
if [ "$#" -ne 3 ]; then
    show_usage
fi

COMMAND=$1
USERNAME=$2
PASSWORD=$3

case "$COMMAND" in
    register)
        echo "Registering user: $USERNAME..."
        curl -s -f -X POST "$API_URL/register" 
            -d "username=$USERNAME" 
            -d "password=$PASSWORD"
        if [ $? -eq 0 ]; then
            echo "Registration successful."
        else
            echo "Registration failed. User might already exist."
            exit 1
        fi
        ;;
    login)
        echo "Logging in user: $USERNAME..."
        # -i to see headers, -s for silent
        RESPONSE=$(curl -s -w "%{http_code}" -X POST "$API_URL/login" 
            -d "username=$USERNAME" 
            -d "password=$PASSWORD")
        
        # Get the HTTP status code (last 3 chars)
        STATUS=$(echo "$RESPONSE" | tail -c 4)
        
        if [ "$STATUS" = "200" ]; then
            echo "Login successful."
        elif [ "$STATUS" = "401" ]; then
            echo "Login failed: Unauthorized."
            exit 1
        else
            echo "Login failed with status: $STATUS"
            exit 1
        fi
        ;;
    *)
        show_usage
        ;;
esac
