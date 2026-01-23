#!/bin/bash
# 06_debug_template.sh
#
# Demonstrates the _debug_render_only flag to see how templates work.
# This shows exactly what prompt is sent to the model after template rendering.

echo "=== Template Debug Demo ==="
echo ""
echo "This shows how Ollama transforms your prompt using model-specific templates."
echo "The _debug_render_only flag returns the rendered template instead of calling the model."
echo ""

MODEL="${1:-gemma3:270m}"

echo "Model: $MODEL"
echo ""

echo "--- Raw prompt ---"
echo "Why is the sky blue?"
echo ""

echo "--- Rendered template ---"
curl -s http://localhost:11434/api/generate -d "{
  \"model\": \"$MODEL\",
  \"prompt\": \"Why is the sky blue?\",
  \"_debug_render_only\": true
}" | jq -r '.response // .error'

echo ""
echo "--- With system prompt ---"
curl -s http://localhost:11434/api/generate -d "{
  \"model\": \"$MODEL\",
  \"prompt\": \"Why is the sky blue?\",
  \"system\": \"You are a helpful physics teacher.\",
  \"_debug_render_only\": true
}" | jq -r '.response // .error'
