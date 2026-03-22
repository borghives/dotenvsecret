# dotenvsecret

`dotenvsecret` is a Python utility that reads key-name pairs from a `.envsecret` file, retrieves the actual secret values from a secret manager (like Google Cloud Secret Manager or a local Keyring), and sets them as environment variables in your application. It serves as a secure alternative or addition to `python-dotenv`.

## Installation

Install using standard python packaging tools:

```bash
pip install dotenvsecret
```

Or using `uv` / `poetry`:

```bash
uv add dotenvsecret
```

## Setup

Create a `.envsecret` file in your project root, mapping your environment variables to secret names.

```dotenv
# .envsecret Example
DB_PASSWORD=my-database-password-secret-id
API_KEY=my-api-key-secret-id
```

## Usage

You can load secrets automatically when your application starts:

```python
import os
from dotenvsecret import load_dotenvsecret

# This will fetch secrets and set them as environment variables
load_dotenvsecret()

print(os.getenv("DB_PASSWORD"))
print(os.getenv("API_KEY"))
```

### Unloading
If you need to remove the secrets from your environment later, you can use:

```python
from dotenvsecret import unload_dotenvsecret

unload_dotenvsecret()
```

## Secret Managers

`dotenvsecret` currently supports two secret managers:
- `SecretManager.GCP_SECRET_MANAGER` (Google Cloud Secret Manager) - default
- `SecretManager.LOCAL_KEYRING` (Local Keyring)

To use a different manager, pass it to `load_dotenvsecret`:

```python
from dotenvsecret import load_dotenvsecret, SecretManager

load_dotenvsecret(manager=SecretManager.LOCAL_KEYRING)
```

### Google Cloud Secret Manager Configuration

If you are using Google Cloud, ensure you have set one of the following environment variables to identify your project:
- `GOOGLE_CLOUD_PROJECT`
- `PROJECT_ID`
- `GOOGLE_CLOUD_PROJECT_NUM`

## Disabling Loading

Like `python-dotenv`, you can disable loading `.envsecret` by setting the environment variable `DOTENVSECRET_DISABLED` to any truthy value (`1`, `true`, `yes`).
