# 🔐 Simple Python Authentication System

A lightweight Python project demonstrating secure password storage and verification using **SQLite** and **cryptographic hashing**.

---

## Features

- Add new users with **secure password hashing** (`PBKDF2-HMAC-SHA256`)  
- Login verification with **constant-time comparison** to prevent timing attacks  
- Unique salts for each user for enhanced security  
- All data stored locally in **SQLite database**  

---

## Usage

```bash
# Add a new user
python3 login.py add <username> <password>

# Login with existing user
python3 login.py login <username> <password>
``` 
Example;
![windows](example.png) 

---
## 🌐 Version 2 – Web-Based Authentication System (Flask)

This project is the **web-based version** of the earlier **terminal-based authentication system** that worked with `sys.argv` and command-line inputs.

In this version, the authentication logic has been migrated to a **Flask web application** with session management and HTML templates.

### ✨ Features
- User registration and login via web forms
- Secure password hashing using **PBKDF2 (SHA-256 + salt)**
- Session-based authentication using Flask sessions
- Protected dashboard page (login required)
- Logout functionality
- SQLite database
- Deployed on **Render**

### 🧠 Differences from Version 1
| Version 1 (CLI) | Version 2 (Web) |
|----------------|----------------|
Terminal-based (`sys.argv`) | Web-based (Flask) |
Command-line input | HTML forms |
No session handling | Flask sessions |
Local execution only | Publicly accessible |

### 🛠 Technologies Used
- Python
- Flask
- SQLite
- hashlib / secrets / hmac
- Gunicorn
- Render

### 🚀 Live Demo
🔗 **Live URL:**  
https://basic-authentication-system-0kau.onrender.com

### 📌 Notes
- This project is created **for educational purposes**
- SQLite is used; data may reset on redeploy (Render free tier)
- Not intended for production use

---

## 🧪 Previous Version (CLI-Based)

The first version of this project was a **command-line authentication system** where users were managed using terminal commands and `sys.argv`.

Both versions share the same core ideas:
- Secure password storage
- Manual authentication logic
- Focus on learning security fundamentals


