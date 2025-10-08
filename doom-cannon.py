#!/usr/bin/env python3

import os

if not os.path.isfile("doom_cannon"):
    print("❌ doom_cannon not found!")
else:
    print("▶ Running doom_cannon...")
    os.system("python3 doom_cannon")
