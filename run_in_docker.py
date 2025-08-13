#!/usr/bin/env python
import subprocess
import pathlib
import argparse
import random

# By appending "|| true" execution is allowed to continue even when 0 isn't returned.

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Run or update a docker container for a bot"
    )

    parser.add_argument(
        "--mount", type=pathlib.Path, metavar='m', required=True,
        help="An absolute path to a folder that contains configuration"
        " and storage to be used by the bot. This becomes the argument that"
        " is passed to the bot when running it.")

    args = parser.parse_args()
    mount_src = args.mount.absolute()
    mount_dest = "/config"
    project_path = pathlib.Path(pathlib.Path(__file__).parent).absolute()
    env_file = str(mount_src / ".env")

    print("Deleting")

    subprocess.run(
        args=["docker", "compose", "--env-file", env_file, "down", "--rmi", "all"],
        check=True,
        cwd=project_path
    )

    print("Running")

    subprocess.run(
        args=["docker", "compose", "--env-file", env_file, "up", "-d"],
        check=True,
        cwd=project_path
    )