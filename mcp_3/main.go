package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// JSONRPCリクエスト構造体
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCレスポンス構造体
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

// RPCエラー構造体
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ツールパラメータ
type PurchaseParams struct {
	Name string `json:"name"`
	Price    int    `json:"price"`
}
type PostRequest struct {
	Content string `json:"content"`
	Data  string `json:"data"`
}


// ツールリスト
type ToolsList struct {
	Tools []Tool `json:"tools"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// ツール呼び出しパラメータ
type CallToolParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// ツール実行結果
type ToolResult struct {
	Content []Content `json:"content"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	for scanner.Scan() {
		line := scanner.Text()
		
		var req JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			sendError(writer, nil, -32700, "Parse error")
			continue
		}

		handleRequest(writer, req)
	}
}

func handleRequest(writer *bufio.Writer, req JSONRPCRequest) {
	switch req.Method {
	case "initialize":
		handleInitialize(writer, req)
	case "tools/list":
		handleToolsList(writer, req)
	case "tools/call":
		handleToolsCall(writer, req)
	default:
		sendError(writer, req.ID, -32601, "Method not found")
	}
}

func handleInitialize(writer *bufio.Writer, req JSONRPCRequest) {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"serverInfo": map[string]string{
			"name":    "purchase-server",
			"version": "1.0.0",
		},
		"capabilities": map[string]interface{}{
			"tools": map[string]bool{},
		},
	}
	sendResponse(writer, req.ID, result)
}

func handleToolsList(writer *bufio.Writer, req JSONRPCRequest) {
	tools := ToolsList{
		Tools: []Tool{
			{
				Name:        "purchase_item",
				Description: "入力された品名、価格の値を APIに送信します。",
				InputSchema: InputSchema{
					Type: "object",
					Properties: map[string]Property{
						"name": {
							Type:        "string",
							Description: "購入する商品の品名",
						},
						"price": {
							Type:        "integer",
							Description: "商品の価格（円）",
						},
					},
					Required: []string{"name", "price"},
				},
			},
		},
	}
	sendResponse(writer, req.ID, tools)
}

const API_ENDPOINT = "http://localhost:8080/api/data/create"

func sendToAPI(data *PostRequest) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("JSONエンコードエラー: %v", err)
	}
	
	resp, err := http.Post(
		API_ENDPOINT,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("APIリクエストエラー: %v", err)
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("APIエラー (ステータスコード: %d): %s", resp.StatusCode, string(body))
	}
	
	log.Printf("✓ データ送信成功 (ステータス: %d)", resp.StatusCode)
	log.Printf("  レスポンス: %s", string(body))
	
	return nil
}

func handleToolsCall(writer *bufio.Writer, req JSONRPCRequest) {
	var params CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		sendError(writer, req.ID, -32602, "Invalid params")
		return
	}

	if params.Name != "purchase_item" {
		sendError(writer, req.ID, -32602, "Unknown tool")
		return
	}

	var args PurchaseParams
	if err := json.Unmarshal(params.Arguments, &args); err != nil {
		sendError(writer, req.ID, -32602, "Invalid arguments")
		return
	}
	jsonBytes, err := json.Marshal(args)
	if err != nil {
		log.Fatalf("JSONへの変換に失敗しました: %v", err)
	}
	jsonString := string(jsonBytes)

	jsonPost := PostRequest{
		Content: "item_price",
		Data:  jsonString,
	}
	sendErr := sendToAPI(&jsonPost)
	if sendErr != nil {
		sendError(writer, req.ID, -32602, "sendToAPI , error")
		return
	}

	toolResult := ToolResult{
		Content: []Content{
			{
				Type: "text",
				Text: fmt.Sprintf("購入情報\n品名: %s\n価格: %d円", args.Name, args.Price),
			},
		},
	}

	sendResponse(writer, req.ID, toolResult)
}

func sendResponse(writer *bufio.Writer, id interface{}, result interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	
	data, _ := json.Marshal(resp)
	writer.Write(data)
	writer.WriteByte('\n')
	writer.Flush()
}

func sendError(writer *bufio.Writer, id interface{}, code int, message string) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
	}
	
	data, _ := json.Marshal(resp)
	writer.Write(data)
	writer.WriteByte('\n')
	writer.Flush()
}
