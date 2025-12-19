import requests
import os

import firebase_admin as fb
import firebase_admin.auth as auth

# init app
CRED_PATH = os.getenv("FIREBASE_KEY_PATH")
if CRED_PATH == None:
    print("Must specify \"FIREBASE_KEY_PATH\"")
    exit()
CREDS = fb.credentials.Certificate(CRED_PATH)
APP = fb.initialize_app(CREDS)

# generate signin link
API_KEY = os.getenv("FIREBASE_API_KEY")
if API_KEY == None:
    print("Must specify \"FIREBASE_API_KEY\"")
    exit()
SIGNIN_URL = f"https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key={API_KEY}"

# TODO: can add custom field here for test
def create_user(uid: str, email: str) -> None:
    auth.create_user(uid=uid, email=email)
    claims = {'test': True}
    auth.set_custom_user_claims(uid=uid, custom_claims=claims)

def get_auth_token(uid: str, email: str) -> str:
    claims = {"email": email}
    token = auth.create_custom_token(uid=uid, developer_claims=claims, app=APP).decode()
    data = {'token': token, 'returnSecureToken': True}
    res = requests.post(SIGNIN_URL, json=data)
    if 'idToken' not in res.json():
        raise ValueError(f"Failed to create Firebase ID token: {res.text}")
    return res.json()['idToken']

def delete_test_accounts():
    uids = []
    page_token = None
    while True:
        cur_page: auth.ListUsersPage = auth.list_users(page_token=page_token, max_results=1000, app=APP)
        
        # delete all test users
        uids = []
        for user in cur_page.users:
            if user.custom_claims and user.custom_claims['test']:
                uids.append(user.uid)
        if len(uids) > 0:
            auth.delete_users(uids=uids, app=APP)

        # iterate next page
        if not cur_page.has_next_page:
            break
        else:
            page_token = cur_page.get_next_page()
    
    # delete any remaining
    if len(uids) > 0:
        auth.delete_users(uids=uids, app=APP)
