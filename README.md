# Input / Output Pairs Dataset Tool 🤖
```go
 ________  ______   ______   ______   _________  
/_______/\/_____/\ /_____/\ /_____/\ /________/\ 
\__.::._\/\:::_ \ \\:::_ \ \\:::_ \ \\__.::.__\/ 
   \::\ \  \:\ \ \ \\:(_) \ \\:\ \ \ \  \::\ \   
   _\::\ \__\:\ \ \ \\: ___\/ \:\ \ \ \  \::\ \  
  /__\::\__/\\:\_\ \ \\ \ \    \:\/.:| |  \::\ \ 
  \________\/ \_____\/ \_\/     \____/_/   \__\/
```

It is a tool that allows you to infer an AI model (using OpenAI API Format) to get back responses on `input/output pairs` or a defined `JSON Schema` from a set of files. It is useful for creating datasets for machine learning models.

## Command

To execute the tool, you need to use the following command:

```sh
iopdt --config <path-to-config.yaml>
```
> NOTE: All the generated files are relative to the configuration file path.

## .Env

From the configuration file, you can use environment variables to store sensitive information like API keys.

You can set them by creating a `.env` file from the binary file is located or by setting them in the environment.

## Config

Definition of the configuration file.

<details>
<summary>All yaml options (click me)</summary>

```yaml
## The target to use for the process
## The target is a set of configurations for the API
## that is compatible with the OpenAI API
targets:
  - name: openai # The name of the target
    api_url: https://api.openai.com/v1/ # The URL of the API
    api_key_env: OPENAI_API_KEY # The environment variable that contains the API key

## The process to execute during the inference
## it may contain multiple processes, each process
## will be executed in order
processes:
  - name: dialog-with-chat-gpt3 # The name of the process
    model: gpt-3.5-turbo-0125 # The model to use
    target: openai # The target configuration to use when calling the API
    temperature: 0.6 # The temperature to use when calling the API. Use `0` to disable
    # The maximum tokens to use when calling the API. 
    ## Max depends on the model and use `0` to disable and 
    # use the default value for the model.
    max_tokens: 2100 
    chunk_size: 2100 # The chunk size in which the input will be split to call the API. default: 4096
    skip: false # Ignore the process when executing the tool. default: false
    # The input to use when calling the API, 
    # it can be a path to a file or a pattern to match multiple files. 
    ## NOTE: All files are relative to the configuration file path.
    documents:
      - corpus/*.txt
    # The output directory to save the results. 
    # This directory is relative to the configuration file path.
    # default: output
    output_dir: output/corpus
    # The instruction to the system (Context) to use when calling the API
    # Try to be as specific as possible to get the best results.
    # The Tool first try to find all JSON inside ````json ... ```` and then if not, parse the whole text.
    system_prompt: |
      You are a smart person that creates questions in 'input' and 'output' pairs from the given document.
      Expected response format:

      ```json
      {
        "questions": [
          {
            "input": "What is the capital of France?",
            "output": "Paris"
          }
        ]
      }
      ```
    # User prompt to use when calling the API.
    # It uses .Document which is the text extracted
    # from the documents splitted by the chunk size.
    user_prompt: |
      Here is the document:
      {{ .Document }}
    # If the response is a JSON object, it will use the JSON schema 
    # to validate the response and extract the data.
    # If the data is not valid, then saves a txt.
    ## NOTE: Follow the JSON Schema format:
    ## https://json-schema.org/learn/getting-started-step-by-step
    json_schema:
      type: object
      properties:
        questions:
          type: array
          items:
            type: object
            properties:
              input:
                type: string
              output:
                type: string
            required: [input, output]
      required: [questions]
```
</details>

### Other tools (Optional)

[jq](https://jqlang.github.io/jq/) A lightweight and flexible command-line JSON processor.

#### JQ Usage

To merge all resultant `json` files after calling the tool.

```sh
jq -s '{questions: add | .questions}' *.json > merged.json
```

Convert it to `sharegpt` format for for [Axolotl](https://github.com/OpenAccess-AI-Collective/axolotl?tab=readme-ov-file#dataset)

```sh
jq -c '.questions[] | {"conversations": [{"from": "human", "value": .input}, {"from": "gpt", "value": .output}]}' merged.json > transformed.jsonl
```
