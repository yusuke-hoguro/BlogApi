# MCPの導入手順

1. CopliotとGitHub Copliot chatを導入 ※拡張機能から

1. Ctrl + Shift + P　でMCPをいれてMCPサーバー追加を選択

1. NPMパッケージを選択

1. @playwright/mcp を入力

1. グローバルで追加する

/C:/Users/保黒悠介/AppData/Roaming/Code/User/mcp.json

```
{
	"servers": {
		"playwright-mcp": {
			"command": "npx",
			"args": [
				"-y",
				"playwright-mcp"
			],
			"type": "stdio"
		}
	},
	"inputs": []
}
```

WSLで使う場合には設定ファイルをWSL側にもつ必要がある