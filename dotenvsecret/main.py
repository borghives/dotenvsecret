import logging
import os
from secret import SecretManager
from typing import IO
from typing import Optional
from dotenv import load_dotenv
from secret import access_secret

load_dotenv()


logger = logging.getLogger(__name__)


def _load_dotenvsecret_disabled() -> bool:
    """
    Determine if dotenv loading has been disabled.
    """
    if "DOTENVSECRET_DISABLED" not in os.environ:
        return False
    value = os.environ["DOTENVSECRET_DISABLED"].casefold()
    return value in {"1", "true", "t", "yes", "y"}
    
def load_dotenvsecret(
    dotenvsecret_path: str = ".envsecret",
    manager: SecretManager = SecretManager.GCP_SECRET_MANAGER,
    encoding: Optional[str] = "utf-8",
) -> bool:
    """Parse a .env file and then load all the variables found as environment variables.

    Parameters:
        dotenvsecret_path: Absolute or relative path to .envsecret file.
        stream: Text stream (such as `io.StringIO`) with .envsecret content, used if
            `dotenvsecret_path` is `None`.
    Returns:
        Bool: True if at least one environment variable is set else False

    If both `dotenvsecret_path` and `stream` are `None`, the default path is `.envsecret`.
    If the environment variable `PYTHON_DOTENV_DISABLED` is set to a truthy value,
    .envsecret loading is disabled.
    """

    if _load_dotenvsecret_disabled():
        logger.debug(
            "dotenvsecret: .envsecret loading disabled by DOTENVSECRET_DISABLED environment variable"
        )
        return False
    
    if not os.path.exists(dotenvsecret_path):
        logger.debug(
            "dotenvsecret: .envsecret file not found at %s", dotenvsecret_path
        )
        return False
    
    with open(dotenvsecret_path, "r", encoding=encoding) as f:
        for line in f:
            line = line.strip()
            if not line or line.startswith("#"):
                continue

            if "=" in line:
                env_var, secret_id = line.split("=", 1)
                env_var = env_var.strip()
                secret_id = secret_id.strip()

                # Optional: strip quotes if present
                if secret_id.startswith('"') and secret_id.endswith('"'):
                    secret_id = secret_id[1:-1]
                elif secret_id.startswith("'") and secret_id.endswith("'"):
                    secret_id = secret_id[1:-1]

                try:
                    secret_value = access_secret(
                        secret_id=secret_id,
                        manager=manager,
                        encoding=encoding
                    )
                    os.environ[env_var] = secret_value
                except Exception as e:
                    print(f"Warning: Failed to load secret '{secret_id}' for environment variable '{env_var}': {e}")

def unload_dotenvsecret(
    dotenvsecret_path: str = ".envsecret",
    encoding: Optional[str] = "utf-8",
):
    """Parse a .envsecret file and then remove all the variables found as environment variables.
    """
    if not os.path.exists(dotenvsecret_path):
        logger.debug(
            "dotenvsecret: .envsecret file not found at %s", dotenvsecret_path
        )
        return
    
    with open(dotenvsecret_path, "r", encoding=encoding) as f:
        for line in f:
            line = line.strip()
            if not line or line.startswith("#"):
                continue

            if "=" in line:
                env_var, secret_id = line.split("=", 1)
                env_var = env_var.strip()

                os.environ.pop(env_var, None)
