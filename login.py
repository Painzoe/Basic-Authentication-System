import sqlite3
import hashlib
import sys
import secrets
import hmac

database = sqlite3.connect('database.db')
cursor = database.cursor()
cursor.execute("CREATE TABLE IF NOT EXISTS 'authentication information'(id INTEGER PRIMARY KEY, username TEXT UNIQUE, salt BLOB, password_hash BLOB)") #BLOB(Binary Large Object)

if len(sys.argv) != 4:
    print("""Usage: python3 login.py <command> <username> <password>  
             command options: add, login""", file=sys.stderr)
    sys.exit(1)

# sys.argv[0] = login.py
command = sys.argv[1]
username = sys.argv[2]
password = sys.argv[3]

if command ==  "add":
    salt = secrets.token_bytes(16)
    password_bytes = password.encode("utf-8")
    key = hashlib.pbkdf2_hmac("sha256", password_bytes, salt, 100000)
    cursor.execute("INSERT OR IGNORE INTO 'authentication information' VALUES(?, ?, ?, ?)",(None, username, salt, key))
    database.commit()

elif command == "login":
    cursor.execute("SELECT salt, password_hash FROM 'authentication information' WHERE username = ?",(username,))
    data = cursor.fetchall() #tupple [(salt, password_hash), (...), ...]
    password_bytes = password.encode("utf-8")
    try:
        control = hashlib.pbkdf2_hmac("sha256", password_bytes, data[0][0], 100000)
    except IndexError:
        print("There is no such user", file=sys.stderr)
        sys.exit(1)
    if not hmac.compare_digest(control, data[0][1]): # Use constant-time comparison to prevent timing attacks instead of control!=data[0][1]
        print("Incorrect password", file=sys.stderr, flush=True) # flush-> Don't wait in the buffer, write immediately
        sys.exit(1)
    else:
        print("Login successful")
        sys.exit(0)
else:
    print("Unknown command", file=sys.stderr)
    sys.exit(1)









