"""Lambda module registering/cleaning users from uploaded .csv in S3"""

import os
import shutil
import boto3
import csv

s3_client = boto3.client('s3')
cognito_client = boto3.client('cognito-idp')

USER_POOL_ID = os.environ['USER_POOL_ID']
DEFAULT_PASSWORD = os.environ['DEFAULT_PASSWORD']
CSV_PATH = os.environ['CSV_PATH']

USERS_TO_CREATE = 0
USERS_TO_DELETE = 0

def lambda_handler(event, context):
    """Lambda module registering/cleaning users from uploaded .csv in S3"""

    print("[INFO] Event: ", event)

    bucket = event['Records'][0]['s3']['bucket']['name']
    key = event['Records'][0]['s3']['object']['key']
    download_path = '/tmp/customers'
    
    # Check if the uploaded file is the correct one
    if key != CSV_PATH:
        return
    
    print("[INFO] Removing download directory if exists..")
    if os.path.exists(download_path):
        shutil.rmtree(download_path)

    print("[INFO] Ensuring download directory exists..")
    os.mkdir(download_path)

    print("[INFO] Download CSV file from S3..")
    download_file_path = f"{download_path}/approved_emails.csv"
    s3_client.download_file(bucket, key, download_file_path)

    print("[INFO] Processing CSV file..")
    with open(download_file_path, 'r') as file:
        reader = csv.DictReader(file)
        approved_emails = [row['email'] for row in reader]

    print("[INFO] Fetching the list of current users in the Cognito User Pool..")
    current_users = get_current_users()

    print("[INFO] Creating new users..")
    for email in approved_emails:
        if email not in current_users:
            create_user(email)

    print("[INFO] Removing old users..")
    for email in current_users:
        if email not in approved_emails:
            delete_user(email)

    print(f"[INFO] Created {USERS_TO_CREATE} users, removed {USERS_TO_DELETE}..")
    print("[INFO] User synchronization completed.")

    return event

def get_current_users():
    current_users = []
    paginator = cognito_client.get_paginator('list_users')
    for response in paginator.paginate(UserPoolId=USER_POOL_ID):
        for user in response['Users']:
            for attribute in user['Attributes']:
                if attribute['Name'] == 'email':
                    current_users.append(attribute['Value'])
    return current_users

def create_user(email):
    try:
        global USERS_TO_CREATE
        USERS_TO_CREATE += 1

        print("[INFO] Creating new user: ", email)
        cognito_client.admin_create_user(
            UserPoolId=USER_POOL_ID,
            Username=email,
            UserAttributes=[
                {'Name': 'email', 'Value': email},
                {'Name': 'email_verified', 'Value': 'true'}
            ],
            TemporaryPassword=DEFAULT_PASSWORD,
            MessageAction='SUPPRESS' # Suppress sending the welcome email
        )
        cognito_client.admin_set_user_password(
            UserPoolId=USER_POOL_ID,
            Username=email,
            Password=DEFAULT_PASSWORD,
            Permanent=True
        )
    except cognito_client.exceptions.UsernameExistsException:
        pass

def delete_user(email):
    try:
        global USERS_TO_DELETE
        USERS_TO_DELETE += 1

        print("[INFO] Removing existing user: ", email)
        cognito_client.admin_delete_user(
            UserPoolId=USER_POOL_ID,
            Username=email
        )
    except cognito_client.exceptions.UserNotFoundException:
        pass
