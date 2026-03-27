import os
import getpass
import keyring
from google.cloud import secretmanager
from typing import Optional, Protocol

class SecretManager(Protocol):
    def access_secret_version(
        self,
        secret_id: str,
        version_id: str = "latest",
        encoding: Optional[str] = "utf-8",
    ) -> str:
        ...

class GCPSecretManager:
    def access_secret_version(
        self,
        secret_id: str,
        version_id: str = "latest",
        encoding: Optional[str] = "utf-8",
    ) -> str:
        project_id = os.getenv("GOOGLE_CLOUD_PROJECT") or os.getenv("PROJECT_ID")
        if (project_id is None) or (len(project_id) == 0):
            raise ValueError("Project ID is missing. Set GOOGLE_CLOUD_PROJECT or PROJECT_ID environment variable.")

        client = secretmanager.SecretManagerServiceClient()
        name = f"projects/{project_id}/secrets/{secret_id}/versions/{version_id}"
        response = client.access_secret_version(request={"name": name})
        return response.payload.data.decode(encoding)

class LocalKeyringManager:
    def access_secret_version(
        self,
        secret_id: str,
        version_id: str = "latest",
        encoding: Optional[str] = "utf-8",
    ) -> str:
        username = (
            os.getenv("LOCAL_KEYRING_USERNAME")
            or getpass.getuser()
        )

        ret = keyring.get_password(secret_id, username)
        if ret is None:
            raise Exception("Secret not found in keyring")
        return ret


def access_secret(
    secret_id: str,
    version_id: str = "latest",
    manager: Optional[SecretManager] = None,
    encoding: Optional[str] = "utf-8",
) -> str:
    """
    Accesses a secret from a specified secret manager.

    Args:
        secret_id (str): The ID of the secret.
        version_id (str, optional): The version of the secret (for GCS).
            Defaults to `"latest"`.
        source_id (str, optional): The source identifier. For GCS, this is the
            project ID (defaults to `GOOGLE_CLOUD_PROJECT_NUM` env var). For
            local keyring, this is the username (defaults to `LOCAL_KEYRING_USERNAME`
            env var or the current user).
        manager (SecretManager, optional): The secret manager instance to use.
            Defaults to `GCPSecretManager`.
        encoding (str, optional): The encoding for the secret. Defaults to "utf-8".

    Returns:
        str: The secret payload.
    """
    if manager is None:
        manager = GCPSecretManager()

    return manager.access_secret_version(
        secret_id=secret_id,
        version_id=version_id,
        encoding=encoding
    )