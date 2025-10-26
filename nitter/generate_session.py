#!/usr/bin/env python3
import re
import requests
import json
import sys
import pyotp
import os

# NOTE: pyotp and requests are dependencies
# > pip install pyotp requests

TW_CONSUMER_KEY = '3nVuSoBZnx6U4vzUxf5w'
TW_CONSUMER_SECRET = 'Bcs59EFbbsdF6Sl9Ng71smgStWEGwXXKSjYvPVt7qys'

def auth(username, password, otp_secret):
    print("Starting authentication process...")
    print(f"Username: {username}")

    print("Requesting bearer token...")
    bearer_token_req = requests.post("https://api.twitter.com/oauth2/token",
        auth=(TW_CONSUMER_KEY, TW_CONSUMER_SECRET),
        headers={"Content-Type": "application/x-www-form-urlencoded"},
        data='grant_type=client_credentials'
    ).json()
    bearer_token = ' '.join(str(x) for x in bearer_token_req.values())
    print(f"Bearer token obtained: {bearer_token[:50]}...")

    print("Requesting guest token...")
    guest_token = requests.post(
        "https://api.twitter.com/1.1/guest/activate.json",
        headers={'Authorization': bearer_token}
    ).json().get('guest_token')

    if not guest_token:
        print("Failed to obtain guest token.")
        sys.exit(1)

    print(f"Guest token obtained: {guest_token}")

    twitter_header = {
        'Authorization': bearer_token,
        "Content-Type": "application/json",
        "User-Agent": "TwitterAndroid/10.21.0-release.0 (310210000-r-0) ONEPLUS+A3010/9 (OnePlus;ONEPLUS+A3010;OnePlus;OnePlus3;0;;1;2016)",
        "X-Twitter-API-Version": '5',
        "X-Twitter-Client": "TwitterAndroid",
        "X-Twitter-Client-Version": "10.21.0-release.0",
        "OS-Version": "28",
        "System-User-Agent": "Dalvik/2.1.0 (Linux; U; Android 9; ONEPLUS A3010 Build/PKQ1.181203.001)",
        "X-Twitter-Active-User": "yes",
        "X-Guest-Token": guest_token,
        "X-Twitter-Client-DeviceID": ""
    }

    session = requests.Session()
    session.headers = twitter_header

    print("Starting login flow (Task 1)...")
    task1 = session.post(
        'https://api.twitter.com/1.1/onboarding/task.json',
        params={
            'flow_name': 'login',
            'api_version': '1',
            'known_device_token': '',
            'sim_country_code': 'us'
        },
        json={
            "flow_token": None,
            "input_flow_data": {
                "country_code": None,
                "flow_context": {
                    "referrer_context": {
                        "referral_details": "utm_source=google-play&utm_medium=organic",
                        "referrer_url": ""
                    },
                    "start_location": {
                        "location": "deeplink"
                    }
                },
                "requested_variant": None,
                "target_user_id": 0
            }
        }
    )

    session.headers['att'] = task1.headers.get('att')
    print(f"Task 1 completed. Flow token: {task1.json().get('flow_token')[:50]}...")

    print("Submitting username (Task 2)...")
    task2 = session.post(
        'https://api.twitter.com/1.1/onboarding/task.json',
        json={
            "flow_token": task1.json().get('flow_token'),
            "subtask_inputs": [{
                "enter_text": {
                    "suggestion_id": None,
                    "text": username,
                    "link": "next_link"
                },
                "subtask_id": "LoginEnterUserIdentifier"
            }]
        }
    ).json()

    print(f"Task 2 completed. Subtask ID: {task2['subtasks'][0]['subtask_id']}")

    if task2['subtasks'][0]['subtask_id'] == "LoginEnterAlternateIdentifierSubtask":
        print("Unusual login activity detected!")
        alt_identifier = input("Unusual login activity detected, enter the account's email address: ")

        task_enter_email = session.post(
            'https://api.twitter.com/1.1/onboarding/task.json',
            json={
                "flow_token": task2.get('flow_token'),
                "subtask_inputs": [{
                    "enter_text": {
                        "suggestion_id": None,
                        "text": alt_identifier,
                        "link": "next_link"
                    },
                    "subtask_id": "LoginEnterAlternateIdentifierSubtask"
                }]
            }
        ).json()
        print(f"Email verification response: {task_enter_email}")

    print("Submitting password (Task 3)...")
    task3 = session.post(
        'https://api.twitter.com/1.1/onboarding/task.json',
        json={
            "flow_token": task2.get('flow_token'),
            "subtask_inputs": [{
                "enter_password": {
                    "password": password,
                    "link": "next_link"
                },
                "subtask_id": "LoginEnterPassword"
            }],
        }
    ).json()

    print(f"Task 3 completed. Subtasks: {[st.get('subtask_id') for st in task3.get('subtasks', [])]}")

    for t3_subtask in task3.get('subtasks', []):
        if "open_account" in t3_subtask:
            print("Login successful! Account opened.")
            return t3_subtask["open_account"]
        elif t3_subtask.get("subtask_id") == "LoginAcid":
            # Prompt the user to enter the confirmation code sent to their email.
            print("Email confirmation required (LoginAcid)")
            confirmation_code = input("A confirmation code has been sent to your email. Enter it: ")
            task4 = session.post(
                "https://api.twitter.com/1.1/onboarding/task.json",
                json={
                    "flow_token": task3.get("flow_token"),
                    "subtask_inputs": [{
                        "enter_text": {
                            "text": confirmation_code,
                            "link": "next_link"
                        },
                        "subtask_id": "LoginAcid"
                    }]
                }
            ).json()
            print(f"Task 4 (Email confirmation) completed.")
            for t4_subtask in task4.get("subtasks", []):
                if "open_account" in t4_subtask:
                    print("Login successful after email confirmation!")
                    return t4_subtask["open_account"]
        elif "enter_text" in t3_subtask:
            response_text = t3_subtask["enter_text"]["hint_text"]
            print(f"2FA required. Hint: {response_text}")
            totp = pyotp.TOTP(otp_secret)
            generated_code = totp.now()
            print(f"Generated 2FA code: {generated_code}")
            task4 = session.post(
                "https://api.twitter.com/1.1/onboarding/task.json",
                json={
                    "flow_token": task3.get("flow_token"),
                    "subtask_inputs": [
                        {
                            "enter_text": {
                                "suggestion_id": None,
                                "text": generated_code,
                                "link": "next_link",
                            },
                            "subtask_id": "LoginTwoFactorAuthChallenge",
                        }
                    ],
                }
            ).json()
            print(f"Task 4 (2FA) completed.")
            for t4_subtask in task4.get("subtasks", []):
                if "open_account" in t4_subtask:
                    print("Login successful after 2FA!")
                    return t4_subtask["open_account"]

    print("Authentication flow completed but no open_account found.")
    return None

import os

if __name__ == "__main__":
    print("=== Twitter Session Generator ===")
    username = os.environ.get("TWITTER_USERNAME")
    password = os.environ.get("TWITTER_PASSWORD")
    otp_secret = os.environ.get("TWITTER_OTP_SECRET")
    path = sys.argv[1]

    print(f"Output path: {path}")
    print(f"Username configured: {bool(username)}")
    print(f"Password configured: {bool(password)}")
    print(f"OTP Secret configured: {bool(otp_secret)}")
    print()

    try:
        result = auth(username, password, otp_secret)
    except Exception as e:
        result = None

    if result is None:
        print("Authentication failed.")

        # read file
        with open(path, "r") as f:
            lines = f.readlines()
            if len(lines) > 0:
                print(f"Sessions loaded: {len(lines)}")
                sys.exit(0)
            else:
                print("No sessions found.")
                sys.exit(1)

    print(f"Auth result: {result}")

    session_entry = {
        "oauth_token": result.get("oauth_token"),
        "oauth_token_secret": result.get("oauth_token_secret")
    }

    print(f"Session entry created: oauth_token={bool(session_entry['oauth_token'])}, oauth_token_secret={bool(session_entry['oauth_token_secret'])}")

    try:
        with open(path, "a") as f:
            f.write(json.dumps(session_entry) + "\n")
        print("Authentication successful. Session appended to", path)
    except Exception as e:
        print(f"Failed to write session information: {e}")
        sys.exit(1)
