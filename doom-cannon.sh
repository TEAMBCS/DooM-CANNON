#!/bin/bash
#code by BLACK ZERO
if [ ! -f "doom-cannon.py" ]; then
    echo "❌ doom-cannon.py not found!"
    exit 1
fi

echo "▶ Running DooM-CANNON..."
python3 doom-cannon.py
