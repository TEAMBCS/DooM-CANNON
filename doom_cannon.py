#!/usr/bin/env python3

import os

if not os.path.isfile("doom-cannon"):
    print("❌ doom_cannon not found!")
else:
    print("▶ Running DooM-CANNON...")
    os.system("python3 doom-cannon")
