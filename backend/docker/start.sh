#!/bin/sh

# Start Postfix and msmtp in the background
postfix start &
sleep 0.1
# Check if both Postfix and msmtp are listening on their respective ports
while ! netstat -tln | grep -qE ':\s*25\s*'; do
    echo "Waiting for Postfix to start listening on port 25..."
    sleep 0.05
done

# while ! netstat -tln | grep -qE ':\s*587\s*'; do
#     echo "Waiting for msmtp to start listening on port 587..."
#     sleep 1
# done

# Run the backend application
./backend
