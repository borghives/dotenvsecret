import os
import getpass
import keyring
from google.cloud import secretmanager, resourcemanager_v3
from enum import Enum
from typing import Optional
class SecretManager(Enum):
    GCP_SECRET_MANAGER = "gcp_secret_manager"
    LOCAL_KEYRING = "keyring"

def get_project_number(project_id: str) -> str:
    client = resourcemanager_v3.ProjectsClient()
    name = f"projects/{project_id}"
    project = client.get_project(request={"name": name})
    # The name returned is 'projects/123456789', so we split it
    project_number = project.name.split('/')[-1]
    return project_number

def access_secret(
    secret_id: str,
    version_id: str = "latest",
    source_id: Optional[str] = None,
    manager: Optional[SecretManager] = SecretManager.GCP_SECRET_MANAGER,
    encoding: Optional[str] = "utf-8",
) -> str:
    """
    Accesses a secret from a specified secret manager.

    Args:
        secret_id (str): The ID of the secret.
        source_id (str, optional): The source identifier. For GCS, this is the
            project ID (defaults to `GOOGLE_CLOUD_PROJECT_NUM` env var). For
            local keyring, this is the username (defaults to `LOCAL_KEYRING_USERNAME`
            env var or the current user).
        version_id (str, optional): The version of the secret (for GCS).
            Defaults to `"latest"`.
        manager (SecretManager, optional): The secret manager to use.
            Defaults to `SecretManager.GCS_SECRET_MANAGER`.

    Returns:
        str: The secret payload.
    """

    match manager:
        case SecretManager.GCP_SECRET_MANAGER:
            project_num = os.getenv("GOOGLE_CLOUD_PROJECT_NUM")
            if (project_num is None) or (len(project_num) == 0):
                project_id = os.getenv("GOOGLE_CLOUD_PROJECT") or os.getenv("PROJECT_ID")
                if (project_id is None) or (len(project_id) == 0):
                    raise ValueError("Project ID is missing. Set GOOGLE_CLOUD_PROJECT or PROJECT_ID environment variable.")

                project_num = get_project_number(project_id)

            client = secretmanager.SecretManagerServiceClient()
            name = f"projects/{project_num}/secrets/{secret_id}/versions/{version_id}"
            response = client.access_secret_version(request={"name": name})
            return response.payload.data.decode(encoding)

        case SecretManager.LOCAL_KEYRING:
            username = (
                source_id
                or os.getenv("LOCAL_KEYRING_USERNAME")
                or getpass.getuser()
            )

            ret = keyring.get_password(secret_id, username)
            if ret is None:
                raise Exception("Secret not found in keyring")
            return ret