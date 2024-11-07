package env

import "testing"

func TestEnv(t *testing.T) {
	if err := Load("./.env"); err != nil {
		t.Fatal(err)
	}

	prints(t, "APP_NAME")
	prints(t, "APP_URL")
	printi(t, "APP_STATUS")
	printf(t, "APP_PRICE")
	printb(t, "APP_OPEN")
}

func printb(t *testing.T, key string) {
	val, err := GetBool(key)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s=%t", key, val)
}

func printf(t *testing.T, key string) {
	val, err := GetFloat(key)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s=%f", key, val)
}

func printi(t *testing.T, key string) {
	val, err := GetInt(key)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s=%d", key, val)
}

func prints(t *testing.T, key string) {
	val, err := Get(key)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s=%s", key, val)
}
