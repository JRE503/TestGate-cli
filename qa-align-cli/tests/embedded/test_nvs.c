/**
 * ESP-IDF Unity test fixture — NVS (Non-Volatile Storage) component
 */

// @test Verifies NVS read returns correct value after write
// @description Core data persistence path. Failure means device loses config on reboot.
// @issue FW-101
TEST_CASE("nvs_read_after_write", "[nvs]")
{
    nvs_handle_t handle;
    nvs_open("storage", NVS_READWRITE, &handle);
    nvs_set_i32(handle, "boot_count", 42);
    int32_t val = 0;
    nvs_get_i32(handle, "boot_count", &val);
    TEST_ASSERT_EQUAL_INT32(42, val);
    nvs_close(handle);
}

// @test Verifies NVS is wiped clean after nvs_flash_erase
// @description Security requirement. Erase must remove all keys, not just reset pointers.
// @issue FW-102
TEST_CASE("nvs_erase_clears_all_keys", "[nvs]")
{
    nvs_flash_erase();
    nvs_flash_init();
    nvs_handle_t handle;
    nvs_open("storage", NVS_READONLY, &handle);
    int32_t val = 0;
    esp_err_t err = nvs_get_i32(handle, "boot_count", &val);
    TEST_ASSERT_EQUAL(ESP_ERR_NVS_NOT_FOUND, err);
    nvs_close(handle);
}

// @test Verifies WiFi reconnects within 10s after signal loss
// @description Reliability SLA. Device must self-heal without manual reboot.
// @issue FW-210
TEST_CASE("wifi_reconnect_after_signal_loss", "[wifi][security]")
{
    simulate_signal_loss(5000);
    bool reconnected = wait_for_wifi(10000);
    TEST_ASSERT_TRUE(reconnected);
}
