# Input / Output Pairs Dataset Tool or (IOPDT) 🤖

## What is IOPDT?

It is a tool that allows you to create a dataset of input / output pairs or a defined `JSON Schema` from a set of files. It is useful for creating datasets for machine learning models.

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
    max_tokens: 2100 # The maximum tokens to use when calling the API. Max depends on the model and use `0` to disable and use the default value for the model.
    chunk_size: 2100 # The chunk size in which the input will be split to call the API.
    steps: 1 # How many iterations of the same input will be used to generate the output.
    skip: false # Ignore the process when executing the tool
    # The input to use when calling the API, it can be a path to a file or a pattern to match multiple files. NOTE: All files are relative to the configuration file path.
    documents:
      - corpus/*.txt
    output_dir: output/corpus # The output directory to save the results. This directory is relative to the configuration file path.
    # The instruction to the system (Context) to use when calling the API
    system_prompt: |
      You are a smart person that creates questions in 'input' and 'output' pairs from the given document.
    ## User prompt to use when calling the API.
    ## It uses .Document which is the text extracted from the documents splitted by the chunk size.
    user_prompt: |
      Here is the document:
      {{ .Document }}
    skip_json_schema: false # Skip the JSON schema validation when merging the results. It is needed when a model doest not support JSONSchema validation.

    ## If the response is a JSON object, it will use the JSON schema to validate the response and extract the data.
    ## NOTE: This might not be supported by all the models. It uses the `function-calling` https://platform.openai.com/docs/guides/function-calling
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
jq -s '{conversations: add | .conversations}' *.json > merged.json
```

Convert it to `sharegpt` format for for [Axolotl](https://github.com/OpenAccess-AI-Collective/axolotl?tab=readme-ov-file#dataset)

```sh
jq -c '.conversations[] | {"conversations": [{"from": "human", "value": .input}, {"from": "gpt", "value": .output}]}' merged.json > transformed.jsonl
```
