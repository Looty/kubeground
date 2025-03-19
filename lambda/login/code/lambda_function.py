"""Lambda module creating customer virtualplatform in repo"""

import os
import shutil
import datetime
import time
import glob
import yaml

from git import Repo

LOCAL_DIR = "/tmp/clone1"
GITUHB_USER = "Lambda"
GITHUB_USER_MAIL = "lambda@gmail.com"
GITHUB_MAIN_BRANCH = "main"

def lambda_handler(event, context):
    """Lambda module creating customer virtualplatform in repo"""

    print("Event: ", event)

    github_token = os.environ['GITHUB_TOKEN']
    repo_url = f"https://{github_token}:x-oauth-basic@github.com/Looty/kubeground-config"

    print("[INFO] Remove the existing directory if it exists..")
    if os.path.exists(LOCAL_DIR):
        shutil.rmtree(LOCAL_DIR)

    print("[INFO] Cloning repo..")
    repo = Repo.clone_from(repo_url, LOCAL_DIR)

    email = event["request"]["userAttributes"]["email"]
    user_id = event["request"]["userAttributes"]["sub"]
    username = email.split('@')[0]
    username_sanitized = username.replace(".", "-")

    creation_date = datetime.datetime.now()
    timestamp = int(datetime.datetime.now().timestamp())
    customer_branch = f"new-customer/{username_sanitized}-{timestamp}"

    print("[INFO] Setting committer data..")
    repo.config_writer().set_value("user", "name", GITUHB_USER).release()
    repo.config_writer().set_value("user", "email", GITHUB_USER_MAIL).release()

    print("[INFO] Creating new customer branch..")
    repo.git.checkout("-b", customer_branch)

    clean_virtual_platforms(repo, customer_branch, email)
    generate_virtual_platform(repo, customer_branch, username, username_sanitized, user_id, email, creation_date, timestamp)

    return event

def generate_virtual_platform(repo, source_branch, username, username_sanitized, user_id, email, creation_date, timestamp):
    print(f"[INFO] Creating new user: {email}")

    virtual_platform_commit = f"""
apiVersion: looty.example.org/v1alpha1
kind: VirtualPlatform
metadata:
  name: '{username_sanitized}'
  namespace: virtual-platform
  labels:
    customer.name: '{username}'
    customer.id: '{user_id}'
  annotations:
    customer.email: '{email}'
    customer.creation: '{creation_date}'
    customer.timestamp: '{timestamp}'
spec:
  vcluster:
    values: {{}}
  platform:
    values: {{}}
""".lstrip()

    print("[INFO] Generating customer yml config..")
    print(virtual_platform_commit)

    print("[INFO] Creating new customer yml..")
    update_file = f"virtualplatforms/{username_sanitized}-{timestamp}.yml"
    full_path = f"{LOCAL_DIR}/{update_file}"
    with open(full_path, "w", encoding="utf-8") as f:
      f.write(virtual_platform_commit)

    print("[INFO] Adding config to staging..")
    repo.index.add([update_file])

    if repo.index.diff(repo.head.commit):
      commit_message=f"Added new customer config for {email} at {datetime.datetime.now()}"
      print(f"[INFO] Merging new customer: {commit_message}")

      repo.index.commit(commit_message)
      merge(repo, source_branch, GITHUB_MAIN_BRANCH)
    else:
      print(f"[INFO] No need adding anything as nothing has changed, skipping..")

def clean_virtual_platforms(repo, source_branch, email):
  print("[INFO] Cleaning duplicate vclusters..")
  virtual_platforms = glob.glob(f"{LOCAL_DIR}/virtualplatforms/*.yml")

  for file in virtual_platforms:
    with open(file, 'r') as stream:
        try:
            file_data = yaml.safe_load(stream)
            file_email = file_data['metadata']['annotations']['customer.email']

            if email == file_email:
              print(f"[INFO] Found existing customer file: {file}")
              repo.git.rm(file)

        except yaml.YAMLError as exc:
            print(f"[INFO] Skipping deleting {file}")

  if repo.index.diff(repo.head.commit):
    commit_message=f"Removed existing customer config for {email} at {datetime.datetime.now()}"
    print(f"[INFO] Merging cleaning: {commit_message}")

    repo.index.commit(commit_message)
    merge(repo, source_branch, GITHUB_MAIN_BRANCH)
  else:
    print(f"[INFO] No need cleaning as nothing has changed, skipping..")

def merge(repo, source_branch, dest_branch):
  repo.git.checkout(dest_branch)
  repo.git.merge(source_branch)
  origin = repo.remote(name='origin')
  origin.push()
