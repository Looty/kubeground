"""Lambda module cleaning old customer virtualplatform in repo"""

import os
import shutil
import datetime
from datetime import timedelta
import glob
import yaml

from git import Repo

LOCAL_DIR = "/tmp/clone2"
GITUHB_USER = "Lambda"
GITHUB_USER_MAIL = "lambda@gmail.com"
GITHUB_MAIN_BRANCH = "main"
GITHUB_TOKEN = os.environ['GITHUB_TOKEN']

TRIGGER_FORMAT = os.environ['TRIGGER_FORMAT']
TRIGGER_TIME = os.environ['TRIGGER_TIME']
DELETE_CUSTOMER_AFTER = os.environ['DELETE_CUSTOMER_AFTER']

def lambda_handler(event, context):
    """Lambda module cleaning old customer virtualplatform in repo"""

    print("Event: ", event)
    print(f"[INFO] This event triggers every {TRIGGER_TIME} {TRIGGER_FORMAT}.")

    repo_url = f"https://{GITHUB_TOKEN}:x-oauth-basic@github.com/Looty/kubeground-config"

    print("[INFO] Remove the existing directory if it exists..")
    if os.path.exists(LOCAL_DIR):
        shutil.rmtree(LOCAL_DIR)

    print("[INFO] Cloning repo..")
    repo = Repo.clone_from(repo_url, LOCAL_DIR)
    cleaning_id_prefix = event["id"].split('-')[0]
    cleaning_branch = f"cleaning/{cleaning_id_prefix}"
    current_timestamp = int(datetime.datetime.now().timestamp())

    print("[INFO] Setting committer data..")
    repo.config_writer().set_value("user", "name", GITUHB_USER).release()
    repo.config_writer().set_value("user", "email", GITHUB_USER_MAIL).release()

    print("[INFO] Creating new customer branch..")
    repo.git.checkout("-b", cleaning_branch)

    clean_old_virtual_platforms(repo, cleaning_branch, current_timestamp)

    return event

def clean_old_virtual_platforms(repo, source_branch, current_timestamp):
  print("[INFO] Cleaning old vclusters..")

  virtual_platforms = glob.glob(f"{LOCAL_DIR}/virtualplatforms/*.yml")
  for file in virtual_platforms:
    with open(file, 'r') as stream:
        try:
            file_data = yaml.safe_load(stream)
            email = file_data['metadata']['annotations']['customer.email']
            id = file_data['metadata']['labels']['customer.id']
            creation_timestamp = int(file_data['metadata']['annotations']['customer.timestamp'])
            creation = file_data['metadata']['annotations']['customer.creation']
            diff = current_timestamp - creation_timestamp

            if diff >= (60 * int(DELETE_CUSTOMER_AFTER)): # in minutes
              print(f"[INFO] Found old customer: {id} / {email} / {creation}")
              print(f"[INFO] Found old customer, timestamp diff - {diff}")
              repo.git.rm(file)

        except yaml.YAMLError as exc:
            print(f"[INFO] Skipping deleting {file}")

  if repo.index.diff(repo.head.commit):
    commit_message=f"Cleaning old customer config for {email} at {current_timestamp}"
    print(f"[INFO] Merging cleaning old customers: {commit_message}")

    repo.index.commit(commit_message)
    merge(repo, source_branch, GITHUB_MAIN_BRANCH)
  else:
    print(f"[INFO] No need cleaning as there are only new customers, skipping..")

def merge(repo, source_branch, dest_branch):
  repo.git.checkout(dest_branch)
  repo.git.merge(source_branch)
  origin = repo.remote(name='origin')
  origin.push()
