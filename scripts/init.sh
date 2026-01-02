CONFIG_FILE=$1
echo "📝 Creating template config at $CONFIG_FILE..."; \
cat > $CONFIG_FILE <<EOF;
spreadsheet_id: "REPLACE_WITH_YOUR_SHEET_ID"
sheet_name: "Sheet1"
timezone: "America/New_York"
auth:
  mode: "service_account"
  key_path: "./credentials.json"
habits:
  - id: "meditate"
    label: "Did you meditate?"
    column: "meditation"
    default: false
  - id: "run"
    label: "Did you run?"
    column: "run_5k"
    default: false
EOF
		echo "✅ Done! Please edit $CONFIG_FILE with your real Sheet ID."; \
