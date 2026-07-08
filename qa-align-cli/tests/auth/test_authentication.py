"""
QA-Align annotated test suite — Authentication module
"""

# [What]: Verifies that a registered user can authenticate with valid credentials
# [Why]: Core security invariant. Session tokens must be correctly issued on success.
# [Reference]: SEC-104
def test_user_authentication_with_valid_credentials():
    user = create_test_user(email="alice@corp.com", password="StrongPass!1")
    response = auth_client.login(user.email, "StrongPass!1")
    assert response.status_code == 200
    assert "access_token" in response.json()


# [What]: Verifies that login fails with an incorrect password
# [Why]: Prevents brute-force bypass. System must reject invalid credential pairs.
# [Reference]: SEC-105
def test_user_authentication_with_invalid_password():
    user = create_test_user(email="bob@corp.com", password="CorrectPass!2")
    response = auth_client.login(user.email, "WrongPassword")
    assert response.status_code == 401
    assert response.json()["error"] == "INVALID_CREDENTIALS"


# [What]: Verifies that a locked account cannot authenticate regardless of credentials
# [Why]: Account lockout policy must be enforced after N failed attempts.
# [Reference]: SEC-110
def test_locked_account_cannot_authenticate():
    user = create_test_user(email="locked@corp.com", locked=True)
    response = auth_client.login(user.email, "AnyPassword1!")
    assert response.status_code == 403
    assert response.json()["error"] == "ACCOUNT_LOCKED"


# [What]: Verifies JWT token expiry is honoured after session timeout
# [Why]: Expired tokens must be rejected to prevent unauthorized session reuse.
# [Reference]: SEC-118
def test_expired_jwt_token_is_rejected():
    token = generate_expired_jwt(user_id=42)
    response = protected_client.get("/profile", headers={"Authorization": f"Bearer {token}"})
    assert response.status_code == 401
