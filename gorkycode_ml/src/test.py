import requests
import json

# Данные для отправки
data = {
    "interests": "парки",
    "time_for_route": 400,
    "coordinates": [56.30981, 44.010701]
}

# Отправляем GET запрос
try:
    response = requests.post(
        'http://localhost:5001/route',
        json=data,
        headers={'Content-Type': 'application/json'}
    )

    print(f"Статус код: {response.status_code}")
    print("Ответ от сервера:")
    print(json.dumps(response.json(), indent=2, ensure_ascii=False))

except requests.exceptions.ConnectionError:
    print("Ошибка: Не удалось подключиться к серверу. Убедитесь, что сервер запущен.")
