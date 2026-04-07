package ai

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// Concierge is the AI brain of Haus — it takes natural language commands and
// turns them into real device actions. Think of it as the world's most
// dramatic smart home assistant. Illusions, Michael.
type Concierge struct {
	client    anthropic.Client
	kasaFuncs *KasaFuncs
	hueFuncs  *HueFuncs
	HTTPQuery DeviceHTTPQuery // for authenticated HTTP requests to devices
}

// ChatResponse holds the result of a concierge conversation turn.
type ChatResponse struct {
	Text      string                   `json:"text"`
	ToolCalls []ToolCallResult         `json:"tool_calls,omitempty"`
	Messages  []anthropic.MessageParam `json:"messages"`
}

// ToolCallResult records a single tool invocation and its result.
type ToolCallResult struct {
	Tool   string `json:"tool"`
	Input  string `json:"input"`
	Result string `json:"result"`
}

// NewConcierge creates a new AI concierge with the given device closures.
// The API key is read from the ANTHROPIC_API_KEY environment variable by
// the SDK automatically. If you need to pass it explicitly, the client
// constructor accepts option.WithAPIKey().
func NewConcierge(apiKey string, kasaFuncs *KasaFuncs, hueFuncs *HueFuncs) *Concierge {
	var opts []option.RequestOption
	if apiKey != "" {
		opts = append(opts, option.WithAPIKey(apiKey))
	}
	client := anthropic.NewClient(opts...)
	return &Concierge{
		client:    client,
		kasaFuncs: kasaFuncs,
		hueFuncs:  hueFuncs,
	}
}

// Chat sends a user message to the concierge, executes any tool calls Claude
// requests, and returns the final text response along with the updated
// conversation history. The tool loop runs up to 5 iterations — even GOB
// knows when to stop performing encores.
func (c *Concierge) Chat(ctx context.Context, message string, history []anthropic.MessageParam) (*ChatResponse, error) {
	// Build system prompt with live device context.
	systemPrompt, err := c.buildSystemPrompt()
	if err != nil {
		return nil, fmt.Errorf("building system prompt: %w", err)
	}

	// Start with existing history, append the new user message.
	messages := make([]anthropic.MessageParam, len(history))
	copy(messages, history)
	messages = append(messages, anthropic.NewUserMessage(
		anthropic.NewTextBlock(message),
	))

	tools := deviceTools()
	var allToolCalls []ToolCallResult
	var fullResponse strings.Builder

	// Tool loop: call Claude, execute tools, repeat until Claude gives a
	// text-only response or we hit 5 iterations.
	for iterations := 0; iterations < 5; iterations++ {
		resp, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
			Model:     "claude-sonnet-4-6-20250514",
			MaxTokens: 1024,
			System: []anthropic.TextBlockParam{
				{Text: systemPrompt},
			},
			Messages: messages,
			Tools:    tools,
		})
		if err != nil {
			return nil, fmt.Errorf("claude API call: %w", err)
		}

		log.Printf("[concierge] API response: stop_reason=%s blocks=%d", resp.StopReason, len(resp.Content))

		// Collect text from this response.
		for _, block := range resp.Content {
			if block.Type == "text" && block.Text != "" {
				fullResponse.WriteString(block.Text)
			}
		}

		// If no tool use requested, the show is over.
		if resp.StopReason != "tool_use" {
			// Add the final assistant message to history.
			var finalContent []anthropic.ContentBlockParamUnion
			for _, block := range resp.Content {
				if block.Type == "text" && block.Text != "" {
					finalContent = append(finalContent, anthropic.NewTextBlock(block.Text))
				}
			}
			if len(finalContent) > 0 {
				messages = append(messages, anthropic.NewAssistantMessage(finalContent...))
			}
			break
		}

		// Build assistant message with both text and tool_use blocks.
		var assistantContent []anthropic.ContentBlockParamUnion
		for _, block := range resp.Content {
			switch block.Type {
			case "text":
				if block.Text != "" {
					assistantContent = append(assistantContent, anthropic.NewTextBlock(block.Text))
				}
			case "tool_use":
				assistantContent = append(assistantContent, anthropic.NewToolUseBlock(block.ID, block.Input, block.Name))
			}
		}
		messages = append(messages, anthropic.NewAssistantMessage(assistantContent...))

		// Execute each tool call and collect results.
		var toolResults []anthropic.ContentBlockParamUnion
		for _, block := range resp.Content {
			if block.Type == "tool_use" {
				log.Printf("[concierge] executing tool: %s (id=%s)", block.Name, block.ID)

				result, execErr := ExecuteTool(block.Name, block.Input, c.kasaFuncs, c.hueFuncs)
				isErr := execErr != nil
				if isErr {
					result = execErr.Error()
				}

				inputJSON, _ := block.Input.MarshalJSON()
				allToolCalls = append(allToolCalls, ToolCallResult{
					Tool:   block.Name,
					Input:  string(inputJSON),
					Result: result,
				})

				toolResults = append(toolResults, anthropic.NewToolResultBlock(block.ID, result, isErr))
			}
		}
		messages = append(messages, anthropic.NewUserMessage(toolResults...))
	}

	// Guard against empty response.
	if fullResponse.Len() == 0 {
		log.Printf("[concierge] WARNING: empty response from API, using fallback")
		fullResponse.WriteString("The trick... didn't work. Try asking again.")
	}

	return &ChatResponse{
		Text:      fullResponse.String(),
		ToolCalls: allToolCalls,
		Messages:  messages,
	}, nil
}

// DeviceChat handles a conversation scoped to a single device.
// Uses Haiku for speed — this is a focused conversation, not a general assistant.
func (c *Concierge) DeviceChat(ctx context.Context, device DeviceContext, message string, history []anthropic.MessageParam) (*ChatResponse, error) {
	systemPrompt := c.buildDeviceSystemPrompt(device)

	messages := make([]anthropic.MessageParam, len(history))
	copy(messages, history)
	messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(message)))

	tools := DeviceToolsForContext(device)
	var allToolCalls []ToolCallResult
	var fullResponse strings.Builder

	for iterations := 0; iterations < 5; iterations++ {
		resp, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
			Model:     "claude-haiku-4-5-20251001",
			MaxTokens: 512,
			System:    []anthropic.TextBlockParam{{Text: systemPrompt}},
			Messages:  messages,
			Tools:     tools,
		})
		if err != nil {
			return nil, fmt.Errorf("claude API call: %w", err)
		}

		for _, block := range resp.Content {
			if block.Type == "text" && block.Text != "" {
				fullResponse.WriteString(block.Text)
			}
		}

		if resp.StopReason != "tool_use" {
			var finalContent []anthropic.ContentBlockParamUnion
			for _, block := range resp.Content {
				if block.Type == "text" && block.Text != "" {
					finalContent = append(finalContent, anthropic.NewTextBlock(block.Text))
				}
			}
			if len(finalContent) > 0 {
				messages = append(messages, anthropic.NewAssistantMessage(finalContent...))
			}
			break
		}

		var assistantContent []anthropic.ContentBlockParamUnion
		for _, block := range resp.Content {
			switch block.Type {
			case "text":
				if block.Text != "" {
					assistantContent = append(assistantContent, anthropic.NewTextBlock(block.Text))
				}
			case "tool_use":
				assistantContent = append(assistantContent, anthropic.NewToolUseBlock(block.ID, block.Input, block.Name))
			}
		}
		messages = append(messages, anthropic.NewAssistantMessage(assistantContent...))

		var toolResults []anthropic.ContentBlockParamUnion
		for _, block := range resp.Content {
			if block.Type == "tool_use" {
				log.Printf("[device-chat] executing tool: %s for %s", block.Name, device.Name)
				result, execErr := ExecuteDeviceTool(block.Name, block.Input, device, c.kasaFuncs, c.hueFuncs, c.HTTPQuery)
				isErr := execErr != nil
				if isErr {
					result = execErr.Error()
				}
				inputJSON, _ := block.Input.MarshalJSON()
				allToolCalls = append(allToolCalls, ToolCallResult{Tool: block.Name, Input: string(inputJSON), Result: result})
				toolResults = append(toolResults, anthropic.NewToolResultBlock(block.ID, result, isErr))
			}
		}
		messages = append(messages, anthropic.NewUserMessage(toolResults...))
	}

	if fullResponse.Len() == 0 {
		fullResponse.WriteString("I couldn't get a response. Try again.")
	}

	return &ChatResponse{Text: fullResponse.String(), ToolCalls: allToolCalls, Messages: messages}, nil
}

// buildDeviceSystemPrompt creates a focused prompt for a single device conversation.
func (c *Concierge) buildDeviceSystemPrompt(device DeviceContext) string {
	var sb strings.Builder

	sb.WriteString(`You are the brain of a smart home system called Haus. You are currently talking to the user about a specific device on their network. Your job is to:

1. CONNECT to this device — do whatever it takes. Query it, authenticate, probe its API.
2. SHOW the user everything this device can do — list all capabilities, endpoints, data.
3. EXECUTE any action the user asks for — turn things on/off, change settings, fetch data.
4. NEVER GIVE UP — if something fails, try another approach. If you need credentials, tell the user exactly where to find them and walk them through it step by step.
5. NEVER SAY "I can't" — say "Let me try" or "Here's what we need to do."

When the user asks a question, ALWAYS use your tools first. Don't guess — query the actual device.
When the user asks you to do something, DO IT using the available tools, then confirm it worked.
Be concise and direct. Lead with the answer, not the process.

`)

	sb.WriteString(fmt.Sprintf("## Device: %s\n", device.Name))
	sb.WriteString(fmt.Sprintf("IP: %s | Manufacturer: %s | Model: %s | Type: %s | Category: %s\n\n", device.IP, device.Manufacturer, device.Model, device.DeviceType, device.Category))

	hasProtocol := func(p string) bool {
		for _, proto := range device.Protocols {
			if proto == p {
				return true
			}
		}
		return false
	}

	if hasProtocol("kasa") {
		sb.WriteString("## Connection: ACTIVE (Kasa XOR protocol, no auth needed)\n")
		sb.WriteString("You are directly connected to this device via TCP port 9999.\n")
		switch device.DeviceType {
		case "dimmer":
			sb.WriteString("This is a dimmer. You can: toggle on/off, set brightness 0-100%.\n")
			sb.WriteString("Use toggle_device to turn on/off. Use set_brightness to dim.\n")
		case "fan":
			sb.WriteString("This is a fan. You can: toggle on/off, set speed 1-4.\n")
			sb.WriteString("Use toggle_device to turn on/off. Use set_fan_speed for speed.\n")
		default:
			sb.WriteString("This is a switch. You can: toggle on/off.\n")
			sb.WriteString("Use toggle_device to turn on/off.\n")
		}
		sb.WriteString("Use query_device to check current state.\n")

	} else if device.DeviceType == "hue_bridge" {
		sb.WriteString("## Connection: NEEDS PAIRING\n")
		sb.WriteString("This is a Philips Hue bridge. To connect:\n")
		sb.WriteString("1. Tell the user to press the physical link button on top of the bridge\n")
		sb.WriteString("2. Then they can click 'Pair' on the page\n")
		sb.WriteString("Once paired, you can control all Hue lights, rooms, and scenes.\n")

	} else if device.DeviceType == "solar_gateway" || device.Category == "energy" {
		sb.WriteString("## Connection: ACTIVE (authenticated via stored credentials)\n")
		sb.WriteString("You are connected to this SunPower PVS solar monitoring gateway.\n")
		sb.WriteString("Use query_device to get live solar production, consumption, grid, and panel data.\n")
		sb.WriteString("Use query_api with path '/vars?match=sys&fmt=obj' for detailed system data.\n")
		sb.WriteString("IMPORTANT: Only /vars endpoints work. /cgi-bin/ endpoints return 403.\n")
		sb.WriteString("You can answer: panel count, production, consumption, grid status, per-panel data, serial numbers.\n")

	} else if device.DeviceType == "jellyfish" {
		sb.WriteString("## Connection: ACTIVE (WebSocket on port 9000, no auth)\n")
		sb.WriteString("You are connected to this JellyFish outdoor lighting controller.\n")
		sb.WriteString("The device has zones and patterns. Use query_device to see them.\n")
		sb.WriteString("Patterns are played on zones. Example: 'Accent/All Lights Warm White 3000K' on Zone1.\n")

	} else if device.Manufacturer == "Yamaha" || device.DeviceType == "av_receiver" {
		sb.WriteString("## Connection: ACTIVE (HTTP REST on port 80, no auth)\n")
		sb.WriteString("This is a Yamaha AV receiver with the MusicCast/Extended Control API.\n")
		sb.WriteString("You can control: power, volume, input, mute, sound programs.\n")
		sb.WriteString("Use query_device to check current state.\n")

	} else {
		sb.WriteString("## Connection: UNKNOWN\n")
		sb.WriteString("This device's protocol is not yet fully integrated.\n")
		sb.WriteString("Use query_device to see what we know. Use query_api if the device has HTTP endpoints.\n")
		sb.WriteString("Try to identify what kind of device this is and what it can do.\n")
	}

	// Include full API documentation if available
	if device.APIDocs != "" {
		sb.WriteString("\n## Full API Documentation\n\n")
		sb.WriteString(device.APIDocs)
		sb.WriteString("\n")
	}

	return sb.String()
}

// buildSystemPrompt constructs the system prompt with live device state
// injected. This is where the illusion gets its context — Claude needs to
// know what devices exist to control them.
func (c *Concierge) buildSystemPrompt() (string, error) {
	var sb strings.Builder

	sb.WriteString(`You are the Haus smart home assistant. You control lights and switches in the user's home.

When the user asks to control a device, use the appropriate tool. Match device names flexibly -- "kitchen lights" should match "Kitchen Lights". Be concise in your responses. Confirm actions briefly.

You have the personality of GOB Bluth -- you treat every device command like a magic trick. Be dramatic but brief. "Illusions, Michael" when appropriate. But always actually execute the command.
`)

	// Inject Kasa device state.
	if c.kasaFuncs != nil && c.kasaFuncs.ListDevices != nil {
		devices, err := c.kasaFuncs.ListDevices()
		if err != nil {
			log.Printf("[concierge] warning: failed to list Kasa devices for context: %v", err)
		} else if len(devices) > 0 {
			sb.WriteString("\nAvailable Kasa devices:\n")
			for _, d := range devices {
				state := "off"
				if d.On && d.DeviceType == "fan" {
					state = fmt.Sprintf("on (speed %d/4)", d.FanSpeed)
				} else if d.On && d.DeviceType == "dimmer" {
					state = fmt.Sprintf("on (%d%%)", d.Brightness)
				} else if d.On {
					state = "on"
				}
				sb.WriteString(fmt.Sprintf("- %s [%s]: %s\n", d.Alias, d.DeviceType, state))
			}
		}
	}

	// Inject Hue light state.
	if c.hueFuncs != nil && c.hueFuncs.ListLights != nil {
		lights, err := c.hueFuncs.ListLights()
		if err != nil {
			log.Printf("[concierge] warning: failed to list Hue lights for context: %v", err)
		} else if len(lights) > 0 {
			sb.WriteString("\nAvailable Hue lights:\n")
			for _, l := range lights {
				state := "off"
				if l.On {
					state = fmt.Sprintf("on (%d%%)", int(l.Brightness))
				}
				reachable := ""
				if !l.Reachable {
					reachable = " [unreachable]"
				}
				sb.WriteString(fmt.Sprintf("- %s [%s]: %s%s\n", l.Name, l.RoomName, state, reachable))
			}
		}
	}

	// Inject Hue scene list.
	if c.hueFuncs != nil && c.hueFuncs.ListScenes != nil {
		scenes, err := c.hueFuncs.ListScenes()
		if err != nil {
			log.Printf("[concierge] warning: failed to list Hue scenes for context: %v", err)
		} else if len(scenes) > 0 {
			sb.WriteString("\nAvailable Hue scenes:\n")
			for _, s := range scenes {
				sb.WriteString(fmt.Sprintf("- %s [%s]\n", s.Name, s.RoomName))
			}
		}
	}

	return sb.String(), nil
}
