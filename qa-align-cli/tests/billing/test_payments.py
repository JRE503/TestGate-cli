"""
QA-Align annotated test suite — Billing & Payments module
"""

# [What]: Verifies that a valid credit card payment is processed and confirmation is returned
# [Why]: Revenue-critical path. Failed payments result in direct revenue loss.
# [Reference]: BILL-201
def test_successful_credit_card_payment():
    order = create_test_order(amount=99.99, currency="USD")
    result = billing_service.charge(order, card=VALID_TEST_CARD)
    assert result.status == "SUCCEEDED"
    assert result.transaction_id is not None


# [What]: Verifies that declined cards produce the correct error response
# [Why]: Must surface actionable error codes to the frontend — silent failures break UX.
# [Reference]: BILL-202
def test_declined_card_returns_error_code():
    order = create_test_order(amount=50.00, currency="USD")
    result = billing_service.charge(order, card=DECLINED_TEST_CARD)
    assert result.status == "DECLINED"
    assert result.error_code == "CARD_DECLINED"


# [What]: Verifies that refund is issued within 24h for qualifying cancellations
# [Why]: SLA compliance. Breach triggers regulatory penalty under EU PSD2.
# [Reference]: BILL-215
def test_refund_issued_on_cancellation_within_sla():
    transaction = create_completed_transaction(amount=149.99)
    refund = billing_service.refund(transaction.id, reason="CUSTOMER_CANCEL")
    assert refund.status == "REFUNDED"
    assert refund.processing_hours <= 24
