#!/bin/bash

if [ ! -f "doom_cannon" ]; then
    echo "❌ doom_cannon not found!"
    exit 1
fi

echo "▶ Running doom_cannon..."
python3 doom_cannon
