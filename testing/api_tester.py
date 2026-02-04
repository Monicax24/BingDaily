import os
import time
import uuid
import json
from typing import Any, Dict

import requests
import firebase_test as firebase_test
import upload_test as upload

# Configurable via env; sensible defaults for local dev
BASE_URL = os.getenv("BACKEND_BASE_URL")
if BASE_URL == None:
    print("Must specify \"BACKEND_BASE_URL\"")
    exit()

COMMUNITY_ID = os.getenv("COMMUNITY_ID")
if COMMUNITY_ID == None:
    print("Must specify \"COMMUNITY_ID\"")
    exit()

def hdrs(token: str) -> Dict[str, str]:
    return {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json",
        "Accept": "application/json",
    }

def req_ok(resp: requests.Response) -> bool:
    try:
        j = resp.json()
    except Exception:
        return False
    return j.get("status") == "success"


def print_test_result(title: str, passed: bool, detail: str = "") -> None:
    status = "PASS" if passed else "FAIL"
    line = f"[ {status} ] {title}"
    if detail:
        line += f" — {detail}"
    print(line)


def get_json(resp: requests.Response) -> Dict[str, Any]:
    try:
        return resp.json()
    except Exception:
        return {"status": "fail", "message": f"non-json response: {resp.text[:200]}"}


def test_profile(token: str) -> bool:
    r = requests.get(f"{BASE_URL}/users/profile", headers=hdrs(token), timeout=15)
    j = get_json(r)
    passed = req_ok(r) and isinstance(j.get("data", {}).get("user"), dict)
    print_test_result("GET /users/profile", passed, j.get("message", ""))
    if passed:
        user = j["data"]["user"]
        print(json.dumps({"userId": user.get("userId"), "username": user.get("username")}, indent=2))
    return passed


def test_join_community(token: str, comm_id: str) -> bool:
    r = requests.get(f"{BASE_URL}/communities/join/{comm_id}", headers=hdrs(token), timeout=15)
    j = get_json(r)
    passed = req_ok(r)
    print_test_result("GET /communities/join/:id", passed, j.get("message", ""))
    return passed


def test_leave_community(token: str, comm_id: str) -> bool:
    r = requests.get(f"{BASE_URL}/communities/leave/{comm_id}", headers=hdrs(token), timeout=15)
    j = get_json(r)
    passed = req_ok(r)
    print_test_result("GET /communities/leave/:id", passed, j.get("message", ""))
    return passed   


def test_get_community(token: str, community_id: str) -> bool:
    r = requests.get(f"{BASE_URL}/communities/{community_id}", headers=hdrs(token), timeout=15)
    j = get_json(r)
    comm = j.get("data", {}).get("community") if req_ok(r) else None
    passed = req_ok(r) and isinstance(comm, dict) and comm.get("communityId") == community_id
    print_test_result("GET /communities/:id", passed, j.get("message", ""))
    if passed:
        print(json.dumps({"name": comm.get("name"), "memberCnt": comm.get("memberCnt")}, indent=2))
    return passed


def test_get_posts(token: str, community_id: str) -> Dict[str, Any]:
    r = requests.get(f"{BASE_URL}/communities/posts/{community_id}", headers=hdrs(token), timeout=20)
    j = get_json(r)
    posts = j.get("data", {}).get("posts") if req_ok(r) else None
    passed = req_ok(r)
    print_test_result("GET /communities/posts/:id", passed, j.get("message", ""))
    count = len(posts) if isinstance(posts, list) else None
    if isinstance(posts, list):
        print(json.dumps({"postCount": count}, indent=2))
        print("POST IMAGES:")
        for i in range(count):
            print(f"URL {i}: {posts[i]["imageUrl"]}\n")
        print("\n\n")
    return {"passed": passed, "count": count}


def test_create_post(token: str, community_id: str) -> Dict[str, Any]:
    body = {"communityId": community_id, "caption": f"Hello from tester {uuid.uuid4().hex[:8]}"}
    r = requests.post(f"{BASE_URL}/communities/posts/upload", headers=hdrs(token), json=body, timeout=20)
    j = get_json(r)
    data = j.get("data", {}) if req_ok(r) else {}
    post_id = data.get("postId")
    upload_url = data.get("uploadUrl")
    passed = req_ok(r) and isinstance(post_id, str) and isinstance(upload_url, str)
    print_test_result("POST /communities/posts/upload", passed, j.get("message", ""))
    return {"passed": passed, "postId": post_id, "uploadUrl": upload_url}


def test_register_user(token: str, uid: str, email: str, username: str) -> Dict[str, Any]:
    body = {"email": email, "username": username, "updatePicture": False}
    r = requests.post(f"{BASE_URL}/users/register", headers=hdrs(token), json=body, timeout=15)
    j = get_json(r)
    passed = req_ok(r)
    print_test_result("POST /users/register", passed, j.get("message", ""))
    return {"passed": passed, "uid": uid, "email": email, "username": username}

def test_fetch_communites(token: str) -> Dict[str, Any]:
    r = requests.get(f"{BASE_URL}/communities/list", headers=hdrs(token), timeout=15)
    j = get_json(r)
    passed = req_ok(r)
    comms = j.get("data", {})["communities"]
    print_test_result("GET /communites/list", passed, j.get('message', ""))
    if passed:
        print(comms)
        print()
    return {"passed": passed}

def test_delete_post(token: str, comm_id: str) -> Dict[str, Any]:
    r = requests.get(
        f"{BASE_URL}/communities/posts/delete/{comm_id}",
        headers=hdrs(token), timeout=15
    ) 
    passed = req_ok(r)
    j = get_json(r)
    print_test_result("GET /communities/posts/delete/:communityId", passed, j.get('message', ""))
    return {"passed": passed}

def main() -> int:
    print(f"Base URL: {BASE_URL}")
    print(f"Community: {COMMUNITY_ID}")
    print("Starting API tests with new user creation flow...\n")

    results = {}

    # Generate test user credentials
    test_uid = str(uuid.uuid4())
    test_email = f"testuser_{uuid.uuid4().hex[:8]}@test.local"
    test_username = f"testuser_{uuid.uuid4().hex[:8]}"

    print(f"Test User: uid={test_uid}, email={test_email}, username={test_username}\n")

    # Step 1: Create user in Firebase
    try:
        print("[*] Creating Firebase user...")
        firebase_test.create_user(test_uid, test_email)
        results["firebase_create"] = True
        print_test_result("Firebase create_user", True)
    except Exception as e:
        print_test_result("Firebase create_user", False, str(e))
        results["firebase_create"] = False
        return 1

    # Step 2: Get auth token for this user (simulating OAuth)
    try:
        print("\n[*] Getting auth token...")
        auth_token = firebase_test.get_auth_token(test_uid, test_email)
        print_test_result("firebase_testing.get_auth_token", True)
        results["get_token"] = True
    except Exception as e:
        print_test_result("firebase_testing.get_auth_token", False, str(e))
        results["get_token"] = False
        return 1

    # Step 3: Register user to backend
    register_result = test_register_user(auth_token, test_uid, test_email, test_username)
    results["register"] = register_result["passed"]
    if not register_result["passed"]:
        return 1

    print(f"\n[*] Using auth token for all subsequent requests...\n")

    # Step 3.1: Fetch communities list
    results["fetch_communities"] = test_fetch_communites(auth_token)

    # Step 4: Fetch user profile
    results["profile"] = test_profile(auth_token)

    # Step 4.5: Join default community
    results["join_community"] = test_join_community(auth_token, COMMUNITY_ID)

    # Step 5: Get community data
    results["community"] = test_get_community(auth_token, COMMUNITY_ID)

    # Step 6: Get posts before upload
    before = test_get_posts(auth_token, COMMUNITY_ID)
    results["posts_before"] = before["passed"]

    # Step 7: Create a post
    created = test_create_post(auth_token, COMMUNITY_ID)
    results["create_post"] = created["passed"]

    # Step 8: Upload media to S3
    if created["passed"]:
        time.sleep(0.5)
        results["upload_media"] = upload.upload_random(created["uploadUrl"])
    else:
        results["upload_media"] = False

    # Step 9: Get posts after upload
    after = test_get_posts(auth_token, COMMUNITY_ID)
    results["posts_after"] = after["passed"]

    # Basic delta check
    if before["passed"] and after["passed"] and results["create_post"]:
        grew = (after["count"] or 0) >= (before["count"] or 0)
        print_test_result("Posts count non-decreasing after upload", grew)
        results["posts_grew"] = grew

    print("\nSummary:")
    print(json.dumps(results, indent=2))

    input("\nPress ENTER to resume testing...\n")
    # Step 10: Delete post 
    results["delete_post"] = test_delete_post(auth_token, COMMUNITY_ID)

    # Non-zero exit if any core step failed
    core = ["firebase_create", "register", "get_token", "profile", "community", "posts_before", "create_post", "upload_media", "posts_after"]
    failed = [k for k in core if not results.get(k, False)]
    if failed:
        print(f"\nOne or more checks failed: {', '.join(failed)}")
        return 1
    return 0

if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print(f"Error occured, aborting tests...\n\n{'='*50}\n")
        firebase_test.delete_test_accounts()
        raise e
    firebase_test.delete_test_accounts()