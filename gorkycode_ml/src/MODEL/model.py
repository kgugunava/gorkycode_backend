from ollama import Client


class Model:
    def __init__(self, model):
        self.model = model
        self.client = Client(host="http://ollama:11434")

    def generate_prompt(self, user_prompt, locations, query):
        user_prompt = user_prompt.format(
            user_request=query,
            route_json=locations
        )
        return user_prompt

    def request_to_model(self, system_prompt, user_prompt, locations, query):
        final_prompt = self.generate_prompt(user_prompt, locations, query)

        response = self.client.chat(
            model=self.model,
            messages=[
                {'role': 'system', 'content': system_prompt},
                {'role': 'user', 'content': final_prompt}
            ]
        )
        locations["description"] = response['message']['content']
