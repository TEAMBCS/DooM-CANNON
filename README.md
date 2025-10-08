<p align="center">
  <img src="https://i.postimg.cc/j2Rx23zp/doom-cannon.jpg" width="600"/>
</p>

<h1 align="center"> DOOM CANNON DDoS Tool</h1>

<p align="center">
  🚀 Powered by <b> BANGLADESH CYBER SQUAD and TEAM SHADOW STRIKER </b><br>
  📆 Year: 2025
</p>


---

## 📖 Overview
**DOOM CANNON** is a hybrid CLI toolkit built with **Python** 🐍 and **Go** 🌀.  
It provides a centralized menu system (Python Rich + PyFiglet) to launch multiple Go-powered modules.  
Each module has customizable options (target, port, threads, duration, proxy, headers, etc.).  

⚠️ **Disclaimer**: This tool is for **educational and research purposes only**.  
The authors take no responsibility for misuse.  

---

## ✨ Features
- Interactive CLI with styled menus
- Python-based launcher with Rich UI
- Multiple Go-based modules for different modes
- Custom input support (target, threads, proxy, headers, wordlists)
- Modular structure – easy to add new Go tools

---
<h1 align="center"> DEVELOPER INFO </h1>

## 👨‍💻 Developers

* **BLACK ZERO**
* **FULL MOON**
* **MR. CODE ERROR**

## 🧑‍💻 HELPED BY 
* **TAUSIF ZAMAN**
* **PAEVES JOY**
* **Kazi Tanvir Mahmud Omi**
* **SHAWON ISLAM SAIM**

## INSPIRED By
* **PARVIS HASAN**

---
<h1 align="center"> Project Languages </h1>



<p align="center">
  <img src="https://i.postimg.cc/FF6y71Ds/python.jpg" alt="Python Logo" width="120"/>
  &nbsp;&nbsp;&nbsp;
  <img src="https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Blue.png" alt="Go Logo" width="140"/>
  &nbsp;&nbsp;&nbsp;
  <img src="https://i.postimg.cc/9XwpPG4q/shell.png" alt="Go Logo" width="140"/>
</p>

---

<h1 align="center"> Installing info </h1>


## 🛠️ Installation for Termux
```bash
pkg update && pkg upgrade
pkg install coreutils -y
pkg install grep -y
pkg install libcaca -y
pkg install python -y
pkg install python3 -y
pkg install golang -y
pkg install git
rm -rf DooM-CANNON 
git clone https://github.com/TEAMBCS/DooM-CANNON.git
cd DooM-CANNON
chmod 777 *
chmod +x *
pip3 install -r requirements.txt
python3 setup.py
```
## 🔌 Installation for Kali Linux 
```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y coreutils grep libcaca python3 python3-pip python3-dev python3-setuptools build-essential golang git
pip3 install --upgrade pip
git clone https://github.com/TEAMBCS/DooM-CANNON.git
cd DooM-CANNON
chmod 777 *
chmod +x *
pip3 install -r requirements.txt
python3 setup-kali.py
```
---

## ▶️ Usage
   *PYTHON language  launcher*
```bash
python3 doom-cannon.py
```
   *GO language launcher*
```bash
go run doom-cannon.go
```
   *SHELL language launcher*
```bash
bash doom-cannon.sh
```    
You will see a CLI menu with available modules.
Select an option (e.g., `01/A` for ORBIT Attack) and follow the prompts.

---
<h1 align="center"> Tool Login Info </h1>

## TOOL USERNAME 
```bash
*as you wish* user name is not fix use random name.
```
## TOOL PASSWOED 
```bash
DOOM CANNON@TEAM BCS
```
---
<h1 align="center"> Important Note </h1>

## HTTP Headers Example

These headers are used when sending requests to a target.  
**Important:** Change the values according to your target server.

```python
headers = {
    Host: victim.com
Origin: https://victim.com
Referer: https://www.google.com/
X-Forwarded-For: 45.76.89.120
X-Forwarded-Host: victim.com
X-Forwarded-Proto: https
X-Real-IP: 203.23.101.55
Client-IP: 149.56.210.87
Forwarded: for=185.12.44.201;proto=https;by=198.51.100.200
CF-Connecting-IP: 91.132.137.45
True-Client-IP: 64.233.160.2
X-Originating-IP: 212.102.44.98
X-Cluster-Client-IP: 103.21.244.15
X-Requested-With: XMLHttpRequest
X-Forwarded-Server: edge-proxy-1
X-Request-ID: 123456789
Via: 2.0 proxy
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
Accept-Language: en-US,en;q=0.8
Accept-Encoding: gzip, deflate, br
Cache-Control: no-cache
Pragma: no-cache
Connection: keep-alive
TE: trailers
Upgrade-Insecure-Requests: 1
DNT: 1
    # Add or modify other headers as needed
}
```
---
<h1 align="center"> PROJECT INFO </h1>

## 🌳 Project Structure (Detailed Tree)

```
DooM-CANNON/
│
├── doom_cannon.py               # Main Python language  launcher (menu, UI, input handling)
├── doom_cannon.go  							#  Go language launcher (menu, UI, input handling)
├── doom_cannon.sh 							#  Shell language launcher (menu, UI, input handling)
│
├── orbit.go                 # Orbit Attack (TLS-based)
│   ├─ Inputs: host, port, method [GET/POST], threads, duration, debug, proxy, header
│
├── axis.go                  # Axis Attack (TLS-based)
│   ├─ Inputs: url, port, method [GET/POST], threads, duration, debug, proxy
│
├── viod.go                  # Viod Attack
│   ├─ Inputs: url, port, threads, method, duration, proxy, wordlist, header
│
├── noise.go                 # Noise Attack
│   ├─ Inputs: url, threads, method, duration, proxy, header
│
├── ghost.go                 # Ghost Attack
│   ├─ Inputs: url, port, threads, method, duration, proxy, wordlist
│
├── xiron.go                 # Xiron Attack
│   ├─ Inputs: url, threads, method, duration, proxy
│
├── orix.go                  # Orix Attack
│   ├─ Inputs: url, port, threads, method, duration, proxy
│
├── orrin.go                 # Orrin Attack
│   ├─ Inputs: site, safe_mode [y/n]
│
├── viont.go                 # Viont Attack (Flood mode)
│   ├─ Inputs: url, port, threads, method, duration, header
│
│
├── proxy.txt                # Optional proxy list
├── header.txt               # Optional custom headers
└── wordlist                 # Optional wordlist for specific modules
```

---
## DOOM CANNON UI 
<p align="center">
  <img src="https://i.postimg.cc/sg05WqgC/doom-ui.jpg" width="600"/></p>

  ---
## Prove
<p align="center">
  <img src="https://i.postimg.cc/RFB9ppbW/prove.jpg" width="800"/>
</p>


---



## 📜 License

MIT License @ copyright 2025

---

