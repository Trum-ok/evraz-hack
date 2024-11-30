from decouple import config
from anthropic import Anthropic

api_key = config("API_KEY")
client = Anthropic(api_key=api_key)

msg = client.messages.create(
    model="claude-3-5-sonnet-20241022",
    max_tokens=1000,
    temperature=0,
    system="ты можешь делаешь мне code-review C#",
    messages=[
        {
            "role": "user",
            "content": [
                {
                    "type": "text",
                    "text": "какой-то C# код"
                }
            ]
        }
    ]
)

print(msg.to_json)
