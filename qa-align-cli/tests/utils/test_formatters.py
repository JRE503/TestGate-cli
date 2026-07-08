"""
QA-Align annotated test suite — Utility / Formatter module
"""

# [What]: Verifies currency string is formatted with correct symbol and 2dp
# [Why]: Display consistency across locales prevents user confusion in invoices.
# [Reference]: UTIL-001
def test_currency_formatter_usd():
    result = format_currency(99.9, currency="USD")
    assert result == "$99.90"


# [What]: Verifies date formatter respects ISO-8601 output
# [Why]: API consumers expect RFC3339. Drift breaks downstream integrations.
# [Reference]: UTIL-002
def test_date_formatter_iso8601():
    dt = datetime(2024, 3, 15, 10, 30, 0)
    result = format_date(dt)
    assert result == "2024-03-15T10:30:00Z"
