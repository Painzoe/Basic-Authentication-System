import sqlite3
import hashlib
import secrets
import hmac
from flask import Flask, request, redirect, session, render_template

app = Flask(__name__)
app.secret_key = "secret_key06"                 # secrets.token_hex(32) -> for real life using


def get_db():                                # To prevent it from exploding in heavy traffic
    conn = sqlite3.connect("database.db")
    return conn

connect = get_db()
cursor = connect.cursor()
cursor.execute("CREATE TABLE IF NOT EXISTS 'authentication information'(id INTEGER PRIMARY KEY, username TEXT UNIQUE, salt BLOB, password_hash BLOB)") #BLOB(Binary Large Object)
connect.commit()
connect.close()


@app.route("/")
def home():
    return redirect("/login")                # First, the login page is displayed.


@app.route("/register", methods=["GET", "POST"])
def create_user():
    if request.method == "POST":
        username = request.form["username"]   # input
        password = request.form["password"]   # input

        conn = get_db()
        cursor = conn.cursor()
        cursor.execute("SELECT * FROM 'authentication information' WHERE username = ?", (username,))
        if cursor.fetchone():
            return render_template("register.html", error="Username already exists!")

        salt = secrets.token_bytes(16)
        password_bytes = password.encode("utf-8")
        key = hashlib.pbkdf2_hmac("sha256", password_bytes, salt, 100000)

        cursor.execute("INSERT OR IGNORE INTO 'authentication information' VALUES(?, ?, ?, ?)",(None, username, salt, key))
        conn.commit()

        # Created new user
        return redirect("/login")
    # Get request
    return render_template("register.html")


@app.route("/login", methods=["GET", "POST"])
def verify_user():
    if request.method == "POST":
        username = request.form["username"]   # input
        password = request.form["password"]   # input

        conn = get_db()
        cursor = conn.cursor()
        cursor.execute("SELECT salt, password_hash, id FROM 'authentication information' WHERE username = ?",(username,))
        data = cursor.fetchone()             # fetchone -> one row,   fetchall -> all rows

        if data is None:
            return render_template("login.html", error="User not found!")

        user_id = int(data[2])

        password_bytes = password.encode("utf-8")
        control = hashlib.pbkdf2_hmac("sha256", password_bytes, data[0], 100000)

         # incorrect password
        if not hmac.compare_digest(control, data[1]): # Use constant-time comparison to prevent timing attacks instead of control!=data[0][1]
            return render_template("login.html", error="Incorrect password!")

        # Login successful
        session["user_id"] = user_id
        return redirect("/dashboard")

        # Get request
    return render_template("login.html")

@app.route("/dashboard")
def dashboard():
    if "user_id" not in session:
        return redirect("/login")
    return render_template("dashboard.html")

@app.route("/logout")
def logout():
    session.clear()
    return redirect("/login")