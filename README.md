# MilkyTeadrop Discord Bot

MilkyTeadrop is a Discord bot designed to provide image generation and Q&A services. It interfaces with a Rust backend server for processing and forwarding questions to a local LLM and also prompts a Python-based AI Service exposing a Rest-API for image generation. It utilizes a custom-written upstream remote UNIX-based fileserver via REST-API to persist and retrieve images.

## Examples: 
![image](https://github.com/CoffeeeAtNight/MilkyTeadrop_Bot/assets/98992091/fe7efe1b-65b4-4d39-b7a2-566216cf5129)

![image](https://github.com/CoffeeeAtNight/MilkyTeadrop_Bot/assets/98992091/666444b1-61af-49d5-98c3-76648d444e9a)

## Features

- **Ask Questions**: Users can ask questions which are processed by a backend server.
- **Image Generation**: Users can request image generation with specific prompts.

## Setup and Installation

### Prerequisites

- Go
- Rust / Cargo
- Python / Pip
- Discord Bot Token
- 'MilkyTeadrop-Fileserver' or any other fileserver (might require changes to code)

### Configuration

1. **Set up the Configuration File**: Create a `config.json` file in the project root with the following content:
  
```
{
"Token": "YOUR_DISCORD_BOT_TOKEN"
}
```

2. **Rust Backend**: Ensure the Rust backend server is running on `127.0.0.1:7878`.

3. **Python REST API**: Make sure the Python REST API is active on `http://127.0.0.1:5000`.
(Note: You might need to donwload the model and specifc pip packages on your own for this to work. Also check if your system is capable to power the models used.)

4. **File Server**: The MilkyTeadrop File Server should be reachable at `http:/{IP}:7676`. 
(Note: This is a custom written open-source standalone fileserver which exposes a REST-API, also avaliable on my github profile [HERE](https://github.com/CoffeeeAtNight/MilkyTeadrop_FileServer))

### Running the Bot

Execute the following commands to start the Discord bot:

```
cd ./scripts
./start.sh && ./start-python.sh
```

## Usage

- **Asking Questions**: Use the `!ask [question]` command in any Discord channel where the bot is present.
- **Requesting Images**: Use the `!img [prompt]` command to generate images based on the provided prompt.

## Notes

- This README is not intended to fully guide you to setup the enviroment on your own machine and requires installation of seperate tools.
- For the Q&A Model, I use ollama/llama2. A docker like system to serve locally hosted LLMs. [llama2](https://ollama.ai/library/llama2)
- The Stable Diffusion model is pulled from hugging face [Stable Diffusion v1.5](https://huggingface.co/runwayml/stable-diffusion-v1-5)
- THE CODE IS NOT INTENDED FOR PRODUCTION USE AND IS ONLY A FUN PROJECT OF MINE. SAME GOES FOR THE [MILKYTEADROP-FILESERVER](https://github.com/CoffeeeAtNight/MilkyTeadrop_FileServer) AS THE SYSTEM EXPOSES APIS THAT ARE VULNERABLE TO CYBER ATTACKS.

## Contributing

Contributions to MilkyTeadrop are welcome!

## Acknowledgements

- Aki, the creator of MilkyTeadrop.
- Mocha, Aki's cat and the unofficial mascot.
  
![image (3)](https://github.com/CoffeeeAtNight/MilkyTeadrop_Bot/assets/98992091/e8b796c2-10d1-4b86-a4a5-e2e6198f732f)


