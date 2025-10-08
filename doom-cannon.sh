#!/bin/bash

if [ ! -f "doom-cannon" ]; then
    echo "❌ doom-cannon not found!"
    exit 1
fi

echo "▶ Running  DooM-CANNON..."
python3 doom_cannon
