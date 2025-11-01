import requests
import json

# Данные для отправки
data = {
    "interests": "история, музеи, парки, архитектура",
    "time_for_route": 400,
    "coordinates": [43.99821, 56.308973]
}

# Отправляем GET запрос
try:
    response = requests.get(
        'http://localhost:5001/route',
        json=data,
        headers={'Content-Type': 'application/json'}
    )

    print(f"Статус код: {response.status_code}")
    print("Ответ от сервера:")
    print(json.dumps(response.json(), indent=2, ensure_ascii=False))

except requests.exceptions.ConnectionError:
    print("Ошибка: Не удалось подключиться к серверу. Убедитесь, что сервер запущен.")
