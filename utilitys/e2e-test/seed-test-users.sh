#!/bin/bash

set -uo pipefail

# 引数でURLを指定可能。なければデフォルト値を使用する
API_URL="${1:-http://localhost:8080}"
# URLの末尾のスラッシュを削除
API_URL="${API_URL%/}"
REGISTER_ENDPOINT="/api/signup"
MAX_TRIES=10
SLEEP_SEC=3

TEST_USERS=( 
	"e2e_test:e2e_password" 
	"e2e_test2:e2e_password2"
)

# テストユーザーをシードする
function seed_user(){
	local username="$1"
	local password="$2"
	local url="${API_URL}${REGISTER_ENDPOINT}"

	echo "Seeding user: $username"

	for i in $(seq 1 "$MAX_TRIES"); do
		# JSONペイロードを作成
		payload=$(printf '{"username":"%s","password":"%s"}' "$username" "$password")
		resp=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${url}" \
			-H "Content-Type: application/json" \
			-d "$payload")

		# HTTPステータスコードに基づいて処理を分岐、成功ステータスなら終了する
		if [[ "$resp" == "200" || "$resp" == "201" || "$resp" == "409" ]]; then
			echo "User '$username' seeded successfully (HTTP $resp)"
			return 0
		else
			echo "Attempt $i: Failed to seed user '$username' (HTTP $resp). Retrying in $SLEEP_SEC seconds"
			sleep "$SLEEP_SEC"
		fi
	done

	echo "ERROR: Failed to seed user '$username' after $MAX_TRIES attempts" >&2
	return 1
}

function main(){
	# 本番環境で実行できないように抑制する
	if [[ "$API_URL" != http://localhost:* && "$API_URL" != http://127.0.0.1:* ]]; then
  		echo "ERROR: Refusing to seed users against non-local URL: $API_URL" >&2
  		exit 1
	fi

	# curlがインストールされているか確認
	if ! command -v curl &> /dev/null; then
		echo "ERROR: curl is not installed. Please install curl to proceed." >&2
		exit 1
	fi

	local fail=0
	for user_cred in "${TEST_USERS[@]}"; do
		# TEST_USERSのフォーマットが正しいか確認
		if [[ "$user_cred" != *:* ]]; then
			echo "ERROR: Invalid user credential format: $user_cred" >&2
			fail=$((fail + 1))
			continue
		fi
		local username="${user_cred%%:*}"
		local password="${user_cred##*:}"

		# 引数に設定してseed_userを呼び出す
		if ! seed_user "$username" "$password"; then
			fail=$((fail + 1))
		fi
	done

	echo "Done seeding users. Failures: $fail"
	((fail == 0)) || exit 1
}

main "$@"
