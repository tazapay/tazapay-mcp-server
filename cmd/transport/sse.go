package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
	register "github.com/tazapay/tazapay-mcp-server/tools/register"
)

type Session struct {
	ID      string
	Writer  http.ResponseWriter
	Flusher http.Flusher
	MsgCh   chan string
	Done    chan struct{}
}

type SessionStore struct {
	sessions map[string]*Session
	mutex    sync.Mutex
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*Session),
	}
}

func (s *SessionStore) Add(sess *Session) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessions[sess.ID] = sess
}

func (s *SessionStore) Get(id string) (*Session, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sess, ok := s.sessions[id]
	return sess, ok
}

func (s *SessionStore) Remove(id string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.sessions, id)
}

type MySSEServer struct {
	SessionStore *SessionStore
	ToolRegistry map[string]func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
	Logger       *slog.Logger
}

func NewMySSEServer() *MySSEServer {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create tool registry and register all tools
	toolRegistry := make(map[string]func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error))
	register.RegisterToolsInMap(toolRegistry, logger)

	// Log all registered tools for debugging
	toolNames := make([]string, 0, len(toolRegistry))
	for name := range toolRegistry {
		toolNames = append(toolNames, name)
	}
	logger.InfoContext(context.Background(), "Registered tools", "count", len(toolRegistry), "tools", toolNames)

	return &MySSEServer{
		SessionStore: NewSessionStore(),
		ToolRegistry: toolRegistry,
		Logger:       logger,
	}
}

// HandleSSE: GET /sse
func (s *MySSEServer) HandleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	sessionID := uuid.New().String()
	session := &Session{
		ID:      sessionID,
		Writer:  w,
		Flusher: w.(http.Flusher),
		MsgCh:   make(chan string, 20),
		Done:    make(chan struct{}),
	}
	s.SessionStore.Add(session)

	fmt.Fprintf(w, "event: endpoint\n")
	fmt.Fprintf(w, "data: /message?sessionId=%s\n\n", sessionID)
	session.Flusher.Flush()

	// Stream messages
	go func(sess *Session) {
		for {
			select {
			case msg := <-sess.MsgCh:
				fmt.Fprintf(sess.Writer, "event: message\n")
				fmt.Fprintf(sess.Writer, "data: %s\n\n", msg)
				sess.Flusher.Flush()
			case <-sess.Done:
				return
			}
		}
	}(session)

	// Wait for client to disconnect
	<-r.Context().Done()
	s.SessionStore.Remove(sessionID)
	close(session.Done)
}

// HandleMessage: POST /message?sessionId=...
func (s *MySSEServer) HandleMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		http.Error(w, "Missing sessionId", http.StatusBadRequest)
		return
	}

	session, ok := s.SessionStore.Get(sessionID)
	if !ok {
		http.Error(w, "Invalid sessionId", http.StatusNotFound)
		return
	}

	// Parse JSON - support multiple formats
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Log the received body for debugging
	bodyBytes, _ := json.Marshal(body)
	s.Logger.InfoContext(r.Context(), "Received request body", "body", string(bodyBytes))

	var toolName string
	var arguments interface{}

	// Check if this is a JSON-RPC MCP request
	if method, ok := body["method"].(string); ok && method == "tools/call" {
		if params, ok := body["params"].(map[string]interface{}); ok {
			if name, ok := params["name"].(string); ok {
				toolName = name
				arguments = params["arguments"]
			}
		}
	} else if method, ok := body["method"].(string); ok && method == "tools/list" {
		// Handle tools/list request - return available tools
		s.handleToolsList(session, body)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Tools list sent"))
		return
	} else if tool, ok := body["tool"].(string); ok {
		toolName = tool
		if args, ok := body["arguments"]; ok {
			arguments = args
		} else if input, ok := body["input"]; ok {
			arguments = input
		}
	} else if name, ok := body["name"].(string); ok {
		// Direct tool call format
		toolName = name
		if args, ok := body["arguments"]; ok {
			arguments = args
		} else if input, ok := body["input"]; ok {
			arguments = input
		}
	}

	if toolName == "" {
		s.Logger.ErrorContext(r.Context(), "Missing tool name in request", "body", string(bodyBytes))
		http.Error(w, "Missing tool name", http.StatusBadRequest)
		return
	}
	// Process MCP tool call asynchronously
	go func(sess *Session, requestBody map[string]any) {
		// Create a background context for tool execution that won't be canceled
		// when the HTTP request completes, but preserve auth context
		ctx := utils.AuthHeaderHTTPContextFunc(context.Background(), r)

		// Find and execute the tool using our registry
		handler, exists := s.ToolRegistry[toolName]
		if !exists {
			s.Logger.ErrorContext(ctx, "Tool not found", "tool", toolName)
			s.sendErrorResponse(sess, requestBody, fmt.Sprintf("Tool [%s] not found", toolName))
			return
		}

		// Create MCP tool call request structure
		mcpRequest := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name:      toolName,
				Arguments: arguments,
			},
		}

		result, err := handler(ctx, mcpRequest)
		if err != nil {
			s.Logger.ErrorContext(ctx, "MCP tool call failed", "tool", toolName, "error", err)
			s.sendErrorResponse(sess, requestBody, fmt.Sprintf("Tool [%s] failed: %v", toolName, err))
			return
		}

		// Send successful result
		s.sendSuccessResponse(sess, requestBody, result)
	}(session, body)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Processing started"))
}

// HandleDirectExecution: POST /execute - Direct tool execution without SSE Channel
func (s *MySSEServer) HandleDirectExecution(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request body
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Log the received body for debugging
	bodyBytes, _ := json.Marshal(body)
	s.Logger.InfoContext(r.Context(), "Direct execution request", "body", string(bodyBytes))

	var toolName string
	var arguments interface{}

	// Parse different request formats
	if method, ok := body["method"].(string); ok && method == "tools/call" {
		if params, ok := body["params"].(map[string]interface{}); ok {
			if name, ok := params["name"].(string); ok {
				toolName = name
				arguments = params["arguments"]
			}
		}
	} else if method, ok := body["method"].(string); ok && method == "tools/list" {
		// Handle tools/list request directly
		tools := make([]map[string]interface{}, 0)
		for name := range s.ToolRegistry {
			tools = append(tools, map[string]interface{}{
				"name":        name,
				"description": fmt.Sprintf("Tazapay %s tool", name),
			})
		}

		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"result": map[string]interface{}{
				"tools": tools,
			},
		}

		if id, ok := body["id"]; ok {
			response["id"] = id
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	} else if tool, ok := body["tool"].(string); ok {
		toolName = tool
		if args, ok := body["arguments"]; ok {
			arguments = args
		} else if input, ok := body["input"]; ok {
			arguments = input
		}
	} else if name, ok := body["name"].(string); ok {
		toolName = name
		if args, ok := body["arguments"]; ok {
			arguments = args
		} else if input, ok := body["input"]; ok {
			arguments = input
		}
	}

	if toolName == "" {
		s.Logger.ErrorContext(r.Context(), "Missing tool name in request", "body", string(bodyBytes))
		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"error": map[string]interface{}{
				"code":    -32602,
				"message": "Missing tool name",
			},
		}
		if id, ok := body["id"]; ok {
			response["id"] = id
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create context with auth token if present
	ctx := utils.AuthHeaderHTTPContextFunc(r.Context(), r)

	// Find the tool handler
	handler, exists := s.ToolRegistry[toolName]
	if !exists {
		s.Logger.ErrorContext(ctx, "Tool not found", "tool", toolName)
		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"error": map[string]interface{}{
				"code":    -32601,
				"message": fmt.Sprintf("Tool [%s] not found", toolName),
			},
		}
		if id, ok := body["id"]; ok {
			response["id"] = id
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create MCP tool call request
	mcpRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      toolName,
			Arguments: arguments,
		},
	}

	// Execute the tool
	result, err := handler(ctx, mcpRequest)
	if err != nil {
		s.Logger.ErrorContext(ctx, "Tool execution failed", "tool", toolName, "error", err)
		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"error": map[string]interface{}{
				"code":    -32000,
				"message": fmt.Sprintf("Tool [%s] failed: %v", toolName, err),
			},
		}
		if id, ok := body["id"]; ok {
			response["id"] = id
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Prepare successful response with structured data
	var responseData interface{}
	if len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(mcp.TextContent); ok {
			// Try to parse the text content as JSON for structured response
			var parsedData interface{}
			if err := json.Unmarshal([]byte(textContent.Text), &parsedData); err == nil {
				responseData = parsedData
			} else {
				// Check if the text contains JSON (common pattern: summary + JSON)
				responseData = s.parseToolResponse(textContent.Text)
			}
		} else {
			responseData = result.Content[0]
		}
	} else {
		responseData = map[string]interface{}{
			"message": "Tool executed successfully",
			"status":  "success",
		}
	}

	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"result":  responseData,
	}

	if id, ok := body["id"]; ok {
		response["id"] = id
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// handleToolsList sends the list of available tools to the client
func (s *MySSEServer) handleToolsList(session *Session, body map[string]any) {
	// Create tools list response
	tools := make([]map[string]interface{}, 0)
	for toolName := range s.ToolRegistry {
		tools = append(tools, map[string]interface{}{
			"name":        toolName,
			"description": fmt.Sprintf("Tazapay %s tool", toolName),
		})
	}

	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"result": map[string]interface{}{
			"tools": tools,
		},
	}

	// Add id if present in request
	if id, ok := body["id"]; ok {
		response["id"] = id
	}

	responseBytes, _ := json.Marshal(response)
	select {
	case session.MsgCh <- string(responseBytes):
	default:
	}
}

// sendErrorResponse sends a JSON-RPC error response
func (s *MySSEServer) sendErrorResponse(session *Session, requestBody map[string]any, errorMsg string) {
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"error": map[string]interface{}{
			"code":    -32000,
			"message": errorMsg,
		},
	}

	// Add id if present in request
	if id, ok := requestBody["id"]; ok {
		response["id"] = id
	}

	responseBytes, _ := json.Marshal(response)
	select {
	case session.MsgCh <- string(responseBytes):
	default:
	}
}

// sendSuccessResponse sends a JSON-RPC success response
func (s *MySSEServer) sendSuccessResponse(session *Session, requestBody map[string]any, result *mcp.CallToolResult) {
	var responseData interface{}
	if len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(mcp.TextContent); ok {
			// Try to parse the text content as JSON for structured response
			var parsedData interface{}
			if err := json.Unmarshal([]byte(textContent.Text), &parsedData); err == nil {
				responseData = parsedData
			} else {
				// Check if the text contains JSON (common pattern: summary + JSON)
				responseData = s.parseToolResponse(textContent.Text)
			}
		} else {
			responseData = result.Content[0]
		}
	} else {
		responseData = map[string]interface{}{
			"message": "Tool executed successfully",
			"status":  "success",
		}
	}

	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"result":  responseData,
	}

	// Add id if present in request
	if id, ok := requestBody["id"]; ok {
		response["id"] = id
	}

	responseBytes, _ := json.Marshal(response)
	select {
	case session.MsgCh <- string(responseBytes):
	default:
	}
}

// parseToolResponse intelligently parses tool responses that may contain
// summary text followed by JSON data
func (s *MySSEServer) parseToolResponse(text string) interface{} {
	// Look for common patterns where JSON follows descriptive text
	jsonStartPatterns := []string{
		"Full Data: {",
		"Data: {",
		"Response: {",
		"{", // Direct JSON
	}

	var jsonStart = -1
	var jsonStartOffset = 0

	// Find where JSON starts using multiple patterns
	for _, pattern := range jsonStartPatterns {
		if idx := strings.Index(text, pattern); idx >= 0 {
			jsonStart = idx
			// If pattern includes prefix text (like "Full Data: "), skip it
			if pattern != "{" {
				jsonStartOffset = len(pattern) - 1 // -1 to include the opening brace
			}
			break
		}
	}

	// If we found JSON, try to parse it
	if jsonStart >= 0 {
		// Extract everything before JSON as summary
		summaryText := text[:jsonStart]

		// Extract JSON part
		jsonText := text[jsonStart+jsonStartOffset:]

		var parsedData interface{}
		if err := json.Unmarshal([]byte(jsonText), &parsedData); err == nil {
			// Successfully parsed JSON, create structured response
			response := map[string]interface{}{
				"data": parsedData,
			}

			// Add summary information if present
			if strings.TrimSpace(summaryText) != "" {
				summary := make(map[string]interface{})
				lines := strings.Split(summaryText, "\n")

				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" && strings.Contains(line, ":") {
						parts := strings.SplitN(line, ":", 2)
						if len(parts) == 2 {
							key := strings.TrimSpace(parts[0])
							value := strings.TrimSpace(parts[1])
							// Convert key to snake_case
							key = strings.ToLower(strings.ReplaceAll(key, " ", "_"))
							summary[key] = value
						}
					}
				}

				if len(summary) > 0 {
					response["summary"] = summary
				}
			}

			return response
		}
	}

	// If no JSON found or parsing failed, return as text
	return map[string]interface{}{
		"content": text,
		"type":    "text",
	}
}
