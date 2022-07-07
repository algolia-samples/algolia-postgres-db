# Algolia Sample application - Python back end

## Requirements

- Python, at least version 3.6
- [Configured .env file](../../README.md)

## How to run

1. Create and activate a new [virtual environment](https://docs.python.org/fr/3/library/venv.html).

**MacOS/Unix**

```
python3 -m venv env
source env/bin/activate
```

**Windows (PowerShell)**

```
python3 -m venv env
.\env\Scripts\activate.bat
```

2. Install dependencies

```
pip install -r requirements.txt
```

3. Run the application

**MacOS/Unix**

```
export FLASK_APP=server.py
python3 -m flask run --port=4242
```

**Windows (PowerShell)**

```
$env:FLASK_APP=“server.py"
python3 -m flask run --port=4242
```

4. Go to [localhost:4242](http://localhost:4242)