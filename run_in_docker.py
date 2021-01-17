#!/usr/bin/env python
import subprocess
import pathlib
import argparse

# By appending "|| true" execution is allowed to continue even when 0 isn't returned.

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Run or update a docker container for a bot"
    )

    parser.add_argument(
        "--name", metavar='N', type=str, required=True,
        help="An identifier for a bot that will be used extensively in docker."
             " If this identifier already exists in docker (including using"
             " this script previously) it will be replaced.")

    parser.add_argument(
        "--mount", type=pathlib.Path, metavar='m', required=True,
        help="An absolute path to a folder that contains configuration"
        " and storage to be used by the bot. This becomes the argument that"
        " is passed to the bot when running it.")

    args = parser.parse_args()
    name = args.name
    tag = "latest"
    mount_src = args.mount.absolute()
    mount_dest = "/config"
    project_path = pathlib.Path(pathlib.Path(__file__).parent).absolute()

    subprocess.run(
        args=["docker", "stop", name],
        check=False
    )

    subprocess.run(
        args=["docker", "rm", name],
        check=False
    )

    subprocess.run(
        args=["docker", "rmi", "{}:{}".format(name, tag)],
        check=False
    )

    subprocess.run(
        args=[
            "docker", "build",
            "-t", "{}:{}".format(name, tag),
            "-f", "dockerfile", project_path],
        check=True
    )

    subprocess.run(
        args=[
            "docker", "run",
            "--name", "{}".format(name),
            "--env", "config_path={}".format(str(mount_dest)),
            "--restart", "always",
            "--mount", "type=bind,source={},target={}".format(str(mount_src),
                                                              str(mount_dest)),
            "-d", "{}:{}".format(name, tag)],
        check=True
    )
