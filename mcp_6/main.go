package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

  "example.com/go-remote-mcp-server6/models"
	"example.com/go-remote-mcp-server6/handler"
	"github.com/joho/godotenv"
)


type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// MCP Server
type MCPServer struct {
	tools map[string]Tool
	mu    sync.RWMutex
}

func NewMCPServer() *MCPServer {
	server := &MCPServer{
		tools: make(map[string]Tool),
	}
	
	// サンプルツールの登録
	server.RegisterTool(Tool{
		Name:        "echo",
		Description: "Echo back the input message",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"message": map[string]interface{}{
					"type":        "string",
					"description": "The message to echo",
				},
			},
			"required": []string{"message"},
		},
	})
	
	server.RegisterTool(Tool{
		Name:        "add",
		Description: "Add two numbers",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"a": map[string]interface{}{
					"type":        "number",
					"description": "First number",
				},
				"b": map[string]interface{}{
					"type":        "number",
					"description": "Second number",
				},
			},
			"required": []string{"a", "b"},
		},
	})

	server.RegisterTool(Tool{
		Name:        "purchase_item",
		Description: "入力された品名、価格の値を APIに送信します。",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "購入する商品の品名",
				},
				"price": map[string]interface{}{
					"type":        "number",
					"description": "商品の価格（円）",
				},				
			},
			"required": []string{"name", "price"},
		},
	})
		
	
	return server
}

func (s *MCPServer) RegisterTool(tool Tool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tools[tool.Name] = tool
}

func (s *MCPServer) HandleRequest(req models.JSONRPCRequest) models.JSONRPCResponse {
	switch req.Method {
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(req)
	case "initialize":
		return s.handleInitialize(req)
	default:
		return models.JSONRPCResponse{
			JSONRPC: "2.0",
			Error: &models.RPCError{
				Code:    -32601,
				Message: "Method not found",
			},
			ID: req.ID,
		}
	}
}

func (s *MCPServer) handleInitialize(req models.JSONRPCRequest) models.JSONRPCResponse {
	return models.JSONRPCResponse{
		JSONRPC: "2.0",
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"serverInfo": map[string]interface{}{
				"name":    "sample-mcp-server",
				"version": "1.0.0",
			},
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
		},
		ID: req.ID,
	}
}

func (s *MCPServer) handleToolsList(req models.JSONRPCRequest) models.JSONRPCResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	tools := make([]Tool, 0, len(s.tools))
	for _, tool := range s.tools {
		tools = append(tools, tool)
	}
	
	return models.JSONRPCResponse{
		JSONRPC: "2.0",
		Result: map[string]interface{}{
			"tools": tools,
		},
		ID: req.ID,
	}
}

func (s *MCPServer) handleToolsCall(req models.JSONRPCRequest) models.JSONRPCResponse {
	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}
	
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return models.JSONRPCResponse{
			JSONRPC: "2.0",
			Error: &models.RPCError{
				Code:    -32602,
				Message: "Invalid params",
				Data:    err.Error(),
			},
			ID: req.ID,
		}
	}
	log.Printf("Arguments %v" , params.Arguments)
	
	// ツールの実行
	result, err := s.executeTool(params.Name, params.Arguments)
	if err != nil {
		return models.JSONRPCResponse{
			JSONRPC: "2.0",
			Error: &models.RPCError{
				Code:    -32603,
				Message: "Internal error",
				Data:    err.Error(),
			},
			ID: req.ID,
		}
	}
	
	return models.JSONRPCResponse{
		JSONRPC: "2.0",
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": result,
				},
			},
		},
		ID: req.ID,
	}
}

func (s *MCPServer) executeTool(name string, args map[string]interface{}) (string, error) {
	s.mu.RLock()
	_, exists := s.tools[name]
	s.mu.RUnlock()
	
	if !exists {
		return "", fmt.Errorf("tool not found: %s", name)
	}
	
	switch name {
	case "echo":
		message, ok := args["message"].(string)
		if !ok {
			return "", fmt.Errorf("invalid message parameter")
		}
		return message, nil
		
	case "add":
		a, ok1 := args["a"].(float64)
		b, ok2 := args["b"].(float64)
		if !ok1 || !ok2 {
			return "", fmt.Errorf("invalid number parameters")
		}
		return fmt.Sprintf("%.2f", a+b), nil

	case "purchase_item":
		log.Printf("# purchase_item-start")
		//log.Printf("args %v" , args)
		name, ok1 := args["name"].(string)
		price, ok2 := args["price"].(float64)
		log.Printf("name=%s", name)
		if !ok1 || !ok2 {
			return "", fmt.Errorf("invalid number parameters")
		}
		price2 := int64(price)
		handler.PurchaseHnadler(name , price2)

		return fmt.Sprintf("name=%s , price= %d 円、登録しました。", name, price2), nil
	
	default:
		return "", fmt.Errorf("tool execution not implemented: %s", name)
	}
}

// HTTP Handler
func (s *MCPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	apiKey := os.Getenv("API_KEY")
	authHeader := r.Header.Get("Authorization")

	fmt.Println("Authorization: %s\n", authHeader)
  if apiKey != authHeader {
		http.Error(w, "error, Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	
	var req models.JSONRPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		resp := models.JSONRPCResponse{
			JSONRPC: "2.0",
			Error: &models.RPCError{
				Code:    -32700,
				Message: "Parse error",
			},
			ID: nil,
		}
		s.writeResponse(w, resp)
		return
	}
	
	resp := s.HandleRequest(req)
	s.writeResponse(w, resp)
}

func (s *MCPServer) writeResponse(w http.ResponseWriter, resp models.JSONRPCResponse) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func main() {
	server := NewMCPServer()
	
	http.Handle("/mcp", server)
	
	port := ":8080"
	log.Printf("MCP Server starting on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
