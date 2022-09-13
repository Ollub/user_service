import pytest
import requests

from faker import Faker

AUTH_HEADER = "x-authentication-token"
BASE_URL = "http://localhost:8080"
REGISTER_URL = f"{BASE_URL}/register"
LOGIN_URL = f"{BASE_URL}/login"
USERS_URL = f"{BASE_URL}/users"


def user_payload(**kwargs):
    faker = Faker()
    return {
        "lastName": "John",
        "firstName": "Doe",
        "email": faker.email(),
        "password": faker.password(),
        **kwargs,
}


@pytest.mark.parametrize(
    "field", ("lastName", "firstName", "email", "password"),
)
def test_user_register_empty_field(field):
    payload = user_payload()
    payload[field] = ""

    resp = requests.post(REGISTER_URL, json=payload)

    assert resp.status_code == 422
    resp_json = resp.json()
    assert resp_json["message"] == f"{field}: may not be empty"


@pytest.mark.parametrize(
    "password, msg", (
            ("Abc123", "password: should contain special characters"),
            ("Aab!!!", "password: should contain numbers"),
            ("AAa1!", "password: length should be greater then 5"),
    ),
    ids=["No special", "No digits", "Length"]
)
def test_user_register_weak_password(password, msg):
    payload = user_payload()
    payload["password"] = password

    resp = requests.post(REGISTER_URL, json=payload)

    assert resp.status_code == 422, resp.json()
    resp_json = resp.json()
    assert resp_json["message"] == msg


@pytest.mark.parametrize(
    "email", ("A", "@", "Aasdf@", "@asd", "asd.com", "@asdf.com"),
)
def test_user_register_bad_email(email):
    payload = user_payload()
    payload["email"] = email

    resp = requests.post(REGISTER_URL, json=payload)

    assert resp.status_code == 422, resp.json()
    resp_json = resp.json()
    assert resp_json["message"] == "email: invalid"


def test_list_users_auth_error():
    resp = requests.get(USERS_URL)
    assert resp.status_code == 401


def test_user_flow():
    """Test all the user flow step by step."""

    # Register User1
    resp = requests.post(REGISTER_URL, json=user_payload())
    assert resp.status_code == 201, resp.json()
    u1_resp = resp.json()

    # Register User2
    u2 = user_payload()
    resp = requests.post(REGISTER_URL, json=u2)
    assert resp.status_code == 201, resp.json()
    u2_resp = resp.json()

    # User2 gets list of users
    resp = requests.get(USERS_URL, headers={AUTH_HEADER: u2_resp["token"]})
    assert resp.status_code == 200, resp.json()
    assert {u1_resp["userId"], u2_resp["userId"]}.intersection(u["id"] for u in resp.json()["users"])

    # User2 tries to update User1 -> 403
    resp = requests.put(f"{USERS_URL}/{u1_resp['userId']}", headers={AUTH_HEADER: u2_resp["token"]})
    assert resp.status_code == 403

    # User2 changes own data
    # After this change user version will change and token will be invalid
    new_firstname = "John"
    new_lastname = "Doe"
    resp = requests.put(
        f"{USERS_URL}/{u2_resp['userId']}",
        headers={AUTH_HEADER: u2_resp["token"]},
        json={"firstName": new_firstname, "lastName": new_lastname},
    )
    assert resp.status_code == 200
    resp_json = resp.json()
    assert resp_json["firstName"] == new_firstname
    assert resp_json["lastName"] == new_lastname

    # User2 calls api with the same token -> 401
    resp = requests.get(USERS_URL, headers={AUTH_HEADER: u2_resp["token"]})
    assert resp.status_code == 401

    # User2 login with wrong password
    resp = requests.post(LOGIN_URL, json={"email": u2["email"], "password": "wrongPass"})
    assert resp.status_code == 400, resp.json()

    # User2 pass login and receive new token
    resp = requests.post(LOGIN_URL, json={"email": u2["email"], "password": u2["password"]})
    assert resp.status_code == 200, resp.json()
    u2_resp = resp.json()

    # Now user2 can call protected api
    resp = requests.get(USERS_URL, headers={AUTH_HEADER: u2_resp["token"]})
    assert resp.status_code == 200, resp.json()
