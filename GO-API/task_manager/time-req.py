import time
import requests

# URL da requisição
URL = "http://localhost:8080/users"

# Tempo de espera (em segundos)
DELAY_SECONDS = 10

print(f"A aguardar {DELAY_SECONDS} segundos...")
time.sleep(DELAY_SECONDS)

try:
    response = requests.get(URL, timeout=10)

    print("Requisição concluída!")
    print("Status Code:", response.status_code)
    print("Resposta:")
    print(response.text)

except requests.RequestException as err:
    print("Erro ao realizar a requisição:")
    print(err)