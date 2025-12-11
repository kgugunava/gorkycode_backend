from ollama import Client
import json
import logging

class Model:
    def __init__(self, model_name):

        self.logger = logging.getLogger("MODEL")
        try:
            self.client = Client(host="http://localhost:11434")
        except:
            self.logger.warning("CLIENT INITIALIZATION PROBLEM")
        self.model = model_name

    def generate_prompt(self, user_prompt, locations, query):
        points = [place["title"] + " - описание :" + place["description"][:100] for place in locations['places']]
        return user_prompt.format(
            user_request=query,
            route_json=points
        )

    def request_to_model(self, user_prompt, locations, query):
        final_prompt = self.generate_prompt(user_prompt, locations, query)

        self.logger.debug(f"REQUEST TO MODEL PROMPT : {final_prompt}")
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

        self.logger.debug(f"REQUEST_TO_MODEL : ANSWER : {answer}")

