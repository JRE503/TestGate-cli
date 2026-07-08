/**
 * QA-Align annotated test suite — User Service API (TypeScript)
 */

// @test Creates a new user account with valid payload
// @description Ensures POST /users returns 201 and a valid user ID on correct input
// @issue USR-301
function test_create_user_with_valid_payload() {
  const res = apiClient.post('/users', { name: 'Alice', email: 'alice@corp.com' });
  expect(res.status).toBe(201);
  expect(res.body.id).toBeDefined();
}

// @test Rejects duplicate email registration
// @description System must enforce email uniqueness at the API layer before DB write
// @issue USR-302
function test_create_user_duplicate_email_rejected() {
  apiClient.post('/users', { name: 'Alice', email: 'alice@corp.com' });
  const res = apiClient.post('/users', { name: 'Alice2', email: 'alice@corp.com' });
  expect(res.status).toBe(409);
  expect(res.body.error).toBe('EMAIL_ALREADY_EXISTS');
}

// @test Returns 404 for non-existent user lookup
// @description GET /users/:id with unknown ID must return structured 404, not 500
// @issue USR-310
function test_get_nonexistent_user_returns_404() {
  const res = apiClient.get('/users/99999999');
  expect(res.status).toBe(404);
}
