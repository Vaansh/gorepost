#!/bin/bash

# Remove based on OS type
if [[ "$OSTYPE" == "linux-gnu" || "$OSTYPE" == "darwin"* ]]; then
  rm ../data/*.mp4
  rm ../log/*.log
elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
  del ..\data\*.mp4
  del ..\log\*.log
else
  echo "Unknown operating system"
fi
