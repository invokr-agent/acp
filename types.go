// Translated once from the ACP v1 JSON schema (schema-v1.26.0) and maintained by hand.

package acp

import "encoding/json"

// Authentication-related capabilities supported by the agent.
type AgentAuthCapabilities struct {
	// Whether the agent supports the logout method.
	Logout *LogoutCapabilities `json:"logout,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Capabilities supported by the agent.
type AgentCapabilities struct {
	// Whether the agent supports `session/load`.
	LoadSession bool `json:"loadSession,omitempty"`
	// Prompt capabilities supported by the agent.
	PromptCapabilities *PromptCapabilities `json:"promptCapabilities,omitempty"`
	// MCP capabilities supported by the agent.
	MCPCapabilities *MCPCapabilities `json:"mcpCapabilities,omitempty"`
	// Session lifecycle and prompt capabilities advertised by the agent.
	SessionCapabilities *SessionCapabilities `json:"sessionCapabilities,omitempty"`
	// Authentication-related capabilities supported by the agent.
	Auth *AgentAuthCapabilities `json:"auth,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Optional annotations for the client. The client can use annotations to inform how objects are used or displayed
type Annotations struct {
	// Intended recipients for this content, such as the user or assistant.
	Audience []Role `json:"audience,omitempty"`
	// Timestamp indicating when the underlying resource was last modified.
	LastModified string `json:"lastModified,omitempty"`
	// Relative importance of this content when clients choose what to surface.
	Priority *float64 `json:"priority,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Audio provided to or from an LLM.
type AudioContent struct {
	// Optional annotations that help clients decide how to display or route this content.
	Annotations *Annotations `json:"annotations,omitempty"`
	// Base64-encoded media payload.
	Data string `json:"data"`
	// MIME type describing the encoded media payload.
	MIMEType string `json:"mimeType"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Authentication capabilities supported by the client.
type AuthCapabilities struct {
	// Whether the client supports `terminal` authentication methods.
	Terminal bool `json:"terminal,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Describes an available authentication method.
type AuthMethod struct {
	// Client runs the configured agent program as a separate interactive process, without passing this method to `authenticate`.
	Terminal *AuthMethodTerminal `tag:"terminal"`
	// Agent handles authentication itself through `authenticate`.
	Agent *AuthMethodAgent `tag:""`
	// Unknown holds a variant this version does not know, verbatim.
	Unknown json.RawMessage
}

func (u AuthMethod) MarshalJSON() ([]byte, error) { return marshalUnion("type", &u) }

func (u *AuthMethod) UnmarshalJSON(data []byte) error { return unmarshalUnion("type", data, u) }

// Agent handles authentication itself through `authenticate`.
type AuthMethodAgent struct {
	// Unique identifier for this authentication method.
	ID AuthMethodID `json:"id"`
	// Human-readable name of the authentication method.
	Name string `json:"name"`
	// Optional description providing more details about this authentication method.
	Description string `json:"description,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Typed identifier used for auth method values on the wire.
type AuthMethodID string

// Terminal-based authentication method.
type AuthMethodTerminal struct {
	// Unique identifier for this authentication method.
	ID AuthMethodID `json:"id"`
	// Human-readable name of the authentication method.
	Name string `json:"name"`
	// Optional description providing more details about this authentication method.
	Description string `json:"description,omitempty"`
	// Additional arguments to append to the configured agent invocation for terminal auth.
	Args []string `json:"args,omitempty"`
	// Additional environment variables to set on the configured agent invocation for terminal auth. These values override same-named variables in the base launch configuration.
	Env map[string]string `json:"env,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Request parameters for the authenticate method.
type AuthenticateRequest struct {
	// The ID of the authentication method to use. Must be one of the methods advertised in the initialize response.
	MethodID AuthMethodID `json:"methodId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response to the `authenticate` method.
type AuthenticateResponse struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Information about a command.
type AvailableCommand struct {
	// Command name (e.g., `create_plan`, `research_codebase`).
	Name string `json:"name"`
	// Human-readable description of what the command does.
	Description string `json:"description"`
	// Input for the command if required
	Input *AvailableCommandInput `json:"input,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// The input specification for a command.
type AvailableCommandInput struct {
	// A hint to display when the input hasn't been provided yet
	Hint string `json:"hint,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Available commands are ready or have changed
type AvailableCommandsUpdate struct {
	// Commands the agent can execute
	AvailableCommands []AvailableCommand `json:"availableCommands"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Binary resource contents.
type BlobResourceContents struct {
	// Base64-encoded bytes for a binary resource payload.
	Blob string `json:"blob"`
	// MIME type describing the encoded media payload.
	MIMEType string `json:"mimeType,omitempty"`
	// URI associated with this resource or media payload.
	URI string `json:"uri"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Capabilities for boolean session configuration options.
type BooleanConfigOptionCapabilities struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Schema for boolean properties in an elicitation form.
type BooleanPropertySchema struct {
	// Optional title for the property.
	Title string `json:"title,omitempty"`
	// Human-readable description.
	Description string `json:"description,omitempty"`
	// Default value.
	Default bool `json:"default,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Notification to cancel ongoing operations for a session.
type CancelNotification struct {
	// The ID of the session to cancel operations for.
	SessionID SessionID `json:"sessionId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Notification to cancel an ongoing request.
type CancelRequestNotification struct {
	// The ID of the request to cancel.
	RequestID json.RawMessage `json:"requestId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Capabilities supported by the client.
type ClientCapabilities struct {
	// File system capabilities supported by the client. Determines which file operations the agent can request.
	FS *FileSystemCapabilities `json:"fs,omitempty"`
	// Whether the Client support all `terminal/*` methods.
	Terminal bool `json:"terminal,omitempty"`
	// Session-related capabilities supported by the client.
	Session *ClientSessionCapabilities `json:"session,omitempty"`
	// Authentication capabilities supported by the client. Determines which authentication method types the agent may include in its `InitializeResponse`.
	Auth *AuthCapabilities `json:"auth,omitempty"`
	// Elicitation capabilities supported by the client. Determines which elicitation modes the agent may use.
	Elicitation *ElicitationCapabilities `json:"elicitation,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Session-related capabilities supported by the client.
type ClientSessionCapabilities struct {
	// Support for ID-addressed context compaction updates. Omitted or `null` means unsupported; `{}` advertises the complete compaction contract.
	Compaction *CompactionCapabilities `json:"compaction,omitempty"`
	// Config option capabilities supported by the client.
	ConfigOptions *SessionConfigOptionsCapabilities `json:"configOptions,omitempty"`
	// Support for live user-facing `notice` session updates.
	Notices *NoticeCapabilities `json:"notices,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Request parameters for closing an active session.
type CloseSessionRequest struct {
	// The ID of the session to close.
	SessionID SessionID `json:"sessionId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response from closing a session.
type CloseSessionResponse struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Client support for ID-addressed context compaction updates.
type CompactionCapabilities struct {
}

// Unique identifier for a context compaction within a session.
type CompactionID string

// Lifecycle state of a context compaction.
type CompactionStatus string

const (
	// Compaction has started and has not finished.
	CompactionStatusInProgress CompactionStatus = "in_progress"
	// Compaction finished successfully.
	CompactionStatusCompleted CompactionStatus = "completed"
	// Compaction finished unsuccessfully.
	CompactionStatusFailed CompactionStatus = "failed"
	// Compaction was cancelled before it finished.
	CompactionStatusCancelled CompactionStatus = "cancelled"
)

// A content block appended to a compaction's summary. A first-seen ID creates an in-progress compaction. Chunks append in receive order. Agents MUST only send this update when the Client advertised the `session.compaction` Client capability.
type CompactionSummaryChunk struct {
	// ID of the compaction whose summary receives this content.
	CompactionID CompactionID `json:"compactionId"`
	// One content block to append.
	Content ContentBlock `json:"content"`
	// Metadata scoped to this chunk. Omission and `null` both mean absent.
	Meta map[string]any `json:"_meta,omitempty"`
}

// A context compaction upsert. The first notification fixes the compaction's timeline position. Later updates with the same ID patch that entity in place. Agents MUST only send this update when the Client advertised the `session.compaction` Client capability.
type CompactionUpdate struct {
	// The Agent-owned ID of this compaction, unique within the session.
	CompactionID CompactionID `json:"compactionId"`
	// Current lifecycle status.
	Status CompactionStatus `json:"status"`
	// Complete replacement user-displayable summary content for the compaction.
	Summary []ContentBlock `json:"summary,omitempty"`
	// Human-readable error details for the compaction.
	Error string `json:"error,omitempty"`
	// Extensible metadata patch for this compaction.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Notification sent by the agent when a URL-based elicitation is complete.
type CompleteElicitationNotification struct {
	// The ID of the elicitation that completed.
	ElicitationID ElicitationID `json:"elicitationId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Session configuration options have been updated.
type ConfigOptionUpdate struct {
	// The full set of configuration options and their current values.
	ConfigOptions []SessionConfigOption `json:"configOptions"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Standard content block (text, images, resources).
type Content struct {
	// The actual content block.
	Content ContentBlock `json:"content"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Content blocks represent displayable information in the Agent Client Protocol.
type ContentBlock struct {
	// Text content. May be plain text or formatted with Markdown.
	Text *TextContent `tag:"text"`
	// Images for visual context or analysis.
	Image *ImageContent `tag:"image"`
	// Audio data for transcription or analysis.
	Audio *AudioContent `tag:"audio"`
	// References to resources that the agent can access.
	ResourceLink *ResourceLink `tag:"resource_link"`
	// Complete resource contents embedded directly in the message.
	Resource *EmbeddedResource `tag:"resource"`
	// Unknown holds a variant this version does not know, verbatim.
	Unknown json.RawMessage
}

func (u ContentBlock) MarshalJSON() ([]byte, error) { return marshalUnion("type", &u) }

func (u *ContentBlock) UnmarshalJSON(data []byte) error { return unmarshalUnion("type", data, u) }

// A streamed item of content
type ContentChunk struct {
	// A single item of content
	Content ContentBlock `json:"content"`
	// A unique identifier for the message this chunk belongs to.
	MessageID MessageID `json:"messageId,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Cost information for a session.
type Cost struct {
	// Total cumulative cost for session.
	Amount float64 `json:"amount"`
	// ISO 4217 currency code (e.g., "USD", "EUR").
	Currency string `json:"currency"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Request from the agent to elicit structured user input.
type CreateElicitationRequest struct {
	// A human-readable message describing what input is needed.
	Message string `json:"message"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
	Mode string         `json:"mode,omitempty"`
	// A JSON Schema describing the form fields to present to the user.
	RequestedSchema *ElicitationSchema `json:"requestedSchema,omitempty"`
	// The session this elicitation is tied to.
	SessionID SessionID `json:"sessionId,omitempty"`
	// Optional tool call within the session.
	ToolCallID ToolCallID `json:"toolCallId,omitempty"`
	// The request this elicitation is tied to.
	RequestID json.RawMessage `json:"requestId,omitempty"`
	// The unique identifier for this elicitation.
	ElicitationID ElicitationID `json:"elicitationId,omitempty"`
	// The URL to direct the user to.
	URL string `json:"url,omitempty"`
}

// Response from the client to an elicitation request.
type CreateElicitationResponse struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta   map[string]any `json:"_meta,omitempty"`
	Action string         `json:"action,omitempty"`
	// The user-provided content, if any, as an object matching the requested schema.
	Content map[string]ElicitationContentValue `json:"content,omitempty"`
}

// Request to create a new terminal and execute a command.
type CreateTerminalRequest struct {
	// The session ID for this request.
	SessionID SessionID `json:"sessionId"`
	// The command to execute.
	Command string `json:"command"`
	// Array of command arguments.
	Args []string `json:"args,omitempty"`
	// Environment variables for the command.
	Env []EnvVariable `json:"env,omitempty"`
	// Working directory for the command. Must be an absolute path.
	Cwd string `json:"cwd,omitempty"`
	// Maximum number of output bytes to retain.
	OutputByteLimit *uint64 `json:"outputByteLimit,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response containing the ID of the created terminal.
type CreateTerminalResponse struct {
	// The unique identifier for the created terminal.
	TerminalID TerminalID `json:"terminalId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// The current mode of the session has changed
type CurrentModeUpdate struct {
	// The ID of the current mode
	CurrentModeID SessionModeID `json:"currentModeId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Request parameters for deleting an existing session from `session/list`.
type DeleteSessionRequest struct {
	// The ID of the session to delete.
	SessionID SessionID `json:"sessionId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response from deleting a session.
type DeleteSessionResponse struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// A diff representing file modifications.
type Diff struct {
	// The absolute file path being modified.
	Path string `json:"path"`
	// The original content (omitted or `null` for new files).
	OldText string `json:"oldText,omitempty"`
	// The new content after modification.
	NewText string `json:"newText"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// The user accepted the elicitation and provided content.
type ElicitationAcceptAction struct {
	// The user-provided content, if any, as an object matching the requested schema.
	Content map[string]ElicitationContentValue `json:"content,omitempty"`
}

// Elicitation capabilities supported by the client.
type ElicitationCapabilities struct {
	// Whether the client supports form-based elicitation.
	Form *ElicitationFormCapabilities `json:"form,omitempty"`
	// Whether the client supports URL-based elicitation.
	URL *ElicitationUrlCapabilities `json:"url,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Allowed wire representations for `ElicitationContentValue`.
type ElicitationContentValue = json.RawMessage

// Form-based elicitation capabilities.
type ElicitationFormCapabilities struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Form-based elicitation mode where the client renders a form from the provided schema.
type ElicitationFormMode struct {
	// A JSON Schema describing the form fields to present to the user.
	RequestedSchema ElicitationSchema `json:"requestedSchema"`
	// The session this elicitation is tied to.
	SessionID SessionID `json:"sessionId,omitempty"`
	// Optional tool call within the session.
	ToolCallID ToolCallID `json:"toolCallId,omitempty"`
	// The request this elicitation is tied to.
	RequestID json.RawMessage `json:"requestId,omitempty"`
}

// Unique identifier for an elicitation.
type ElicitationID string

// Property schema for elicitation form fields.
type ElicitationPropertySchema = json.RawMessage

// Request-scoped elicitation, tied to a specific JSON-RPC request outside of a session (e.g., during auth/configuration phases before any session is started).
type ElicitationRequestScope struct {
	// The request this elicitation is tied to.
	RequestID json.RawMessage `json:"requestId"`
}

// Type-safe elicitation schema for requesting structured user input.
type ElicitationSchema struct {
	// Type discriminator. Always `"object"`.
	Type ElicitationSchemaType `json:"type,omitempty"`
	// Optional title for the schema.
	Title string `json:"title,omitempty"`
	// Property definitions (must be primitive types).
	Properties map[string]ElicitationPropertySchema `json:"properties,omitempty"`
	// List of required property names.
	Required []string `json:"required,omitempty"`
	// Optional description of what this schema represents.
	Description string `json:"description,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Type discriminator for elicitation schemas.
type ElicitationSchemaType string

const (
	// Object schema type.
	ElicitationSchemaTypeObject ElicitationSchemaType = "object"
)

// Session-scoped elicitation, optionally tied to a specific tool call.
type ElicitationSessionScope struct {
	// The session this elicitation is tied to.
	SessionID SessionID `json:"sessionId"`
	// Optional tool call within the session.
	ToolCallID ToolCallID `json:"toolCallId,omitempty"`
}

// URL-based elicitation capabilities.
type ElicitationUrlCapabilities struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// URL-based elicitation mode where the client directs the user to a URL.
type ElicitationUrlMode struct {
	// The unique identifier for this elicitation.
	ElicitationID ElicitationID `json:"elicitationId"`
	// The URL to direct the user to.
	URL string `json:"url"`
	// The session this elicitation is tied to.
	SessionID SessionID `json:"sessionId,omitempty"`
	// Optional tool call within the session.
	ToolCallID ToolCallID `json:"toolCallId,omitempty"`
	// The request this elicitation is tied to.
	RequestID json.RawMessage `json:"requestId,omitempty"`
}

// The contents of a resource, embedded into a prompt or tool call result.
type EmbeddedResource struct {
	// Optional annotations that help clients decide how to display or route this content.
	Annotations *Annotations `json:"annotations,omitempty"`
	// Embedded resource payload, either text or binary data.
	Resource EmbeddedResourceResource `json:"resource"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Resource content that can be embedded in a message.
type EmbeddedResourceResource struct {
	// MIME type describing the encoded media payload.
	MIMEType string `json:"mimeType,omitempty"`
	// Text payload carried by this content block.
	Text string `json:"text,omitempty"`
	// URI associated with this resource or media payload.
	URI string `json:"uri,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
	// Base64-encoded bytes for a binary resource payload.
	Blob string `json:"blob,omitempty"`
}

// A titled enum option with a const value, human-readable title, and optional description.
type EnumOption struct {
	// The constant value for this option.
	Const string `json:"const"`
	// Human-readable title for this option.
	Title string `json:"title"`
	// Human-readable description.
	Description string `json:"description,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// An environment variable to set when launching an MCP server.
type EnvVariable struct {
	// The name of the environment variable.
	Name string `json:"name"`
	// The value to set for the environment variable.
	Value string `json:"value"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// JSON-RPC error object.
type Error struct {
	// A number indicating the error type that occurred. This must be an integer as defined in the JSON-RPC specification.
	Code ErrorCode `json:"code"`
	// A string providing a short description of the error. The message should be limited to a concise single sentence.
	Message string `json:"message"`
	// Optional primitive or structured value that contains additional information about the error. This may include debugging information or context-specific details.
	Data json.RawMessage `json:"data,omitempty"`
}

// Predefined error codes for common JSON-RPC and ACP-specific errors.
type ErrorCode int

const (
	// **Parse error**: Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text.
	ErrParseError ErrorCode = -32700
	// **Invalid request**: The JSON sent is not a valid Request object.
	ErrInvalidRequest ErrorCode = -32600
	// **Method not found**: The method does not exist or is not available.
	ErrMethodNotFound ErrorCode = -32601
	// **Invalid params**: Invalid method parameter(s).
	ErrInvalidParams ErrorCode = -32602
	// **Internal error**: Internal JSON-RPC error. Reserved for implementation-defined server errors.
	ErrInternalError ErrorCode = -32603
	// **Request cancelled**: Execution of the method was aborted either due to a cancellation request from the caller or because of resource constraints or shutdown.
	ErrRequestCancelled ErrorCode = -32800
	// **Authentication required**: Authentication is required before this operation can be performed.
	ErrAuthRequired ErrorCode = -32000
	// **Resource not found**: A given resource, such as a session or file, was not found.
	ErrResourceNotFound ErrorCode = -32002
)

// File system capabilities that a client may support.
type FileSystemCapabilities struct {
	// Whether the Client supports `fs/read_text_file` requests.
	ReadTextFile bool `json:"readTextFile,omitempty"`
	// Whether the Client supports `fs/write_text_file` requests.
	WriteTextFile bool `json:"writeTextFile,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// An HTTP header to set when making requests to the MCP server.
type HttpHeader struct {
	// The name of the HTTP header.
	Name string `json:"name"`
	// The value to set for the HTTP header.
	Value string `json:"value"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// An image provided to or from an LLM.
type ImageContent struct {
	// Optional annotations that help clients decide how to display or route this content.
	Annotations *Annotations `json:"annotations,omitempty"`
	// Base64-encoded media payload.
	Data string `json:"data"`
	// MIME type describing the encoded media payload.
	MIMEType string `json:"mimeType"`
	// URI associated with this resource or media payload.
	URI string `json:"uri,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Metadata about the implementation of the client or agent. Describes the name and version of an ACP implementation, with an optional title for UI representation.
type Implementation struct {
	// Intended for programmatic or logical use, but can be used as a display name fallback if title isn’t present.
	Name string `json:"name"`
	// Intended for UI and end-user contexts — optimized to be human-readable and easily understood.
	Title string `json:"title,omitempty"`
	// Version of the implementation. Can be displayed to the user or used for debugging or metrics purposes. (e.g. "1.0.0").
	Version string `json:"version"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Request parameters for the initialize method.
type InitializeRequest struct {
	// The latest protocol version supported by the client.
	ProtocolVersion ProtocolVersion `json:"protocolVersion"`
	// Capabilities supported by the client.
	ClientCapabilities *ClientCapabilities `json:"clientCapabilities,omitempty"`
	// Information about the Client name and version sent to the Agent.
	ClientInfo *Implementation `json:"clientInfo,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response to the `initialize` method.
type InitializeResponse struct {
	// The protocol version the client specified if supported by the agent, or the latest protocol version supported by the agent.
	ProtocolVersion ProtocolVersion `json:"protocolVersion"`
	// Capabilities supported by the agent.
	AgentCapabilities *AgentCapabilities `json:"agentCapabilities,omitempty"`
	// Authentication methods supported by the agent.
	AuthMethods []AuthMethod `json:"authMethods,omitempty"`
	// Information about the Agent name and version sent to the Client.
	AgentInfo *Implementation `json:"agentInfo,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Schema for integer properties in an elicitation form.
type IntegerPropertySchema struct {
	// Optional title for the property.
	Title string `json:"title,omitempty"`
	// Human-readable description.
	Description string `json:"description,omitempty"`
	// Minimum value (inclusive).
	Minimum *int64 `json:"minimum,omitempty"`
	// Maximum value (inclusive).
	Maximum *int64 `json:"maximum,omitempty"`
	// Default value.
	Default *int64 `json:"default,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Request to kill a terminal without releasing it.
type KillTerminalRequest struct {
	// The session ID for this request.
	SessionID SessionID `json:"sessionId"`
	// The ID of the terminal to kill.
	TerminalID TerminalID `json:"terminalId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response to `terminal/kill` method
type KillTerminalResponse struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Request parameters for listing existing sessions.
type ListSessionsRequest struct {
	// Filter sessions by working directory. Must be an absolute path.
	Cwd string `json:"cwd,omitempty"`
	// Opaque cursor token from a previous response's nextCursor field for cursor-based pagination
	Cursor string `json:"cursor,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response from listing sessions.
type ListSessionsResponse struct {
	// Array of session information objects
	Sessions []SessionInfo `json:"sessions"`
	// Opaque cursor token. If present, pass this in the next request's cursor parameter to fetch the next page. If absent, there are no more results.
	NextCursor string `json:"nextCursor,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Request parameters for loading an existing session.
type LoadSessionRequest struct {
	// List of MCP servers to connect to for this session.
	MCPServers []MCPServer `json:"mcpServers"`
	// The working directory for this session. Must be an absolute path.
	Cwd string `json:"cwd"`
	// Additional workspace roots to activate for this session. Each path must be absolute.
	AdditionalDirectories []string `json:"additionalDirectories,omitempty"`
	// The ID of the session to load.
	SessionID SessionID `json:"sessionId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response from loading an existing session.
type LoadSessionResponse struct {
	// Initial mode state if supported by the Agent
	Modes *SessionModeState `json:"modes,omitempty"`
	// Initial session configuration options if supported by the Agent.
	ConfigOptions []SessionConfigOption `json:"configOptions,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Logout capabilities supported by the agent.
type LogoutCapabilities struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Request parameters for the logout method.
type LogoutRequest struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response to the `logout` method.
type LogoutResponse struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// MCP capabilities supported by the agent
type MCPCapabilities struct {
	// Agent supports `http` MCP servers.
	HTTP bool `json:"http,omitempty"`
	// Agent supports `sse` MCP servers.
	SSE bool `json:"sse,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Configuration for connecting to an MCP (Model Context Protocol) server.
type MCPServer struct {
	// HTTP transport configuration
	HTTP *MCPServerHTTP `tag:"http"`
	// SSE transport configuration
	SSE *MCPServerSSE `tag:"sse"`
	// Stdio transport configuration
	Stdio *MCPServerStdio `tag:""`
	// Unknown holds a variant this version does not know, verbatim.
	Unknown json.RawMessage
}

func (u MCPServer) MarshalJSON() ([]byte, error) { return marshalUnion("type", &u) }

func (u *MCPServer) UnmarshalJSON(data []byte) error { return unmarshalUnion("type", data, u) }

// HTTP transport configuration for MCP.
type MCPServerHTTP struct {
	// Human-readable name identifying this MCP server.
	Name string `json:"name"`
	// URL to the MCP server.
	URL string `json:"url"`
	// HTTP headers to set when making requests to the MCP server.
	Headers []HttpHeader `json:"headers"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// SSE transport configuration for MCP.
type MCPServerSSE struct {
	// Human-readable name identifying this MCP server.
	Name string `json:"name"`
	// URL to the MCP server.
	URL string `json:"url"`
	// HTTP headers to set when making requests to the MCP server.
	Headers []HttpHeader `json:"headers"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Stdio transport configuration for MCP.
type MCPServerStdio struct {
	// Human-readable name identifying this MCP server.
	Name string `json:"name"`
	// Absolute path to the MCP server executable.
	Command string `json:"command"`
	// Command-line arguments to pass to the MCP server.
	Args []string `json:"args"`
	// Environment variables to set when launching the MCP server.
	Env []EnvVariable `json:"env"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Unique identifier for a message within a session.
type MessageID string

// Items for a multi-select (array) property schema.
type MultiSelectItems = json.RawMessage

// Schema for multi-select (array) properties in an elicitation form.
type MultiSelectPropertySchema struct {
	// Optional title for the property.
	Title string `json:"title,omitempty"`
	// Human-readable description.
	Description string `json:"description,omitempty"`
	// Minimum number of items to select.
	MinItems *uint64 `json:"minItems,omitempty"`
	// Maximum number of items to select.
	MaxItems *uint64 `json:"maxItems,omitempty"`
	// The items definition describing allowed values.
	Items MultiSelectItems `json:"items"`
	// Default selected values.
	Default []string `json:"default,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Request parameters for creating a new session.
type NewSessionRequest struct {
	// The working directory for this session. Must be an absolute path.
	Cwd string `json:"cwd"`
	// Additional workspace roots for this session. Each path must be absolute.
	AdditionalDirectories []string `json:"additionalDirectories,omitempty"`
	// List of MCP (Model Context Protocol) servers the agent should connect to.
	MCPServers []MCPServer `json:"mcpServers"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response from creating a new session.
type NewSessionResponse struct {
	// Unique identifier for the created session.
	SessionID SessionID `json:"sessionId"`
	// Initial mode state if supported by the Agent
	Modes *SessionModeState `json:"modes,omitempty"`
	// Initial session configuration options if supported by the Agent.
	ConfigOptions []SessionConfigOption `json:"configOptions,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Fire-and-forget information for the user.
type Notice struct {
	// Presentation severity hint.
	Severity NoticeSeverity `json:"severity"`
	// Required non-empty plain-text title that can stand alone.
	Title string `json:"title"`
	// Optional plain-text detail or guidance.
	Description string `json:"description,omitempty"`
	// Metadata scoped to this notice.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Client support for presenting live notices to the user.
type NoticeCapabilities struct {
}

// Severity hint for a session notice.
type NoticeSeverity string

const (
	// Informational notice.
	NoticeSeverityInfo NoticeSeverity = "info"
	// Warning notice.
	NoticeSeverityWarning NoticeSeverity = "warning"
	// Error notice.
	NoticeSeverityError NoticeSeverity = "error"
)

// Schema for number (floating-point) properties in an elicitation form.
type NumberPropertySchema struct {
	// Optional title for the property.
	Title string `json:"title,omitempty"`
	// Human-readable description.
	Description string `json:"description,omitempty"`
	// Minimum value (inclusive).
	Minimum *float64 `json:"minimum,omitempty"`
	// Maximum value (inclusive).
	Maximum *float64 `json:"maximum,omitempty"`
	// Default value.
	Default *float64 `json:"default,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// An option presented to the user when requesting permission.
type PermissionOption struct {
	// Unique identifier for this permission option.
	OptionID PermissionOptionID `json:"optionId"`
	// Human-readable label to display to the user.
	Name string `json:"name"`
	// Hint about the nature of this permission option.
	Kind PermissionOptionKind `json:"kind"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Unique identifier for a permission option.
type PermissionOptionID string

// The type of permission option being presented to the user.
type PermissionOptionKind string

const (
	// Allow this operation only this time.
	PermissionOptionKindAllowOnce PermissionOptionKind = "allow_once"
	// Allow this operation and remember the choice.
	PermissionOptionKindAllowAlways PermissionOptionKind = "allow_always"
	// Reject this operation only this time.
	PermissionOptionKindRejectOnce PermissionOptionKind = "reject_once"
	// Reject this operation and remember the choice.
	PermissionOptionKindRejectAlways PermissionOptionKind = "reject_always"
)

// An execution plan for accomplishing complex tasks.
type Plan struct {
	// The list of tasks to be accomplished.
	Entries []PlanEntry `json:"entries"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// A single entry in the execution plan.
type PlanEntry struct {
	// Human-readable description of what this task aims to accomplish.
	Content string `json:"content"`
	// The relative importance of this task. Used to indicate which tasks are most critical to the overall goal.
	Priority PlanEntryPriority `json:"priority"`
	// Current execution status of this task.
	Status PlanEntryStatus `json:"status"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Priority levels for plan entries.
type PlanEntryPriority string

const (
	// High priority task - critical to the overall goal.
	PlanEntryPriorityHigh PlanEntryPriority = "high"
	// Medium priority task - important but not critical.
	PlanEntryPriorityMedium PlanEntryPriority = "medium"
	// Low priority task - nice to have but not essential.
	PlanEntryPriorityLow PlanEntryPriority = "low"
)

// Status of a plan entry in the execution flow.
type PlanEntryStatus string

const (
	// The task has not started yet.
	PlanEntryStatusPending PlanEntryStatus = "pending"
	// The task is currently being worked on.
	PlanEntryStatusInProgress PlanEntryStatus = "in_progress"
	// The task has been successfully completed.
	PlanEntryStatusCompleted PlanEntryStatus = "completed"
)

// Prompt capabilities supported by the agent in `session/prompt` requests.
type PromptCapabilities struct {
	// Agent supports `image` content blocks.
	Image bool `json:"image,omitempty"`
	// Agent supports `audio` content blocks.
	Audio bool `json:"audio,omitempty"`
	// Agent supports embedded context in `session/prompt` requests.
	EmbeddedContext bool `json:"embeddedContext,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Request parameters for sending a user prompt to the agent.
type PromptRequest struct {
	// The ID of the session to send this user message to
	SessionID SessionID `json:"sessionId"`
	// The blocks of content that compose the user's message.
	Prompt []ContentBlock `json:"prompt"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response from processing a user prompt.
type PromptResponse struct {
	// Indicates why the agent stopped processing the turn.
	StopReason StopReason `json:"stopReason"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Protocol version identifier.
type ProtocolVersion uint16

// Request to read content from a text file.
type ReadTextFileRequest struct {
	// The session ID for this request.
	SessionID SessionID `json:"sessionId"`
	// Absolute path to the file to read.
	Path string `json:"path"`
	// Line number to start reading from (1-based).
	Line *uint32 `json:"line,omitempty"`
	// Maximum number of lines to read.
	Limit *uint32 `json:"limit,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response containing the contents of a text file.
type ReadTextFileResponse struct {
	// Content payload returned by this response.
	Content string `json:"content"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Request to release a terminal and free its resources.
type ReleaseTerminalRequest struct {
	// The session ID for this request.
	SessionID SessionID `json:"sessionId"`
	// The ID of the terminal to release.
	TerminalID TerminalID `json:"terminalId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response to terminal/release method
type ReleaseTerminalResponse struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// The outcome of a permission request.
type RequestPermissionOutcome struct {
	// The prompt turn was cancelled before the user responded.
	Cancelled *RequestPermissionOutcomeCancelled `tag:"cancelled"`
	// The user selected one of the provided options.
	Selected *SelectedPermissionOutcome `tag:"selected"`
	// Unknown holds a variant this version does not know, verbatim.
	Unknown json.RawMessage
}

func (u RequestPermissionOutcome) MarshalJSON() ([]byte, error) { return marshalUnion("outcome", &u) }

func (u *RequestPermissionOutcome) UnmarshalJSON(data []byte) error {
	return unmarshalUnion("outcome", data, u)
}

// Request for user permission to execute a tool call.
type RequestPermissionRequest struct {
	// The session ID for this request.
	SessionID SessionID `json:"sessionId"`
	// Details about the tool call requiring permission.
	ToolCall ToolCallUpdate `json:"toolCall"`
	// Available permission options for the user to choose from.
	Options []PermissionOption `json:"options"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response to a permission request.
type RequestPermissionResponse struct {
	// The user's decision on the permission request.
	Outcome RequestPermissionOutcome `json:"outcome"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// A resource that the server is capable of reading, included in a prompt or tool call result.
type ResourceLink struct {
	// Optional annotations that help clients decide how to display or route this content.
	Annotations *Annotations `json:"annotations,omitempty"`
	// Optional human-readable details shown with this protocol object.
	Description string `json:"description,omitempty"`
	// MIME type describing the encoded media payload.
	MIMEType string `json:"mimeType,omitempty"`
	// Human-readable name shown for this protocol object.
	Name string `json:"name"`
	// Optional size of the linked resource in bytes, if known.
	Size *int64 `json:"size,omitempty"`
	// Optional display title for end-user UI.
	Title string `json:"title,omitempty"`
	// URI associated with this resource or media payload.
	URI string `json:"uri"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Request parameters for resuming an existing session.
type ResumeSessionRequest struct {
	// The ID of the session to resume.
	SessionID SessionID `json:"sessionId"`
	// The working directory for this session. Must be an absolute path.
	Cwd string `json:"cwd"`
	// Additional workspace roots to activate for this session. Each path must be absolute.
	AdditionalDirectories []string `json:"additionalDirectories,omitempty"`
	// List of MCP servers to connect to for this session.
	MCPServers []MCPServer `json:"mcpServers,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response from resuming an existing session.
type ResumeSessionResponse struct {
	// Initial mode state if supported by the Agent
	Modes *SessionModeState `json:"modes,omitempty"`
	// Initial session configuration options if supported by the Agent.
	ConfigOptions []SessionConfigOption `json:"configOptions,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// The sender or recipient of messages and data in a conversation.
type Role string

const (
	// The assistant side of a conversation.
	RoleAssistant Role = "assistant"
	// The user side of a conversation.
	RoleUser Role = "user"
)

// The user selected one of the provided options.
type SelectedPermissionOutcome struct {
	// The ID of the option the user selected.
	OptionID PermissionOptionID `json:"optionId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Capabilities for additional session directories support.
type SessionAdditionalDirectoriesCapabilities struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Session capabilities supported by the agent.
type SessionCapabilities struct {
	// Whether the agent supports `session/list`.
	List *SessionListCapabilities `json:"list,omitempty"`
	// Whether the agent supports `session/delete`.
	Delete *SessionDeleteCapabilities `json:"delete,omitempty"`
	// Whether the agent supports `additionalDirectories` on supported session lifecycle requests.
	AdditionalDirectories *SessionAdditionalDirectoriesCapabilities `json:"additionalDirectories,omitempty"`
	// Whether the agent supports `session/resume`.
	Resume *SessionResumeCapabilities `json:"resume,omitempty"`
	// Whether the agent supports `session/close`.
	Close *SessionCloseCapabilities `json:"close,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Capabilities for the `session/close` method.
type SessionCloseCapabilities struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// A boolean on/off toggle session configuration option payload.
type SessionConfigBoolean struct {
	// The current value of the boolean option.
	CurrentValue bool `json:"currentValue"`
}

// Unique identifier for a session configuration option value group.
type SessionConfigGroupID string

// Unique identifier for a session configuration option.
type SessionConfigID string

// A session configuration option selector and its current state.
type SessionConfigOption struct {
	// Unique identifier for the configuration option.
	ID SessionConfigID `json:"id"`
	// Human-readable label for the option.
	Name string `json:"name"`
	// Optional description for the Client to display to the user.
	Description string `json:"description,omitempty"`
	// Optional semantic category for this option (UX only).
	Category SessionConfigOptionCategory `json:"category,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
	Type string         `json:"type,omitempty"`
	// The currently selected value.
	CurrentValue json.RawMessage `json:"currentValue,omitempty"`
	// The set of selectable options.
	Options SessionConfigSelectOptions `json:"options,omitempty"`
}

// Semantic category for a session configuration option.
type SessionConfigOptionCategory string

const (
	// Session mode selector.
	SessionConfigOptionCategoryMode SessionConfigOptionCategory = "mode"
	// Model selector.
	SessionConfigOptionCategoryModel SessionConfigOptionCategory = "model"
	// Model-related configuration parameter.
	SessionConfigOptionCategoryModelConfig SessionConfigOptionCategory = "model_config"
	// Thought/reasoning level selector.
	SessionConfigOptionCategoryThoughtLevel SessionConfigOptionCategory = "thought_level"
)

// Session configuration option capabilities supported by the client.
type SessionConfigOptionsCapabilities struct {
	// Whether the client supports boolean session configuration options.
	Boolean *BooleanConfigOptionCapabilities `json:"boolean,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// A single-value selector (dropdown) session configuration option payload.
type SessionConfigSelect struct {
	// The currently selected value.
	CurrentValue SessionConfigValueID `json:"currentValue"`
	// The set of selectable options.
	Options SessionConfigSelectOptions `json:"options"`
}

// A group of possible values for a session configuration option.
type SessionConfigSelectGroup struct {
	// Unique identifier for this group.
	Group SessionConfigGroupID `json:"group"`
	// Human-readable label for this group.
	Name string `json:"name"`
	// The set of option values in this group.
	Options []SessionConfigSelectOption `json:"options"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// A possible value for a session configuration option.
type SessionConfigSelectOption struct {
	// Unique identifier for this option value.
	Value SessionConfigValueID `json:"value"`
	// Human-readable label for this option value.
	Name string `json:"name"`
	// Optional description for this option value.
	Description string `json:"description,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Possible values for a session configuration option.
type SessionConfigSelectOptions = json.RawMessage

// Unique identifier for a session configuration option value.
type SessionConfigValueID string

// Capabilities for the `session/delete` method.
type SessionDeleteCapabilities struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// A unique identifier for a conversation session between a client and agent.
type SessionID string

// Information about a session returned by session/list
type SessionInfo struct {
	// Unique identifier for the session
	SessionID SessionID `json:"sessionId"`
	// The working directory for this session. Must be an absolute path.
	Cwd string `json:"cwd"`
	// Additional workspace roots reported for this session. Each path must be absolute.
	AdditionalDirectories []string `json:"additionalDirectories,omitempty"`
	// Human-readable title for the session
	Title string `json:"title,omitempty"`
	// ISO 8601 timestamp of last activity
	UpdatedAt string `json:"updatedAt,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Update to session metadata. All fields are optional to support partial updates.
type SessionInfoUpdate struct {
	// Human-readable title for the session. Set to null to clear.
	Title string `json:"title,omitempty"`
	// ISO 8601 timestamp of last activity. Set to null to clear.
	UpdatedAt string `json:"updatedAt,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Capabilities for the `session/list` method.
type SessionListCapabilities struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// A mode the agent can operate in.
type SessionMode struct {
	// Stable identifier used to refer to this protocol object in later messages.
	ID SessionModeID `json:"id"`
	// Human-readable name shown for this protocol object.
	Name string `json:"name"`
	// Optional human-readable details shown with this protocol object.
	Description string `json:"description,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Unique identifier for a Session Mode.
type SessionModeID string

// The set of modes and the one currently active.
type SessionModeState struct {
	// The current mode the Agent is in.
	CurrentModeID SessionModeID `json:"currentModeId"`
	// The set of modes that the Agent can operate in
	AvailableModes []SessionMode `json:"availableModes"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Notification containing a session update from the agent.
type SessionNotification struct {
	// The ID of the session this update pertains to.
	SessionID SessionID `json:"sessionId"`
	// The actual update content.
	Update SessionUpdate `json:"update"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Capabilities for the `session/resume` method.
type SessionResumeCapabilities struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Different types of updates that can be sent during session processing.
type SessionUpdate struct {
	// A chunk of the user's message being streamed.
	UserMessageChunk *ContentChunk `tag:"user_message_chunk"`
	// A chunk of the agent's response being streamed.
	AgentMessageChunk *ContentChunk `tag:"agent_message_chunk"`
	// A chunk of the agent's internal reasoning being streamed.
	AgentThoughtChunk *ContentChunk `tag:"agent_thought_chunk"`
	// Notification that a new tool call has been initiated.
	ToolCall *ToolCall `tag:"tool_call"`
	// Update on the status or results of a tool call.
	ToolCallUpdate *ToolCallUpdate `tag:"tool_call_update"`
	// The agent's execution plan for complex tasks. See protocol docs: [Agent Plan](https://agentclientprotocol.com/protocol/agent-plan)
	Plan *Plan `tag:"plan"`
	// Available commands are ready or have changed
	AvailableCommandsUpdate *AvailableCommandsUpdate `tag:"available_commands_update"`
	// The current mode of the session has changed
	CurrentModeUpdate *CurrentModeUpdate `tag:"current_mode_update"`
	// Session configuration options have been updated.
	ConfigOptionUpdate *ConfigOptionUpdate `tag:"config_option_update"`
	// Session metadata has been updated (title, timestamps, custom metadata)
	SessionInfoUpdate *SessionInfoUpdate `tag:"session_info_update"`
	// Context window and cost update for the session.
	UsageUpdate *UsageUpdate `tag:"usage_update"`
	// Information for the user that is not part of session history.
	Notice *Notice `tag:"notice"`
	// A context compaction has been created or updated.
	CompactionUpdate *CompactionUpdate `tag:"compaction_update"`
	// A content block appended to a context compaction's retained summary.
	CompactionSummaryChunk *CompactionSummaryChunk `tag:"compaction_summary_chunk"`
	// Unknown holds a variant this version does not know, verbatim.
	Unknown json.RawMessage
}

func (u SessionUpdate) MarshalJSON() ([]byte, error) { return marshalUnion("sessionUpdate", &u) }

func (u *SessionUpdate) UnmarshalJSON(data []byte) error {
	return unmarshalUnion("sessionUpdate", data, u)
}

// Request parameters for setting a session configuration option.
type SetSessionConfigOptionRequest struct {
	// The ID of the session to set the configuration option for.
	SessionID SessionID `json:"sessionId"`
	// The ID of the configuration option to set.
	ConfigID SessionConfigID `json:"configId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
	Type string         `json:"type,omitempty"`
	// The boolean value.
	Value json.RawMessage `json:"value,omitempty"`
}

// Response to `session/set_config_option` method.
type SetSessionConfigOptionResponse struct {
	// The full set of configuration options and their current values.
	ConfigOptions []SessionConfigOption `json:"configOptions"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Request parameters for setting a session mode.
type SetSessionModeRequest struct {
	// The ID of the session to set the mode for.
	SessionID SessionID `json:"sessionId"`
	// The ID of the mode to set.
	ModeID SessionModeID `json:"modeId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response to `session/set_mode` method.
type SetSessionModeResponse struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Reasons why an agent stops processing a prompt turn.
type StopReason string

const (
	// The turn ended successfully.
	StopReasonEndTurn StopReason = "end_turn"
	// The turn ended because the agent reached the maximum number of tokens.
	StopReasonMaxTokens StopReason = "max_tokens"
	// The turn ended because the agent reached the maximum number of allowed agent requests between user turns.
	StopReasonMaxTurnRequests StopReason = "max_turn_requests"
	// The turn ended because the agent refused to continue. The user prompt and everything that comes after it won't be included in the next prompt, so this should be reflected in the UI.
	StopReasonRefusal StopReason = "refusal"
	// The turn was cancelled by the client via `session/cancel`.
	StopReasonCancelled StopReason = "cancelled"
)

// String format types for string properties in elicitation schemas.
type StringFormat string

const (
	// Email address format.
	StringFormatEmail StringFormat = "email"
	// URI format.
	StringFormatURI StringFormat = "uri"
	// Date format (YYYY-MM-DD).
	StringFormatDate StringFormat = "date"
	// Date-time format (ISO 8601).
	StringFormatDateminusTime StringFormat = "date-time"
)

// String item schema for multi-select enum properties.
type StringMultiSelectItems struct {
	// Allowed enum values.
	Enum []string `json:"enum"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Schema for string properties in an elicitation form.
type StringPropertySchema struct {
	// Optional title for the property.
	Title string `json:"title,omitempty"`
	// Human-readable description.
	Description string `json:"description,omitempty"`
	// Minimum string length.
	MinLength *uint32 `json:"minLength,omitempty"`
	// Maximum string length.
	MaxLength *uint32 `json:"maxLength,omitempty"`
	// Pattern the string must match.
	Pattern string `json:"pattern,omitempty"`
	// String format.
	Format StringFormat `json:"format,omitempty"`
	// Default value.
	Default string `json:"default,omitempty"`
	// Enum values for untitled single-select enums. Optional. Omitted and `null` are equivalent and mean no untitled single-select choices are declared by `enum`.
	Enum []string `json:"enum,omitempty"`
	// Titled enum options for titled single-select enums. Optional. Omitted and `null` are equivalent and mean no titled single-select choices are declared by `oneOf`.
	OneOf []EnumOption `json:"oneOf,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Embed a terminal created with `terminal/create` by its id.
type Terminal struct {
	// Identifier of the terminal instance to embed in the content stream.
	TerminalID TerminalID `json:"terminalId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Exit status of a terminal command.
type TerminalExitStatus struct {
	// The process exit code (may be null if terminated by signal).
	ExitCode *uint32 `json:"exitCode,omitempty"`
	// The signal that terminated the process (may be null if exited normally).
	Signal string `json:"signal,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Typed identifier used for terminal values on the wire.
type TerminalID string

// Request to get the current output and status of a terminal.
type TerminalOutputRequest struct {
	// The session ID for this request.
	SessionID SessionID `json:"sessionId"`
	// The ID of the terminal to get output from.
	TerminalID TerminalID `json:"terminalId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response containing the terminal output and exit status.
type TerminalOutputResponse struct {
	// The terminal output captured so far.
	Output string `json:"output"`
	// Whether the output was truncated due to byte limits.
	Truncated bool `json:"truncated"`
	// Exit status if the command has completed.
	ExitStatus *TerminalExitStatus `json:"exitStatus,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Text provided to or from an LLM.
type TextContent struct {
	// Optional annotations that help clients decide how to display or route this content.
	Annotations *Annotations `json:"annotations,omitempty"`
	// Text payload carried by this content block.
	Text string `json:"text"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Text-based resource contents.
type TextResourceContents struct {
	// MIME type describing the encoded media payload.
	MIMEType string `json:"mimeType,omitempty"`
	// Text payload carried by this content block.
	Text string `json:"text"`
	// URI associated with this resource or media payload.
	URI string `json:"uri"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Items definition for titled multi-select enum properties.
type TitledMultiSelectItems struct {
	// Titled enum options.
	AnyOf []EnumOption `json:"anyOf"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Represents a tool call that the language model has requested.
type ToolCall struct {
	// Unique identifier for this tool call within the session.
	ToolCallID ToolCallID `json:"toolCallId"`
	// Human-readable title describing what the tool is doing.
	Title string `json:"title"`
	// Programmatic name of the tool being invoked.
	Name string `json:"name,omitempty"`
	// The category of tool being invoked. Helps clients choose appropriate icons and UI treatment.
	Kind ToolKind `json:"kind,omitempty"`
	// Current execution status of the tool call.
	Status ToolCallStatus `json:"status,omitempty"`
	// Content produced by the tool call.
	Content []ToolCallContent `json:"content,omitempty"`
	// File locations affected by this tool call. Enables "follow-along" features in clients.
	Locations []ToolCallLocation `json:"locations,omitempty"`
	// Raw input parameters sent to the tool.
	RawInput json.RawMessage `json:"rawInput,omitempty"`
	// Raw output returned by the tool.
	RawOutput json.RawMessage `json:"rawOutput,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Content produced by a tool call.
type ToolCallContent struct {
	// Standard content block (text, images, resources).
	Content *Content `tag:"content"`
	// File modification shown as a diff.
	Diff *Diff `tag:"diff"`
	// Embed a terminal created with `terminal/create` by its id.
	Terminal *Terminal `tag:"terminal"`
	// Unknown holds a variant this version does not know, verbatim.
	Unknown json.RawMessage
}

func (u ToolCallContent) MarshalJSON() ([]byte, error) { return marshalUnion("type", &u) }

func (u *ToolCallContent) UnmarshalJSON(data []byte) error { return unmarshalUnion("type", data, u) }

// Unique identifier for a tool call within a session.
type ToolCallID string

// A file location being accessed or modified by a tool.
type ToolCallLocation struct {
	// The absolute file path being accessed or modified.
	Path string `json:"path"`
	// Optional line number within the file.
	Line *uint32 `json:"line,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Execution status of a tool call.
type ToolCallStatus string

const (
	// The tool call hasn't started running yet because the input is either streaming or we're awaiting approval.
	ToolCallStatusPending ToolCallStatus = "pending"
	// The tool call is currently running.
	ToolCallStatusInProgress ToolCallStatus = "in_progress"
	// The tool call completed successfully.
	ToolCallStatusCompleted ToolCallStatus = "completed"
	// The tool call failed with an error.
	ToolCallStatusFailed ToolCallStatus = "failed"
)

// An update to an existing tool call.
type ToolCallUpdate struct {
	// The ID of the tool call being updated.
	ToolCallID ToolCallID `json:"toolCallId"`
	// Update the tool kind.
	Kind ToolKind `json:"kind,omitempty"`
	// Update the execution status.
	Status ToolCallStatus `json:"status,omitempty"`
	// Update the human-readable title.
	Title string `json:"title,omitempty"`
	// Update the programmatic name of the tool being invoked.
	Name string `json:"name,omitempty"`
	// Replace the content collection.
	Content []ToolCallContent `json:"content,omitempty"`
	// Replace the locations collection.
	Locations []ToolCallLocation `json:"locations,omitempty"`
	// Update the raw input.
	RawInput json.RawMessage `json:"rawInput,omitempty"`
	// Update the raw output.
	RawOutput json.RawMessage `json:"rawOutput,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Categories of tools that can be invoked.
type ToolKind string

const (
	// Reading files or data.
	ToolKindRead ToolKind = "read"
	// Modifying files or content.
	ToolKindEdit ToolKind = "edit"
	// Removing files or data.
	ToolKindDelete ToolKind = "delete"
	// Moving or renaming files.
	ToolKindMove ToolKind = "move"
	// Searching for information.
	ToolKindSearch ToolKind = "search"
	// Running commands or code.
	ToolKindExecute ToolKind = "execute"
	// Internal reasoning or planning.
	ToolKindThink ToolKind = "think"
	// Retrieving external data.
	ToolKindFetch ToolKind = "fetch"
	// Switching the current session mode.
	ToolKindSwitchMode ToolKind = "switch_mode"
	// Other tool types (default).
	ToolKindOther ToolKind = "other"
)

// All text that was typed after the command name is provided as input.
type UnstructuredCommandInput struct {
	// A hint to display when the input hasn't been provided yet
	Hint string `json:"hint"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Context window and cost update for a session.
type UsageUpdate struct {
	// Tokens currently in context.
	Used uint64 `json:"used"`
	// Total context window size in tokens.
	Size uint64 `json:"size"`
	// Cumulative session cost (optional).
	Cost *Cost `json:"cost,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Request to wait for a terminal command to exit.
type WaitForTerminalExitRequest struct {
	// The session ID for this request.
	SessionID SessionID `json:"sessionId"`
	// The ID of the terminal to wait for.
	TerminalID TerminalID `json:"terminalId"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response containing the exit status of a terminal command.
type WaitForTerminalExitResponse struct {
	// The process exit code (may be null if terminated by signal).
	ExitCode *uint32 `json:"exitCode,omitempty"`
	// The signal that terminated the process (may be null if exited normally).
	Signal string `json:"signal,omitempty"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Request to write content to a text file.
type WriteTextFileRequest struct {
	// The session ID for this request.
	SessionID SessionID `json:"sessionId"`
	// Absolute path to the file to write.
	Path string `json:"path"`
	// The text content to write to the file.
	Content string `json:"content"`
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// Response to `fs/write_text_file`
type WriteTextFileResponse struct {
	// The _meta property is reserved by ACP to allow clients and agents to attach additional metadata to their interactions. Implementations MUST NOT make assumptions about values at these keys.
	Meta map[string]any `json:"_meta,omitempty"`
}

// The prompt turn was cancelled before the user responded.
type RequestPermissionOutcomeCancelled struct {
}
