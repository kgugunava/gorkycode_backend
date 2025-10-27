from ollama import chat, ChatResponse


class Model:
    def __init__(self, model):
        self.model = model

    def generate_prompt(self, user_prompt, locations):
        return user_prompt+str(locations)

    def request_to_model(self, system_prompt, user_prompt, locations):
        final_prompt = self.generate_prompt(user_prompt, locations)
        response: ChatResponse = chat(model=self.model, message=[
            {
                'role': 'system',
                'content': system_prompt
            },
            {
                'role': 'user',
                'content': final_prompt
            }
        ])
        return response['message']['content']
