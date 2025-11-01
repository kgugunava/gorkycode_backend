from ollama import Client
import json


class Model:
    def __init__(self, model_name):
        self.client = Client(host="http://localhost:11434")
        self.model = model_name

    def generate_prompt(self, user_prompt, locations, query):
        points = [place["title"] + " - описание :" + place["description"][:100] for place in locations['places']]
        return user_prompt.format(
            user_request=query,
            route_json=points
        )

    def request_to_model(self, user_prompt, locations, query):
       
        final_prompt = self.generate_prompt(user_prompt, locations, query)
        # print(final_prompt)
        response = self.client.chat(
            model=self.model,
            messages=[
                {"role": "system", "content": "Ты — эксперт по культурному и туристическому планированию маршрутов."},
                {"role": "user", "content": final_prompt}
            ],
            options={
                "temperature": 0.6,
                "num_predict": 400 
            }
        )

        answer = response["message"]["content"].strip()
        locations["description"] = str(answer)

