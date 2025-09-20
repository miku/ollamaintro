package main

import (
	"flag"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/ollama/ollama/fs/gguf"
)

func main() {
	var noColor bool
	flag.BoolVar(&noColor, "no-color", false, "disable colored output")
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("path to model required")
	}

	if noColor {
		color.NoColor = true
	}

	modelPath := flag.Arg(0)
	f, err := gguf.Open(modelPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	displayGGUFInfo(f, modelPath)
}

func displayGGUFInfo(f *gguf.File, modelPath string) {
	// Color definitions
	headerColor := color.New(color.FgCyan, color.Bold)
	sectionColor := color.New(color.FgYellow, color.Bold)
	keyColor := color.New(color.FgGreen)
	valueColor := color.New(color.FgWhite)
	dimColor := color.New(color.FgHiBlack)

	// Header
	headerColor.Println("╭─────────────────────────────────────────────────────────────────╮")
	headerColor.Printf("│ %-63s │\n", "GGUF Model Metadata")
	headerColor.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Printf("│ %-15s %-47s │\n", "File:", truncateString(modelPath, 47))
	fmt.Printf("│ %-15s %-47s │\n", "Magic:", f.Magic)
	fmt.Printf("│ %-15s %-47d │\n", "Version:", f.Version)
	fmt.Printf("│ %-15s %-47d │\n", "Tensors:", f.NumTensors())
	headerColor.Println("╰─────────────────────────────────────────────────────────────────╯")
	fmt.Println()

	// Model metadata
	sectionColor.Println("📋 Model Metadata")
	dimColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Collect key-value pairs from iterator
	kvMap := make(map[string]any)
	kvCount := 0
	for _, kv := range f.KeyValues() {
		kvMap[kv.Key] = kv.Value
		kvCount++
	}

	if kvCount > 0 {
		categories := categorizeMetadata(kvMap)

		for category, keys := range categories {
			if len(keys) == 0 {
				continue
			}

			fmt.Printf("\n%s:\n", category)
			maxKeyLen := getMaxKeyLength(keys, kvMap)

			for _, key := range keys {
				value := kvMap[key]
				formattedValue := formatValue(value)

				keyColor.Printf("  %-*s ", maxKeyLen, key)
				dimColor.Print("│ ")
				valueColor.Println(formattedValue)
			}
		}
	} else {
		dimColor.Println("  No metadata found")
	}

	fmt.Println()

	// Tensor information
	sectionColor.Println("🧠 Tensor Information")
	dimColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Collect tensor infos from iterator
	tensorMap := make(map[string]gguf.TensorInfo)
	tensorCount := 0
	for _, tensor := range f.TensorInfos() {
		tensorMap[tensor.Name] = tensor
		tensorCount++
	}

	if tensorCount > 0 {
		displayTensorTable(tensorMap)
	} else {
		dimColor.Println("  No tensor information found")
	}
}

func categorizeMetadata(kvs map[string]any) map[string][]string {
	categories := map[string][]string{
		"General":      {},
		"Tokenizer":    {},
		"Architecture": {},
		"Training":     {},
		"Other":        {},
	}

	for key := range kvs {
		keyLower := strings.ToLower(key)
		switch {
		case strings.Contains(keyLower, "tokenizer") || strings.Contains(keyLower, "token"):
			categories["Tokenizer"] = append(categories["Tokenizer"], key)
		case strings.Contains(keyLower, "arch") || strings.Contains(keyLower, "layer") ||
			strings.Contains(keyLower, "head") || strings.Contains(keyLower, "dim") ||
			strings.Contains(keyLower, "embedding") || strings.Contains(keyLower, "vocab"):
			categories["Architecture"] = append(categories["Architecture"], key)
		case strings.Contains(keyLower, "train") || strings.Contains(keyLower, "learning") ||
			strings.Contains(keyLower, "epoch"):
			categories["Training"] = append(categories["Training"], key)
		case strings.Contains(keyLower, "name") || strings.Contains(keyLower, "version") ||
			strings.Contains(keyLower, "author") || strings.Contains(keyLower, "description") ||
			strings.Contains(keyLower, "license") || strings.Contains(keyLower, "url"):
			categories["General"] = append(categories["General"], key)
		default:
			categories["Other"] = append(categories["Other"], key)
		}
	}

	// Sort keys within each category
	for category := range categories {
		sort.Strings(categories[category])
	}

	return categories
}

func displayTensorTable(tensorInfos map[string]gguf.TensorInfo) {
	if len(tensorInfos) == 0 {
		return
	}

	// Sort tensor names
	var names []string
	for name := range tensorInfos {
		names = append(names, name)
	}
	sort.Strings(names)

	// Calculate column widths
	maxNameLen := 20
	maxShapeLen := 15
	maxTypeLen := 10

	for name, info := range tensorInfos {
		if len(name) > maxNameLen {
			maxNameLen = len(name)
		}
		shapeStr := formatShape(info.Shape)
		if len(shapeStr) > maxShapeLen {
			maxShapeLen = len(shapeStr)
		}
		typeStr := fmt.Sprintf("%v", info.Type)
		if len(typeStr) > maxTypeLen {
			maxTypeLen = len(typeStr)
		}
	}

	// Limit max widths for readability
	if maxNameLen > 50 {
		maxNameLen = 50
	}
	if maxShapeLen > 25 {
		maxShapeLen = 25
	}

	// Table header
	headerColor := color.New(color.FgHiBlue, color.Bold)
	headerColor.Printf("  %-*s │ %-*s │ %-*s │ %s\n",
		maxNameLen, "Name",
		maxShapeLen, "Shape",
		maxTypeLen, "Type",
		"Offset")

	dimColor := color.New(color.FgHiBlack)
	dimColor.Printf("  %s─┼─%s─┼─%s─┼─%s\n",
		strings.Repeat("─", maxNameLen),
		strings.Repeat("─", maxShapeLen),
		strings.Repeat("─", maxTypeLen),
		strings.Repeat("─", 12))

	// Group tensors by layer/type for better organization
	layerGroups := groupTensorsByLayer(names)

	keyColor := color.New(color.FgGreen)
	valueColor := color.New(color.FgWhite)

	for i, group := range layerGroups {
		if i > 0 {
			fmt.Println() // Add spacing between groups
		}

		for _, name := range group {
			info := tensorInfos[name]
			displayName := truncateString(name, maxNameLen)
			shapeStr := truncateString(formatShape(info.Shape), maxShapeLen)
			typeStr := fmt.Sprintf("%v", info.Type)

			keyColor.Printf("  %-*s ", maxNameLen, displayName)
			dimColor.Print("│ ")
			valueColor.Printf("%-*s ", maxShapeLen, shapeStr)
			dimColor.Print("│ ")
			valueColor.Printf("%-*s ", maxTypeLen, typeStr)
			dimColor.Print("│ ")
			valueColor.Printf("%d", info.Offset)
			fmt.Println()
		}
	}
}

func groupTensorsByLayer(names []string) [][]string {
	groups := make(map[string][]string)

	for _, name := range names {
		// Extract layer prefix (e.g., "blk.0", "output", "token_embd")
		parts := strings.Split(name, ".")
		groupKey := parts[0]
		if len(parts) > 1 && (parts[0] == "blk" || parts[0] == "layers") {
			if len(parts) > 2 {
				groupKey = strings.Join(parts[:2], ".")
			}
		}

		if _, exists := groups[groupKey]; !exists {
			groups[groupKey] = []string{}
		}
		groups[groupKey] = append(groups[groupKey], name)
	}

	// Define architectural order for common layer types
	layerOrder := []string{
		"token_embd", // Input embeddings
		"pos_embd",   // Positional embeddings
		"input_norm", // Input normalization
		"embed_norm", // Embedding normalization
	}

	// Add transformer blocks in numerical order
	var blockKeys []string
	for key := range groups {
		if strings.HasPrefix(key, "blk.") || strings.HasPrefix(key, "layers.") {
			blockKeys = append(blockKeys, key)
		}
	}
	sort.Slice(blockKeys, func(i, j int) bool {
		return naturalLess(blockKeys[i], blockKeys[j])
	})
	layerOrder = append(layerOrder, blockKeys...)

	// Add output layers in architectural order
	outputOrder := []string{
		"output_norm", // Final normalization
		"norm",        // Alternative final norm name
		"ln_f",        // GPT-style final layer norm
		"output",      // Final output projection
		"lm_head",     // Alternative output head name
		"head",        // Generic head name
	}
	layerOrder = append(layerOrder, outputOrder...)

	// Collect results in architectural order
	var result [][]string
	usedKeys := make(map[string]bool)

	for _, key := range layerOrder {
		if tensors, exists := groups[key]; exists {
			// Sort tensors within each group
			sort.Strings(tensors)
			result = append(result, tensors)
			usedKeys[key] = true
		}
	}

	// Add any remaining groups that weren't in our predefined order
	var remainingKeys []string
	for key := range groups {
		if !usedKeys[key] {
			remainingKeys = append(remainingKeys, key)
		}
	}
	sort.Slice(remainingKeys, func(i, j int) bool {
		return naturalLess(remainingKeys[i], remainingKeys[j])
	})

	for _, key := range remainingKeys {
		sort.Strings(groups[key])
		result = append(result, groups[key])
	}

	return result
}

// naturalLess compares strings with natural number ordering
func naturalLess(a, b string) bool {
	// Handle simple cases first
	if a == b {
		return false
	}

	// Extract parts and compare
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	minLen := len(aParts)
	if len(bParts) < minLen {
		minLen = len(bParts)
	}

	for i := 0; i < minLen; i++ {
		aPart := aParts[i]
		bPart := bParts[i]

		// Try to parse as numbers
		aNum, aErr := strconv.Atoi(aPart)
		bNum, bErr := strconv.Atoi(bPart)

		if aErr == nil && bErr == nil {
			// Both are numbers, compare numerically
			if aNum != bNum {
				return aNum < bNum
			}
		} else {
			// At least one is not a number, compare as strings
			if aPart != bPart {
				return aPart < bPart
			}
		}
	}

	// If all compared parts are equal, shorter string comes first
	return len(aParts) < len(bParts)
}

func formatShape(dims []uint64) string {
	if len(dims) == 0 {
		return "[]"
	}

	strs := make([]string, len(dims))
	for i, dim := range dims {
		strs[i] = strconv.FormatUint(dim, 10)
	}
	return "[" + strings.Join(strs, ", ") + "]"
}

func formatValue(value any) string {
	if value == nil {
		return "<nil>"
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			return "[]"
		}
		if v.Len() > 10 {
			return fmt.Sprintf("[%d items]", v.Len())
		}
		// For small arrays, show the values
		var items []string
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			items = append(items, fmt.Sprintf("%v", item))
		}
		result := "[" + strings.Join(items, ", ") + "]"
		if len(result) > 60 {
			return fmt.Sprintf("[%d items]", v.Len())
		}
		return result
	case reflect.String:
		str := value.(string)
		if len(str) > 60 {
			return str[:57] + "..."
		}
		return str
	default:
		result := fmt.Sprintf("%v", value)
		if len(result) > 60 {
			return result[:57] + "..."
		}
		return result
	}
}

func getMaxKeyLength(keys []string, kvs map[string]any) int {
	maxLen := 0
	for _, key := range keys {
		if len(key) > maxLen {
			maxLen = len(key)
		}
	}
	if maxLen > 30 {
		maxLen = 30
	}
	return maxLen
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
