package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	inputPathJSON  = "docs/swagger.json"
	outputPathJSON = "docs/swagger.json"
	outputPathYAML = "docs/swagger.yaml"
	outputPathGo   = "docs/docs.go"
)

const (
	envServerURLs = "SWAGGER_SERVER_URLS"
	envHost       = "SWAGGER_HOST"
	envBasePath   = "SWAGGER_BASE_PATH"
)

const docsTemplate = `package docs

import (
	_ "embed"
	"sync"

	"github.com/swaggo/swag"
)

//go:embed swagger.json
var openAPIDoc string

var registerOnce sync.Once

type embeddedSwagger struct {
	doc string
}

func (s *embeddedSwagger) ReadDoc() string {
	return s.doc
}

func init() {
	registerOnce.Do(func() {
		swag.Register(swag.Name, &embeddedSwagger{doc: openAPIDoc})
	})
}
`

func main() {
	loadEnvFile(".env")

	raw, err := os.ReadFile(inputPathJSON)
	if err != nil {
		log.Fatalf("read swagger json: %v", err)
	}

	var spec2 map[string]interface{}
	if err := json.Unmarshal(raw, &spec2); err != nil {
		log.Fatalf("parse swagger json: %v", err)
	}

	spec3 := convertSpec(spec2)

	data, err := json.MarshalIndent(spec3, "", "  ")
	if err != nil {
		log.Fatalf("encode openapi json: %v", err)
	}
	if err := os.WriteFile(outputPathJSON, append(data, '\n'), 0o644); err != nil {
		log.Fatalf("write openapi json: %v", err)
	}

	yamlData, err := yaml.Marshal(spec3)
	if err != nil {
		log.Fatalf("encode openapi yaml: %v", err)
	}
	if err := os.WriteFile(outputPathYAML, yamlData, 0o644); err != nil {
		log.Fatalf("write openapi yaml: %v", err)
	}

	if err := os.WriteFile(outputPathGo, []byte(docsTemplate), 0o644); err != nil {
		log.Fatalf("write docs.go: %v", err)
	}

	fmt.Println("OpenAPI 3 documentation generated successfully.")
}

func convertSpec(spec map[string]interface{}) map[string]interface{} {
	info := spec["info"]
	tags := spec["tags"]
	externalDocs := spec["externalDocs"]

	host, _ := spec["host"].(string)
	basePath, _ := spec["basePath"].(string)

	if envHostOverride := strings.TrimSpace(os.Getenv(envHost)); envHostOverride != "" {
		host = envHostOverride
	}

	if envBasePathOverride := strings.TrimSpace(os.Getenv(envBasePath)); envBasePathOverride != "" {
		basePath = envBasePathOverride
	}

	basePath = standardizeBasePath(basePath)
	host = strings.TrimSuffix(strings.TrimSpace(host), "/")

	servers := buildServers(host, basePath, envServerList())

	consumesDefault := toStringSlice(spec["consumes"])
	if len(consumesDefault) == 0 {
		consumesDefault = []string{"application/json"}
	}
	producesDefault := toStringSlice(spec["produces"])
	if len(producesDefault) == 0 {
		producesDefault = []string{"application/json"}
	}

	paths, _ := spec["paths"].(map[string]interface{})
	paths3 := convertPaths(paths, consumesDefault, producesDefault)

	definitions, _ := spec["definitions"].(map[string]interface{})
	securityDefinitions, _ := spec["securityDefinitions"].(map[string]interface{})

	components := map[string]interface{}{}
	if len(definitions) > 0 {
		fixRefs(definitions)
		components["schemas"] = definitions
	}
	if len(securityDefinitions) > 0 {
		components["securitySchemes"] = securityDefinitions
	}

	spec3 := map[string]interface{}{
		"openapi":    "3.0.3",
		"info":       info,
		"servers":    servers,
		"paths":      paths3,
		"components": components,
	}

	if tags != nil {
		spec3["tags"] = tags
	}

	if externalDocs != nil {
		spec3["externalDocs"] = externalDocs
	}

	fixRefs(spec3)
	return spec3
}

func envServerList() []string {
	raw := os.Getenv(envServerURLs)
	if raw == "" {
		return nil
	}
	return splitAndTrim(raw)
}

func splitAndTrim(input string) []string {
	fields := strings.FieldsFunc(input, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r'
	})
	result := make([]string, 0, len(fields))
	for _, field := range fields {
		trimmed := strings.TrimSpace(field)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func standardizeBasePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" || path == "/" {
		return ""
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return strings.TrimRight(path, "/")
}

func loadEnvFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		} else if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
			value = strings.Trim(value, "'")
		}

		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}

func buildServers(host, basePath string, extras []string) []map[string]interface{} {
	urls := map[string]struct{}{}

	if host != "" {
		urls[fmt.Sprintf("https://%s%s", host, basePath)] = struct{}{}
		urls[fmt.Sprintf("http://%s%s", host, basePath)] = struct{}{}
	}

	if basePath != "" {
		urls[basePath] = struct{}{}
	}

	// add sensible development default
	urls[fmt.Sprintf("http://localhost:8080%s", basePath)] = struct{}{}

	for _, raw := range extras {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		raw = strings.TrimSuffix(raw, "/")
		urls[raw] = struct{}{}
		if basePath != "" && !strings.HasSuffix(raw, basePath) {
			urls[raw+basePath] = struct{}{}
		}
	}

	keys := make([]string, 0, len(urls))
	for k := range urls {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	servers := make([]map[string]interface{}, 0, len(keys))
	for _, url := range keys {
		servers = append(servers, map[string]interface{}{"url": url})
	}
	return servers
}

func convertPaths(paths map[string]interface{}, consumesDefault, producesDefault []string) map[string]interface{} {
	if len(paths) == 0 {
		return map[string]interface{}{}
	}

	result := make(map[string]interface{}, len(paths))
	for route, value := range paths {
		pathItem, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		newPathItem := make(map[string]interface{})
		for method, opValue := range pathItem {
			opMap, ok := opValue.(map[string]interface{})
			if !ok {
				newPathItem[method] = opValue
				continue
			}

			operation := make(map[string]interface{})
			for k, v := range opMap {
				switch k {
				case "consumes", "produces", "parameters":
					continue
				default:
					operation[k] = v
				}
			}

			consumes := toStringSlice(opMap["consumes"])
			if len(consumes) == 0 {
				consumes = consumesDefault
			}
			produces := toStringSlice(opMap["produces"])
			if len(produces) == 0 {
				produces = producesDefault
			}

			newParams, requestBody := convertParameters(opMap["parameters"], consumes)
			if len(newParams) > 0 {
				operation["parameters"] = newParams
			}
			if requestBody != nil {
				operation["requestBody"] = requestBody
			}

			if responses, ok := opMap["responses"].(map[string]interface{}); ok {
				operation["responses"] = convertResponses(responses, produces)
			}

			fixRefs(operation)
			newPathItem[method] = operation
		}
		result[route] = newPathItem
	}
	return result
}

func convertParameters(value interface{}, consumes []string) ([]interface{}, map[string]interface{}) {
	paramsSlice, ok := value.([]interface{})
	if !ok || len(paramsSlice) == 0 {
		return nil, nil
	}

	var parameters []interface{}
	var requestBody map[string]interface{}

	for _, item := range paramsSlice {
		param, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		switch paramIn, _ := param["in"].(string); paramIn {
		case "body":
			schema, _ := param["schema"].(map[string]interface{})
			if schema == nil {
				schema = map[string]interface{}{}
			}
			fixRefs(schema)
			description, _ := param["description"].(string)
			required, _ := param["required"].(bool)

			content := map[string]interface{}{}
			for _, mime := range consumes {
				content[mime] = map[string]interface{}{"schema": schema}
			}

			requestBody = map[string]interface{}{
				"description": description,
				"required":    required,
				"content":     content,
			}
		case "formData":
			name, _ := param["name"].(string)
			if name == "" {
				continue
			}

			schema := map[string]interface{}{}
			if typ, _ := param["type"].(string); typ != "" {
				if typ == "file" {
					schema["type"] = "string"
					schema["format"] = "binary"
				} else {
					schema["type"] = typ
				}
			}
			if items, ok := param["items"].(map[string]interface{}); ok {
				schema["items"] = items
			}
			if enumVals, ok := param["enum"]; ok {
				schema["enum"] = enumVals
			}

			description, _ := param["description"].(string)
			required, _ := param["required"].(bool)

			req := map[string]interface{}{
				"description": description,
				"required":    required,
				"content": map[string]interface{}{
					"multipart/form-data": map[string]interface{}{
						"schema": map[string]interface{}{
							"type":       "object",
							"properties": map[string]interface{}{name: schema},
						},
					},
				},
			}
			if required {
				reqSchema := req["content"].(map[string]interface{})["multipart/form-data"].(map[string]interface{})["schema"].(map[string]interface{})
				reqSchema["required"] = []string{name}
			}
			requestBody = req
		default:
			fixRefs(param)
			parameters = append(parameters, param)
		}
	}

	return parameters, requestBody
}

func convertResponses(responses map[string]interface{}, produces []string) map[string]interface{} {
	result := make(map[string]interface{}, len(responses))
	for code, respVal := range responses {
		response, ok := respVal.(map[string]interface{})
		if !ok {
			result[code] = respVal
			continue
		}

		schema, _ := response["schema"].(map[string]interface{})
		examples, _ := response["examples"].(map[string]interface{})
		delete(response, "schema")
		delete(response, "examples")

		if schema != nil {
			fixRefs(schema)
			content := map[string]interface{}{}
			for _, mime := range produces {
				media := map[string]interface{}{"schema": schema}
				if examples != nil {
					if exampleVal, ok := examples[mime]; ok {
						media["example"] = exampleVal
					}
				}
				content[mime] = media
			}
			response["content"] = content
		} else if examples != nil && len(examples) > 0 {
			content := map[string]interface{}{}
			for mime, exampleVal := range examples {
				content[toString(mime)] = map[string]interface{}{"example": exampleVal}
			}
			if len(content) > 0 {
				response["content"] = content
			}
		}

		result[code] = response
	}
	return result
}

func fixRefs(value interface{}) {
	switch v := value.(type) {
	case map[string]interface{}:
		for key, val := range v {
			if key == "$ref" {
				if refStr, ok := val.(string); ok && strings.HasPrefix(refStr, "#/definitions/") {
					v[key] = strings.Replace(refStr, "#/definitions/", "#/components/schemas/", 1)
				}
				continue
			}
			fixRefs(val)
		}
	case []interface{}:
		for _, item := range v {
			fixRefs(item)
		}
	}
}

func toStringSlice(value interface{}) []string {
	switch vals := value.(type) {
	case []interface{}:
		result := make([]string, 0, len(vals))
		for _, item := range vals {
			if s := toString(item); s != "" {
				result = append(result, s)
			}
		}
		return result
	case []string:
		return vals
	default:
		return nil
	}
}

func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
